package monetisasi

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrIklanNotFound      = errors.New("iklan OPOP tidak ditemukan")
	ErrKomisiNotFound     = errors.New("komisi OPOP tidak ditemukan")
	ErrNominalHarusPosistif = errors.New("nominal iklan harus lebih dari 0")
	ErrTanggalTidakValid  = errors.New("tanggal selesai iklan harus setelah tanggal mulai")
)

type TipeIklan string

const (
	TipeFeaturedToko  TipeIklan = "FEATURED_TOKO"
	TipeBannerProduk  TipeIklan = "BANNER_PRODUK"
	TipeTopSearch     TipeIklan = "TOP_SEARCH"
)

type StatusIklan string

const (
	StatusIklanAktif    StatusIklan = "AKTIF"
	StatusIklanNonAktif StatusIklan = "NONAKTIF"
	StatusIklanExpired  StatusIklan = "EXPIRED"
)

type StatusKomisi string

const (
	StatusKomisiPending StatusKomisi = "PENDING"
	StatusKomisiDitagih StatusKomisi = "DITAGIH"
	StatusKomisiLunas   StatusKomisi = "LUNAS"
)

// OPOPIklan represents a premium advertising slot purchased by a toko in the OPOP marketplace
type OPOPIklan struct {
	ID             uuid.UUID   `json:"id"`
	TokoID         uuid.UUID   `json:"toko_id"`
	BMTID          uuid.UUID   `json:"bmt_id"`
	Tipe           TipeIklan   `json:"tipe"`
	TanggalMulai   time.Time   `json:"tanggal_mulai"`
	TanggalSelesai time.Time   `json:"tanggal_selesai"`
	Nominal        int64       `json:"nominal"`
	Status         StatusIklan `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// OPOPKomisi records the platform commission for each completed OPOP transaction
type OPOPKomisi struct {
	ID            uuid.UUID    `json:"id"`
	PesananID     uuid.UUID    `json:"pesanan_id"`
	BMTSellerID   uuid.UUID    `json:"bmt_seller_id"`
	NilaiPesanan  int64        `json:"nilai_pesanan"`
	PersenKomisi  float64      `json:"persen_komisi"`
	NominalKomisi int64        `json:"nominal_komisi"`
	Periode       string       `json:"periode"` // "2025-01"
	Status        StatusKomisi `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type CreateIklanInput struct {
	TokoID         uuid.UUID
	BMTID          uuid.UUID
	Tipe           TipeIklan
	TanggalMulai   time.Time
	TanggalSelesai time.Time
	Nominal        int64
}

type CreateKomisiInput struct {
	PesananID    uuid.UUID
	BMTSellerID  uuid.UUID
	NilaiPesanan int64
	PersenKomisi float64
	Periode      string
}

type ListIklanFilter struct {
	BMTID    *uuid.UUID
	TokoID   *uuid.UUID
	Tipe     *TipeIklan
	Status   *StatusIklan
	Page     int
	PerPage  int
}

type ListKomisiFilter struct {
	BMTSellerID *uuid.UUID
	Periode     string
	Status      *StatusKomisi
	Page        int
	PerPage     int
}

type Repository interface {
	// Iklan
	CreateIklan(ctx context.Context, i *OPOPIklan) error
	GetIklanByID(ctx context.Context, id uuid.UUID) (*OPOPIklan, error)
	ListIklan(ctx context.Context, filter ListIklanFilter) ([]*OPOPIklan, int64, error)
	ListIklanAktifByTipe(ctx context.Context, tipe TipeIklan) ([]*OPOPIklan, error)
	UpdateStatusIklan(ctx context.Context, id uuid.UUID, status StatusIklan) error
	ExpireIklanLewatTanggal(ctx context.Context) (int64, error)

	// Komisi
	CreateKomisi(ctx context.Context, k *OPOPKomisi) error
	GetKomisiByID(ctx context.Context, id uuid.UUID) (*OPOPKomisi, error)
	GetKomisiByPesanan(ctx context.Context, pesananID uuid.UUID) (*OPOPKomisi, error)
	ListKomisi(ctx context.Context, filter ListKomisiFilter) ([]*OPOPKomisi, int64, error)
	SumKomisiByPeriode(ctx context.Context, bmtSellerID uuid.UUID, periode string) (int64, error)
	UpdateStatusKomisi(ctx context.Context, id uuid.UUID, status StatusKomisi) error
}

func NewIklan(input CreateIklanInput) (*OPOPIklan, error) {
	if input.Nominal <= 0 {
		return nil, ErrNominalHarusPosistif
	}
	if input.TanggalSelesai.Before(input.TanggalMulai) {
		return nil, ErrTanggalTidakValid
	}
	now := time.Now()
	return &OPOPIklan{
		ID:             uuid.New(),
		TokoID:         input.TokoID,
		BMTID:          input.BMTID,
		Tipe:           input.Tipe,
		TanggalMulai:   input.TanggalMulai,
		TanggalSelesai: input.TanggalSelesai,
		Nominal:        input.Nominal,
		Status:         StatusIklanAktif,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func NewKomisi(input CreateKomisiInput) (*OPOPKomisi, error) {
	if input.NilaiPesanan <= 0 {
		return nil, errors.New("nilai pesanan harus lebih dari 0")
	}
	if input.PersenKomisi < 0 {
		return nil, errors.New("persen komisi tidak boleh negatif")
	}
	if input.Periode == "" {
		return nil, errors.New("periode komisi wajib diisi")
	}

	// Use integer arithmetic: avoid float precision issues
	nominalKomisi := int64(float64(input.NilaiPesanan) * input.PersenKomisi / 100)

	now := time.Now()
	return &OPOPKomisi{
		ID:            uuid.New(),
		PesananID:     input.PesananID,
		BMTSellerID:   input.BMTSellerID,
		NilaiPesanan:  input.NilaiPesanan,
		PersenKomisi:  input.PersenKomisi,
		NominalKomisi: nominalKomisi,
		Periode:       input.Periode,
		Status:        StatusKomisiPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (i *OPOPIklan) IsExpired() bool {
	return time.Now().After(i.TanggalSelesai)
}
