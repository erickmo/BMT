package ppdb

import (
	"context"
	"errors"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPendaftarNotFound      = errors.New("pendaftar PPDB tidak ditemukan")
	ErrNomorPendaftaranDuplikat = errors.New("nomor pendaftaran sudah ada")
	ErrStatusTransisiTidakValid = errors.New("transisi status PPDB tidak valid")
	ErrGelombangTidakAktif    = errors.New("gelombang PPDB tidak aktif atau sudah ditutup")
)

type StatusPendaftar string

const (
	StatusDaftar    StatusPendaftar = "DAFTAR"
	StatusSeleksi   StatusPendaftar = "SELEKSI"
	StatusDiterima  StatusPendaftar = "DITERIMA"
	StatusDitolak   StatusPendaftar = "DITOLAK"
	StatusMundur    StatusPendaftar = "MUNDUR"
)

type StatusBayar string

const (
	StatusBelumBayar StatusBayar = "BELUM_BAYAR"
	StatusSudahBayar StatusBayar = "SUDAH_BAYAR"
)

// GelombangPPDB represents an enrollment wave/batch for a school year
type GelombangPPDB struct {
	ID             uuid.UUID `json:"id"`
	BMTID          uuid.UUID `json:"bmt_id"`
	CabangID       uuid.UUID `json:"cabang_id"`
	Nama           string    `json:"nama"`
	TahunAjaran    string    `json:"tahun_ajaran"`
	TanggalBuka    time.Time `json:"tanggal_buka"`
	TanggalTutup   time.Time `json:"tanggal_tutup"`
	Kuota          *int16    `json:"kuota,omitempty"`
	BiayaDaftar    int64     `json:"biaya_daftar"`
	Persyaratan    json.RawMessage `json:"persyaratan"` // ["KK", "Akta", "Ijazah"]
	IsAktif        bool      `json:"is_aktif"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PPDBPendaftar represents a prospective new student registration
type PPDBPendaftar struct {
	ID                  uuid.UUID       `json:"id"`
	BMTID               uuid.UUID       `json:"bmt_id"`
	CabangID            uuid.UUID       `json:"cabang_id"`
	GelombangID         uuid.UUID       `json:"gelombang_id"`
	NomorPendaftaran    string          `json:"nomor_pendaftaran"`
	NamaLengkap         string          `json:"nama_lengkap"`
	NIK                 string          `json:"nik,omitempty"`
	TanggalLahir        *time.Time      `json:"tanggal_lahir,omitempty"`
	NamaWali            string          `json:"nama_wali"`
	TeleponWali         string          `json:"telepon_wali"`
	EmailWali           string          `json:"email_wali,omitempty"`
	PilihanTingkat      string          `json:"pilihan_tingkat,omitempty"`
	Status              StatusPendaftar `json:"status"`
	StatusBayar         StatusBayar     `json:"status_bayar"`
	MidtransOrderID     *string         `json:"midtrans_order_id,omitempty"`
	Dokumen             json.RawMessage `json:"dokumen"` // {"kk": "url", "akta": "url"}
	Catatan             string          `json:"catatan,omitempty"`
	// If accepted — linked to newly created santri and nasabah
	SantriID            *uuid.UUID      `json:"santri_id,omitempty"`
	NasabahID           *uuid.UUID      `json:"nasabah_id,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type CreateGelombangInput struct {
	BMTID          uuid.UUID
	CabangID       uuid.UUID
	Nama           string
	TahunAjaran    string
	TanggalBuka    time.Time
	TanggalTutup   time.Time
	Kuota          *int16
	BiayaDaftar    int64
	Persyaratan    []string
}

type CreatePendaftarInput struct {
	BMTID          uuid.UUID
	CabangID       uuid.UUID
	GelombangID    uuid.UUID
	NamaLengkap    string
	NIK            string
	TanggalLahir   *time.Time
	NamaWali       string
	TeleponWali    string
	EmailWali      string
	PilihanTingkat string
}

type ListPendaftarFilter struct {
	BMTID       *uuid.UUID
	CabangID    *uuid.UUID
	GelombangID *uuid.UUID
	Status      *StatusPendaftar
	TahunAjaran string
	Page        int
	PerPage     int
}

// validTransisi defines allowed status transitions for PPDB
var validTransisi = map[StatusPendaftar][]StatusPendaftar{
	StatusDaftar:   {StatusSeleksi, StatusMundur},
	StatusSeleksi:  {StatusDiterima, StatusDitolak, StatusMundur},
	StatusDiterima: {StatusMundur},
	StatusDitolak:  {},
	StatusMundur:   {},
}

type Repository interface {
	// Gelombang
	CreateGelombang(ctx context.Context, g *GelombangPPDB) error
	GetGelombangByID(ctx context.Context, id uuid.UUID) (*GelombangPPDB, error)
	ListGelombang(ctx context.Context, bmtID, cabangID uuid.UUID, tahunAjaran string) ([]*GelombangPPDB, error)
	UpdateGelombang(ctx context.Context, g *GelombangPPDB) error

	// Pendaftar
	Create(ctx context.Context, p *PPDBPendaftar) error
	GetByID(ctx context.Context, id uuid.UUID) (*PPDBPendaftar, error)
	GetByNomor(ctx context.Context, nomor string) (*PPDBPendaftar, error)
	List(ctx context.Context, filter ListPendaftarFilter) ([]*PPDBPendaftar, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusPendaftar, catatan string) error
	UpdateStatusBayar(ctx context.Context, id uuid.UUID, statusBayar StatusBayar, midtransOrderID *string) error
	LinkSantri(ctx context.Context, id uuid.UUID, santriID, nasabahID uuid.UUID) error
	GenerateNomor(ctx context.Context, bmtID, cabangID uuid.UUID, tahunAjaran string) (string, error)
}

func NewGelombang(input CreateGelombangInput) (*GelombangPPDB, error) {
	if input.Nama == "" {
		return nil, errors.New("nama gelombang PPDB wajib diisi")
	}
	if input.TahunAjaran == "" {
		return nil, errors.New("tahun ajaran wajib diisi")
	}
	if input.TanggalTutup.Before(input.TanggalBuka) {
		return nil, errors.New("tanggal tutup harus setelah tanggal buka")
	}

	syarat := json.RawMessage(`[]`)
	if len(input.Persyaratan) > 0 {
		b, err := json.Marshal(input.Persyaratan)
		if err != nil {
			return nil, errors.New("gagal memproses persyaratan")
		}
		syarat = b
	}

	now := time.Now()
	return &GelombangPPDB{
		ID:           uuid.New(),
		BMTID:        input.BMTID,
		CabangID:     input.CabangID,
		Nama:         input.Nama,
		TahunAjaran:  input.TahunAjaran,
		TanggalBuka:  input.TanggalBuka,
		TanggalTutup: input.TanggalTutup,
		Kuota:        input.Kuota,
		BiayaDaftar:  input.BiayaDaftar,
		Persyaratan:  syarat,
		IsAktif:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func NewPendaftar(input CreatePendaftarInput, nomor string) (*PPDBPendaftar, error) {
	if input.NamaLengkap == "" {
		return nil, errors.New("nama lengkap pendaftar wajib diisi")
	}
	if input.NamaWali == "" {
		return nil, errors.New("nama wali wajib diisi")
	}
	if input.TeleponWali == "" {
		return nil, errors.New("telepon wali wajib diisi")
	}
	now := time.Now()
	return &PPDBPendaftar{
		ID:               uuid.New(),
		BMTID:            input.BMTID,
		CabangID:         input.CabangID,
		GelombangID:      input.GelombangID,
		NomorPendaftaran: nomor,
		NamaLengkap:      input.NamaLengkap,
		NIK:              input.NIK,
		TanggalLahir:     input.TanggalLahir,
		NamaWali:         input.NamaWali,
		TeleponWali:      input.TeleponWali,
		EmailWali:        input.EmailWali,
		PilihanTingkat:   input.PilihanTingkat,
		Status:           StatusDaftar,
		StatusBayar:      StatusBelumBayar,
		Dokumen:          json.RawMessage(`{}`),
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (p *PPDBPendaftar) TransisiStatus(statusBaru StatusPendaftar) error {
	allowed, ok := validTransisi[p.Status]
	if !ok {
		return ErrStatusTransisiTidakValid
	}
	for _, s := range allowed {
		if s == statusBaru {
			return nil
		}
	}
	return ErrStatusTransisiTidakValid
}
