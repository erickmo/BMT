package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSnapshotNotFound  = errors.New("snapshot analytics tidak ditemukan")
	ErrTemplateNotFound  = errors.New("template laporan tidak ditemukan")
)

// MetrikAnalytics defines the type of metric being tracked
type MetrikAnalytics string

const (
	MetrikTotalNasabah     MetrikAnalytics = "TOTAL_NASABAH"
	MetrikTotalDPK         MetrikAnalytics = "TOTAL_DPK"
	MetrikTotalPembiayaan  MetrikAnalytics = "TOTAL_PEMBIAYAAN"
	MetrikTransaksiHariIni MetrikAnalytics = "TRANSAKSI_HARI_INI"
	MetrikNPFRatio         MetrikAnalytics = "NPF_RATIO"
	MetrikOPOPPenjualan    MetrikAnalytics = "OPOP_PENJUALAN"
	MetrikKolektibilitas1  MetrikAnalytics = "KOLEKTIBILITAS_1"
	MetrikKolektibilitas2  MetrikAnalytics = "KOLEKTIBILITAS_2"
	MetrikKolektibilitas3  MetrikAnalytics = "KOLEKTIBILITAS_3"
	MetrikKolektibilitas4  MetrikAnalytics = "KOLEKTIBILITAS_4"
	MetrikKolektibilitas5  MetrikAnalytics = "KOLEKTIBILITAS_5"
)

// DomainLaporan defines the domain/subject of a report template
type DomainLaporan string

const (
	DomainNasabah    DomainLaporan = "NASABAH"
	DomainRekening   DomainLaporan = "REKENING"
	DomainTransaksi  DomainLaporan = "TRANSAKSI"
	DomainPembiayaan DomainLaporan = "PEMBIAYAAN"
	DomainAbsensi    DomainLaporan = "ABSENSI"
	DomainNilai      DomainLaporan = "NILAI"
	DomainOPOP       DomainLaporan = "OPOP"
)

// AnalyticsSnapshot stores a daily point-in-time metric value
type AnalyticsSnapshot struct {
	ID        uuid.UUID       `json:"id"`
	BMTID     uuid.UUID       `json:"bmt_id"`
	CabangID  *uuid.UUID      `json:"cabang_id,omitempty"` // NULL = BMT consolidation
	Tanggal   time.Time       `json:"tanggal"`
	Metrik    MetrikAnalytics `json:"metrik"`
	Nilai     float64         `json:"nilai"`
	CreatedAt time.Time       `json:"created_at"`
}

// LaporanTemplate stores a saved custom report configuration
type LaporanTemplate struct {
	ID         uuid.UUID       `json:"id"`
	BMTID      uuid.UUID       `json:"bmt_id"`
	Nama       string          `json:"nama"`
	Domain     DomainLaporan   `json:"domain"`
	Kolom      json.RawMessage `json:"kolom"`  // ["field1", "field2", ...]
	Filter     json.RawMessage `json:"filter"` // {"status": "AKTIF"}
	Urutan     json.RawMessage `json:"urutan"` // [{"field": "tanggal", "dir": "DESC"}]
	DibuatOleh uuid.UUID       `json:"dibuat_oleh"`
	IsPublik   bool            `json:"is_publik"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// OPOPAnalyticsHarian stores daily analytics snapshot for each OPOP toko
type OPOPAnalyticsHarian struct {
	Tanggal           time.Time `json:"tanggal"`
	TokoID            uuid.UUID `json:"toko_id"`
	BMTID             uuid.UUID `json:"bmt_id"`
	TotalPesanan      int       `json:"total_pesanan"`
	TotalPendapatan   int64     `json:"total_pendapatan"`
	TotalItemTerjual  int       `json:"total_item_terjual"`
	PengunjungUnik    int       `json:"pengunjung_unik"`
}

type CreateSnapshotInput struct {
	BMTID    uuid.UUID
	CabangID *uuid.UUID
	Tanggal  time.Time
	Metrik   MetrikAnalytics
	Nilai    float64
}

type CreateLaporanTemplateInput struct {
	BMTID      uuid.UUID
	Nama       string
	Domain     DomainLaporan
	Kolom      []string
	Filter     json.RawMessage
	Urutan     json.RawMessage
	DibuatOleh uuid.UUID
	IsPublik   bool
}

type ListSnapshotFilter struct {
	BMTID    uuid.UUID
	CabangID *uuid.UUID
	Metrik   *MetrikAnalytics
	DariTgl  time.Time
	SampaiTgl time.Time
}

type Repository interface {
	// Snapshot
	UpsertSnapshot(ctx context.Context, s *AnalyticsSnapshot) error
	ListSnapshot(ctx context.Context, filter ListSnapshotFilter) ([]*AnalyticsSnapshot, error)
	GetLatestSnapshot(ctx context.Context, bmtID uuid.UUID, cabangID *uuid.UUID, metrik MetrikAnalytics) (*AnalyticsSnapshot, error)

	// Laporan Template
	CreateTemplate(ctx context.Context, t *LaporanTemplate) error
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*LaporanTemplate, error)
	ListTemplate(ctx context.Context, bmtID uuid.UUID, isPublik *bool, page, perPage int) ([]*LaporanTemplate, int64, error)
	UpdateTemplate(ctx context.Context, t *LaporanTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error

	// OPOP Analytics
	UpsertOPOPHarian(ctx context.Context, a *OPOPAnalyticsHarian) error
	ListOPOPHarian(ctx context.Context, tokoID uuid.UUID, dari, sampai time.Time) ([]*OPOPAnalyticsHarian, error)
	SumOPOPByBMT(ctx context.Context, bmtID uuid.UUID, dari, sampai time.Time) (int64, int, error)
}

func NewLaporanTemplate(input CreateLaporanTemplateInput) (*LaporanTemplate, error) {
	if input.Nama == "" {
		return nil, errors.New("nama template laporan wajib diisi")
	}
	if len(input.Kolom) == 0 {
		return nil, errors.New("kolom laporan wajib diisi minimal satu")
	}

	kolomBytes, err := json.Marshal(input.Kolom)
	if err != nil {
		return nil, errors.New("gagal memproses kolom laporan")
	}

	filter := input.Filter
	if filter == nil {
		filter = json.RawMessage(`{}`)
	}
	urutan := input.Urutan
	if urutan == nil {
		urutan = json.RawMessage(`[]`)
	}

	now := time.Now()
	return &LaporanTemplate{
		ID:         uuid.New(),
		BMTID:      input.BMTID,
		Nama:       input.Nama,
		Domain:     input.Domain,
		Kolom:      kolomBytes,
		Filter:     filter,
		Urutan:     urutan,
		DibuatOleh: input.DibuatOleh,
		IsPublik:   input.IsPublik,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
