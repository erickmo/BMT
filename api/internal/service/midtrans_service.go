package service

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bmt-saas/api/internal/domain/payment"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

// MidtransService menangani webhook dan cek status pending dari Midtrans.
type MidtransService struct {
	paymentRepo       payment.Repository
	settingsResolver  *settings.Resolver
	platformServerKey string // fallback jika BMT tidak punya server key
}

func NewMidtransService(
	paymentRepo payment.Repository,
	settingsResolver *settings.Resolver,
	platformServerKey string,
) *MidtransService {
	return &MidtransService{
		paymentRepo:       paymentRepo,
		settingsResolver:  settingsResolver,
		platformServerKey: platformServerKey,
	}
}

// MidtransNotification adalah body webhook dari Midtrans.
type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	SettlementTime    string `json:"settlement_time"`
	FraudStatus       string `json:"fraud_status"`
}

// HandleWebhook memverifikasi signature SHA512 dan memperbarui status payment.
// Verifikasi: SHA512(order_id + status_code + gross_amount + server_key) == signature_key
// Jika SETTLEMENT atau CAPTURE: UpdateStatus ke SETTLEMENT + set settledAt
// Idempotent: jika sudah SETTLEMENT, return nil tanpa re-proses
func (s *MidtransService) HandleWebhook(ctx context.Context, body []byte) error {
	var notif MidtransNotification
	if err := json.Unmarshal(body, &notif); err != nil {
		return fmt.Errorf("gagal parse body webhook Midtrans: %w", err)
	}

	existing, err := s.paymentRepo.GetByOrderID(ctx, notif.OrderID)
	if err != nil {
		return fmt.Errorf("payment dengan order_id %s tidak ditemukan: %w", notif.OrderID, err)
	}

	// Idempotent: sudah SETTLEMENT, tidak perlu re-proses
	if existing.Status == payment.StatusSettlement {
		return nil
	}

	serverKey := s.GetServerKey(ctx, existing.BMTID)

	// Verifikasi signature: SHA512(order_id + status_code + gross_amount + server_key)
	raw := notif.OrderID + notif.StatusCode + notif.GrossAmount + serverKey
	hash := sha512.Sum512([]byte(raw))
	expectedSig := fmt.Sprintf("%x", hash)

	if expectedSig != notif.SignatureKey {
		return payment.ErrSignatureInvalid
	}

	switch notif.TransactionStatus {
	case "settlement", "capture":
		var settledAt *time.Time
		if notif.SettlementTime != "" {
			t, err := time.Parse("2006-01-02 15:04:05", notif.SettlementTime)
			if err == nil {
				settledAt = &t
			}
		}
		if settledAt == nil {
			now := time.Now()
			settledAt = &now
		}

		if err := s.paymentRepo.UpdateStatus(ctx, existing.ID, payment.StatusSettlement, settledAt, string(body)); err != nil {
			return fmt.Errorf("gagal update status payment ke SETTLEMENT: %w", err)
		}

	case "deny":
		if err := s.paymentRepo.UpdateStatus(ctx, existing.ID, payment.StatusDeny, nil, string(body)); err != nil {
			return fmt.Errorf("gagal update status payment ke DENY: %w", err)
		}

	case "cancel":
		if err := s.paymentRepo.UpdateStatus(ctx, existing.ID, payment.StatusCancel, nil, string(body)); err != nil {
			return fmt.Errorf("gagal update status payment ke CANCEL: %w", err)
		}

	case "expire":
		if err := s.paymentRepo.UpdateStatus(ctx, existing.ID, payment.StatusExpire, nil, string(body)); err != nil {
			return fmt.Errorf("gagal update status payment ke EXPIRE: %w", err)
		}
	}

	return nil
}

// GetServerKey mengambil server key untuk BMT tertentu dari settings, fallback ke platform key.
func (s *MidtransService) GetServerKey(ctx context.Context, bmtID uuid.UUID) string {
	key := s.settingsResolver.Resolve(ctx, bmtID, uuid.Nil, "midtrans.server_key")
	if key == "" {
		return s.platformServerKey
	}
	return key
}

// midtransStatusResponse adalah partial response dari Midtrans status API.
type midtransStatusResponse struct {
	TransactionStatus string `json:"transaction_status"`
	StatusCode        string `json:"status_code"`
	OrderID           string `json:"order_id"`
	FraudStatus       string `json:"fraud_status"`
	SettlementTime    string `json:"settlement_time"`
}

// CekPending memverifikasi status transaksi PENDING > 30 menit ke Midtrans API.
// Menggunakan HTTP GET ke https://api.midtrans.com/v2/{order_id}/status
// Basic Auth: base64(serverKey + ":")
func (s *MidtransService) CekPending(ctx context.Context) error {
	pendingList, err := s.paymentRepo.ListPending(ctx, 30*time.Minute)
	if err != nil {
		return fmt.Errorf("gagal ambil daftar payment pending: %w", err)
	}

	for _, p := range pendingList {
		serverKey := s.GetServerKey(ctx, p.BMTID)
		isSandbox := s.settingsResolver.ResolveBool(ctx, p.BMTID, uuid.Nil, "midtrans.sandbox", true)

		baseURL := "https://api.midtrans.com/v2"
		if isSandbox {
			baseURL = "https://api.sandbox.midtrans.com/v2"
		}

		statusURL := fmt.Sprintf("%s/%s/status", baseURL, p.OrderID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
		if err != nil {
			return fmt.Errorf("gagal buat request ke Midtrans untuk order %s: %w", p.OrderID, err)
		}

		credentials := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
		req.Header.Set("Authorization", "Basic "+credentials)
		req.Header.Set("Accept", "application/json")

		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		if err != nil {
			// Lanjut ke payment berikutnya jika satu gagal
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var statusResp midtransStatusResponse
		if err := json.Unmarshal(respBody, &statusResp); err != nil {
			continue
		}

		switch statusResp.TransactionStatus {
		case "settlement", "capture":
			var settledAt *time.Time
			if statusResp.SettlementTime != "" {
				t, err := time.Parse("2006-01-02 15:04:05", statusResp.SettlementTime)
				if err == nil {
					settledAt = &t
				}
			}
			if settledAt == nil {
				now := time.Now()
				settledAt = &now
			}
			_ = s.paymentRepo.UpdateStatus(ctx, p.ID, payment.StatusSettlement, settledAt, string(respBody))

		case "deny":
			_ = s.paymentRepo.UpdateStatus(ctx, p.ID, payment.StatusDeny, nil, string(respBody))

		case "cancel":
			_ = s.paymentRepo.UpdateStatus(ctx, p.ID, payment.StatusCancel, nil, string(respBody))

		case "expire":
			_ = s.paymentRepo.UpdateStatus(ctx, p.ID, payment.StatusExpire, nil, string(respBody))
		}
	}

	return nil
}
