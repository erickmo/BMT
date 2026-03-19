package ulasan

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUlasanNotFound     = errors.New("ulasan tidak ditemukan")
	ErrUlasanSudahAda     = errors.New("ulasan untuk pesanan ini sudah diberikan")
	ErrRatingTidakValid   = errors.New("rating harus antara 1 sampai 5")
)

type UlasanProduk struct {
	ID        uuid.UUID       `json:"id"`
	ProdukID  uuid.UUID       `json:"produk_id"`
	PesananID uuid.UUID       `json:"pesanan_id"`
	NasabahID uuid.UUID       `json:"nasabah_id"`
	Rating    int16           `json:"rating"`
	Komentar  string          `json:"komentar,omitempty"`
	FotoURLs  json.RawMessage `json:"foto_urls"`
	CreatedAt time.Time       `json:"created_at"`
}

type CreateUlasanInput struct {
	ProdukID  uuid.UUID
	PesananID uuid.UUID
	NasabahID uuid.UUID
	Rating    int16
	Komentar  string
	FotoURLs  []string
}

type ListUlasanFilter struct {
	ProdukID  *uuid.UUID
	NasabahID *uuid.UUID
	Rating    *int16
	Page      int
	PerPage   int
}

type Repository interface {
	Create(ctx context.Context, u *UlasanProduk) error
	GetByID(ctx context.Context, id uuid.UUID) (*UlasanProduk, error)
	GetByPesananAndProduk(ctx context.Context, pesananID, produkID uuid.UUID) (*UlasanProduk, error)
	List(ctx context.Context, filter ListUlasanFilter) ([]*UlasanProduk, int64, error)
	AverageRatingByProduk(ctx context.Context, produkID uuid.UUID) (float64, error)
}

func NewUlasan(input CreateUlasanInput) (*UlasanProduk, error) {
	if input.Rating < 1 || input.Rating > 5 {
		return nil, ErrRatingTidakValid
	}

	fotoBytes := json.RawMessage(`[]`)
	if len(input.FotoURLs) > 0 {
		b, err := json.Marshal(input.FotoURLs)
		if err != nil {
			return nil, errors.New("gagal memproses foto URLs")
		}
		fotoBytes = b
	}

	return &UlasanProduk{
		ID:        uuid.New(),
		ProdukID:  input.ProdukID,
		PesananID: input.PesananID,
		NasabahID: input.NasabahID,
		Rating:    input.Rating,
		Komentar:  input.Komentar,
		FotoURLs:  fotoBytes,
		CreatedAt: time.Now(),
	}, nil
}
