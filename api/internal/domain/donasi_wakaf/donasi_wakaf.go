package donasi_wakaf

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProgramNotFound      = errors.New("program donasi tidak ditemukan")
	ErrTransaksiNotFound    = errors.New("transaksi donasi tidak ditemukan")
	ErrAsetWakafNotFound    = errors.New("aset wakaf tidak ditemukan")
	ErrNominalHarusPosistif = errors.New("nominal harus lebih dari 0")
	ErrProgramTidakAktif    = errors.New("program donasi tidak aktif")
)

type TipeProgram string

const (
	TipeDonasi   TipeProgram = "DONASI"
	TipeWakaf    TipeProgram = "WAKAF"
	TipeInfaq    TipeProgram = "INFAQ"
	TipeZakat    TipeProgram = "ZAKAT"
)

type StatusProgram string

const (
	StatusProgramAktif   StatusProgram = "AKTIF"
	StatusProgramSelesai StatusProgram = "SELESAI"
	StatusProgramDitutup StatusProgram = "DITUTUP"
)

type MetodePembayaran string

const (
	MetodeMidtrans    MetodePembayaran = "MIDTRANS"
	MetodeRekeningBMT MetodePembayaran = "REKENING_BMT"
	MetodeNFC         MetodePembayaran = "NFC"
)

type StatusTransaksiDonasi string

const (
	StatusTxPending  StatusTransaksiDonasi = "PENDING"
	StatusTxSettled  StatusTransaksiDonasi = "SETTLED"
	StatusTxGagal    StatusTransaksiDonasi = "GAGAL"
)

type JenisAsetWakaf string

const (
	AsetTanah     JenisAsetWakaf = "TANAH"
	AsetBangunan  JenisAsetWakaf = "BANGUNAN"
	AsetUang      JenisAsetWakaf = "UANG"
	AsetKendaraan JenisAsetWakaf = "KENDARAAN"
	AsetLainnya   JenisAsetWakaf = "LAINNYA"
)

// ProgramDonasi represents a fundraising campaign (donasi, wakaf, infaq, or zakat)
type ProgramDonasi struct {
	ID             uuid.UUID     `json:"id"`
	BMTID          uuid.UUID     `json:"bmt_id"`
	CabangID       uuid.UUID     `json:"cabang_id"`
	Nama           string        `json:"nama"`
	Deskripsi      string        `json:"deskripsi,omitempty"`
	Tipe           TipeProgram   `json:"tipe"`
	TargetNominal  *int64        `json:"target_nominal,omitempty"` // NULL = no target
	Terkumpul      int64         `json:"terkumpul"`
	TanggalMulai   time.Time     `json:"tanggal_mulai"`
	TanggalSelesai *time.Time    `json:"tanggal_selesai,omitempty"`
	FotoURL        string        `json:"foto_url,omitempty"`
	Status         StatusProgram `json:"status"`
	RekeningID     uuid.UUID     `json:"rekening_id"` // dedicated rekening for this program
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// TransaksiDonasi records a single donation transaction
type TransaksiDonasi struct {
	ID               uuid.UUID              `json:"id"`
	BMTID            uuid.UUID              `json:"bmt_id"`
	ProgramID        uuid.UUID              `json:"program_id"`
	NasabahID        *uuid.UUID             `json:"nasabah_id,omitempty"`
	NamaDonatur      string                 `json:"nama_donatur,omitempty"`
	Nominal          int64                  `json:"nominal"`
	IsAnonim         bool                   `json:"is_anonim"`
	Pesan            string                 `json:"pesan,omitempty"`
	Metode           MetodePembayaran       `json:"metode"`
	MidtransOrderID  *string                `json:"midtrans_order_id,omitempty"`
	RekeningID       *uuid.UUID             `json:"rekening_id,omitempty"`
	IdempotencyKey   *uuid.UUID             `json:"idempotency_key,omitempty"`
	Status           StatusTransaksiDonasi  `json:"status"`
	CreatedAt        time.Time              `json:"created_at"`
}

// AsetWakaf represents a wakaf asset managed by the BMT as nazhir
type AsetWakaf struct {
	ID           uuid.UUID      `json:"id"`
	BMTID        uuid.UUID      `json:"bmt_id"`
	Nama         string         `json:"nama"`
	Deskripsi    string         `json:"deskripsi,omitempty"`
	Jenis        JenisAsetWakaf `json:"jenis"`
	NilaiAwal    int64          `json:"nilai_awal"`
	Wakif        string         `json:"wakif,omitempty"`
	Nazhir       string         `json:"nazhir,omitempty"`
	Peruntukan   string         `json:"peruntukan"`
	DokumenURL   string         `json:"dokumen_url,omitempty"`
	TanggalWakaf time.Time      `json:"tanggal_wakaf"`
	Status       string         `json:"status"` // AKTIF | DIKELOLA | SELESAI
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// HasilWakaf records the productive output of a wakaf asset per period
type HasilWakaf struct {
	ID           uuid.UUID       `json:"id"`
	AsetID       uuid.UUID       `json:"aset_id"`
	Periode      string          `json:"periode"` // "2025-01"
	Pendapatan   int64           `json:"pendapatan"`
	Beban        int64           `json:"beban"`
	HasilBersih  int64           `json:"hasil_bersih"`
	Distribusi   interface{}     `json:"distribusi,omitempty"` // JSON — where results are channeled
	CreatedAt    time.Time       `json:"created_at"`
}

type CreateProgramInput struct {
	BMTID          uuid.UUID
	CabangID       uuid.UUID
	Nama           string
	Deskripsi      string
	Tipe           TipeProgram
	TargetNominal  *int64
	TanggalMulai   time.Time
	TanggalSelesai *time.Time
	FotoURL        string
	RekeningID     uuid.UUID
}

type CreateTransaksiInput struct {
	BMTID           uuid.UUID
	ProgramID       uuid.UUID
	NasabahID       *uuid.UUID
	NamaDonatur     string
	Nominal         int64
	IsAnonim        bool
	Pesan           string
	Metode          MetodePembayaran
	MidtransOrderID *string
	RekeningID      *uuid.UUID
	IdempotencyKey  *uuid.UUID
}

type ListProgramFilter struct {
	BMTID    *uuid.UUID
	CabangID *uuid.UUID
	Tipe     *TipeProgram
	Status   *StatusProgram
	Page     int
	PerPage  int
}

type Repository interface {
	// Program Donasi
	CreateProgram(ctx context.Context, p *ProgramDonasi) error
	GetProgramByID(ctx context.Context, id uuid.UUID) (*ProgramDonasi, error)
	ListProgram(ctx context.Context, filter ListProgramFilter) ([]*ProgramDonasi, int64, error)
	UpdateProgram(ctx context.Context, p *ProgramDonasi) error
	TambahTerkumpul(ctx context.Context, programID uuid.UUID, nominal int64) error

	// Transaksi Donasi
	CreateTransaksi(ctx context.Context, t *TransaksiDonasi) error
	GetTransaksiByID(ctx context.Context, id uuid.UUID) (*TransaksiDonasi, error)
	GetTransaksiByIdempotency(ctx context.Context, key uuid.UUID) (*TransaksiDonasi, error)
	ListTransaksiByProgram(ctx context.Context, programID uuid.UUID, page, perPage int) ([]*TransaksiDonasi, int64, error)
	UpdateStatusTransaksi(ctx context.Context, id uuid.UUID, status StatusTransaksiDonasi) error

	// Aset Wakaf
	CreateAset(ctx context.Context, a *AsetWakaf) error
	GetAsetByID(ctx context.Context, id uuid.UUID) (*AsetWakaf, error)
	ListAset(ctx context.Context, bmtID uuid.UUID, page, perPage int) ([]*AsetWakaf, int64, error)
	UpdateAset(ctx context.Context, a *AsetWakaf) error

	// Hasil Wakaf
	CreateHasil(ctx context.Context, h *HasilWakaf) error
	ListHasilByAset(ctx context.Context, asetID uuid.UUID, page, perPage int) ([]*HasilWakaf, int64, error)
}

func NewProgram(input CreateProgramInput) (*ProgramDonasi, error) {
	if input.Nama == "" {
		return nil, errors.New("nama program donasi wajib diisi")
	}
	if input.RekeningID == uuid.Nil {
		return nil, errors.New("rekening program donasi wajib diisi")
	}
	now := time.Now()
	return &ProgramDonasi{
		ID:             uuid.New(),
		BMTID:          input.BMTID,
		CabangID:       input.CabangID,
		Nama:           input.Nama,
		Deskripsi:      input.Deskripsi,
		Tipe:           input.Tipe,
		TargetNominal:  input.TargetNominal,
		Terkumpul:      0,
		TanggalMulai:   input.TanggalMulai,
		TanggalSelesai: input.TanggalSelesai,
		FotoURL:        input.FotoURL,
		Status:         StatusProgramAktif,
		RekeningID:     input.RekeningID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func NewTransaksiDonasi(input CreateTransaksiInput) (*TransaksiDonasi, error) {
	if input.Nominal <= 0 {
		return nil, ErrNominalHarusPosistif
	}
	return &TransaksiDonasi{
		ID:              uuid.New(),
		BMTID:           input.BMTID,
		ProgramID:       input.ProgramID,
		NasabahID:       input.NasabahID,
		NamaDonatur:     input.NamaDonatur,
		Nominal:         input.Nominal,
		IsAnonim:        input.IsAnonim,
		Pesan:           input.Pesan,
		Metode:          input.Metode,
		MidtransOrderID: input.MidtransOrderID,
		RekeningID:      input.RekeningID,
		IdempotencyKey:  input.IdempotencyKey,
		Status:          StatusTxPending,
		CreatedAt:       time.Now(),
	}, nil
}
