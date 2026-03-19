package platform

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBMTNotFound     = errors.New("bmt tidak ditemukan")
	ErrCabangNotFound  = errors.New("cabang tidak ditemukan")
	ErrBMTSudahAda     = errors.New("kode bmt sudah digunakan")
	ErrFiturTidakAktif = errors.New("fitur tidak diaktifkan di kontrak BMT")
	ErrKontrakExpired  = errors.New("kontrak BMT sudah expired")
)

type StatusBMT string

const (
	StatusBMTAktif    StatusBMT = "AKTIF"
	StatusBMTSuspend  StatusBMT = "SUSPEND"
	StatusBMTNonaktif StatusBMT = "NONAKTIF"
)

type BMT struct {
	ID          uuid.UUID         `json:"id"`
	Kode        string            `json:"kode"`
	Nama        string            `json:"nama"`
	Alamat      string            `json:"alamat"`
	Telepon     string            `json:"telepon"`
	Email       string            `json:"email"`
	LogoURL     string            `json:"logo_url"`
	Status      StatusBMT         `json:"status"`
	Whitelabel  map[string]string `json:"whitelabel"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type Cabang struct {
	ID        uuid.UUID `json:"id"`
	BMTID     uuid.UUID `json:"bmt_id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	Alamat    string    `json:"alamat"`
	Telepon   string    `json:"telepon"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KontrakBMT struct {
	ID             uuid.UUID              `json:"id"`
	BMTID          uuid.UUID              `json:"bmt_id"`
	TanggalMulai   time.Time              `json:"tanggal_mulai"`
	TanggalSelesai time.Time              `json:"tanggal_selesai"`
	Fitur          map[string]interface{} `json:"fitur"`
	Tarif          map[string]interface{} `json:"tarif"`
	PICNama        string                 `json:"pic_nama"`
	PICTelepon     string                 `json:"pic_telepon"`
	PICEmail       string                 `json:"pic_email"`
	Status         string                 `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type PecahanUang struct {
	ID           uuid.UUID  `json:"id"`
	Nominal      int64      `json:"nominal"`
	Jenis        string     `json:"jenis"`
	Label        string     `json:"label"`
	IsAktif      bool       `json:"is_aktif"`
	Urutan       int        `json:"urutan"`
	BerlakuSejak time.Time  `json:"berlaku_sejak"`
	DitarikPada  *time.Time `json:"ditarik_pada,omitempty"`
}

type CreateBMTInput struct {
	Kode    string `json:"kode" validate:"required,min=2,max=20"`
	Nama    string `json:"nama" validate:"required,min=3,max=255"`
	Alamat  string `json:"alamat"`
	Telepon string `json:"telepon"`
	Email   string `json:"email"`
}

type CreateCabangInput struct {
	BMTID   uuid.UUID `json:"bmt_id" validate:"required"`
	Kode    string    `json:"kode" validate:"required,min=2,max=20"`
	Nama    string    `json:"nama" validate:"required,min=3,max=255"`
	Alamat  string    `json:"alamat"`
	Telepon string    `json:"telepon"`
}

type Repository interface {
	CreateBMT(ctx context.Context, bmt *BMT) error
	GetBMT(ctx context.Context, id uuid.UUID) (*BMT, error)
	GetBMTByKode(ctx context.Context, kode string) (*BMT, error)
	ListBMT(ctx context.Context) ([]*BMT, error)
	UpdateBMTStatus(ctx context.Context, id uuid.UUID, status StatusBMT) error

	CreateCabang(ctx context.Context, cabang *Cabang) error
	GetCabang(ctx context.Context, id uuid.UUID) (*Cabang, error)
	ListCabangByBMT(ctx context.Context, bmtID uuid.UUID) ([]*Cabang, error)

	CreateKontrak(ctx context.Context, kontrak *KontrakBMT) error
	GetKontrakAktif(ctx context.Context, bmtID uuid.UUID) (*KontrakBMT, error)

	GetPecahanAktif(ctx context.Context) ([]*PecahanUang, error)
	CreatePecahan(ctx context.Context, p *PecahanUang) error
	UpdatePecahan(ctx context.Context, p *PecahanUang) error
}

func NewBMT(input CreateBMTInput) (*BMT, error) {
	if input.Kode == "" {
		return nil, errors.New("kode bmt wajib diisi")
	}
	if input.Nama == "" {
		return nil, errors.New("nama bmt wajib diisi")
	}
	return &BMT{
		ID:        uuid.New(),
		Kode:      input.Kode,
		Nama:      input.Nama,
		Alamat:    input.Alamat,
		Telepon:   input.Telepon,
		Email:     input.Email,
		Status:    StatusBMTAktif,
		Whitelabel: map[string]string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func NewCabang(input CreateCabangInput) (*Cabang, error) {
	if input.BMTID == uuid.Nil {
		return nil, errors.New("bmt_id wajib diisi — cabang harus terikat ke tenant BMT")
	}
	if input.Kode == "" {
		return nil, errors.New("kode cabang wajib diisi")
	}
	return &Cabang{
		ID:        uuid.New(),
		BMTID:     input.BMTID,
		Kode:      input.Kode,
		Nama:      input.Nama,
		Alamat:    input.Alamat,
		Telepon:   input.Telepon,
		Status:    "AKTIF",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
