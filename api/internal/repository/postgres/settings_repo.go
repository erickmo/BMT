package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepository struct {
	db *pgxpool.Pool
}

func NewSettingsRepository(db *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) GetPlatform(ctx context.Context, kunci string) (string, error) {
	var nilai string
	err := r.db.QueryRow(ctx, `SELECT nilai FROM platform_settings WHERE kunci = $1`, kunci).Scan(&nilai)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("settings %s tidak ditemukan", kunci)
		}
		return "", err
	}
	return nilai, nil
}

func (r *SettingsRepository) GetBMT(ctx context.Context, bmtID uuid.UUID, kunci string) (string, error) {
	var nilai string
	err := r.db.QueryRow(ctx, `SELECT nilai FROM bmt_settings WHERE bmt_id = $1 AND kunci = $2`, bmtID, kunci).Scan(&nilai)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("bmt settings %s tidak ditemukan", kunci)
		}
		return "", err
	}
	return nilai, nil
}

func (r *SettingsRepository) GetCabang(ctx context.Context, cabangID uuid.UUID, kunci string) (string, error) {
	var nilai string
	err := r.db.QueryRow(ctx, `SELECT nilai FROM cabang_settings WHERE cabang_id = $1 AND kunci = $2`, cabangID, kunci).Scan(&nilai)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("cabang settings %s tidak ditemukan", kunci)
		}
		return "", err
	}
	return nilai, nil
}

func (r *SettingsRepository) SetPlatform(ctx context.Context, kunci, nilai, tipe, updatedBy string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO platform_settings (kunci, nilai, tipe, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (kunci) DO UPDATE SET nilai = $2, tipe = $3, updated_by = $4, updated_at = NOW()
	`, kunci, nilai, tipe, updatedBy)
	return err
}

func (r *SettingsRepository) SetBMT(ctx context.Context, bmtID uuid.UUID, kunci, nilai, tipe string, updatedBy uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO bmt_settings (bmt_id, kunci, nilai, tipe, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (bmt_id, kunci) DO UPDATE SET nilai = $3, tipe = $4, updated_by = $5, updated_at = NOW()
	`, bmtID, kunci, nilai, tipe, updatedBy)
	return err
}

func (r *SettingsRepository) SetCabang(ctx context.Context, cabangID, bmtID uuid.UUID, kunci, nilai, tipe string, updatedBy uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO cabang_settings (cabang_id, bmt_id, kunci, nilai, tipe, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (cabang_id, kunci) DO UPDATE SET nilai = $4, tipe = $5, updated_by = $6, updated_at = NOW()
	`, cabangID, bmtID, kunci, nilai, tipe, updatedBy)
	return err
}
