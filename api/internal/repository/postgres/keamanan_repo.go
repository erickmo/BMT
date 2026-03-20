package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bmt-saas/api/internal/domain/keamanan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KeamananRepository struct {
	db *pgxpool.Pool
}

func NewKeamananRepository(db *pgxpool.Pool) *KeamananRepository {
	return &KeamananRepository{db: db}
}

// ── OTP ──────────────────────────────────────────────────────────────────────

func (r *KeamananRepository) CreateOTP(ctx context.Context, o *keamanan.OTPLog) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO otp_log (id, tujuan, channel, kode_hash, tipe, referensi_id, is_digunakan, expired_at, ip_address, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, o.ID, o.Tujuan, o.Channel, o.KodeHash, o.Tipe, o.ReferensiID,
		o.IsDigunakan, o.ExpiredAt, o.IPAddress, o.CreatedAt)
	return err
}

func (r *KeamananRepository) GetOTPByTujuanAndTipe(ctx context.Context, tujuan string, tipe keamanan.TipeOTP) (*keamanan.OTPLog, error) {
	o := &keamanan.OTPLog{}
	err := r.db.QueryRow(ctx, `
		SELECT id, tujuan, channel, kode_hash, tipe, referensi_id, is_digunakan, expired_at, ip_address, created_at
		FROM otp_log
		WHERE tujuan = $1 AND tipe = $2 AND is_digunakan = false
		ORDER BY created_at DESC LIMIT 1
	`, tujuan, tipe).Scan(&o.ID, &o.Tujuan, &o.Channel, &o.KodeHash, &o.Tipe,
		&o.ReferensiID, &o.IsDigunakan, &o.ExpiredAt, &o.IPAddress, &o.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, keamanan.ErrOTPNotFound
		}
		return nil, err
	}
	return o, nil
}

func (r *KeamananRepository) MarkOTPDigunakan(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE otp_log SET is_digunakan = true WHERE id = $1`, id)
	return err
}

// ── Sesi ─────────────────────────────────────────────────────────────────────

func (r *KeamananRepository) CreateSesi(ctx context.Context, s *keamanan.SesiAktif) error {
	deviceInfo, _ := json.Marshal(s.DeviceInfo)
	_, err := r.db.Exec(ctx, `
		INSERT INTO sesi_aktif (id, subjek_id, subjek_tipe, refresh_token_hash, device_info, ip_address,
		last_active_at, expired_at, is_aktif, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, s.ID, s.SubjekID, s.SubjekTipe, s.RefreshTokenHash, deviceInfo,
		s.IPAddress, s.LastActiveAt, s.ExpiredAt, s.IsAktif, s.CreatedAt)
	return err
}

func (r *KeamananRepository) GetSesiByRefreshHash(ctx context.Context, hash string) (*keamanan.SesiAktif, error) {
	s := &keamanan.SesiAktif{}
	var deviceInfo []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, subjek_id, subjek_tipe, refresh_token_hash, device_info, ip_address,
		last_active_at, expired_at, is_aktif, created_at
		FROM sesi_aktif WHERE refresh_token_hash = $1 AND is_aktif = true
	`, hash).Scan(&s.ID, &s.SubjekID, &s.SubjekTipe, &s.RefreshTokenHash,
		&deviceInfo, &s.IPAddress, &s.LastActiveAt, &s.ExpiredAt, &s.IsAktif, &s.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, keamanan.ErrSesiNotFound
		}
		return nil, err
	}
	s.DeviceInfo = json.RawMessage(deviceInfo)
	return s, nil
}

func (r *KeamananRepository) ListSesiBySubjek(ctx context.Context, subjekID uuid.UUID) ([]*keamanan.SesiAktif, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, subjek_id, subjek_tipe, refresh_token_hash, device_info, ip_address,
		last_active_at, expired_at, is_aktif, created_at
		FROM sesi_aktif WHERE subjek_id = $1 ORDER BY created_at DESC
	`, subjekID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*keamanan.SesiAktif
	for rows.Next() {
		s := &keamanan.SesiAktif{}
		var deviceInfo []byte
		if err := rows.Scan(&s.ID, &s.SubjekID, &s.SubjekTipe, &s.RefreshTokenHash,
			&deviceInfo, &s.IPAddress, &s.LastActiveAt, &s.ExpiredAt, &s.IsAktif, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.DeviceInfo = json.RawMessage(deviceInfo)
		result = append(result, s)
	}
	return result, nil
}

func (r *KeamananRepository) NonaktifkanSesi(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE sesi_aktif SET is_aktif = false WHERE id = $1`, id)
	return err
}

func (r *KeamananRepository) NonaktifkanSemuaSesi(ctx context.Context, subjekID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE sesi_aktif SET is_aktif = false WHERE subjek_id = $1`, subjekID)
	return err
}

func (r *KeamananRepository) UpdateLastActive(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE sesi_aktif SET last_active_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *KeamananRepository) DeleteExpiredSesi(ctx context.Context) (int64, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM sesi_aktif WHERE expired_at < NOW()`)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// ── Fraud ────────────────────────────────────────────────────────────────────

func (r *KeamananRepository) CreateFraudRule(ctx context.Context, rule *keamanan.FraudRule) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO fraud_rule (id, bmt_id, nama, tipe, kondisi, aksi, is_aktif, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, rule.ID, rule.BMTID, rule.Nama, rule.Tipe, []byte(rule.Kondisi), rule.Aksi, rule.IsAktif, rule.CreatedAt)
	return err
}

func (r *KeamananRepository) GetFraudRuleByID(ctx context.Context, id uuid.UUID) (*keamanan.FraudRule, error) {
	rule := &keamanan.FraudRule{}
	var kondisi []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, nama, tipe, kondisi, aksi, is_aktif, created_at
		FROM fraud_rule WHERE id = $1
	`, id).Scan(&rule.ID, &rule.BMTID, &rule.Nama, &rule.Tipe, &kondisi, &rule.Aksi, &rule.IsAktif, &rule.CreatedAt)
	if err != nil {
		return nil, err
	}
	rule.Kondisi = json.RawMessage(kondisi)
	return rule, nil
}

func (r *KeamananRepository) ListFraudRuleAktif(ctx context.Context, bmtID *uuid.UUID) ([]*keamanan.FraudRule, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, nama, tipe, kondisi, aksi, is_aktif, created_at
		FROM fraud_rule WHERE is_aktif = true AND (bmt_id IS NULL OR bmt_id = $1)
		ORDER BY created_at
	`, bmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*keamanan.FraudRule
	for rows.Next() {
		rule := &keamanan.FraudRule{}
		var kondisi []byte
		if err := rows.Scan(&rule.ID, &rule.BMTID, &rule.Nama, &rule.Tipe, &kondisi, &rule.Aksi, &rule.IsAktif, &rule.CreatedAt); err != nil {
			return nil, err
		}
		rule.Kondisi = json.RawMessage(kondisi)
		result = append(result, rule)
	}
	return result, nil
}

func (r *KeamananRepository) CreateFraudAlert(ctx context.Context, a *keamanan.FraudAlert) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO fraud_alert (id, bmt_id, rule_id, nasabah_id, transaksi_id, deskripsi, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, a.ID, a.BMTID, a.RuleID, a.NasabahID, a.TransaksiID, a.Deskripsi, a.Status, a.CreatedAt)
	return err
}

func (r *KeamananRepository) GetFraudAlertByID(ctx context.Context, id uuid.UUID) (*keamanan.FraudAlert, error) {
	a := &keamanan.FraudAlert{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, rule_id, nasabah_id, transaksi_id, deskripsi, status, direview_oleh, direview_at, created_at
		FROM fraud_alert WHERE id = $1
	`, id).Scan(&a.ID, &a.BMTID, &a.RuleID, &a.NasabahID, &a.TransaksiID,
		&a.Deskripsi, &a.Status, &a.DireviewOleh, &a.DireviewAt, &a.CreatedAt)
	return a, err
}

func (r *KeamananRepository) ListFraudAlert(ctx context.Context, bmtID uuid.UUID, status *keamanan.StatusFraud, page, perPage int) ([]*keamanan.FraudAlert, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, rule_id, nasabah_id, transaksi_id, deskripsi, status, direview_oleh, direview_at, created_at
		FROM fraud_alert WHERE bmt_id = $1 AND ($2::text IS NULL OR status = $2)
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, bmtID, status, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []*keamanan.FraudAlert
	for rows.Next() {
		a := &keamanan.FraudAlert{}
		if err := rows.Scan(&a.ID, &a.BMTID, &a.RuleID, &a.NasabahID, &a.TransaksiID,
			&a.Deskripsi, &a.Status, &a.DireviewOleh, &a.DireviewAt, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, a)
	}
	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM fraud_alert WHERE bmt_id = $1 AND ($2::text IS NULL OR status = $2)`,
		bmtID, status).Scan(&total)
	return result, total, nil
}

func (r *KeamananRepository) UpdateStatusFraudAlert(ctx context.Context, id uuid.UUID, status keamanan.StatusFraud, reviewOleh uuid.UUID) error {
	now := time.Now()
	_, err := r.db.Exec(ctx, `
		UPDATE fraud_alert SET status=$1, direview_oleh=$2, direview_at=$3 WHERE id=$4
	`, status, reviewOleh, now, id)
	return err
}

// ── Audit ─────────────────────────────────────────────────────────────────────

func (r *KeamananRepository) CreateAuditLog(ctx context.Context, l *keamanan.AuditLog) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO audit_log (id, bmt_id, subjek_id, subjek_tipe, aksi, resource_tipe, resource_id,
		data_sebelum, data_sesudah, ip_address, user_agent, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, l.ID, l.BMTID, l.SubjekID, l.SubjekTipe, l.Aksi, l.ResourceTipe, l.ResourceID,
		[]byte(l.DataSebelum), []byte(l.DataSesudah), l.IPAddress, l.UserAgent, l.CreatedAt)
	return err
}

func (r *KeamananRepository) ListAuditLog(ctx context.Context, bmtID *uuid.UUID, subjekID *uuid.UUID, resourceTipe string, page, perPage int) ([]*keamanan.AuditLog, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, subjek_id, subjek_tipe, aksi, resource_tipe, resource_id, data_sebelum, data_sesudah, ip_address, user_agent, created_at
		FROM audit_log
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		AND ($2::uuid IS NULL OR subjek_id = $2)
		AND ($3 = '' OR resource_tipe = $3)
		ORDER BY created_at DESC LIMIT $4 OFFSET $5
	`, bmtID, subjekID, resourceTipe, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []*keamanan.AuditLog
	for rows.Next() {
		l := &keamanan.AuditLog{}
		var sebelum, sesudah []byte
		if err := rows.Scan(&l.ID, &l.BMTID, &l.SubjekID, &l.SubjekTipe, &l.Aksi,
			&l.ResourceTipe, &l.ResourceID, &sebelum, &sesudah, &l.IPAddress, &l.UserAgent, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		l.DataSebelum = json.RawMessage(sebelum)
		l.DataSesudah = json.RawMessage(sesudah)
		result = append(result, l)
	}
	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM audit_log WHERE ($1::uuid IS NULL OR bmt_id = $1) AND ($2::uuid IS NULL OR subjek_id = $2) AND ($3 = '' OR resource_tipe = $3)`,
		bmtID, subjekID, resourceTipe).Scan(&total)
	return result, total, nil
}

func (r *KeamananRepository) DeleteAuditLogOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM audit_log WHERE created_at < $1`, cutoff)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
