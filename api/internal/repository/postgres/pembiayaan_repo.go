package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/pembiayaan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PembiayaanRepository struct {
	db *pgxpool.Pool
}

func NewPembiayaanRepository(db *pgxpool.Pool) *PembiayaanRepository {
	return &PembiayaanRepository{db: db}
}

func (r *PembiayaanRepository) Create(ctx context.Context, p *pembiayaan.Pembiayaan) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pembiayaan (id, bmt_id, cabang_id, nasabah_id, produk_pembiayaan_id, nomor_pembiayaan,
		akad, pokok, margin_persen, nisbah_nasabah, jangka_bulan, angsuran_per_bulan, total_kewajiban,
		ada_beasiswa, beasiswa_persen, beasiswa_nominal, beasiswa_sumber, beasiswa_ditetapkan_oleh, beasiswa_ditetapkan_at,
		status, kolektibilitas, hari_tunggak, saldo_pokok, saldo_margin,
		created_at, updated_at, created_by, updated_by, is_voided)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29)
	`, p.ID, p.BMTID, p.CabangID, p.NasabahID, p.ProdukPembiayaanID, p.NomorPembiayaan,
		p.Akad, p.Pokok, p.MarginPersen, p.NisbahNasabah, p.JangkaBulan, p.AngsuranPerBulan, p.TotalKewajiban,
		p.AdaBeasiswa, p.BeasiswaPersen, p.BeasiswaNominal, p.BeasiswaSumber, p.BeasiswaDitetapkanOleh, p.BeasiswaDitetapkanAt,
		p.Status, p.Kolektibilitas, p.HariTunggak, p.SaldoPokok, p.SaldoMargin,
		p.CreatedAt, p.UpdatedAt, p.CreatedBy, p.UpdatedBy, p.IsVoided)
	return err
}

func (r *PembiayaanRepository) GetByID(ctx context.Context, id uuid.UUID) (*pembiayaan.Pembiayaan, error) {
	return r.scanPembiayaan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, produk_pembiayaan_id, nomor_pembiayaan,
		akad, pokok, margin_persen, nisbah_nasabah, jangka_bulan, angsuran_per_bulan, total_kewajiban,
		ada_beasiswa, beasiswa_persen, beasiswa_nominal, beasiswa_sumber, beasiswa_ditetapkan_oleh, beasiswa_ditetapkan_at,
		status, kolektibilitas, hari_tunggak, saldo_pokok, saldo_margin,
		created_at, updated_at, created_by, updated_by, is_voided
		FROM pembiayaan WHERE id = $1
	`, id))
}

func (r *PembiayaanRepository) GetByNomor(ctx context.Context, nomor string) (*pembiayaan.Pembiayaan, error) {
	return r.scanPembiayaan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, produk_pembiayaan_id, nomor_pembiayaan,
		akad, pokok, margin_persen, nisbah_nasabah, jangka_bulan, angsuran_per_bulan, total_kewajiban,
		ada_beasiswa, beasiswa_persen, beasiswa_nominal, beasiswa_sumber, beasiswa_ditetapkan_oleh, beasiswa_ditetapkan_at,
		status, kolektibilitas, hari_tunggak, saldo_pokok, saldo_margin,
		created_at, updated_at, created_by, updated_by, is_voided
		FROM pembiayaan WHERE nomor_pembiayaan = $1
	`, nomor))
}

func (r *PembiayaanRepository) List(ctx context.Context, filter pembiayaan.ListPembiayaanFilter) ([]*pembiayaan.Pembiayaan, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, produk_pembiayaan_id, nomor_pembiayaan,
		akad, pokok, margin_persen, nisbah_nasabah, jangka_bulan, angsuran_per_bulan, total_kewajiban,
		ada_beasiswa, beasiswa_persen, beasiswa_nominal, beasiswa_sumber, beasiswa_ditetapkan_oleh, beasiswa_ditetapkan_at,
		status, kolektibilitas, hari_tunggak, saldo_pokok, saldo_margin,
		created_at, updated_at, created_by, updated_by, is_voided
		FROM pembiayaan
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR cabang_id = $2)
		  AND ($3::uuid IS NULL OR nasabah_id = $3)
		  AND ($4::text IS NULL OR status = $4)
		  AND is_voided = false
		ORDER BY created_at DESC LIMIT $5 OFFSET $6
	`, filter.BMTID, filter.CabangID, filter.NasabahID, filter.Status, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*pembiayaan.Pembiayaan
	for rows.Next() {
		p, err := r.scanPembiayaan(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pembiayaan WHERE ($1::uuid IS NULL OR bmt_id = $1) AND ($2::uuid IS NULL OR cabang_id = $2) AND is_voided = false`,
		filter.BMTID, filter.CabangID).Scan(&total)

	return result, total, nil
}

func (r *PembiayaanRepository) Update(ctx context.Context, p *pembiayaan.Pembiayaan) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pembiayaan SET status=$1, kolektibilitas=$2, hari_tunggak=$3, saldo_pokok=$4, saldo_margin=$5,
		ada_beasiswa=$6, beasiswa_persen=$7, beasiswa_nominal=$8, beasiswa_sumber=$9,
		beasiswa_ditetapkan_oleh=$10, beasiswa_ditetapkan_at=$11, updated_at=NOW(), updated_by=$12
		WHERE id=$13
	`, p.Status, p.Kolektibilitas, p.HariTunggak, p.SaldoPokok, p.SaldoMargin,
		p.AdaBeasiswa, p.BeasiswaPersen, p.BeasiswaNominal, p.BeasiswaSumber,
		p.BeasiswaDitetapkanOleh, p.BeasiswaDitetapkanAt, p.UpdatedBy, p.ID)
	return err
}

func (r *PembiayaanRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status pembiayaan.StatusPembiayaan) error {
	_, err := r.db.Exec(ctx, `UPDATE pembiayaan SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *PembiayaanRepository) UpdateKolektibilitas(ctx context.Context, id uuid.UUID, kolektibilitas int16, hariTunggak int) error {
	_, err := r.db.Exec(ctx, `UPDATE pembiayaan SET kolektibilitas=$1, hari_tunggak=$2, updated_at=NOW() WHERE id=$3`,
		kolektibilitas, hariTunggak, id)
	return err
}

func (r *PembiayaanRepository) UpdateSaldo(ctx context.Context, id uuid.UUID, saldoPokok, saldoMargin int64) error {
	_, err := r.db.Exec(ctx, `UPDATE pembiayaan SET saldo_pokok=$1, saldo_margin=$2, updated_at=NOW() WHERE id=$3`,
		saldoPokok, saldoMargin, id)
	return err
}

func (r *PembiayaanRepository) SetBeasiswa(ctx context.Context, id uuid.UUID, persen float64, nominal int64, sumber string, oleh uuid.UUID) error {
	now := time.Now()
	_, err := r.db.Exec(ctx, `
		UPDATE pembiayaan SET ada_beasiswa=true, beasiswa_persen=$1, beasiswa_nominal=$2,
		beasiswa_sumber=$3, beasiswa_ditetapkan_oleh=$4, beasiswa_ditetapkan_at=$5, updated_at=NOW()
		WHERE id=$6
	`, persen, nominal, sumber, oleh, now, id)
	return err
}

func (r *PembiayaanRepository) LockForUpdate(ctx context.Context, id uuid.UUID) (*pembiayaan.Pembiayaan, error) {
	return r.scanPembiayaan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, produk_pembiayaan_id, nomor_pembiayaan,
		akad, pokok, margin_persen, nisbah_nasabah, jangka_bulan, angsuran_per_bulan, total_kewajiban,
		ada_beasiswa, beasiswa_persen, beasiswa_nominal, beasiswa_sumber, beasiswa_ditetapkan_oleh, beasiswa_ditetapkan_at,
		status, kolektibilitas, hari_tunggak, saldo_pokok, saldo_margin,
		created_at, updated_at, created_by, updated_by, is_voided
		FROM pembiayaan WHERE id = $1 FOR UPDATE
	`, id))
}

func (r *PembiayaanRepository) GenerateNomor(ctx context.Context, bmtID, cabangID uuid.UUID) (string, error) {
	var kodeBMT, kodeCabang string
	r.db.QueryRow(ctx, `SELECT kode FROM bmt WHERE id = $1`, bmtID).Scan(&kodeBMT)
	r.db.QueryRow(ctx, `SELECT kode FROM cabang WHERE id = $1`, cabangID).Scan(&kodeCabang)

	var seq int
	r.db.QueryRow(ctx, `SELECT COALESCE(COUNT(*), 0) + 1 FROM pembiayaan WHERE bmt_id = $1 AND cabang_id = $2`, bmtID, cabangID).Scan(&seq)

	now := time.Now()
	return fmt.Sprintf("%s-%s-PB-%d%02d-%06d", kodeBMT, kodeCabang, now.Year(), now.Month(), seq), nil
}

func (r *PembiayaanRepository) CreateAngsuran(ctx context.Context, a *pembiayaan.AngsuranPembiayaan) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO angsuran_pembiayaan (id, bmt_id, pembiayaan_id, periode_bulan, nominal_pokok, nominal_margin,
		total_angsuran, tanggal_jatuh_tempo, nominal_terbayar, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, a.ID, a.BMTID, a.PembiayaanID, a.PeriodeBulan, a.NominalPokok, a.NominalMargin,
		a.TotalAngsuran, a.TanggalJatuhTempo, a.NominalTerbayar, a.Status, a.CreatedAt)
	return err
}

func (r *PembiayaanRepository) GetAngsuranByID(ctx context.Context, id uuid.UUID) (*pembiayaan.AngsuranPembiayaan, error) {
	a := &pembiayaan.AngsuranPembiayaan{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, pembiayaan_id, periode_bulan, nominal_pokok, nominal_margin,
		total_angsuran, tanggal_jatuh_tempo, tanggal_bayar, nominal_terbayar, status, transaksi_id, created_at
		FROM angsuran_pembiayaan WHERE id = $1
	`, id).Scan(&a.ID, &a.BMTID, &a.PembiayaanID, &a.PeriodeBulan, &a.NominalPokok, &a.NominalMargin,
		&a.TotalAngsuran, &a.TanggalJatuhTempo, &a.TanggalBayar, &a.NominalTerbayar, &a.Status, &a.TransaksiID, &a.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pembiayaan.ErrPembiayaanNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *PembiayaanRepository) ListAngsuran(ctx context.Context, pembiayaanID uuid.UUID) ([]*pembiayaan.AngsuranPembiayaan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, pembiayaan_id, periode_bulan, nominal_pokok, nominal_margin,
		total_angsuran, tanggal_jatuh_tempo, tanggal_bayar, nominal_terbayar, status, transaksi_id, created_at
		FROM angsuran_pembiayaan WHERE pembiayaan_id = $1 ORDER BY periode_bulan
	`, pembiayaanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pembiayaan.AngsuranPembiayaan
	for rows.Next() {
		a := &pembiayaan.AngsuranPembiayaan{}
		err := rows.Scan(&a.ID, &a.BMTID, &a.PembiayaanID, &a.PeriodeBulan, &a.NominalPokok, &a.NominalMargin,
			&a.TotalAngsuran, &a.TanggalJatuhTempo, &a.TanggalBayar, &a.NominalTerbayar, &a.Status, &a.TransaksiID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *PembiayaanRepository) GetAngsuranJatuhTempo(ctx context.Context, bmtID uuid.UUID, tanggal time.Time) ([]*pembiayaan.AngsuranPembiayaan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.bmt_id, a.pembiayaan_id, a.periode_bulan, a.nominal_pokok, a.nominal_margin,
		a.total_angsuran, a.tanggal_jatuh_tempo, a.tanggal_bayar, a.nominal_terbayar, a.status, a.transaksi_id, a.created_at
		FROM angsuran_pembiayaan a
		JOIN pembiayaan p ON p.id = a.pembiayaan_id
		WHERE p.bmt_id = $1 AND a.tanggal_jatuh_tempo::date = $2::date AND a.status = 'MENUNGGU'
	`, bmtID, tanggal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pembiayaan.AngsuranPembiayaan
	for rows.Next() {
		a := &pembiayaan.AngsuranPembiayaan{}
		err := rows.Scan(&a.ID, &a.BMTID, &a.PembiayaanID, &a.PeriodeBulan, &a.NominalPokok, &a.NominalMargin,
			&a.TotalAngsuran, &a.TanggalJatuhTempo, &a.TanggalBayar, &a.NominalTerbayar, &a.Status, &a.TransaksiID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *PembiayaanRepository) UpdateAngsuranTerbayar(ctx context.Context, id uuid.UUID, nominal int64, tanggalBayar time.Time, transaksiID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE angsuran_pembiayaan SET nominal_terbayar=$1, tanggal_bayar=$2, transaksi_id=$3,
		status=CASE WHEN nominal_terbayar + $1 >= total_angsuran THEN 'TERBAYAR' ELSE 'SEBAGIAN' END
		WHERE id=$4
	`, nominal, tanggalBayar, transaksiID, id)
	return err
}

func (r *PembiayaanRepository) CreateBeasiswaRiwayat(ctx context.Context, riwayat *pembiayaan.BeasiswaRiwayat) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO beasiswa_riwayat (id, pembiayaan_id, persen_sebelum, persen_sesudah, nominal_sebelum, nominal_sesudah, alasan, ditetapkan_oleh, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, riwayat.ID, riwayat.PembiayaanID, riwayat.PersenSebelum, riwayat.PersenSesudah,
		riwayat.NominalSebelum, riwayat.NominalSesudah, riwayat.Alasan, riwayat.DitetapkanOleh, riwayat.CreatedAt)
	return err
}

func (r *PembiayaanRepository) ListBeasiswaRiwayat(ctx context.Context, pembiayaanID uuid.UUID) ([]*pembiayaan.BeasiswaRiwayat, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, pembiayaan_id, persen_sebelum, persen_sesudah, nominal_sebelum, nominal_sesudah, alasan, ditetapkan_oleh, created_at
		FROM beasiswa_riwayat WHERE pembiayaan_id = $1 ORDER BY created_at DESC
	`, pembiayaanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pembiayaan.BeasiswaRiwayat
	for rows.Next() {
		rw := &pembiayaan.BeasiswaRiwayat{}
		err := rows.Scan(&rw.ID, &rw.PembiayaanID, &rw.PersenSebelum, &rw.PersenSesudah,
			&rw.NominalSebelum, &rw.NominalSesudah, &rw.Alasan, &rw.DitetapkanOleh, &rw.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, rw)
	}
	return result, nil
}

func (r *PembiayaanRepository) scanPembiayaan(s scanner) (*pembiayaan.Pembiayaan, error) {
	p := &pembiayaan.Pembiayaan{}
	err := s.Scan(&p.ID, &p.BMTID, &p.CabangID, &p.NasabahID, &p.ProdukPembiayaanID, &p.NomorPembiayaan,
		&p.Akad, &p.Pokok, &p.MarginPersen, &p.NisbahNasabah, &p.JangkaBulan, &p.AngsuranPerBulan, &p.TotalKewajiban,
		&p.AdaBeasiswa, &p.BeasiswaPersen, &p.BeasiswaNominal, &p.BeasiswaSumber, &p.BeasiswaDitetapkanOleh, &p.BeasiswaDitetapkanAt,
		&p.Status, &p.Kolektibilitas, &p.HariTunggak, &p.SaldoPokok, &p.SaldoMargin,
		&p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy, &p.IsVoided)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pembiayaan.ErrPembiayaanNotFound
		}
		return nil, err
	}
	return p, nil
}
