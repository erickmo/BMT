package pengguna

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPenggunaNotFound   = errors.New("pengguna tidak ditemukan")
	ErrPasswordSalah      = errors.New("username atau password salah")
	ErrPenggunaNonAktif   = errors.New("akun tidak aktif")
	ErrPenggunaBlokir     = errors.New("akun diblokir")
	ErrUsernameSudahAda   = errors.New("username sudah digunakan")
)

type StatusPengguna string

const (
	StatusAktif   StatusPengguna = "AKTIF"
	StatusNonAktif StatusPengguna = "NONAKTIF"
	StatusBlokir  StatusPengguna = "BLOKIR"
)

type Pengguna struct {
	ID           uuid.UUID      `json:"id"`
	BMTID        uuid.UUID      `json:"bmt_id"`
	CabangID     *uuid.UUID     `json:"cabang_id,omitempty"` // nil = akses semua cabang
	Username     string         `json:"username"`
	PasswordHash string         `json:"-"`
	NamaLengkap  string         `json:"nama_lengkap"`
	Email        string         `json:"email"`
	Telepon      string         `json:"telepon"`
	Role         string         `json:"role"`
	Status       StatusPengguna `json:"status"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CreatePenggunaInput struct {
	BMTID       uuid.UUID  `json:"bmt_id" validate:"required"`
	CabangID    *uuid.UUID `json:"cabang_id"`
	Username    string     `json:"username" validate:"required,min=4,max=50"`
	Password    string     `json:"password" validate:"required,min=8"`
	NamaLengkap string     `json:"nama_lengkap" validate:"required"`
	Email       string     `json:"email"`
	Telepon     string     `json:"telepon"`
	Role        string     `json:"role" validate:"required"`
}

type Repository interface {
	Create(ctx context.Context, p *Pengguna) error
	GetByID(ctx context.Context, id uuid.UUID) (*Pengguna, error)
	GetByUsername(ctx context.Context, bmtID uuid.UUID, username string) (*Pengguna, error)
	ListByBMT(ctx context.Context, bmtID uuid.UUID) ([]*Pengguna, error)
	Update(ctx context.Context, p *Pengguna) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusPengguna) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

func New(input CreatePenggunaInput, passwordHash string) (*Pengguna, error) {
	if input.Username == "" {
		return nil, errors.New("username wajib diisi")
	}
	if input.Role == "" {
		return nil, errors.New("role wajib diisi")
	}
	return &Pengguna{
		ID:           uuid.New(),
		BMTID:        input.BMTID,
		CabangID:     input.CabangID,
		Username:     input.Username,
		PasswordHash: passwordHash,
		NamaLengkap:  input.NamaLengkap,
		Email:        input.Email,
		Telepon:      input.Telepon,
		Role:         input.Role,
		Status:       StatusAktif,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
