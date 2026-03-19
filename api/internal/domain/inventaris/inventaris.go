package inventaris

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAsetNotFound       = errors.New("aset tidak ditemukan")
	ErrPeminjamanNotFound = errors.New("peminjaman aset tidak ditemukan")
	ErrAsetSedangDipinjam = errors.New("aset sedang dipinjam")
	ErrAsetTidakLayak     = errors.New("kondisi aset tidak memungkinkan untuk dipinjam")
)

type KategoriAset string

const (
	KategoriGedung     KategoriAset = "GEDUNG"
	KategoriKendaraan  KategoriAset = "KENDARAAN"
	KategoriPeralatan  KategoriAset = "PERALATAN"
	KategoriFurnitur   KategoriAset = "FURNITUR"
	KategoriElektronik KategoriAset = "ELEKTRONIK"
	KategoriLainnya    KategoriAset = "LAINNYA"
)

type KondisiAset string

const (
	KondsiBaik        KondisiAset = "BAIK"
	KondisiRusakRingan KondisiAset = "RUSAK_RINGAN"
	KondisiRusakBerat KondisiAset = "RUSAK_BERAT"
	KondisiTidakAktif KondisiAset = "TIDAK_AKTIF"
)

type StatusPeminjaman string

const (
	StatusDipinjam      StatusPeminjaman = "DIPINJAM"
	StatusDikembalikan  StatusPeminjaman = "DIKEMBALIKAN"
	StatusTerlambat     StatusPeminjaman = "TERLAMBAT"
)

type TipePeminjam string

const (
	PeminjamSantri        TipePeminjam = "SANTRI"
	PeminjamPenggunaPondok TipePeminjam = "PENGGUNA_PONDOK"
	PeminjamExternal      TipePeminjam = "EXTERNAL"
)

// InventarisAset represents a fixed asset of the pondok
type InventarisAset struct {
	ID                uuid.UUID    `json:"id"`
	BMTID             uuid.UUID    `json:"bmt_id"`
	CabangID          uuid.UUID    `json:"cabang_id"`
	KodeAset          string       `json:"kode_aset"`
	Nama              string       `json:"nama"`
	Kategori          KategoriAset `json:"kategori"`
	Lokasi            string       `json:"lokasi,omitempty"`
	TanggalPerolehan  time.Time    `json:"tanggal_perolehan"`
	NilaiPerolehan    int64        `json:"nilai_perolehan"`
	NilaiBuku         int64        `json:"nilai_buku"` // after depreciation
	UmurEkonomis      *int16       `json:"umur_ekonomis,omitempty"` // in years
	Kondisi           KondisiAset  `json:"kondisi"`
	FotoURL           string       `json:"foto_url,omitempty"`
	KodeAkun          string       `json:"kode_akun"` // default "131"
	IsAktif           bool         `json:"is_aktif"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// InventarisPeminjaman records borrowing of an asset or room
type InventarisPeminjaman struct {
	ID                     uuid.UUID        `json:"id"`
	BMTID                  uuid.UUID        `json:"bmt_id"`
	AsetID                 uuid.UUID        `json:"aset_id"`
	PeminjamID             uuid.UUID        `json:"peminjam_id"`
	TipePeminjam           TipePeminjam     `json:"tipe_peminjam"`
	Keperluan              string           `json:"keperluan"`
	TanggalPinjam          time.Time        `json:"tanggal_pinjam"`
	TanggalKembaliRencana  time.Time        `json:"tanggal_kembali_rencana"`
	TanggalKembaliAktual   *time.Time       `json:"tanggal_kembali_aktual,omitempty"`
	Status                 StatusPeminjaman `json:"status"`
	KondisiKembali         string           `json:"kondisi_kembali,omitempty"`
	Catatan                string           `json:"catatan,omitempty"`
	DisetujuiOleh          *uuid.UUID       `json:"disetujui_oleh,omitempty"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
}

type CreateAsetInput struct {
	BMTID            uuid.UUID
	CabangID         uuid.UUID
	KodeAset         string
	Nama             string
	Kategori         KategoriAset
	Lokasi           string
	TanggalPerolehan time.Time
	NilaiPerolehan   int64
	UmurEkonomis     *int16
	Kondisi          KondisiAset
	FotoURL          string
	KodeAkun         string
}

type CreatePeminjamanInput struct {
	BMTID                 uuid.UUID
	AsetID                uuid.UUID
	PeminjamID            uuid.UUID
	TipePeminjam          TipePeminjam
	Keperluan             string
	TanggalPinjam         time.Time
	TanggalKembaliRencana time.Time
	DisetujuiOleh         *uuid.UUID
}

type ListAsetFilter struct {
	BMTID    *uuid.UUID
	CabangID *uuid.UUID
	Kategori *KategoriAset
	Kondisi  *KondisiAset
	IsAktif  *bool
	Page     int
	PerPage  int
}

type Repository interface {
	// Aset
	CreateAset(ctx context.Context, a *InventarisAset) error
	GetAsetByID(ctx context.Context, id uuid.UUID) (*InventarisAset, error)
	GetAsetByKode(ctx context.Context, kode string) (*InventarisAset, error)
	ListAset(ctx context.Context, filter ListAsetFilter) ([]*InventarisAset, int64, error)
	UpdateAset(ctx context.Context, a *InventarisAset) error
	UpdateNilaiBuku(ctx context.Context, id uuid.UUID, nilaiBuku int64) error
	UpdateKondisi(ctx context.Context, id uuid.UUID, kondisi KondisiAset) error

	// Peminjaman
	CreatePeminjaman(ctx context.Context, p *InventarisPeminjaman) error
	GetPeminjamanByID(ctx context.Context, id uuid.UUID) (*InventarisPeminjaman, error)
	ListPeminjamanByAset(ctx context.Context, asetID uuid.UUID, status *StatusPeminjaman) ([]*InventarisPeminjaman, error)
	ListPeminjamanAktif(ctx context.Context, bmtID uuid.UUID) ([]*InventarisPeminjaman, error)
	KembalikanAset(ctx context.Context, id uuid.UUID, kondisi string, aktual time.Time) error
}

func NewAset(input CreateAsetInput) (*InventarisAset, error) {
	if input.KodeAset == "" {
		return nil, errors.New("kode aset wajib diisi")
	}
	if input.Nama == "" {
		return nil, errors.New("nama aset wajib diisi")
	}
	if input.NilaiPerolehan <= 0 {
		return nil, errors.New("nilai perolehan aset harus lebih dari 0")
	}
	kodeAkun := input.KodeAkun
	if kodeAkun == "" {
		kodeAkun = "131"
	}
	kondisi := input.Kondisi
	if kondisi == "" {
		kondisi = KondsiBaik
	}
	now := time.Now()
	return &InventarisAset{
		ID:               uuid.New(),
		BMTID:            input.BMTID,
		CabangID:         input.CabangID,
		KodeAset:         input.KodeAset,
		Nama:             input.Nama,
		Kategori:         input.Kategori,
		Lokasi:           input.Lokasi,
		TanggalPerolehan: input.TanggalPerolehan,
		NilaiPerolehan:   input.NilaiPerolehan,
		NilaiBuku:        input.NilaiPerolehan,
		UmurEkonomis:     input.UmurEkonomis,
		Kondisi:          kondisi,
		FotoURL:          input.FotoURL,
		KodeAkun:         kodeAkun,
		IsAktif:          true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func NewPeminjaman(input CreatePeminjamanInput) (*InventarisPeminjaman, error) {
	if input.Keperluan == "" {
		return nil, errors.New("keperluan peminjaman wajib diisi")
	}
	if input.TanggalKembaliRencana.Before(input.TanggalPinjam) {
		return nil, errors.New("tanggal kembali rencana harus setelah tanggal pinjam")
	}
	now := time.Now()
	return &InventarisPeminjaman{
		ID:                    uuid.New(),
		BMTID:                 input.BMTID,
		AsetID:                input.AsetID,
		PeminjamID:            input.PeminjamID,
		TipePeminjam:          input.TipePeminjam,
		Keperluan:             input.Keperluan,
		TanggalPinjam:         input.TanggalPinjam,
		TanggalKembaliRencana: input.TanggalKembaliRencana,
		Status:                StatusDipinjam,
		DisetujuiOleh:         input.DisetujuiOleh,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}
