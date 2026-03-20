package postgres

import (
	"context"

	"github.com/bmt-saas/api/internal/domain/pengguna"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PenggunaRepository struct {
	db *pgxpool.Pool
}

func NewPenggunaRepository(db *pgxpool.Pool) *PenggunaRepository {
	return &PenggunaRepository{db: db}
}

func (r *PenggunaRepository) Create(ctx context.Context, p *pengguna.Pengguna) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pengguna (id, bmt_id, cabang_id, username, password_hash, nama_lengkap,
		email, telepon, role, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, p.ID, p.BMTID, p.CabangID, p.Username, p.PasswordHash, p.NamaLengkap,
		p.Email, p.Telepon, p.Role, p.Status, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PenggunaRepository) GetByID(ctx context.Context, id uuid.UUID) (*pengguna.Pengguna, error) {
	return r.scan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, username, password_hash, nama_lengkap,
		email, telepon, role, status, last_login_at, created_at, updated_at
		FROM pengguna WHERE id = $1
	`, id))
}

func (r *PenggunaRepository) GetByUsername(ctx context.Context, bmtID uuid.UUID, username string) (*pengguna.Pengguna, error) {
	return r.scan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, username, password_hash, nama_lengkap,
		email, telepon, role, status, last_login_at, created_at, updated_at
		FROM pengguna WHERE bmt_id = $1 AND username = $2
	`, bmtID, username))
}

func (r *PenggunaRepository) ListByBMT(ctx context.Context, bmtID uuid.UUID) ([]*pengguna.Pengguna, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, username, password_hash, nama_lengkap,
		email, telepon, role, status, last_login_at, created_at, updated_at
		FROM pengguna WHERE bmt_id = $1 ORDER BY nama_lengkap
	`, bmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*pengguna.Pengguna
	for rows.Next() {
		p, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *PenggunaRepository) Update(ctx context.Context, p *pengguna.Pengguna) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pengguna SET nama_lengkap=$1, email=$2, telepon=$3, role=$4, updated_at=NOW()
		WHERE id=$5
	`, p.NamaLengkap, p.Email, p.Telepon, p.Role, p.ID)
	return err
}

func (r *PenggunaRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status pengguna.StatusPengguna) error {
	_, err := r.db.Exec(ctx, `UPDATE pengguna SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *PenggunaRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE pengguna SET last_login_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *PenggunaRepository) scan(s scanner) (*pengguna.Pengguna, error) {
	p := &pengguna.Pengguna{}
	err := s.Scan(&p.ID, &p.BMTID, &p.CabangID, &p.Username, &p.PasswordHash,
		&p.NamaLengkap, &p.Email, &p.Telepon, &p.Role, &p.Status,
		&p.LastLoginAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pengguna.ErrPenggunaNotFound
		}
		return nil, err
	}
	return p, nil
}
