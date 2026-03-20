package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/payment"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// ── MidtransPayment ───────────────────────────────────────────────────────────

func (r *PaymentRepository) Create(ctx context.Context, p *payment.MidtransPayment) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO midtrans_payment (id, bmt_id, order_id, referensi_id, referensi_tipe, nominal,
		status, metode_bayar, snap_token, snap_url, midtrans_response, settled_at, expired_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, p.ID, p.BMTID, p.OrderID, p.ReferensiID, p.ReferensiTipe, p.Nominal,
		p.Status, p.MetodeBayar, p.SnapToken, p.SnapURL, p.MidtransResponse,
		p.SettledAt, p.ExpiredAt, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create midtrans payment: %w", err)
	}
	return nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*payment.MidtransPayment, error) {
	return r.scanPayment(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, order_id, referensi_id, referensi_tipe, nominal, status, metode_bayar,
		snap_token, snap_url, midtrans_response, settled_at, expired_at, created_at, updated_at
		FROM midtrans_payment WHERE id = $1
	`, id))
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*payment.MidtransPayment, error) {
	return r.scanPayment(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, order_id, referensi_id, referensi_tipe, nominal, status, metode_bayar,
		snap_token, snap_url, midtrans_response, settled_at, expired_at, created_at, updated_at
		FROM midtrans_payment WHERE order_id = $1
	`, orderID))
}

func (r *PaymentRepository) GetByReferensi(ctx context.Context, referensiID uuid.UUID) (*payment.MidtransPayment, error) {
	return r.scanPayment(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, order_id, referensi_id, referensi_tipe, nominal, status, metode_bayar,
		snap_token, snap_url, midtrans_response, settled_at, expired_at, created_at, updated_at
		FROM midtrans_payment WHERE referensi_id = $1
		ORDER BY created_at DESC LIMIT 1
	`, referensiID))
}

func (r *PaymentRepository) List(ctx context.Context, filter payment.ListPaymentFilter) ([]*payment.MidtransPayment, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage

	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, order_id, referensi_id, referensi_tipe, nominal, status, metode_bayar,
		snap_token, snap_url, midtrans_response, settled_at, expired_at, created_at, updated_at
		FROM midtrans_payment
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR referensi_id = $2)
		  AND ($3::text = '' OR referensi_tipe = $3)
		  AND ($4::text IS NULL OR status = $4)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`, filter.BMTID, filter.ReferensiID, filter.ReferensiTipe, filter.Status, filter.PerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list midtrans payment: %w", err)
	}
	defer rows.Close()

	var result []*payment.MidtransPayment
	for rows.Next() {
		p, err := r.scanPayment(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	var total int64
	r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM midtrans_payment
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR referensi_id = $2)
		  AND ($3::text = '' OR referensi_tipe = $3)
		  AND ($4::text IS NULL OR status = $4)
	`, filter.BMTID, filter.ReferensiID, filter.ReferensiTipe, filter.Status).Scan(&total)

	return result, total, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status payment.StatusMidtrans, settledAt *time.Time, rawResponse string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE midtrans_payment
		SET status = $1, settled_at = $2, midtrans_response = $3, updated_at = NOW()
		WHERE id = $4
	`, status, settledAt, rawResponse, id)
	if err != nil {
		return fmt.Errorf("update status midtrans payment: %w", err)
	}
	return nil
}

func (r *PaymentRepository) UpdateSnapToken(ctx context.Context, id uuid.UUID, token, url string, expiredAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE midtrans_payment
		SET snap_token = $1, snap_url = $2, expired_at = $3, updated_at = NOW()
		WHERE id = $4
	`, token, url, expiredAt, id)
	if err != nil {
		return fmt.Errorf("update snap token midtrans payment: %w", err)
	}
	return nil
}

func (r *PaymentRepository) ListPending(ctx context.Context, olderThan time.Duration) ([]*payment.MidtransPayment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, order_id, referensi_id, referensi_tipe, nominal, status, metode_bayar,
		snap_token, snap_url, midtrans_response, settled_at, expired_at, created_at, updated_at
		FROM midtrans_payment
		WHERE status = 'PENDING' AND created_at < NOW() - $1::interval
		ORDER BY created_at
	`, olderThan.String())
	if err != nil {
		return nil, fmt.Errorf("list pending midtrans payment: %w", err)
	}
	defer rows.Close()

	var result []*payment.MidtransPayment
	for rows.Next() {
		p, err := r.scanPayment(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// ── UsageLog ──────────────────────────────────────────────────────────────────

func (r *PaymentRepository) CreateUsageLog(ctx context.Context, l *payment.UsageLog) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO usage_log (id, bmt_id, cabang_id, jenis, referensi_id, nominal, biaya_admin, periode, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, l.ID, l.BMTID, l.CabangID, l.Jenis, l.ReferensiID, l.Nominal, l.BiayaAdmin, l.Periode, l.CreatedAt)
	if err != nil {
		return fmt.Errorf("create usage log: %w", err)
	}
	return nil
}

func (r *PaymentRepository) SumBiayaAdminByPeriode(ctx context.Context, bmtID uuid.UUID, periode string) (int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(biaya_admin), 0) FROM usage_log WHERE bmt_id = $1 AND periode = $2
	`, bmtID, periode).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum biaya admin by periode: %w", err)
	}
	return total, nil
}

func (r *PaymentRepository) ListUsageLog(ctx context.Context, bmtID uuid.UUID, periode string, page, perPage int) ([]*payment.UsageLog, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, jenis, referensi_id, nominal, biaya_admin, periode, created_at
		FROM usage_log
		WHERE bmt_id = $1 AND periode = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, bmtID, periode, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list usage log: %w", err)
	}
	defer rows.Close()

	var result []*payment.UsageLog
	for rows.Next() {
		l := &payment.UsageLog{}
		err := rows.Scan(&l.ID, &l.BMTID, &l.CabangID, &l.Jenis, &l.ReferensiID, &l.Nominal, &l.BiayaAdmin, &l.Periode, &l.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan usage log: %w", err)
		}
		result = append(result, l)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM usage_log WHERE bmt_id = $1 AND periode = $2`, bmtID, periode).Scan(&total)

	return result, total, nil
}

// ── scanner helper ────────────────────────────────────────────────────────────

func (r *PaymentRepository) scanPayment(s scanner) (*payment.MidtransPayment, error) {
	p := &payment.MidtransPayment{}
	err := s.Scan(&p.ID, &p.BMTID, &p.OrderID, &p.ReferensiID, &p.ReferensiTipe, &p.Nominal,
		&p.Status, &p.MetodeBayar, &p.SnapToken, &p.SnapURL, &p.MidtransResponse,
		&p.SettledAt, &p.ExpiredAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, payment.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("scan midtrans payment: %w", err)
	}
	return p, nil
}
