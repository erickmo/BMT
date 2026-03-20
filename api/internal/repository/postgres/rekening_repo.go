package postgres

import (
	"context"
	"fmt"

	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RekeningRepository struct {
	db *pgxpool.Pool
}

func NewRekeningRepository(db *pgxpool.Pool) *RekeningRepository {
	return &RekeningRepository{db: db}
}

func (r *RekeningRepository) Create(ctx context.Context, rek *rekening.Rekening) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO rekening (id, bmt_id, cabang_id, nasabah_id, jenis_rekening_id, nomor_rekening,
		saldo, status, biaya_admin_bulanan, nominal_deposito, nisbah_nasabah, tanggal_buka,
		tanggal_jatuh_tempo, tanggal_tutup, created_at, updated_at, created_by_form_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`, rek.ID, rek.BMTID, rek.CabangID, rek.NasabahID, rek.JenisRekeningID,
		rek.NomorRekening, rek.Saldo.Int64(), rek.Status, rek.BiayaAdminBulanan,
		rek.NominalDeposito, rek.NisbahNasabah, rek.TanggalBuka,
		rek.TanggalJatuhTempo, rek.TanggalTutup, rek.CreatedAt, rek.UpdatedAt,
		rek.CreatedByFormID)
	return err
}

func (r *RekeningRepository) GetByID(ctx context.Context, id uuid.UUID) (*rekening.Rekening, error) {
	return r.scanRekening(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, jenis_rekening_id, nomor_rekening,
		saldo, status, alasan_blokir, biaya_admin_bulanan, nominal_deposito, nisbah_nasabah,
		tanggal_buka, tanggal_jatuh_tempo, tanggal_tutup, created_at, updated_at, created_by_form_id
		FROM rekening WHERE id = $1
	`, id))
}

func (r *RekeningRepository) GetByNomor(ctx context.Context, nomor string) (*rekening.Rekening, error) {
	return r.scanRekening(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, jenis_rekening_id, nomor_rekening,
		saldo, status, alasan_blokir, biaya_admin_bulanan, nominal_deposito, nisbah_nasabah,
		tanggal_buka, tanggal_jatuh_tempo, tanggal_tutup, created_at, updated_at, created_by_form_id
		FROM rekening WHERE nomor_rekening = $1
	`, nomor))
}

func (r *RekeningRepository) ListByNasabah(ctx context.Context, nasabahID uuid.UUID) ([]*rekening.Rekening, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, jenis_rekening_id, nomor_rekening,
		saldo, status, alasan_blokir, biaya_admin_bulanan, nominal_deposito, nisbah_nasabah,
		tanggal_buka, tanggal_jatuh_tempo, tanggal_tutup, created_at, updated_at, created_by_form_id
		FROM rekening WHERE nasabah_id = $1 AND status != 'TUTUP' ORDER BY tanggal_buka
	`, nasabahID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*rekening.Rekening
	for rows.Next() {
		rek, err := r.scanRekening(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, rek)
	}
	return result, nil
}

func (r *RekeningRepository) ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, page, perPage int) ([]*rekening.Rekening, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, jenis_rekening_id, nomor_rekening,
		saldo, status, alasan_blokir, biaya_admin_bulanan, nominal_deposito, nisbah_nasabah,
		tanggal_buka, tanggal_jatuh_tempo, tanggal_tutup, created_at, updated_at, created_by_form_id
		FROM rekening WHERE bmt_id = $1 AND cabang_id = $2
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, bmtID, cabangID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*rekening.Rekening
	for rows.Next() {
		rek, err := r.scanRekening(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, rek)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM rekening WHERE bmt_id = $1 AND cabang_id = $2`, bmtID, cabangID).Scan(&total)

	return result, total, nil
}

func (r *RekeningRepository) UpdateSaldo(ctx context.Context, id uuid.UUID, saldoBaru int64) error {
	tag, err := r.db.Exec(ctx, `UPDATE rekening SET saldo=$1, updated_at=NOW() WHERE id=$2`, saldoBaru, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return rekening.ErrRekeningNotFound
	}
	return nil
}

func (r *RekeningRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status rekening.StatusRekening, alasan string) error {
	_, err := r.db.Exec(ctx, `UPDATE rekening SET status=$1, alasan_blokir=$2, updated_at=NOW() WHERE id=$3`, status, alasan, id)
	return err
}

func (r *RekeningRepository) LockForUpdate(ctx context.Context, id uuid.UUID) (*rekening.Rekening, error) {
	return r.scanRekening(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nasabah_id, jenis_rekening_id, nomor_rekening,
		saldo, status, alasan_blokir, biaya_admin_bulanan, nominal_deposito, nisbah_nasabah,
		tanggal_buka, tanggal_jatuh_tempo, tanggal_tutup, created_at, updated_at, created_by_form_id
		FROM rekening WHERE id = $1 FOR UPDATE
	`, id))
}

func (r *RekeningRepository) CreateTransaksi(ctx context.Context, t *rekening.TransaksiRekening) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transaksi_rekening (id, bmt_id, cabang_id, rekening_id, jenis, posisi, nominal,
		saldo_sebelum, saldo_sesudah, keterangan, referensi_id, referensi_tipe, idempotency_key, created_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, t.ID, t.BMTID, t.CabangID, t.RekeningID, t.Jenis, t.Posisi, t.Nominal,
		t.SaldoSebelum, t.SaldoSesudah, t.Keterangan, t.ReferensiID, t.ReferensiTipe,
		t.IdempotencyKey, t.CreatedBy, t.CreatedAt)
	return err
}

func (r *RekeningRepository) ListTransaksi(ctx context.Context, rekeningID uuid.UUID, limit, offset int) ([]*rekening.TransaksiRekening, int64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, rekening_id, jenis, posisi, nominal,
		saldo_sebelum, saldo_sesudah, keterangan, referensi_id, referensi_tipe, idempotency_key, created_by, created_at
		FROM transaksi_rekening WHERE rekening_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, rekeningID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*rekening.TransaksiRekening
	for rows.Next() {
		t := &rekening.TransaksiRekening{}
		err := rows.Scan(&t.ID, &t.BMTID, &t.CabangID, &t.RekeningID, &t.Jenis, &t.Posisi,
			&t.Nominal, &t.SaldoSebelum, &t.SaldoSesudah, &t.Keterangan,
			&t.ReferensiID, &t.ReferensiTipe, &t.IdempotencyKey, &t.CreatedBy, &t.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM transaksi_rekening WHERE rekening_id = $1`, rekeningID).Scan(&total)

	return result, total, nil
}

func (r *RekeningRepository) GetTransaksiByIdempotency(ctx context.Context, key uuid.UUID) (*rekening.TransaksiRekening, error) {
	t := &rekening.TransaksiRekening{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, rekening_id, jenis, posisi, nominal,
		saldo_sebelum, saldo_sesudah, keterangan, referensi_id, referensi_tipe, idempotency_key, created_by, created_at
		FROM transaksi_rekening WHERE idempotency_key = $1
	`, key).Scan(&t.ID, &t.BMTID, &t.CabangID, &t.RekeningID, &t.Jenis, &t.Posisi,
		&t.Nominal, &t.SaldoSebelum, &t.SaldoSesudah, &t.Keterangan,
		&t.ReferensiID, &t.ReferensiTipe, &t.IdempotencyKey, &t.CreatedBy, &t.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rekening.ErrRekeningNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *RekeningRepository) GenerateNomorRekening(ctx context.Context, bmtID, cabangID uuid.UUID, kodeJenis string) (string, error) {
	var kodeBMT, kodeCabang string
	r.db.QueryRow(ctx, `SELECT kode FROM bmt WHERE id = $1`, bmtID).Scan(&kodeBMT)
	r.db.QueryRow(ctx, `SELECT kode FROM cabang WHERE id = $1`, cabangID).Scan(&kodeCabang)

	var seq int
	r.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(CAST(RIGHT(nomor_rekening, 8) AS INTEGER)), 0) + 1
		FROM rekening WHERE bmt_id = $1 AND cabang_id = $2
	`, bmtID, cabangID).Scan(&seq)

	return fmt.Sprintf("%s-%s-%s-%08d", kodeBMT, kodeCabang, kodeJenis, seq), nil
}

func (r *RekeningRepository) CreateJenis(ctx context.Context, jr *rekening.JenisRekening) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO jenis_rekening (id, bmt_id, kode, nama, tipe_dasar, akad, deskripsi, setoran_awal_min,
		setoran_min, bisa_ditarik, syarat_penarikan, nisbah_nasabah, jangka_hari, biaya_admin_bulanan,
		bisa_nfc, bisa_autodebet, biaya_admin_buka, is_aktif, urutan_tampil, created_at, updated_at, created_by, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)
	`, jr.ID, jr.BMTID, jr.Kode, jr.Nama, jr.TipeDasar, jr.Akad, jr.Deskripsi,
		jr.SetoranAwalMin, jr.SetoranMin, jr.BisaDitarik, jr.SyaratPenarikan,
		jr.NisbahNasabah, jr.JangkaHari, jr.BiayaAdminBulanan, jr.BisaNFC,
		jr.BisaAutodebet, jr.BiayaAdminBuka, jr.IsAktif, jr.UrutanTampil,
		jr.CreatedAt, jr.UpdatedAt, jr.CreatedBy, jr.UpdatedBy)
	return err
}

func (r *RekeningRepository) GetJenisByID(ctx context.Context, id uuid.UUID) (*rekening.JenisRekening, error) {
	jr := &rekening.JenisRekening{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, kode, nama, tipe_dasar, akad, deskripsi, setoran_awal_min,
		setoran_min, bisa_ditarik, syarat_penarikan, nisbah_nasabah, jangka_hari, biaya_admin_bulanan,
		bisa_nfc, bisa_autodebet, biaya_admin_buka, is_aktif, urutan_tampil, created_at, updated_at, created_by, updated_by
		FROM jenis_rekening WHERE id = $1
	`, id).Scan(&jr.ID, &jr.BMTID, &jr.Kode, &jr.Nama, &jr.TipeDasar, &jr.Akad, &jr.Deskripsi,
		&jr.SetoranAwalMin, &jr.SetoranMin, &jr.BisaDitarik, &jr.SyaratPenarikan,
		&jr.NisbahNasabah, &jr.JangkaHari, &jr.BiayaAdminBulanan, &jr.BisaNFC,
		&jr.BisaAutodebet, &jr.BiayaAdminBuka, &jr.IsAktif, &jr.UrutanTampil,
		&jr.CreatedAt, &jr.UpdatedAt, &jr.CreatedBy, &jr.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rekening.ErrRekeningNotFound
		}
		return nil, err
	}
	return jr, nil
}

func (r *RekeningRepository) ListJenisByBMT(ctx context.Context, bmtID uuid.UUID) ([]*rekening.JenisRekening, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, kode, nama, tipe_dasar, akad, deskripsi, setoran_awal_min,
		setoran_min, bisa_ditarik, syarat_penarikan, nisbah_nasabah, jangka_hari, biaya_admin_bulanan,
		bisa_nfc, bisa_autodebet, biaya_admin_buka, is_aktif, urutan_tampil, created_at, updated_at, created_by, updated_by
		FROM jenis_rekening WHERE bmt_id = $1 AND is_aktif = true ORDER BY urutan_tampil
	`, bmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*rekening.JenisRekening
	for rows.Next() {
		jr := &rekening.JenisRekening{}
		err := rows.Scan(&jr.ID, &jr.BMTID, &jr.Kode, &jr.Nama, &jr.TipeDasar, &jr.Akad, &jr.Deskripsi,
			&jr.SetoranAwalMin, &jr.SetoranMin, &jr.BisaDitarik, &jr.SyaratPenarikan,
			&jr.NisbahNasabah, &jr.JangkaHari, &jr.BiayaAdminBulanan, &jr.BisaNFC,
			&jr.BisaAutodebet, &jr.BiayaAdminBuka, &jr.IsAktif, &jr.UrutanTampil,
			&jr.CreatedAt, &jr.UpdatedAt, &jr.CreatedBy, &jr.UpdatedBy)
		if err != nil {
			return nil, err
		}
		result = append(result, jr)
	}
	return result, nil
}

func (r *RekeningRepository) UpdateJenis(ctx context.Context, jr *rekening.JenisRekening) error {
	_, err := r.db.Exec(ctx, `
		UPDATE jenis_rekening SET nama=$1, deskripsi=$2, setoran_awal_min=$3, setoran_min=$4,
		biaya_admin_bulanan=$5, biaya_admin_buka=$6, bisa_ditarik=$7, bisa_nfc=$8,
		bisa_autodebet=$9, is_aktif=$10, urutan_tampil=$11, updated_at=NOW(), updated_by=$12
		WHERE id=$13
	`, jr.Nama, jr.Deskripsi, jr.SetoranAwalMin, jr.SetoranMin, jr.BiayaAdminBulanan,
		jr.BiayaAdminBuka, jr.BisaDitarik, jr.BisaNFC, jr.BisaAutodebet, jr.IsAktif,
		jr.UrutanTampil, jr.UpdatedBy, jr.ID)
	return err
}

func (r *RekeningRepository) ListDepositoAktif(ctx context.Context, bmtID uuid.UUID) ([]*rekening.Rekening, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.bmt_id, r.cabang_id, r.nasabah_id, r.jenis_rekening_id,
		r.nomor_rekening, r.saldo, r.status, r.alasan_blokir, r.biaya_admin_bulanan,
		r.nominal_deposito, r.nisbah_nasabah, r.tanggal_buka, r.tanggal_jatuh_tempo,
		r.tanggal_tutup, r.created_at, r.updated_at, r.created_by_form_id
		FROM rekening r
		JOIN jenis_rekening jr ON jr.id = r.jenis_rekening_id
		WHERE r.bmt_id = $1 AND r.status = 'AKTIF' AND jr.tipe_dasar = 'DEPOSITO'
	`, bmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*rekening.Rekening
	for rows.Next() {
		rek, err := r.scanRekening(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, rek)
	}
	return result, nil
}

func (r *RekeningRepository) scanRekening(s scanner) (*rekening.Rekening, error) {
	rek := &rekening.Rekening{}
	var saldo int64
	err := s.Scan(&rek.ID, &rek.BMTID, &rek.CabangID, &rek.NasabahID, &rek.JenisRekeningID,
		&rek.NomorRekening, &saldo, &rek.Status, &rek.AlasanBlokir, &rek.BiayaAdminBulanan,
		&rek.NominalDeposito, &rek.NisbahNasabah, &rek.TanggalBuka, &rek.TanggalJatuhTempo,
		&rek.TanggalTutup, &rek.CreatedAt, &rek.UpdatedAt, &rek.CreatedByFormID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, rekening.ErrRekeningNotFound
		}
		return nil, err
	}
	rek.Saldo = money.New(saldo)
	return rek, nil
}
