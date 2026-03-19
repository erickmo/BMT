package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/ecommerce/pesanan"
	"github.com/bmt-saas/api/internal/domain/ecommerce/produk"
	"github.com/bmt-saas/api/internal/domain/ecommerce/toko"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ── Toko Repository ───────────────────────────────────────────────────────────

type TokoRepository struct {
	db *pgxpool.Pool
}

func NewTokoRepository(db *pgxpool.Pool) *TokoRepository {
	return &TokoRepository{db: db}
}

func (r *TokoRepository) Create(ctx context.Context, t *toko.Toko) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO toko (id, bmt_id, cabang_id, nama, slug, deskripsi, logo_url, banner_url, kategori_toko, is_opop, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`, t.ID, t.BMTID, t.CabangID, t.Nama, t.Slug, t.Deskripsi, t.LogoURL, t.BannerURL,
		t.KategoriToko, t.IsOPOP, t.Status, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *TokoRepository) GetByID(ctx context.Context, id uuid.UUID) (*toko.Toko, error) {
	return r.scanToko(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nama, slug, deskripsi, logo_url, banner_url, kategori_toko, is_opop, status, rating, created_at, updated_at
		FROM toko WHERE id = $1
	`, id))
}

func (r *TokoRepository) GetBySlug(ctx context.Context, slug string) (*toko.Toko, error) {
	return r.scanToko(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nama, slug, deskripsi, logo_url, banner_url, kategori_toko, is_opop, status, rating, created_at, updated_at
		FROM toko WHERE slug = $1
	`, slug))
}

func (r *TokoRepository) GetByCabang(ctx context.Context, bmtID, cabangID uuid.UUID) (*toko.Toko, error) {
	return r.scanToko(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, nama, slug, deskripsi, logo_url, banner_url, kategori_toko, is_opop, status, rating, created_at, updated_at
		FROM toko WHERE bmt_id = $1 AND cabang_id = $2 LIMIT 1
	`, bmtID, cabangID))
}

func (r *TokoRepository) List(ctx context.Context, filter toko.ListTokoFilter) ([]*toko.Toko, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nama, slug, deskripsi, logo_url, banner_url, kategori_toko, is_opop, status, rating, created_at, updated_at
		FROM toko
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR cabang_id = $2)
		  AND ($3::bool IS NULL OR is_opop = $3)
		  AND ($4::text IS NULL OR status = $4)
		ORDER BY nama LIMIT $5 OFFSET $6
	`, filter.BMTID, filter.CabangID, filter.IsOPOP, filter.Status, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*toko.Toko
	for rows.Next() {
		t, err := r.scanToko(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM toko`).Scan(&total)

	return result, total, nil
}

func (r *TokoRepository) ListOPOP(ctx context.Context, page, perPage int) ([]*toko.Toko, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, nama, slug, deskripsi, logo_url, banner_url, kategori_toko, is_opop, status, rating, created_at, updated_at
		FROM toko WHERE is_opop = true AND status = 'AKTIF'
		ORDER BY nama LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*toko.Toko
	for rows.Next() {
		t, err := r.scanToko(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM toko WHERE is_opop = true AND status = 'AKTIF'`).Scan(&total)

	return result, total, nil
}

func (r *TokoRepository) Update(ctx context.Context, t *toko.Toko) error {
	_, err := r.db.Exec(ctx, `
		UPDATE toko SET nama=$1, deskripsi=$2, logo_url=$3, banner_url=$4, is_opop=$5, status=$6, updated_at=NOW()
		WHERE id=$7
	`, t.Nama, t.Deskripsi, t.LogoURL, t.BannerURL, t.IsOPOP, t.Status, t.ID)
	return err
}

func (r *TokoRepository) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	_, err := r.db.Exec(ctx, `UPDATE toko SET rating=$1, updated_at=NOW() WHERE id=$2`, rating, id)
	return err
}

func (r *TokoRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status toko.StatusToko) error {
	_, err := r.db.Exec(ctx, `UPDATE toko SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *TokoRepository) scanToko(s scanner) (*toko.Toko, error) {
	t := &toko.Toko{}
	err := s.Scan(&t.ID, &t.BMTID, &t.CabangID, &t.Nama, &t.Slug, &t.Deskripsi,
		&t.LogoURL, &t.BannerURL, &t.KategoriToko, &t.IsOPOP, &t.Status, &t.Rating, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, toko.ErrTokoNotFound
		}
		return nil, err
	}
	return t, nil
}

// ── Produk Repository ─────────────────────────────────────────────────────────

type ProdukRepository struct {
	db *pgxpool.Pool
}

func NewProdukRepository(db *pgxpool.Pool) *ProdukRepository {
	return &ProdukRepository{db: db}
}

func (r *ProdukRepository) Create(ctx context.Context, p *produk.Produk) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO produk (id, bmt_id, toko_id, nama, slug, deskripsi, kategori, harga, harga_b2b,
		stok, satuan, berat_gram, foto_urls, is_opop, is_aktif, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`, p.ID, p.BMTID, p.TokoID, p.Nama, p.Slug, p.Deskripsi, p.Kategori, p.Harga, p.HargaB2B,
		p.Stok, p.Satuan, p.BeratGram, p.FotoURLs, p.IsOPOP, p.IsAktif, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *ProdukRepository) GetByID(ctx context.Context, id uuid.UUID) (*produk.Produk, error) {
	return r.scanProduk(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, toko_id, nama, slug, deskripsi, kategori, harga, harga_b2b,
		stok, satuan, berat_gram, foto_urls, is_opop, is_aktif, rating, total_terjual, created_at, updated_at
		FROM produk WHERE id = $1
	`, id))
}

func (r *ProdukRepository) GetBySlugAndToko(ctx context.Context, tokoID uuid.UUID, slug string) (*produk.Produk, error) {
	return r.scanProduk(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, toko_id, nama, slug, deskripsi, kategori, harga, harga_b2b,
		stok, satuan, berat_gram, foto_urls, is_opop, is_aktif, rating, total_terjual, created_at, updated_at
		FROM produk WHERE toko_id = $1 AND slug = $2
	`, tokoID, slug))
}

func (r *ProdukRepository) List(ctx context.Context, filter produk.ListProdukFilter) ([]*produk.Produk, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, toko_id, nama, slug, deskripsi, kategori, harga, harga_b2b,
		stok, satuan, berat_gram, foto_urls, is_opop, is_aktif, rating, total_terjual, created_at, updated_at
		FROM produk
		WHERE ($1::uuid IS NULL OR bmt_id = $1)
		  AND ($2::uuid IS NULL OR toko_id = $2)
		  AND ($3 = '' OR kategori = $3)
		  AND ($4::bool IS NULL OR is_opop = $4)
		  AND ($5::bool IS NULL OR is_aktif = $5)
		  AND ($6 = '' OR nama ILIKE '%' || $6 || '%')
		ORDER BY nama LIMIT $7 OFFSET $8
	`, filter.BMTID, filter.TokoID, filter.Kategori, filter.IsOPOP, filter.IsAktif, filter.Search, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*produk.Produk
	for rows.Next() {
		p, err := r.scanProduk(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM produk WHERE ($1::uuid IS NULL OR bmt_id = $1) AND ($2::uuid IS NULL OR toko_id = $2)`,
		filter.BMTID, filter.TokoID).Scan(&total)

	return result, total, nil
}

func (r *ProdukRepository) Update(ctx context.Context, p *produk.Produk) error {
	_, err := r.db.Exec(ctx, `
		UPDATE produk SET nama=$1, deskripsi=$2, kategori=$3, harga=$4, harga_b2b=$5,
		satuan=$6, berat_gram=$7, foto_urls=$8, is_opop=$9, is_aktif=$10, updated_at=NOW()
		WHERE id=$11
	`, p.Nama, p.Deskripsi, p.Kategori, p.Harga, p.HargaB2B, p.Satuan, p.BeratGram,
		p.FotoURLs, p.IsOPOP, p.IsAktif, p.ID)
	return err
}

func (r *ProdukRepository) UpdateStok(ctx context.Context, id uuid.UUID, stokBaru int) error {
	_, err := r.db.Exec(ctx, `UPDATE produk SET stok=$1, updated_at=NOW() WHERE id=$2`, stokBaru, id)
	return err
}

func (r *ProdukRepository) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	_, err := r.db.Exec(ctx, `UPDATE produk SET rating=$1, updated_at=NOW() WHERE id=$2`, rating, id)
	return err
}

func (r *ProdukRepository) IncrementTerjual(ctx context.Context, id uuid.UUID, jumlah int) error {
	_, err := r.db.Exec(ctx, `UPDATE produk SET total_terjual = total_terjual + $1, updated_at=NOW() WHERE id=$2`, jumlah, id)
	return err
}

func (r *ProdukRepository) LockForUpdate(ctx context.Context, id uuid.UUID) (*produk.Produk, error) {
	return r.scanProduk(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, toko_id, nama, slug, deskripsi, kategori, harga, harga_b2b,
		stok, satuan, berat_gram, foto_urls, is_opop, is_aktif, rating, total_terjual, created_at, updated_at
		FROM produk WHERE id = $1 FOR UPDATE
	`, id))
}

func (r *ProdukRepository) scanProduk(s scanner) (*produk.Produk, error) {
	p := &produk.Produk{}
	err := s.Scan(&p.ID, &p.BMTID, &p.TokoID, &p.Nama, &p.Slug, &p.Deskripsi, &p.Kategori,
		&p.Harga, &p.HargaB2B, &p.Stok, &p.Satuan, &p.BeratGram, &p.FotoURLs,
		&p.IsOPOP, &p.IsAktif, &p.Rating, &p.TotalTerjual, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, produk.ErrProdukNotFound
		}
		return nil, err
	}
	return p, nil
}

// ── Pesanan Repository ────────────────────────────────────────────────────────

type PesananRepository struct {
	db *pgxpool.Pool
}

func NewPesananRepository(db *pgxpool.Pool) *PesananRepository {
	return &PesananRepository{db: db}
}

func (r *PesananRepository) Create(ctx context.Context, p *pesanan.Pesanan) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pesanan (id, buyer_tipe, nasabah_id, bmt_buyer_id, toko_id, bmt_seller_id,
		nomor_pesanan, status, subtotal, ongkir, total, alamat_kirim, metode_bayar, catatan, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`, p.ID, p.BuyerTipe, p.NasabahID, p.BMTBuyerID, p.TokoID, p.BMTSellerID,
		p.NomorPesanan, p.Status, p.Subtotal, p.Ongkir, p.Total, p.AlamatKirim,
		p.MetodeBayar, p.Catatan, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PesananRepository) CreateItems(ctx context.Context, items []*pesanan.PesananItem) error {
	for _, item := range items {
		_, err := r.db.Exec(ctx, `
			INSERT INTO pesanan_item (id, pesanan_id, produk_id, nama_produk, harga, jumlah, subtotal)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`, item.ID, item.PesananID, item.ProdukID, item.NamaProduk, item.Harga, item.Jumlah, item.Subtotal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PesananRepository) GetByID(ctx context.Context, id uuid.UUID) (*pesanan.Pesanan, error) {
	return r.scanPesanan(r.db.QueryRow(ctx, `
		SELECT id, buyer_tipe, nasabah_id, bmt_buyer_id, toko_id, bmt_seller_id,
		nomor_pesanan, status, subtotal, ongkir, total, alamat_kirim, kurir, nomor_resi, metode_bayar, catatan, created_at, updated_at
		FROM pesanan WHERE id = $1
	`, id))
}

func (r *PesananRepository) GetByNomor(ctx context.Context, nomor string) (*pesanan.Pesanan, error) {
	return r.scanPesanan(r.db.QueryRow(ctx, `
		SELECT id, buyer_tipe, nasabah_id, bmt_buyer_id, toko_id, bmt_seller_id,
		nomor_pesanan, status, subtotal, ongkir, total, alamat_kirim, kurir, nomor_resi, metode_bayar, catatan, created_at, updated_at
		FROM pesanan WHERE nomor_pesanan = $1
	`, nomor))
}

func (r *PesananRepository) GetWithItems(ctx context.Context, id uuid.UUID) (*pesanan.Pesanan, error) {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, pesanan_id, produk_id, nama_produk, harga, jumlah, subtotal
		FROM pesanan_item WHERE pesanan_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &pesanan.PesananItem{}
		err := rows.Scan(&item.ID, &item.PesananID, &item.ProdukID, &item.NamaProduk,
			&item.Harga, &item.Jumlah, &item.Subtotal)
		if err != nil {
			return nil, err
		}
		p.Items = append(p.Items, item)
	}

	return p, nil
}

func (r *PesananRepository) List(ctx context.Context, filter pesanan.ListPesananFilter) ([]*pesanan.Pesanan, int64, error) {
	offset := (filter.Page - 1) * filter.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, buyer_tipe, nasabah_id, bmt_buyer_id, toko_id, bmt_seller_id,
		nomor_pesanan, status, subtotal, ongkir, total, alamat_kirim, kurir, nomor_resi, metode_bayar, catatan, created_at, updated_at
		FROM pesanan
		WHERE ($1::uuid IS NULL OR nasabah_id = $1)
		  AND ($2::uuid IS NULL OR toko_id = $2)
		  AND ($3::uuid IS NULL OR bmt_seller_id = $3)
		  AND ($4::text IS NULL OR status = $4)
		ORDER BY created_at DESC LIMIT $5 OFFSET $6
	`, filter.NasabahID, filter.TokoID, filter.BMTSellerID, filter.Status, filter.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*pesanan.Pesanan
	for rows.Next() {
		p, err := r.scanPesanan(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM pesanan`).Scan(&total)

	return result, total, nil
}

func (r *PesananRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status pesanan.StatusPesanan) error {
	_, err := r.db.Exec(ctx, `UPDATE pesanan SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *PesananRepository) UpdatePengiriman(ctx context.Context, id uuid.UUID, kurir, nomorResi string) error {
	_, err := r.db.Exec(ctx, `UPDATE pesanan SET kurir=$1, nomor_resi=$2, updated_at=NOW() WHERE id=$3`, kurir, nomorResi, id)
	return err
}

func (r *PesananRepository) UpdateMetodeBayar(ctx context.Context, id uuid.UUID, metode pesanan.MetodeBayar) error {
	_, err := r.db.Exec(ctx, `UPDATE pesanan SET metode_bayar=$1, updated_at=NOW() WHERE id=$2`, metode, id)
	return err
}

func (r *PesananRepository) GenerateNomor(ctx context.Context, bmtID uuid.UUID) (string, error) {
	var seq int
	r.db.QueryRow(ctx, `SELECT COALESCE(COUNT(*), 0) + 1 FROM pesanan WHERE bmt_seller_id = $1`, bmtID).Scan(&seq)
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%d%02d%02d-%06d",
		bmtID.String()[:8],
		now.Year(), now.Month(), now.Day(),
		seq), nil
}

func (r *PesananRepository) scanPesanan(s scanner) (*pesanan.Pesanan, error) {
	p := &pesanan.Pesanan{}
	var alamat json.RawMessage
	err := s.Scan(&p.ID, &p.BuyerTipe, &p.NasabahID, &p.BMTBuyerID, &p.TokoID, &p.BMTSellerID,
		&p.NomorPesanan, &p.Status, &p.Subtotal, &p.Ongkir, &p.Total, &alamat,
		&p.Kurir, &p.NomorResi, &p.MetodeBayar, &p.Catatan, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pesanan.ErrPesananNotFound
		}
		return nil, err
	}
	p.AlamatKirim = alamat
	return p, nil
}
