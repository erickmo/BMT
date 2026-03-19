package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/platform"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlatformRepository struct {
	db *pgxpool.Pool
}

func NewPlatformRepository(db *pgxpool.Pool) *PlatformRepository {
	return &PlatformRepository{db: db}
}

func (r *PlatformRepository) CreateBMT(ctx context.Context, bmt *platform.BMT) error {
	wl, _ := json.Marshal(bmt.Whitelabel)
	_, err := r.db.Exec(ctx, `
		INSERT INTO bmt (id, kode, nama, alamat, telepon, email, logo_url, status, whitelabel, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, bmt.ID, bmt.Kode, bmt.Nama, bmt.Alamat, bmt.Telepon, bmt.Email,
		bmt.LogoURL, bmt.Status, wl, bmt.CreatedAt, bmt.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal create bmt: %w", err)
	}
	return nil
}

func (r *PlatformRepository) GetBMT(ctx context.Context, id uuid.UUID) (*platform.BMT, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, kode, nama, alamat, telepon, email, logo_url, status, whitelabel, created_at, updated_at
		FROM bmt WHERE id = $1
	`, id)
	return scanBMT(row)
}

func (r *PlatformRepository) GetBMTByKode(ctx context.Context, kode string) (*platform.BMT, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, kode, nama, alamat, telepon, email, logo_url, status, whitelabel, created_at, updated_at
		FROM bmt WHERE kode = $1
	`, kode)
	return scanBMT(row)
}

func (r *PlatformRepository) ListBMT(ctx context.Context) ([]*platform.BMT, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, kode, nama, alamat, telepon, email, logo_url, status, whitelabel, created_at, updated_at
		FROM bmt ORDER BY nama
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*platform.BMT
	for rows.Next() {
		bmt, err := scanBMT(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, bmt)
	}
	return result, nil
}

func (r *PlatformRepository) UpdateBMTStatus(ctx context.Context, id uuid.UUID, status platform.StatusBMT) error {
	_, err := r.db.Exec(ctx, `UPDATE bmt SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *PlatformRepository) CreateCabang(ctx context.Context, cabang *platform.Cabang) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO cabang (id, bmt_id, kode, nama, alamat, telepon, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, cabang.ID, cabang.BMTID, cabang.Kode, cabang.Nama, cabang.Alamat,
		cabang.Telepon, cabang.Status, cabang.CreatedAt, cabang.UpdatedAt)
	return err
}

func (r *PlatformRepository) GetCabang(ctx context.Context, id uuid.UUID) (*platform.Cabang, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, kode, nama, alamat, telepon, status, created_at, updated_at
		FROM cabang WHERE id = $1
	`, id)
	return scanCabang(row)
}

func (r *PlatformRepository) ListCabangByBMT(ctx context.Context, bmtID uuid.UUID) ([]*platform.Cabang, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, kode, nama, alamat, telepon, status, created_at, updated_at
		FROM cabang WHERE bmt_id = $1 ORDER BY nama
	`, bmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*platform.Cabang
	for rows.Next() {
		c, err := scanCabang(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *PlatformRepository) CreateKontrak(ctx context.Context, k *platform.KontrakBMT) error {
	fitur, _ := json.Marshal(k.Fitur)
	tarif, _ := json.Marshal(k.Tarif)
	_, err := r.db.Exec(ctx, `
		INSERT INTO kontrak_bmt (id, bmt_id, tanggal_mulai, tanggal_selesai, fitur, tarif, pic_nama, pic_telepon, pic_email, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, k.ID, k.BMTID, k.TanggalMulai, k.TanggalSelesai, fitur, tarif,
		k.PICNama, k.PICTelepon, k.PICEmail, k.Status, k.CreatedAt, k.UpdatedAt)
	return err
}

func (r *PlatformRepository) GetKontrakAktif(ctx context.Context, bmtID uuid.UUID) (*platform.KontrakBMT, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, tanggal_mulai, tanggal_selesai, fitur, tarif, pic_nama, pic_telepon, pic_email, status, created_at, updated_at
		FROM kontrak_bmt
		WHERE bmt_id = $1 AND status = 'AKTIF' AND tanggal_selesai >= CURRENT_DATE
		ORDER BY tanggal_mulai DESC LIMIT 1
	`, bmtID)

	var k platform.KontrakBMT
	var fitur, tarif []byte
	err := row.Scan(&k.ID, &k.BMTID, &k.TanggalMulai, &k.TanggalSelesai,
		&fitur, &tarif, &k.PICNama, &k.PICTelepon, &k.PICEmail,
		&k.Status, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, platform.ErrBMTNotFound
		}
		return nil, err
	}
	json.Unmarshal(fitur, &k.Fitur)
	json.Unmarshal(tarif, &k.Tarif)
	return &k, nil
}

func (r *PlatformRepository) GetPecahanAktif(ctx context.Context) ([]*platform.PecahanUang, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, nominal, jenis, label, is_aktif, urutan, berlaku_sejak, ditarik_pada, created_at, updated_at
		FROM pecahan_uang
		WHERE is_aktif = true AND (ditarik_pada IS NULL OR ditarik_pada > CURRENT_DATE)
		ORDER BY urutan
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*platform.PecahanUang
	for rows.Next() {
		p := &platform.PecahanUang{}
		err := rows.Scan(&p.ID, &p.Nominal, &p.Jenis, &p.Label,
			&p.IsAktif, &p.Urutan, &p.BerlakuSejak, &p.DitarikPada,
			new(time.Time), new(time.Time))
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *PlatformRepository) CreatePecahan(ctx context.Context, p *platform.PecahanUang) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pecahan_uang (id, nominal, jenis, label, is_aktif, urutan, berlaku_sejak)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, p.ID, p.Nominal, p.Jenis, p.Label, p.IsAktif, p.Urutan, p.BerlakuSejak)
	return err
}

func (r *PlatformRepository) UpdatePecahan(ctx context.Context, p *platform.PecahanUang) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pecahan_uang SET label = $1, is_aktif = $2, urutan = $3, ditarik_pada = $4, updated_at = NOW()
		WHERE id = $5
	`, p.Label, p.IsAktif, p.Urutan, p.DitarikPada, p.ID)
	return err
}

// Helper scanner functions
type scanner interface {
	Scan(dest ...any) error
}

func scanBMT(s scanner) (*platform.BMT, error) {
	bmt := &platform.BMT{}
	var wl []byte
	err := s.Scan(&bmt.ID, &bmt.Kode, &bmt.Nama, &bmt.Alamat, &bmt.Telepon,
		&bmt.Email, &bmt.LogoURL, &bmt.Status, &wl, &bmt.CreatedAt, &bmt.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, platform.ErrBMTNotFound
		}
		return nil, err
	}
	bmt.Whitelabel = make(map[string]string)
	json.Unmarshal(wl, &bmt.Whitelabel)
	return bmt, nil
}

func scanCabang(s scanner) (*platform.Cabang, error) {
	c := &platform.Cabang{}
	err := s.Scan(&c.ID, &c.BMTID, &c.Kode, &c.Nama, &c.Alamat,
		&c.Telepon, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, platform.ErrCabangNotFound
		}
		return nil, err
	}
	return c, nil
}
