package payment

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPaymentNotFound       = errors.New("pembayaran Midtrans tidak ditemukan")
	ErrSignatureInvalid      = errors.New("signature Midtrans tidak valid")
	ErrUsageLogNotFound      = errors.New("usage log tidak ditemukan")
)

type StatusMidtrans string

const (
	StatusPending    StatusMidtrans = "PENDING"
	StatusSettlement StatusMidtrans = "SETTLEMENT"
	StatusCapture    StatusMidtrans = "CAPTURE"
	StatusDeny       StatusMidtrans = "DENY"
	StatusCancel     StatusMidtrans = "CANCEL"
	StatusExpire     StatusMidtrans = "EXPIRE"
	StatusRefund     StatusMidtrans = "REFUND"
)

type JenisUsage string

const (
	UsageTransaksiRekening JenisUsage = "TRANSAKSI_REKENING"
	UsageTransaksiNFC      JenisUsage = "TRANSAKSI_NFC"
	UsagePesananOPOP       JenisUsage = "PESANAN_OPOP"
	UsagePesananEcommerce  JenisUsage = "PESANAN_ECOMMERCE"
	UsageTransaksiDonasi   JenisUsage = "TRANSAKSI_DONASI"
)

// MidtransPayment records a payment transaction via Midtrans gateway
type MidtransPayment struct {
	ID              uuid.UUID      `json:"id"`
	BMTID           uuid.UUID      `json:"bmt_id"`
	OrderID         string         `json:"order_id"`      // sent to Midtrans
	ReferensiID     uuid.UUID      `json:"referensi_id"`  // pesanan_id / donasi_id / etc.
	ReferensiTipe   string         `json:"referensi_tipe"` // PESANAN | DONASI | SPP
	Nominal         int64          `json:"nominal"`
	Status          StatusMidtrans `json:"status"`
	MetodeBayar     string         `json:"metode_bayar,omitempty"` // gopay, qris, va_bni, etc.
	SnapToken       string         `json:"snap_token,omitempty"`
	SnapURL         string         `json:"snap_url,omitempty"`
	MidtransResponse string        `json:"midtrans_response,omitempty"` // raw JSON response
	SettledAt       *time.Time     `json:"settled_at,omitempty"`
	ExpiredAt       *time.Time     `json:"expired_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// UsageLog records every billable usage event for platform revenue tracking
type UsageLog struct {
	ID            uuid.UUID  `json:"id"`
	BMTID         uuid.UUID  `json:"bmt_id"`
	CabangID      *uuid.UUID `json:"cabang_id,omitempty"`
	Jenis         JenisUsage `json:"jenis"`
	ReferensiID   *uuid.UUID `json:"referensi_id,omitempty"`
	Nominal       int64      `json:"nominal"`        // transaction value
	BiayaAdmin    int64      `json:"biaya_admin"`    // platform fee
	Periode       string     `json:"periode"`        // "2025-01"
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateMidtransInput struct {
	BMTID         uuid.UUID
	OrderID       string
	ReferensiID   uuid.UUID
	ReferensiTipe string
	Nominal       int64
}

type CreateUsageLogInput struct {
	BMTID       uuid.UUID
	CabangID    *uuid.UUID
	Jenis       JenisUsage
	ReferensiID *uuid.UUID
	Nominal     int64
	BiayaAdmin  int64
	Periode     string
}

type ListPaymentFilter struct {
	BMTID         *uuid.UUID
	ReferensiID   *uuid.UUID
	ReferensiTipe string
	Status        *StatusMidtrans
	Page          int
	PerPage       int
}

type Repository interface {
	// Midtrans Payment
	Create(ctx context.Context, p *MidtransPayment) error
	GetByID(ctx context.Context, id uuid.UUID) (*MidtransPayment, error)
	GetByOrderID(ctx context.Context, orderID string) (*MidtransPayment, error)
	GetByReferensi(ctx context.Context, referensiID uuid.UUID) (*MidtransPayment, error)
	List(ctx context.Context, filter ListPaymentFilter) ([]*MidtransPayment, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusMidtrans, settledAt *time.Time, rawResponse string) error
	UpdateSnapToken(ctx context.Context, id uuid.UUID, token, url string, expiredAt time.Time) error
	ListPending(ctx context.Context, olderThan time.Duration) ([]*MidtransPayment, error)

	// Usage Log
	CreateUsageLog(ctx context.Context, l *UsageLog) error
	SumBiayaAdminByPeriode(ctx context.Context, bmtID uuid.UUID, periode string) (int64, error)
	ListUsageLog(ctx context.Context, bmtID uuid.UUID, periode string, page, perPage int) ([]*UsageLog, int64, error)
}

func NewMidtransPayment(input CreateMidtransInput) (*MidtransPayment, error) {
	if input.OrderID == "" {
		return nil, errors.New("order ID Midtrans wajib diisi")
	}
	if input.Nominal <= 0 {
		return nil, errors.New("nominal pembayaran harus lebih dari 0")
	}
	if input.ReferensiTipe == "" {
		return nil, errors.New("referensi tipe wajib diisi")
	}
	now := time.Now()
	return &MidtransPayment{
		ID:            uuid.New(),
		BMTID:         input.BMTID,
		OrderID:       input.OrderID,
		ReferensiID:   input.ReferensiID,
		ReferensiTipe: input.ReferensiTipe,
		Nominal:       input.Nominal,
		Status:        StatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func NewUsageLog(input CreateUsageLogInput) (*UsageLog, error) {
	if input.Periode == "" {
		return nil, errors.New("periode usage log wajib diisi")
	}
	return &UsageLog{
		ID:          uuid.New(),
		BMTID:       input.BMTID,
		CabangID:    input.CabangID,
		Jenis:       input.Jenis,
		ReferensiID: input.ReferensiID,
		Nominal:     input.Nominal,
		BiayaAdmin:  input.BiayaAdmin,
		Periode:     input.Periode,
		CreatedAt:   time.Now(),
	}, nil
}
