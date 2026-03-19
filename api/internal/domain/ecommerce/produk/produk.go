package produk

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProdukNotFound    = errors.New("produk tidak ditemukan")
	ErrStokTidakCukup    = errors.New("stok produk tidak mencukupi")
	ErrHargaHarusPosistif = errors.New("harga produk harus lebih dari 0")
	ErrProdukTidakAktif  = errors.New("produk tidak aktif")
)

type Produk struct {
	ID           uuid.UUID       `json:"id"`
	BMTID        uuid.UUID       `json:"bmt_id"`
	TokoID       uuid.UUID       `json:"toko_id"`
	Nama         string          `json:"nama"`
	Slug         string          `json:"slug"`
	Deskripsi    string          `json:"deskripsi,omitempty"`
	Kategori     string          `json:"kategori"`
	Harga        int64           `json:"harga"`
	HargaB2B     *int64          `json:"harga_b2b,omitempty"`
	Stok         int             `json:"stok"`
	Satuan       string          `json:"satuan"`
	BeratGram    *int            `json:"berat_gram,omitempty"`
	FotoURLs     json.RawMessage `json:"foto_urls"`
	IsOPOP       bool            `json:"is_opop"`
	IsAktif      bool            `json:"is_aktif"`
	Rating       *float64        `json:"rating,omitempty"`
	TotalTerjual int             `json:"total_terjual"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CreateProdukInput struct {
	BMTID     uuid.UUID `json:"bmt_id"`
	TokoID    uuid.UUID `json:"toko_id"`
	Nama      string    `json:"nama"`
	Slug      string    `json:"slug"`
	Deskripsi string    `json:"deskripsi"`
	Kategori  string    `json:"kategori"`
	Harga     int64     `json:"harga"`
	HargaB2B  *int64    `json:"harga_b2b"`
	Stok      int       `json:"stok"`
	Satuan    string    `json:"satuan"`
	BeratGram *int      `json:"berat_gram"`
	FotoURLs  []string  `json:"foto_urls"`
	IsOPOP    bool      `json:"is_opop"`
}

type UpdateStokInput struct {
	ProdukID uuid.UUID `json:"produk_id"`
	Delta    int       `json:"delta"` // positif = tambah, negatif = kurang
}

type ListProdukFilter struct {
	BMTID    *uuid.UUID
	TokoID   *uuid.UUID
	Kategori string
	IsOPOP   *bool
	IsAktif  *bool
	Search   string
	Page     int
	PerPage  int
}

type Repository interface {
	Create(ctx context.Context, p *Produk) error
	GetByID(ctx context.Context, id uuid.UUID) (*Produk, error)
	GetBySlugAndToko(ctx context.Context, tokoID uuid.UUID, slug string) (*Produk, error)
	List(ctx context.Context, filter ListProdukFilter) ([]*Produk, int64, error)
	Update(ctx context.Context, p *Produk) error
	UpdateStok(ctx context.Context, id uuid.UUID, stokBaru int) error
	UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error
	IncrementTerjual(ctx context.Context, id uuid.UUID, jumlah int) error
	LockForUpdate(ctx context.Context, id uuid.UUID) (*Produk, error)
}

func NewProduk(input CreateProdukInput) (*Produk, error) {
	if input.Nama == "" {
		return nil, errors.New("nama produk wajib diisi")
	}
	if input.Harga <= 0 {
		return nil, ErrHargaHarusPosistif
	}
	if input.Satuan == "" {
		input.Satuan = "pcs"
	}

	fotoBytes := json.RawMessage(`[]`)
	if len(input.FotoURLs) > 0 {
		b, err := json.Marshal(input.FotoURLs)
		if err != nil {
			return nil, errors.New("gagal memproses foto URLs")
		}
		fotoBytes = b
	}

	now := time.Now()
	return &Produk{
		ID:        uuid.New(),
		BMTID:     input.BMTID,
		TokoID:    input.TokoID,
		Nama:      input.Nama,
		Slug:      input.Slug,
		Deskripsi: input.Deskripsi,
		Kategori:  input.Kategori,
		Harga:     input.Harga,
		HargaB2B:  input.HargaB2B,
		Stok:      input.Stok,
		Satuan:    input.Satuan,
		BeratGram: input.BeratGram,
		FotoURLs:  fotoBytes,
		IsOPOP:    input.IsOPOP,
		IsAktif:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p *Produk) ValidasiPembelian(jumlah int) error {
	if !p.IsAktif {
		return ErrProdukTidakAktif
	}
	if p.Stok < jumlah {
		return ErrStokTidakCukup
	}
	return nil
}

func (p *Produk) HargaEfektif(isB2B bool) int64 {
	if isB2B && p.HargaB2B != nil && *p.HargaB2B > 0 {
		return *p.HargaB2B
	}
	return p.Harga
}
