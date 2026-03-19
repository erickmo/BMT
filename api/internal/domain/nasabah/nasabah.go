package nasabah

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNasabahNotFound   = errors.New("nasabah tidak ditemukan")
	ErrNIKSudahTerdaftar = errors.New("NIK sudah terdaftar")
	ErrNasabahNonAktif   = errors.New("nasabah tidak aktif")
	ErrPINSalah          = errors.New("PIN nasabah tidak sesuai")
)

type StatusNasabah string

const (
	StatusAktif    StatusNasabah = "AKTIF"
	StatusNonAktif StatusNasabah = "NONAKTIF"
	StatusBlokir   StatusNasabah = "BLOKIR"
)

type Nasabah struct {
	ID           uuid.UUID     `json:"id"`
	BMTID        uuid.UUID     `json:"bmt_id"`
	CabangID     uuid.UUID     `json:"cabang_id"`
	NomorNasabah string        `json:"nomor_nasabah"`
	NIK          string        `json:"nik,omitempty"`
	NamaLengkap  string        `json:"nama_lengkap"`
	TempatLahir  string        `json:"tempat_lahir"`
	TanggalLahir *time.Time    `json:"tanggal_lahir,omitempty"`
	JenisKelamin string        `json:"jenis_kelamin"`
	Alamat       string        `json:"alamat"`
	Telepon      string        `json:"telepon"`
	Email        string        `json:"email"`
	FotoURL      string        `json:"foto_url,omitempty"`
	Pekerjaan    string        `json:"pekerjaan"`
	Status       StatusNasabah `json:"status"`
	PINHash      string        `json:"-"`
	PasswordHash string        `json:"-"`
	LastLoginAt  *time.Time    `json:"last_login_at,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type KartuNFC struct {
	ID                uuid.UUID  `json:"id"`
	BMTID             uuid.UUID  `json:"bmt_id"`
	NasabahID         uuid.UUID  `json:"nasabah_id"`
	UID               string     `json:"uid"`
	PINHash           string     `json:"-"`
	LimitPerTransaksi int64      `json:"limit_per_transaksi"`
	LimitHarian       int64      `json:"limit_harian"`
	SaldoNFC          int64      `json:"saldo_nfc"`
	Status            string     `json:"status"`
	ExpiredAt         *time.Time `json:"expired_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateNasabahInput struct {
	BMTID        uuid.UUID  `json:"bmt_id" validate:"required"`
	CabangID     uuid.UUID  `json:"cabang_id" validate:"required"`
	NIK          string     `json:"nik"`
	NamaLengkap  string     `json:"nama_lengkap" validate:"required,min=3"`
	TempatLahir  string     `json:"tempat_lahir"`
	TanggalLahir *time.Time `json:"tanggal_lahir"`
	JenisKelamin string     `json:"jenis_kelamin"`
	Alamat       string     `json:"alamat"`
	Telepon      string     `json:"telepon" validate:"required"`
	Email        string     `json:"email"`
	Pekerjaan    string     `json:"pekerjaan"`
}

type Repository interface {
	Create(ctx context.Context, n *Nasabah) error
	GetByID(ctx context.Context, id uuid.UUID) (*Nasabah, error)
	GetByNomor(ctx context.Context, bmtID uuid.UUID, nomor string) (*Nasabah, error)
	GetByNIK(ctx context.Context, bmtID uuid.UUID, nik string) (*Nasabah, error)
	Search(ctx context.Context, bmtID uuid.UUID, query string, limit, offset int) ([]*Nasabah, int64, error)
	Update(ctx context.Context, n *Nasabah) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusNasabah) error

	CreateKartuNFC(ctx context.Context, k *KartuNFC) error
	GetKartuNFCByUID(ctx context.Context, uid string) (*KartuNFC, error)
	GetKartuNFCByNasabah(ctx context.Context, nasabahID uuid.UUID) ([]*KartuNFC, error)
	UpdateKartuNFC(ctx context.Context, k *KartuNFC) error

	GenerateNomorNasabah(ctx context.Context, bmtID uuid.UUID) (string, error)
}

func New(input CreateNasabahInput, nomorNasabah string) (*Nasabah, error) {
	if input.NamaLengkap == "" {
		return nil, errors.New("nama lengkap wajib diisi")
	}
	if input.Telepon == "" {
		return nil, errors.New("telepon wajib diisi")
	}
	return &Nasabah{
		ID:           uuid.New(),
		BMTID:        input.BMTID,
		CabangID:     input.CabangID,
		NomorNasabah: nomorNasabah,
		NIK:          input.NIK,
		NamaLengkap:  input.NamaLengkap,
		TempatLahir:  input.TempatLahir,
		TanggalLahir: input.TanggalLahir,
		JenisKelamin: input.JenisKelamin,
		Alamat:       input.Alamat,
		Telepon:      input.Telepon,
		Email:        input.Email,
		Pekerjaan:    input.Pekerjaan,
		Status:       StatusAktif,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
