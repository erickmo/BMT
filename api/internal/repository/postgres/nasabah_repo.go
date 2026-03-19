package postgres

import (
	"context"
	"fmt"

	"github.com/bmt-saas/api/internal/domain/nasabah"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NasabahRepository struct {
	db *pgxpool.Pool
}

func NewNasabahRepository(db *pgxpool.Pool) *NasabahRepository {
	return &NasabahRepository{db: db}
}

func (r *NasabahRepository) Create(ctx context.Context, n *nasabah.Nasabah) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO nasabah (id, bmt_id, cabang_id, nomor_nasabah, nik, nama_lengkap, tempat_lahir,
		tanggal_lahir, jenis_kelamin, alamat, telepon, email, foto_url, pekerjaan, status,
		pin_hash, password_hash, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
	`, n.ID, n.BMTID, n.CabangID, n.NomorNasabah, n.NIK, n.NamaLengkap,
		n.TempatLahir, n.TanggalLahir, n.JenisKelamin, n.Alamat, n.Telepon,
		n.Email, n.FotoURL, n.Pekerjaan, n.Status, n.PINHash, n.PasswordHash,
		n.CreatedAt, n.UpdatedAt)
	return err
}

func (r *NasabahRepository) GetByID(ctx context.Context, id uuid.UUID) (*nasabah.Nasabah, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_nasabah, nik, nama_lengkap, tempat_lahir,
		tanggal_lahir, jenis_kelamin, alamat, telepon, email, foto_url, pekerjaan, status,
		pin_hash, password_hash, last_login_at, created_at, updated_at
		FROM nasabah WHERE id = $1
	`, id)
	return scanNasabah(row)
}

func (r *NasabahRepository) GetByNomor(ctx context.Context, bmtID uuid.UUID, nomor string) (*nasabah.Nasabah, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_nasabah, nik, nama_lengkap, tempat_lahir,
		tanggal_lahir, jenis_kelamin, alamat, telepon, email, foto_url, pekerjaan, status,
		pin_hash, password_hash, last_login_at, created_at, updated_at
		FROM nasabah WHERE bmt_id = $1 AND nomor_nasabah = $2
	`, bmtID, nomor)
	return scanNasabah(row)
}

func (r *NasabahRepository) GetByNIK(ctx context.Context, bmtID uuid.UUID, nik string) (*nasabah.Nasabah, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_nasabah, nik, nama_lengkap, tempat_lahir,
		tanggal_lahir, jenis_kelamin, alamat, telepon, email, foto_url, pekerjaan, status,
		pin_hash, password_hash, last_login_at, created_at, updated_at
		FROM nasabah WHERE bmt_id = $1 AND nik = $2
	`, bmtID, nik)
	return scanNasabah(row)
}

func (r *NasabahRepository) Search(ctx context.Context, bmtID uuid.UUID, query string, limit, offset int) ([]*nasabah.Nasabah, int64, error) {
	// Full text search
	likeQuery := "%" + query + "%"
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_nasabah, nik, nama_lengkap, tempat_lahir,
		tanggal_lahir, jenis_kelamin, alamat, telepon, email, foto_url, pekerjaan, status,
		pin_hash, password_hash, last_login_at, created_at, updated_at
		FROM nasabah
		WHERE bmt_id = $1 AND (nama_lengkap ILIKE $2 OR nomor_nasabah ILIKE $2 OR telepon ILIKE $2 OR nik ILIKE $2)
		ORDER BY nama_lengkap
		LIMIT $3 OFFSET $4
	`, bmtID, likeQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*nasabah.Nasabah
	for rows.Next() {
		n, err := scanNasabah(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, n)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM nasabah WHERE bmt_id = $1 AND (nama_lengkap ILIKE $2 OR nomor_nasabah ILIKE $2)`, bmtID, likeQuery).Scan(&total)

	return result, total, nil
}

func (r *NasabahRepository) Update(ctx context.Context, n *nasabah.Nasabah) error {
	_, err := r.db.Exec(ctx, `
		UPDATE nasabah SET nik=$1, nama_lengkap=$2, tempat_lahir=$3, tanggal_lahir=$4,
		jenis_kelamin=$5, alamat=$6, telepon=$7, email=$8, foto_url=$9, pekerjaan=$10,
		updated_at=NOW() WHERE id=$11
	`, n.NIK, n.NamaLengkap, n.TempatLahir, n.TanggalLahir, n.JenisKelamin,
		n.Alamat, n.Telepon, n.Email, n.FotoURL, n.Pekerjaan, n.ID)
	return err
}

func (r *NasabahRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status nasabah.StatusNasabah) error {
	_, err := r.db.Exec(ctx, `UPDATE nasabah SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *NasabahRepository) CreateKartuNFC(ctx context.Context, k *nasabah.KartuNFC) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO kartu_nfc (id, bmt_id, nasabah_id, uid, pin_hash, limit_per_transaksi, limit_harian, saldo_nfc, status, expired_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, k.ID, k.BMTID, k.NasabahID, k.UID, k.PINHash, k.LimitPerTransaksi,
		k.LimitHarian, k.SaldoNFC, k.Status, k.ExpiredAt, k.CreatedAt, k.UpdatedAt)
	return err
}

func (r *NasabahRepository) GetKartuNFCByUID(ctx context.Context, uid string) (*nasabah.KartuNFC, error) {
	k := &nasabah.KartuNFC{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, nasabah_id, uid, pin_hash, limit_per_transaksi, limit_harian, saldo_nfc, status, expired_at, created_at, updated_at
		FROM kartu_nfc WHERE uid = $1
	`, uid).Scan(&k.ID, &k.BMTID, &k.NasabahID, &k.UID, &k.PINHash, &k.LimitPerTransaksi,
		&k.LimitHarian, &k.SaldoNFC, &k.Status, &k.ExpiredAt, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nasabah.ErrNasabahNotFound
		}
		return nil, err
	}
	return k, nil
}

func (r *NasabahRepository) GetKartuNFCByNasabah(ctx context.Context, nasabahID uuid.UUID) ([]*nasabah.KartuNFC, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, nasabah_id, uid, pin_hash, limit_per_transaksi, limit_harian, saldo_nfc, status, expired_at, created_at, updated_at
		FROM kartu_nfc WHERE nasabah_id = $1
	`, nasabahID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*nasabah.KartuNFC
	for rows.Next() {
		k := &nasabah.KartuNFC{}
		err := rows.Scan(&k.ID, &k.BMTID, &k.NasabahID, &k.UID, &k.PINHash,
			&k.LimitPerTransaksi, &k.LimitHarian, &k.SaldoNFC, &k.Status,
			&k.ExpiredAt, &k.CreatedAt, &k.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	return result, nil
}

func (r *NasabahRepository) UpdateKartuNFC(ctx context.Context, k *nasabah.KartuNFC) error {
	_, err := r.db.Exec(ctx, `
		UPDATE kartu_nfc SET pin_hash=$1, limit_per_transaksi=$2, limit_harian=$3, saldo_nfc=$4, status=$5, updated_at=NOW()
		WHERE id=$6
	`, k.PINHash, k.LimitPerTransaksi, k.LimitHarian, k.SaldoNFC, k.Status, k.ID)
	return err
}

func (r *NasabahRepository) GenerateNomorNasabah(ctx context.Context, bmtID uuid.UUID) (string, error) {
	// Get BMT kode
	var kode string
	err := r.db.QueryRow(ctx, `SELECT kode FROM bmt WHERE id = $1`, bmtID).Scan(&kode)
	if err != nil {
		return "", fmt.Errorf("gagal ambil kode bmt: %w", err)
	}

	// Get next sequence
	var seq int
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(CAST(SUBSTRING(nomor_nasabah FROM '(\d+)$') AS INTEGER)), 0) + 1
		FROM nasabah WHERE bmt_id = $1
	`, bmtID).Scan(&seq)
	if err != nil {
		return "", fmt.Errorf("gagal generate nomor: %w", err)
	}

	return fmt.Sprintf("%s-%08d", kode, seq), nil
}

func scanNasabah(s scanner) (*nasabah.Nasabah, error) {
	n := &nasabah.Nasabah{}
	err := s.Scan(&n.ID, &n.BMTID, &n.CabangID, &n.NomorNasabah, &n.NIK,
		&n.NamaLengkap, &n.TempatLahir, &n.TanggalLahir, &n.JenisKelamin,
		&n.Alamat, &n.Telepon, &n.Email, &n.FotoURL, &n.Pekerjaan, &n.Status,
		&n.PINHash, &n.PasswordHash, &n.LastLoginAt, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nasabah.ErrNasabahNotFound
		}
		return nil, err
	}
	return n, nil
}
