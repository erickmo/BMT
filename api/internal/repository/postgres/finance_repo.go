package postgres

import (
	"context"
	"time"

	"github.com/bmt-saas/api/internal/domain/finance"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FinanceRepository mengimplementasikan finance.Repository menggunakan PostgreSQL.
type FinanceRepository struct {
	db *pgxpool.Pool
}

// NewFinanceRepository membuat instance baru FinanceRepository.
func NewFinanceRepository(db *pgxpool.Pool) *FinanceRepository {
	return &FinanceRepository{db: db}
}

// ─── Jurnal Manual ───────────────────────────────────────────────────────────

// CreateJurnal menyimpan jurnal manual beserta seluruh entri-nya.
func (r *FinanceRepository) CreateJurnal(ctx context.Context, j *finance.JurnalManual) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO jurnal_manual
		  (id, bmt_id, cabang_id, tanggal, keterangan, referensi, status,
		   dibuat_oleh, disetujui_oleh, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, j.ID, j.BMTID, j.CabangID, j.Tanggal, j.Keterangan, j.Referensi,
		j.Status, j.DibuatOleh, j.DisetujuiOleh, j.CreatedAt, j.UpdatedAt)
	if err != nil {
		return err
	}

	for _, e := range j.Entries {
		_, err = r.db.Exec(ctx, `
			INSERT INTO jurnal_manual_entri
			  (id, jurnal_id, kode_akun, nama_akun, posisi, nominal)
			VALUES ($1,$2,$3,$4,$5,$6)
		`, e.ID, e.JurnalID, e.KodeAkun, e.NamaAkun, e.Posisi, e.Nominal)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetJurnalByID mengambil satu jurnal manual beserta entri-nya berdasarkan ID.
func (r *FinanceRepository) GetJurnalByID(ctx context.Context, id uuid.UUID) (*finance.JurnalManual, error) {
	j := &finance.JurnalManual{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, tanggal, keterangan, referensi,
		       status, dibuat_oleh, disetujui_oleh, created_at, updated_at
		FROM jurnal_manual
		WHERE id = $1
	`, id).Scan(
		&j.ID, &j.BMTID, &j.CabangID, &j.Tanggal, &j.Keterangan, &j.Referensi,
		&j.Status, &j.DibuatOleh, &j.DisetujuiOleh, &j.CreatedAt, &j.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, finance.ErrJurnalNotFound
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, jurnal_id, kode_akun, nama_akun, posisi, nominal
		FROM jurnal_manual_entri
		WHERE jurnal_id = $1
		ORDER BY id
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := &finance.EntriJurnal{}
		if err := rows.Scan(&e.ID, &e.JurnalID, &e.KodeAkun, &e.NamaAkun, &e.Posisi, &e.Nominal); err != nil {
			return nil, err
		}
		j.Entries = append(j.Entries, e)
	}
	return j, nil
}

// ListJurnal mengembalikan daftar jurnal manual dengan filter opsional.
func (r *FinanceRepository) ListJurnal(ctx context.Context, filter finance.ListJurnalFilter) ([]*finance.JurnalManual, int64, error) {
	offset := 0
	if filter.Page > 1 {
		offset = (filter.Page - 1) * filter.PerPage
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, tanggal, keterangan, referensi,
		       status, dibuat_oleh, disetujui_oleh, created_at, updated_at
		FROM jurnal_manual
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR cabang_id = $2)
		  AND ($3::timestamptz IS NULL OR tanggal >= $3)
		  AND ($4::timestamptz IS NULL OR tanggal <= $4)
		  AND ($5::text IS NULL OR status = $5)
		ORDER BY tanggal DESC
		LIMIT $6 OFFSET $7
	`, filter.BMTID, filter.CabangID, filter.TanggalDari, filter.TanggalSampai,
		filter.Status, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*finance.JurnalManual
	for rows.Next() {
		j := &finance.JurnalManual{}
		if err := rows.Scan(
			&j.ID, &j.BMTID, &j.CabangID, &j.Tanggal, &j.Keterangan, &j.Referensi,
			&j.Status, &j.DibuatOleh, &j.DisetujuiOleh, &j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, j)
	}

	var total int64
	r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM jurnal_manual
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR cabang_id = $2)
		  AND ($3::timestamptz IS NULL OR tanggal >= $3)
		  AND ($4::timestamptz IS NULL OR tanggal <= $4)
		  AND ($5::text IS NULL OR status = $5)
	`, filter.BMTID, filter.CabangID, filter.TanggalDari, filter.TanggalSampai, filter.Status).Scan(&total)

	return result, total, nil
}

// PostJurnal mengubah status jurnal menjadi POSTED dan menyimpan siapa yang menyetujui.
func (r *FinanceRepository) PostJurnal(ctx context.Context, id uuid.UUID, disetujuiOleh uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE jurnal_manual
		SET status = 'POSTED', disetujui_oleh = $2, updated_at = NOW()
		WHERE id = $1
	`, id, disetujuiOleh)
	return err
}

// VoidJurnal mengubah status jurnal menjadi VOID.
func (r *FinanceRepository) VoidJurnal(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE jurnal_manual
		SET status = 'VOID', updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// ─── Vendor ──────────────────────────────────────────────────────────────────

// CreateVendor menyimpan data vendor baru.
func (r *FinanceRepository) CreateVendor(ctx context.Context, v *finance.Vendor) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO vendor
		  (id, bmt_id, nama, npwp, alamat, telepon, email,
		   rekening_bank, nama_bank, is_aktif, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, v.ID, v.BMTID, v.Nama, v.NPWP, v.Alamat, v.Telepon, v.Email,
		v.RekeningBank, v.NamaBank, v.IsAktif, v.CreatedAt, v.UpdatedAt)
	return err
}

// GetVendorByID mengambil satu vendor berdasarkan ID.
func (r *FinanceRepository) GetVendorByID(ctx context.Context, id uuid.UUID) (*finance.Vendor, error) {
	v := &finance.Vendor{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, nama, npwp, alamat, telepon, email,
		       rekening_bank, nama_bank, is_aktif, created_at, updated_at
		FROM vendor
		WHERE id = $1
	`, id).Scan(
		&v.ID, &v.BMTID, &v.Nama, &v.NPWP, &v.Alamat, &v.Telepon, &v.Email,
		&v.RekeningBank, &v.NamaBank, &v.IsAktif, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, finance.ErrVendorNotFound
		}
		return nil, err
	}
	return v, nil
}

// ListVendor mengembalikan daftar vendor milik BMT tertentu dengan paginasi.
func (r *FinanceRepository) ListVendor(ctx context.Context, bmtID uuid.UUID, page, perPage int) ([]*finance.Vendor, int64, error) {
	offset := 0
	if page > 1 {
		offset = (page - 1) * perPage
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, nama, npwp, alamat, telepon, email,
		       rekening_bank, nama_bank, is_aktif, created_at, updated_at
		FROM vendor
		WHERE bmt_id = $1
		ORDER BY nama ASC
		LIMIT $2 OFFSET $3
	`, bmtID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*finance.Vendor
	for rows.Next() {
		v := &finance.Vendor{}
		if err := rows.Scan(
			&v.ID, &v.BMTID, &v.Nama, &v.NPWP, &v.Alamat, &v.Telepon, &v.Email,
			&v.RekeningBank, &v.NamaBank, &v.IsAktif, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, v)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM vendor WHERE bmt_id = $1`, bmtID).Scan(&total)

	return result, total, nil
}

// UpdateVendor memperbarui data vendor.
func (r *FinanceRepository) UpdateVendor(ctx context.Context, v *finance.Vendor) error {
	_, err := r.db.Exec(ctx, `
		UPDATE vendor
		SET nama = $2, npwp = $3, alamat = $4, telepon = $5, email = $6,
		    rekening_bank = $7, nama_bank = $8, is_aktif = $9, updated_at = NOW()
		WHERE id = $1
	`, v.ID, v.Nama, v.NPWP, v.Alamat, v.Telepon, v.Email,
		v.RekeningBank, v.NamaBank, v.IsAktif)
	return err
}

// ─── Transaksi Operasional ───────────────────────────────────────────────────

// CreateTransaksiOperasional menyimpan satu transaksi operasional.
func (r *FinanceRepository) CreateTransaksiOperasional(ctx context.Context, t *finance.TransaksiOperasional) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transaksi_operasional
		  (id, bmt_id, cabang_id, vendor_id, tanggal, jenis, kategori, keterangan,
		   nominal, kode_akun_debit, kode_akun_kredit, lampiran, jurnal_id, dibuat_oleh, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, t.ID, t.BMTID, t.CabangID, t.VendorID, t.Tanggal, t.Jenis, t.Kategori, t.Keterangan,
		t.Nominal, t.KodeAkunDebit, t.KodeAkunKredit, t.Lampiran, t.JurnalID, t.DibuatOleh, t.CreatedAt)
	return err
}

// GetTransaksiOperasionalByID mengambil satu transaksi operasional berdasarkan ID.
func (r *FinanceRepository) GetTransaksiOperasionalByID(ctx context.Context, id uuid.UUID) (*finance.TransaksiOperasional, error) {
	t := &finance.TransaksiOperasional{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, vendor_id, tanggal, jenis, kategori, keterangan,
		       nominal, kode_akun_debit, kode_akun_kredit, lampiran, jurnal_id, dibuat_oleh, created_at
		FROM transaksi_operasional
		WHERE id = $1
	`, id).Scan(
		&t.ID, &t.BMTID, &t.CabangID, &t.VendorID, &t.Tanggal, &t.Jenis, &t.Kategori, &t.Keterangan,
		&t.Nominal, &t.KodeAkunDebit, &t.KodeAkunKredit, &t.Lampiran, &t.JurnalID, &t.DibuatOleh, &t.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, finance.ErrTransaksiOperasionalNotFound
		}
		return nil, err
	}
	return t, nil
}

// ListTransaksiOperasional mengembalikan daftar transaksi operasional dengan filter rentang tanggal.
func (r *FinanceRepository) ListTransaksiOperasional(
	ctx context.Context,
	bmtID, cabangID uuid.UUID,
	dari, sampai time.Time,
	page, perPage int,
) ([]*finance.TransaksiOperasional, int64, error) {
	offset := 0
	if page > 1 {
		offset = (page - 1) * perPage
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, vendor_id, tanggal, jenis, kategori, keterangan,
		       nominal, kode_akun_debit, kode_akun_kredit, lampiran, jurnal_id, dibuat_oleh, created_at
		FROM transaksi_operasional
		WHERE bmt_id = $1
		  AND cabang_id = $2
		  AND tanggal >= $3
		  AND tanggal <= $4
		ORDER BY tanggal DESC
		LIMIT $5 OFFSET $6
	`, bmtID, cabangID, dari, sampai, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*finance.TransaksiOperasional
	for rows.Next() {
		t := &finance.TransaksiOperasional{}
		if err := rows.Scan(
			&t.ID, &t.BMTID, &t.CabangID, &t.VendorID, &t.Tanggal, &t.Jenis, &t.Kategori, &t.Keterangan,
			&t.Nominal, &t.KodeAkunDebit, &t.KodeAkunKredit, &t.Lampiran, &t.JurnalID, &t.DibuatOleh, &t.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}

	var total int64
	r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM transaksi_operasional
		WHERE bmt_id = $1 AND cabang_id = $2 AND tanggal >= $3 AND tanggal <= $4
	`, bmtID, cabangID, dari, sampai).Scan(&total)

	return result, total, nil
}

// ─── Laporan (non-interface, dipakai service) ────────────────────────────────

// SumByAkunPrefix menghitung total debit dan kredit per prefix kode akun (1 karakter)
// untuk rentang tanggal tertentu — digunakan untuk kalkulasi neraca sederhana.
// Hasil map berisi kunci seperti "1_debit", "2_kredit", dst.
func (r *FinanceRepository) SumByAkunPrefix(
	ctx context.Context,
	bmtID uuid.UUID,
	dari, sampai time.Time,
) (map[string]int64, error) {
	rows, err := r.db.Query(ctx, `
		SELECT LEFT(e.kode_akun, 1)   AS prefix,
		       e.posisi,
		       COALESCE(SUM(e.nominal), 0) AS total
		FROM jurnal_manual_entri e
		JOIN jurnal_manual j ON j.id = e.jurnal_id
		WHERE j.bmt_id = $1
		  AND j.tanggal BETWEEN $2 AND $3
		  AND j.status = 'POSTED'
		GROUP BY LEFT(e.kode_akun, 1), e.posisi
	`, bmtID, dari, sampai)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var prefix, posisi string
		var total int64
		if err := rows.Scan(&prefix, &posisi, &total); err != nil {
			return nil, err
		}
		key := prefix + "_" + posisi // contoh: "1_DEBIT", "2_KREDIT"
		result[key] += total
	}
	return result, nil
}
