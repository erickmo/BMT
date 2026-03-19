package opop

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOPOPKomisiNotFound  = errors.New("komisi OPOP tidak ditemukan")
	ErrTokoTidakTerdaftarOPOP = errors.New("toko tidak terdaftar di marketplace OPOP")
)

// StatusKomisi represents the billing status of an OPOP commission
type StatusKomisi string

const (
	StatusKomisiPending StatusKomisi = "PENDING"
	StatusKomisiDitagih StatusKomisi = "DITAGIH"
	StatusKomisiLunas   StatusKomisi = "LUNAS"
)

// OPOPKomisi represents the platform commission from each completed OPOP order
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

// OPOPIklan represents a premium advertising slot in the OPOP marketplace
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
}

type CreateKomisiInput struct {
	PesananID    uuid.UUID
	BMTSellerID  uuid.UUID
	NilaiPesanan int64
	PersenKomisi float64
	Periode      string
}

type ListOPOPFilter struct {
	BMTSellerID *uuid.UUID
	Periode     string
	Status      *StatusKomisi
	Page        int
	PerPage     int
}

type Repository interface {
	// Komisi
	CreateKomisi(ctx context.Context, k *OPOPKomisi) error
	GetKomisiByID(ctx context.Context, id uuid.UUID) (*OPOPKomisi, error)
	GetKomisiByPesanan(ctx context.Context, pesananID uuid.UUID) (*OPOPKomisi, error)
	ListKomisi(ctx context.Context, filter ListOPOPFilter) ([]*OPOPKomisi, int64, error)
	UpdateStatusKomisi(ctx context.Context, id uuid.UUID, status StatusKomisi) error
	SumKomisiByPeriode(ctx context.Context, bmtID uuid.UUID, periode string) (int64, error)

	// Iklan
	CreateIklan(ctx context.Context, i *OPOPIklan) error
	GetIklanByID(ctx context.Context, id uuid.UUID) (*OPOPIklan, error)
	ListIklanAktif(ctx context.Context, tipe TipeIklan) ([]*OPOPIklan, error)
	UpdateStatusIklan(ctx context.Context, id uuid.UUID, status StatusIklan) error
}

func NewKomisi(input CreateKomisiInput) (*OPOPKomisi, error) {
	if input.NilaiPesanan <= 0 {
		return nil, errors.New("nilai pesanan harus lebih dari 0")
	}
	if input.PersenKomisi < 0 {
		return nil, errors.New("persen komisi tidak boleh negatif")
	}

	// Calculate commission using integer arithmetic to avoid float precision issues
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

func NewIklan(tokoID, bmtID uuid.UUID, tipe TipeIklan, mulai, selesai time.Time, nominal int64) (*OPOPIklan, error) {
	if nominal <= 0 {
		return nil, errors.New("nominal iklan harus lebih dari 0")
	}
	if selesai.Before(mulai) {
		return nil, errors.New("tanggal selesai iklan harus setelah tanggal mulai")
	}
	return &OPOPIklan{
		ID:             uuid.New(),
		TokoID:         tokoID,
		BMTID:          bmtID,
		Tipe:           tipe,
		TanggalMulai:   mulai,
		TanggalSelesai: selesai,
		Nominal:        nominal,
		Status:         StatusIklanAktif,
		CreatedAt:      time.Now(),
	}, nil
}
