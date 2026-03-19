package sdm

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrKontrakNotFound  = errors.New("kontrak SDM tidak ditemukan")
	ErrSlipGajiNotFound = errors.New("slip gaji tidak ditemukan")
	ErrSlipSudahDibayar = errors.New("slip gaji sudah dalam status dibayar")
	ErrGajiPokokNol     = errors.New("gaji pokok harus lebih dari 0")
)

type TipePegawai string

const (
	TipePengajar  TipePegawai = "PENGAJAR"
	TipeKaryawan  TipePegawai = "KARYAWAN"
)

type TipeKontrak string

const (
	KontrakTetap     TipeKontrak = "TETAP"
	KontrakKontrak   TipeKontrak = "KONTRAK"
	KontrakHonorer   TipeKontrak = "HONORER"
	KontrakMagang    TipeKontrak = "MAGANG"
)

type StatusSlipGaji string

const (
	StatusSlipDraft     StatusSlipGaji = "DRAFT"
	StatusSlipDisetujui StatusSlipGaji = "DISETUJUI"
	StatusSlipDibayar   StatusSlipGaji = "DIBAYAR"
)

// SDMKontrak represents an employment contract for a teacher or staff member
type SDMKontrak struct {
	ID                   uuid.UUID       `json:"id"`
	BMTID                uuid.UUID       `json:"bmt_id"`
	CabangID             uuid.UUID       `json:"cabang_id"`
	PegawaiID            uuid.UUID       `json:"pegawai_id"`
	TipePegawai          TipePegawai     `json:"tipe_pegawai"`
	NomorKontrak         string          `json:"nomor_kontrak"`
	TipeKontrak          TipeKontrak     `json:"tipe_kontrak"`
	TanggalMulai         time.Time       `json:"tanggal_mulai"`
	TanggalSelesai       *time.Time      `json:"tanggal_selesai,omitempty"`
	GajiPokok            int64           `json:"gaji_pokok"`
	Tunjangan            json.RawMessage `json:"tunjangan"`             // {"transport": 200000, "makan": 150000}
	PotonganTetap        json.RawMessage `json:"potongan_tetap"`        // {"bpjs_kesehatan": 50000}
	RekeningGajiID       *uuid.UUID      `json:"rekening_gaji_id,omitempty"`
	Status               string          `json:"status"` // AKTIF | BERAKHIR | DIBATALKAN
	DokumenURL           string          `json:"dokumen_url,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// SDMSlipGaji represents a monthly payslip generated for a contract
type SDMSlipGaji struct {
	ID               uuid.UUID      `json:"id"`
	BMTID            uuid.UUID      `json:"bmt_id"`
	KontrakID        uuid.UUID      `json:"kontrak_id"`
	Periode          string         `json:"periode"` // "2025-01"
	// Pendapatan
	GajiPokok        int64          `json:"gaji_pokok"`
	TunjanganTotal   int64          `json:"tunjangan_total"`
	TunjanganDetail  json.RawMessage `json:"tunjangan_detail"`
	// Potongan
	PotonganAbsensi  int64          `json:"potongan_absensi"`
	PotonganTetap    int64          `json:"potongan_tetap"`
	PotonganLain     int64          `json:"potongan_lain"`
	// Total
	GajiBersih       int64          `json:"gaji_bersih"`
	// Rekap absensi
	HariKerja        int16          `json:"hari_kerja"`
	HariHadir        int16          `json:"hari_hadir"`
	HariSakit        int16          `json:"hari_sakit"`
	HariIzin         int16          `json:"hari_izin"`
	HariAlfa         int16          `json:"hari_alfa"`
	// Status payroll
	Status           StatusSlipGaji `json:"status"`
	DibayarAt        *time.Time     `json:"dibayar_at,omitempty"`
	TransaksiID      *uuid.UUID     `json:"transaksi_id,omitempty"`
	FileURL          string         `json:"file_url,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type CreateKontrakInput struct {
	BMTID          uuid.UUID
	CabangID       uuid.UUID
	PegawaiID      uuid.UUID
	TipePegawai    TipePegawai
	NomorKontrak   string
	TipeKontrak    TipeKontrak
	TanggalMulai   time.Time
	TanggalSelesai *time.Time
	GajiPokok      int64
	Tunjangan      json.RawMessage
	PotonganTetap  json.RawMessage
	RekeningGajiID *uuid.UUID
}

type GenerateSlipInput struct {
	BMTID           uuid.UUID
	KontrakID       uuid.UUID
	Periode         string
	GajiPokok       int64
	TunjanganTotal  int64
	TunjanganDetail json.RawMessage
	PotonganAbsensi int64
	PotonganTetap   int64
	HariKerja       int16
	HariHadir       int16
	HariSakit       int16
	HariIzin        int16
	HariAlfa        int16
}

type ListKontrakFilter struct {
	BMTID       *uuid.UUID
	CabangID    *uuid.UUID
	PegawaiID   *uuid.UUID
	TipePegawai *TipePegawai
	Status      string
	Page        int
	PerPage     int
}

type Repository interface {
	// Kontrak
	CreateKontrak(ctx context.Context, k *SDMKontrak) error
	GetKontrakByID(ctx context.Context, id uuid.UUID) (*SDMKontrak, error)
	GetKontrakByNomor(ctx context.Context, nomor string) (*SDMKontrak, error)
	ListKontrak(ctx context.Context, filter ListKontrakFilter) ([]*SDMKontrak, int64, error)
	ListKontrakAktif(ctx context.Context, bmtID uuid.UUID) ([]*SDMKontrak, error)
	UpdateKontrak(ctx context.Context, k *SDMKontrak) error
	UpdateStatusKontrak(ctx context.Context, id uuid.UUID, status string) error

	// Slip Gaji
	CreateSlip(ctx context.Context, s *SDMSlipGaji) error
	GetSlipByID(ctx context.Context, id uuid.UUID) (*SDMSlipGaji, error)
	GetSlipByKontrakAndPeriode(ctx context.Context, kontrakID uuid.UUID, periode string) (*SDMSlipGaji, error)
	ListSlipByBMT(ctx context.Context, bmtID uuid.UUID, periode string, status *StatusSlipGaji, page, perPage int) ([]*SDMSlipGaji, int64, error)
	UpdateStatusSlip(ctx context.Context, id uuid.UUID, status StatusSlipGaji, dibayarAt *time.Time, transaksiID *uuid.UUID) error
	UpdateFileSlip(ctx context.Context, id uuid.UUID, fileURL string) error
}

func NewKontrak(input CreateKontrakInput) (*SDMKontrak, error) {
	if input.GajiPokok <= 0 {
		return nil, ErrGajiPokokNol
	}
	if input.NomorKontrak == "" {
		return nil, errors.New("nomor kontrak wajib diisi")
	}

	tunjangan := input.Tunjangan
	if tunjangan == nil {
		tunjangan = json.RawMessage(`{}`)
	}
	potongan := input.PotonganTetap
	if potongan == nil {
		potongan = json.RawMessage(`{}`)
	}

	now := time.Now()
	return &SDMKontrak{
		ID:             uuid.New(),
		BMTID:          input.BMTID,
		CabangID:       input.CabangID,
		PegawaiID:      input.PegawaiID,
		TipePegawai:    input.TipePegawai,
		NomorKontrak:   input.NomorKontrak,
		TipeKontrak:    input.TipeKontrak,
		TanggalMulai:   input.TanggalMulai,
		TanggalSelesai: input.TanggalSelesai,
		GajiPokok:      input.GajiPokok,
		Tunjangan:      tunjangan,
		PotonganTetap:  potongan,
		RekeningGajiID: input.RekeningGajiID,
		Status:         "AKTIF",
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func NewSlipGaji(input GenerateSlipInput) (*SDMSlipGaji, error) {
	if input.Periode == "" {
		return nil, errors.New("periode slip gaji wajib diisi")
	}
	detail := input.TunjanganDetail
	if detail == nil {
		detail = json.RawMessage(`{}`)
	}
	gajiBersih := input.GajiPokok + input.TunjanganTotal - input.PotonganAbsensi - input.PotonganTetap

	now := time.Now()
	return &SDMSlipGaji{
		ID:              uuid.New(),
		BMTID:           input.BMTID,
		KontrakID:       input.KontrakID,
		Periode:         input.Periode,
		GajiPokok:       input.GajiPokok,
		TunjanganTotal:  input.TunjanganTotal,
		TunjanganDetail: detail,
		PotonganAbsensi: input.PotonganAbsensi,
		PotonganTetap:   input.PotonganTetap,
		GajiBersih:      gajiBersih,
		HariKerja:       input.HariKerja,
		HariHadir:       input.HariHadir,
		HariSakit:       input.HariSakit,
		HariIzin:        input.HariIzin,
		HariAlfa:        input.HariAlfa,
		Status:          StatusSlipDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}
