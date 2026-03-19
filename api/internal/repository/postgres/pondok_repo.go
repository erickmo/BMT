package postgres

import (
	"context"
	"time"

	"github.com/bmt-saas/api/internal/domain/pondok/administrasi"
	"github.com/bmt-saas/api/internal/domain/pondok/keuangan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ── Santri Repository ─────────────────────────────────────────────────────────

type SantriRepository struct {
	db *pgxpool.Pool
}

func NewSantriRepository(db *pgxpool.Pool) *SantriRepository {
	return &SantriRepository{db: db}
}

func (r *SantriRepository) Create(ctx context.Context, s *administrasi.Santri) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pondok_santri (id, bmt_id, cabang_id, nomor_induk_santri, nama_lengkap, nasabah_id,
		tingkat, kelas_id, asrama, kamar, angkatan, status_aktif, tanggal_masuk, foto_url,
		nama_wali, telepon_wali, nasabah_wali_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
	`, s.ID, s.BMTID, s.CabangID, s.NomorIndukSantri, s.NamaLengkap, s.NasabahID,
		s.Tingkat, s.KelasID, s.Asrama, s.Kamar, s.Angkatan, s.StatusAktif, s.TanggalMasuk, s.FotoURL,
		s.NamaWali, s.TeleponWali, s.NasabahWaliID, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SantriRepository) GetByID(ctx context.Context, id uuid.UUID) (*administrasi.Santri, error) {
	return r.scanSantri(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_induk_santri, nama_lengkap, nasabah_id,
		tingkat, kelas_id, asrama, kamar, angkatan, status_aktif, tanggal_masuk, tanggal_keluar, foto_url,
		nama_wali, telepon_wali, nasabah_wali_id, created_at, updated_at
		FROM pondok_santri WHERE id = $1
	`, id))
}

func (r *SantriRepository) GetByNIS(ctx context.Context, bmtID uuid.UUID, nis string) (*administrasi.Santri, error) {
	return r.scanSantri(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_induk_santri, nama_lengkap, nasabah_id,
		tingkat, kelas_id, asrama, kamar, angkatan, status_aktif, tanggal_masuk, tanggal_keluar, foto_url,
		nama_wali, telepon_wali, nasabah_wali_id, created_at, updated_at
		FROM pondok_santri WHERE bmt_id = $1 AND nomor_induk_santri = $2
	`, bmtID, nis))
}

func (r *SantriRepository) List(ctx context.Context, filter administrasi.ListSantriFilter) ([]*administrasi.Santri, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nomor_induk_santri, nama_lengkap, nasabah_id,
		tingkat, kelas_id, asrama, kamar, angkatan, status_aktif, tanggal_masuk, tanggal_keluar, foto_url,
		nama_wali, telepon_wali, nasabah_wali_id, created_at, updated_at
		FROM pondok_santri
		WHERE bmt_id = $1 AND cabang_id = $2
		  AND ($3::uuid IS NULL OR kelas_id = $3)
		  AND ($4 = '' OR tingkat = $4)
		  AND ($5::bool IS NULL OR status_aktif = $5)
		  AND ($6 = '' OR nama_lengkap ILIKE '%' || $6 || '%' OR nomor_induk_santri ILIKE '%' || $6 || '%')
		ORDER BY nama_lengkap LIMIT $7 OFFSET $8
	`, filter.BMTID, filter.CabangID, filter.KelasID, filter.Tingkat, filter.StatusAktif,
		filter.Keyword, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*administrasi.Santri
	for rows.Next() {
		s, err := r.scanSantri(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, s)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pondok_santri WHERE bmt_id = $1 AND cabang_id = $2`, filter.BMTID, filter.CabangID).Scan(&total)

	return result, total, nil
}

func (r *SantriRepository) Update(ctx context.Context, s *administrasi.Santri) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pondok_santri SET nama_lengkap=$1, tingkat=$2, kelas_id=$3, asrama=$4, kamar=$5,
		status_aktif=$6, tanggal_keluar=$7, foto_url=$8, nama_wali=$9, telepon_wali=$10, updated_at=NOW()
		WHERE id=$11
	`, s.NamaLengkap, s.Tingkat, s.KelasID, s.Asrama, s.Kamar,
		s.StatusAktif, s.TanggalKeluar, s.FotoURL, s.NamaWali, s.TeleponWali, s.ID)
	return err
}

func (r *SantriRepository) UpdateKelas(ctx context.Context, santriID, kelasID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE pondok_santri SET kelas_id=$1, updated_at=NOW() WHERE id=$2`, kelasID, santriID)
	return err
}

func (r *SantriRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pondok_santri WHERE id=$1`, id)
	return err
}

func (r *SantriRepository) scanSantri(s scanner) (*administrasi.Santri, error) {
	santri := &administrasi.Santri{}
	err := s.Scan(&santri.ID, &santri.BMTID, &santri.CabangID, &santri.NomorIndukSantri, &santri.NamaLengkap,
		&santri.NasabahID, &santri.Tingkat, &santri.KelasID, &santri.Asrama, &santri.Kamar, &santri.Angkatan,
		&santri.StatusAktif, &santri.TanggalMasuk, &santri.TanggalKeluar, &santri.FotoURL,
		&santri.NamaWali, &santri.TeleponWali, &santri.NasabahWaliID, &santri.CreatedAt, &santri.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, administrasi.ErrSantriNotFound
		}
		return nil, err
	}
	return santri, nil
}

// ── Kelas Repository ──────────────────────────────────────────────────────────

type KelasRepository struct {
	db *pgxpool.Pool
}

func NewKelasRepository(db *pgxpool.Pool) *KelasRepository {
	return &KelasRepository{db: db}
}

func (r *KelasRepository) Create(ctx context.Context, k *administrasi.Kelas) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pondok_kelas (id, bmt_id, cabang_id, nama, tingkat, tahun_ajaran, wali_kelas_id, kapasitas, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, k.ID, k.BMTID, k.CabangID, k.Nama, k.Tingkat, k.TahunAjaran, k.WaliKelasID, k.Kapasitas, k.CreatedAt, k.UpdatedAt)
	return err
}

func (r *KelasRepository) GetByID(ctx context.Context, id uuid.UUID) (*administrasi.Kelas, error) {
	k := &administrasi.Kelas{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nama, tingkat, tahun_ajaran, wali_kelas_id, kapasitas, created_at, updated_at
		FROM pondok_kelas WHERE id = $1
	`, id).Scan(&k.ID, &k.BMTID, &k.CabangID, &k.Nama, &k.Tingkat, &k.TahunAjaran,
		&k.WaliKelasID, &k.Kapasitas, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, administrasi.ErrKelasNotFound
		}
		return nil, err
	}
	return k, nil
}

func (r *KelasRepository) List(ctx context.Context, bmtID, cabangID uuid.UUID, tahunAjaran string) ([]*administrasi.Kelas, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nama, tingkat, tahun_ajaran, wali_kelas_id, kapasitas, created_at, updated_at
		FROM pondok_kelas WHERE bmt_id = $1 AND cabang_id = $2 AND ($3 = '' OR tahun_ajaran = $3)
		ORDER BY nama
	`, bmtID, cabangID, tahunAjaran)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*administrasi.Kelas
	for rows.Next() {
		k := &administrasi.Kelas{}
		err := rows.Scan(&k.ID, &k.BMTID, &k.CabangID, &k.Nama, &k.Tingkat, &k.TahunAjaran,
			&k.WaliKelasID, &k.Kapasitas, &k.CreatedAt, &k.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	return result, nil
}

func (r *KelasRepository) Update(ctx context.Context, k *administrasi.Kelas) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pondok_kelas SET nama=$1, tingkat=$2, wali_kelas_id=$3, kapasitas=$4, updated_at=NOW()
		WHERE id=$5
	`, k.Nama, k.Tingkat, k.WaliKelasID, k.Kapasitas, k.ID)
	return err
}

func (r *KelasRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pondok_kelas WHERE id=$1`, id)
	return err
}

func (r *KelasRepository) CountSantri(ctx context.Context, kelasID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pondok_santri WHERE kelas_id = $1 AND status_aktif = true`, kelasID).Scan(&count)
	return count, err
}

// ── Pengajar Repository ───────────────────────────────────────────────────────

type PengajarRepository struct {
	db *pgxpool.Pool
}

func NewPengajarRepository(db *pgxpool.Pool) *PengajarRepository {
	return &PengajarRepository{db: db}
}

func (r *PengajarRepository) Create(ctx context.Context, p *administrasi.Pengajar) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pondok_pengajar (id, bmt_id, cabang_id, nip, nama_lengkap, jabatan, spesialisasi, nasabah_id, status_aktif, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, p.ID, p.BMTID, p.CabangID, p.NIP, p.NamaLengkap, p.Jabatan, p.Spesialisasi, p.NasabahID, p.StatusAktif, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PengajarRepository) GetByID(ctx context.Context, id uuid.UUID) (*administrasi.Pengajar, error) {
	p := &administrasi.Pengajar{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nip, nama_lengkap, jabatan, spesialisasi, nasabah_id, status_aktif, created_at, updated_at
		FROM pondok_pengajar WHERE id = $1
	`, id).Scan(&p.ID, &p.BMTID, &p.CabangID, &p.NIP, &p.NamaLengkap, &p.Jabatan,
		&p.Spesialisasi, &p.NasabahID, &p.StatusAktif, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, administrasi.ErrPengajarNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *PengajarRepository) GetByNIP(ctx context.Context, bmtID uuid.UUID, nip string) (*administrasi.Pengajar, error) {
	p := &administrasi.Pengajar{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nip, nama_lengkap, jabatan, spesialisasi, nasabah_id, status_aktif, created_at, updated_at
		FROM pondok_pengajar WHERE bmt_id = $1 AND nip = $2
	`, bmtID, nip).Scan(&p.ID, &p.BMTID, &p.CabangID, &p.NIP, &p.NamaLengkap, &p.Jabatan,
		&p.Spesialisasi, &p.NasabahID, &p.StatusAktif, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, administrasi.ErrPengajarNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *PengajarRepository) List(ctx context.Context, filter administrasi.ListPengajarFilter) ([]*administrasi.Pengajar, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nip, nama_lengkap, jabatan, spesialisasi, nasabah_id, status_aktif, created_at, updated_at
		FROM pondok_pengajar
		WHERE bmt_id = $1 AND cabang_id = $2
		  AND ($3::bool IS NULL OR status_aktif = $3)
		ORDER BY nama_lengkap LIMIT $4 OFFSET $5
	`, filter.BMTID, filter.CabangID, filter.StatusAktif, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*administrasi.Pengajar
	for rows.Next() {
		p := &administrasi.Pengajar{}
		err := rows.Scan(&p.ID, &p.BMTID, &p.CabangID, &p.NIP, &p.NamaLengkap, &p.Jabatan,
			&p.Spesialisasi, &p.NasabahID, &p.StatusAktif, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pondok_pengajar WHERE bmt_id = $1 AND cabang_id = $2`, filter.BMTID, filter.CabangID).Scan(&total)

	return result, total, nil
}

func (r *PengajarRepository) Update(ctx context.Context, p *administrasi.Pengajar) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pondok_pengajar SET nama_lengkap=$1, jabatan=$2, spesialisasi=$3, status_aktif=$4, updated_at=NOW()
		WHERE id=$5
	`, p.NamaLengkap, p.Jabatan, p.Spesialisasi, p.StatusAktif, p.ID)
	return err
}

func (r *PengajarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pondok_pengajar WHERE id=$1`, id)
	return err
}

// ── JenisTagihan Repository ───────────────────────────────────────────────────

type JenisTagihanRepository struct {
	db *pgxpool.Pool
}

func NewJenisTagihanRepository(db *pgxpool.Pool) *JenisTagihanRepository {
	return &JenisTagihanRepository{db: db}
}

func (r *JenisTagihanRepository) Create(ctx context.Context, j *keuangan.JenisTagihan) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pondok_jenis_tagihan (id, bmt_id, kode, nama, nominal, frekuensi, is_aktif, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, j.ID, j.BMTID, j.Kode, j.Nama, j.Nominal, j.Frekuensi, j.IsAktif, j.CreatedAt, j.UpdatedAt)
	return err
}

func (r *JenisTagihanRepository) GetByID(ctx context.Context, id uuid.UUID) (*keuangan.JenisTagihan, error) {
	j := &keuangan.JenisTagihan{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, kode, nama, nominal, frekuensi, is_aktif, created_at, updated_at
		FROM pondok_jenis_tagihan WHERE id = $1
	`, id).Scan(&j.ID, &j.BMTID, &j.Kode, &j.Nama, &j.Nominal, &j.Frekuensi, &j.IsAktif, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, keuangan.ErrJenisTagihanNotFound
		}
		return nil, err
	}
	return j, nil
}

func (r *JenisTagihanRepository) GetByKode(ctx context.Context, bmtID uuid.UUID, kode string) (*keuangan.JenisTagihan, error) {
	j := &keuangan.JenisTagihan{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, kode, nama, nominal, frekuensi, is_aktif, created_at, updated_at
		FROM pondok_jenis_tagihan WHERE bmt_id = $1 AND kode = $2
	`, bmtID, kode).Scan(&j.ID, &j.BMTID, &j.Kode, &j.Nama, &j.Nominal, &j.Frekuensi, &j.IsAktif, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, keuangan.ErrJenisTagihanNotFound
		}
		return nil, err
	}
	return j, nil
}

func (r *JenisTagihanRepository) List(ctx context.Context, bmtID uuid.UUID, aktifSaja bool) ([]*keuangan.JenisTagihan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, kode, nama, nominal, frekuensi, is_aktif, created_at, updated_at
		FROM pondok_jenis_tagihan WHERE bmt_id = $1 AND ($2 = false OR is_aktif = true)
		ORDER BY nama
	`, bmtID, aktifSaja)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*keuangan.JenisTagihan
	for rows.Next() {
		j := &keuangan.JenisTagihan{}
		err := rows.Scan(&j.ID, &j.BMTID, &j.Kode, &j.Nama, &j.Nominal, &j.Frekuensi, &j.IsAktif, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, j)
	}
	return result, nil
}

func (r *JenisTagihanRepository) Update(ctx context.Context, j *keuangan.JenisTagihan) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pondok_jenis_tagihan SET nama=$1, nominal=$2, frekuensi=$3, is_aktif=$4, updated_at=NOW()
		WHERE id=$5
	`, j.Nama, j.Nominal, j.Frekuensi, j.IsAktif, j.ID)
	return err
}

func (r *JenisTagihanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pondok_jenis_tagihan WHERE id=$1`, id)
	return err
}

// ── TagihanSPP Repository ─────────────────────────────────────────────────────

type TagihanSPPRepository struct {
	db *pgxpool.Pool
}

func NewTagihanSPPRepository(db *pgxpool.Pool) *TagihanSPPRepository {
	return &TagihanSPPRepository{db: db}
}

func (r *TagihanSPPRepository) Create(ctx context.Context, t *keuangan.TagihanSPP) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tagihan_spp (id, bmt_id, cabang_id, santri_id, jenis_tagihan_id, periode,
		nominal, nominal_terbayar, nominal_sisa, beasiswa_persen, beasiswa_nominal, nominal_efektif,
		status, tanggal_jatuh_tempo, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`, t.ID, t.BMTID, t.CabangID, t.SantriID, t.JenisTagihanID, t.Periode,
		t.Nominal, t.NominalTerbayar, t.NominalSisa, t.BeasiswaPersen, t.BeasiswaNominal, t.NominalEfektif,
		t.Status, t.TanggalJatuhTempo, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *TagihanSPPRepository) GetByID(ctx context.Context, id uuid.UUID) (*keuangan.TagihanSPP, error) {
	return r.scanTagihan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, santri_id, jenis_tagihan_id, periode,
		nominal, nominal_terbayar, nominal_sisa, beasiswa_persen, beasiswa_nominal, nominal_efektif,
		status, tanggal_jatuh_tempo, tanggal_lunas, created_at, updated_at
		FROM tagihan_spp WHERE id = $1
	`, id))
}

func (r *TagihanSPPRepository) GetBySantriPeriode(ctx context.Context, santriID uuid.UUID, periode string) (*keuangan.TagihanSPP, error) {
	return r.scanTagihan(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, santri_id, jenis_tagihan_id, periode,
		nominal, nominal_terbayar, nominal_sisa, beasiswa_persen, beasiswa_nominal, nominal_efektif,
		status, tanggal_jatuh_tempo, tanggal_lunas, created_at, updated_at
		FROM tagihan_spp WHERE santri_id = $1 AND periode = $2
	`, santriID, periode))
}

func (r *TagihanSPPRepository) List(ctx context.Context, filter keuangan.ListTagihanFilter) ([]*keuangan.TagihanSPP, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, santri_id, jenis_tagihan_id, periode,
		nominal, nominal_terbayar, nominal_sisa, beasiswa_persen, beasiswa_nominal, nominal_efektif,
		status, tanggal_jatuh_tempo, tanggal_lunas, created_at, updated_at
		FROM tagihan_spp
		WHERE bmt_id = $1 AND cabang_id = $2
		  AND ($3::uuid IS NULL OR santri_id = $3)
		  AND ($4 = '' OR periode = $4)
		  AND ($5 = '' OR status = $5)
		ORDER BY tanggal_jatuh_tempo DESC LIMIT $6 OFFSET $7
	`, filter.BMTID, filter.CabangID, filter.SantriID, filter.Periode, filter.Status, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*keuangan.TagihanSPP
	for rows.Next() {
		t, err := r.scanTagihan(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM tagihan_spp WHERE bmt_id = $1 AND cabang_id = $2`, filter.BMTID, filter.CabangID).Scan(&total)

	return result, total, nil
}

func (r *TagihanSPPRepository) ListBelumLunas(ctx context.Context, bmtID uuid.UUID, tanggalJatuhTempo time.Time) ([]*keuangan.TagihanSPP, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, santri_id, jenis_tagihan_id, periode,
		nominal, nominal_terbayar, nominal_sisa, beasiswa_persen, beasiswa_nominal, nominal_efektif,
		status, tanggal_jatuh_tempo, tanggal_lunas, created_at, updated_at
		FROM tagihan_spp
		WHERE bmt_id = $1 AND status != 'LUNAS' AND tanggal_jatuh_tempo <= $2
	`, bmtID, tanggalJatuhTempo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*keuangan.TagihanSPP
	for rows.Next() {
		t, err := r.scanTagihan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TagihanSPPRepository) Update(ctx context.Context, t *keuangan.TagihanSPP) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tagihan_spp SET nominal_terbayar=$1, nominal_sisa=$2, status=$3, tanggal_lunas=$4, updated_at=NOW()
		WHERE id=$5
	`, t.NominalTerbayar, t.NominalSisa, t.Status, t.TanggalLunas, t.ID)
	return err
}

func (r *TagihanSPPRepository) UpdateBeasiswa(ctx context.Context, id uuid.UUID, persen float64, nominal, efektif, sisa int64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tagihan_spp SET beasiswa_persen=$1, beasiswa_nominal=$2, nominal_efektif=$3, nominal_sisa=$4, updated_at=NOW()
		WHERE id=$5
	`, persen, nominal, efektif, sisa, id)
	return err
}

func (r *TagihanSPPRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tagihan_spp WHERE id=$1`, id)
	return err
}

func (r *TagihanSPPRepository) scanTagihan(s scanner) (*keuangan.TagihanSPP, error) {
	t := &keuangan.TagihanSPP{}
	err := s.Scan(&t.ID, &t.BMTID, &t.CabangID, &t.SantriID, &t.JenisTagihanID, &t.Periode,
		&t.Nominal, &t.NominalTerbayar, &t.NominalSisa, &t.BeasiswaPersen, &t.BeasiswaNominal, &t.NominalEfektif,
		&t.Status, &t.TanggalJatuhTempo, &t.TanggalLunas, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, keuangan.ErrTagihanSPPNotFound
		}
		return nil, err
	}
	return t, nil
}
