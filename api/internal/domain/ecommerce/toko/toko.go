package toko

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokoNotFound      = errors.New("toko tidak ditemukan")
	ErrSlugSudahDipakai  = errors.New("slug toko sudah digunakan")
	ErrTokoTidakAktif    = errors.New("toko tidak aktif")
)

type KategoriToko string

const (
	KategoriPondok     KategoriToko = "PONDOK"
	KategoriBMTKoperasi KategoriToko = "BMT_KOPERASI"
)

type StatusToko string

const (
	StatusTokoAktif    StatusToko = "AKTIF"
	StatusTokoNonAktif StatusToko = "NONAKTIF"
	StatusTokoSuspend  StatusToko = "SUSPEND"
)

type Toko struct {
	ID           uuid.UUID    `json:"id"`
	BMTID        uuid.UUID    `json:"bmt_id"`
	CabangID     uuid.UUID    `json:"cabang_id"`
	Nama         string       `json:"nama"`
	Slug         string       `json:"slug"`
	Deskripsi    string       `json:"deskripsi,omitempty"`
	LogoURL      string       `json:"logo_url,omitempty"`
	BannerURL    string       `json:"banner_url,omitempty"`
	KategoriToko KategoriToko `json:"kategori_toko"`
	IsOPOP       bool         `json:"is_opop"`
	Status       StatusToko   `json:"status"`
	Rating       *float64     `json:"rating,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type CreateTokoInput struct {
	BMTID        uuid.UUID    `json:"bmt_id"`
	CabangID     uuid.UUID    `json:"cabang_id"`
	Nama         string       `json:"nama"`
	Slug         string       `json:"slug"`
	Deskripsi    string       `json:"deskripsi"`
	LogoURL      string       `json:"logo_url"`
	BannerURL    string       `json:"banner_url"`
	KategoriToko KategoriToko `json:"kategori_toko"`
	IsOPOP       bool         `json:"is_opop"`
}

type UpdateTokoInput struct {
	Nama      string     `json:"nama"`
	Deskripsi string     `json:"deskripsi"`
	LogoURL   string     `json:"logo_url"`
	BannerURL string     `json:"banner_url"`
	IsOPOP    bool       `json:"is_opop"`
	Status    StatusToko `json:"status"`
}

type ListTokoFilter struct {
	BMTID    *uuid.UUID
	CabangID *uuid.UUID
	IsOPOP   *bool
	Status   *StatusToko
	Page     int
	PerPage  int
}

type Repository interface {
	Create(ctx context.Context, t *Toko) error
	GetByID(ctx context.Context, id uuid.UUID) (*Toko, error)
	GetBySlug(ctx context.Context, slug string) (*Toko, error)
	GetByCabang(ctx context.Context, bmtID, cabangID uuid.UUID) (*Toko, error)
	List(ctx context.Context, filter ListTokoFilter) ([]*Toko, int64, error)
	ListOPOP(ctx context.Context, page, perPage int) ([]*Toko, int64, error)
	Update(ctx context.Context, t *Toko) error
	UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusToko) error
}

func NewToko(input CreateTokoInput) (*Toko, error) {
	if input.Nama == "" {
		return nil, errors.New("nama toko wajib diisi")
	}
	if input.Slug == "" {
		return nil, errors.New("slug toko wajib diisi")
	}
	if input.KategoriToko == "" {
		return nil, errors.New("kategori toko wajib diisi")
	}
	now := time.Now()
	return &Toko{
		ID:           uuid.New(),
		BMTID:        input.BMTID,
		CabangID:     input.CabangID,
		Nama:         input.Nama,
		Slug:         input.Slug,
		Deskripsi:    input.Deskripsi,
		LogoURL:      input.LogoURL,
		BannerURL:    input.BannerURL,
		KategoriToko: input.KategoriToko,
		IsOPOP:       input.IsOPOP,
		Status:       StatusTokoAktif,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (t *Toko) ValidasiAktif() error {
	if t.Status != StatusTokoAktif {
		return ErrTokoTidakAktif
	}
	return nil
}
