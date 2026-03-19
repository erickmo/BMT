package merchant

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMerchantNotFound      = errors.New("merchant tidak ditemukan")
	ErrMerchantTidakAktif    = errors.New("merchant tidak aktif")
	ErrTerminalKioskNotFound = errors.New("terminal kiosk tidak ditemukan")
)

type StatusMerchant string

const (
	StatusMerchantAktif    StatusMerchant = "AKTIF"
	StatusMerchantNonAktif StatusMerchant = "NONAKTIF"
	StatusMerchantSuspend  StatusMerchant = "SUSPEND"
)

// Merchant represents an NFC payment merchant (pondok shop or canteen)
type Merchant struct {
	ID           uuid.UUID      `json:"id"`
	BMTID        uuid.UUID      `json:"bmt_id"`
	CabangID     uuid.UUID      `json:"cabang_id"`
	Nama         string         `json:"nama"`
	Deskripsi    string         `json:"deskripsi,omitempty"`
	Kategori     string         `json:"kategori,omitempty"` // KANTIN | KOPERASI | TOKO | dll.
	AlamatLokasi string         `json:"alamat_lokasi,omitempty"`
	FotoURL      string         `json:"foto_url,omitempty"`
	// Owner/Kasir reference
	OwnerID      *uuid.UUID     `json:"owner_id,omitempty"`
	// BMT rekening for receiving payments
	RekeningID   *uuid.UUID     `json:"rekening_id,omitempty"`
	Status       StatusMerchant `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// TerminalKiosk represents an NFC kiosk terminal with IP whitelist
// (also stored here for merchant context; NFC domain has its own for low-level ops)
type TerminalKiosk struct {
	ID        uuid.UUID `json:"id"`
	BMTID     uuid.UUID `json:"bmt_id"`
	CabangID  uuid.UUID `json:"cabang_id"`
	Nama      string    `json:"nama"`
	IPAddress string    `json:"ip_address"`
	Lokasi    string    `json:"lokasi,omitempty"`
	IsAktif   bool      `json:"is_aktif"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateMerchantInput struct {
	BMTID        uuid.UUID
	CabangID     uuid.UUID
	Nama         string
	Deskripsi    string
	Kategori     string
	AlamatLokasi string
	FotoURL      string
	OwnerID      *uuid.UUID
	RekeningID   *uuid.UUID
}

type ListMerchantFilter struct {
	BMTID    *uuid.UUID
	CabangID *uuid.UUID
	Status   *StatusMerchant
	Page     int
	PerPage  int
}

type Repository interface {
	// Merchant
	Create(ctx context.Context, m *Merchant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	List(ctx context.Context, filter ListMerchantFilter) ([]*Merchant, int64, error)
	Update(ctx context.Context, m *Merchant) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusMerchant) error

	// TerminalKiosk
	CreateTerminal(ctx context.Context, t *TerminalKiosk) error
	GetTerminalByID(ctx context.Context, id uuid.UUID) (*TerminalKiosk, error)
	GetTerminalByIP(ctx context.Context, ip string) (*TerminalKiosk, error)
	ListTerminalByBMT(ctx context.Context, bmtID uuid.UUID) ([]*TerminalKiosk, error)
	UpdateTerminalStatus(ctx context.Context, id uuid.UUID, isAktif bool) error
}

func NewMerchant(input CreateMerchantInput) (*Merchant, error) {
	if input.Nama == "" {
		return nil, errors.New("nama merchant wajib diisi")
	}
	now := time.Now()
	return &Merchant{
		ID:           uuid.New(),
		BMTID:        input.BMTID,
		CabangID:     input.CabangID,
		Nama:         input.Nama,
		Deskripsi:    input.Deskripsi,
		Kategori:     input.Kategori,
		AlamatLokasi: input.AlamatLokasi,
		FotoURL:      input.FotoURL,
		OwnerID:      input.OwnerID,
		RekeningID:   input.RekeningID,
		Status:       StatusMerchantAktif,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func NewTerminalKiosk(bmtID, cabangID uuid.UUID, nama, ip, lokasi string) (*TerminalKiosk, error) {
	if nama == "" {
		return nil, errors.New("nama terminal kiosk wajib diisi")
	}
	if ip == "" {
		return nil, errors.New("IP address terminal kiosk wajib diisi")
	}
	now := time.Now()
	return &TerminalKiosk{
		ID:        uuid.New(),
		BMTID:     bmtID,
		CabangID:  cabangID,
		Nama:      nama,
		IPAddress: ip,
		Lokasi:    lokasi,
		IsAktif:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
