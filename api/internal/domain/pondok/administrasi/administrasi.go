package administrasi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrSantriNotFound    = errors.New("santri tidak ditemukan")
	ErrKelasNotFound     = errors.New("kelas tidak ditemukan")
	ErrPengajarNotFound  = errors.New("pengajar tidak ditemukan")
	ErrKaryawanNotFound  = errors.New("karyawan tidak ditemukan")
	ErrKelasFull         = errors.New("kapasitas kelas sudah penuh")
	ErrNISAlreadyExists  = errors.New("nomor induk santri sudah terdaftar")
	ErrNIPAlreadyExists  = errors.New("NIP sudah terdaftar")
	ErrSantriTidakAktif  = errors.New("santri tidak aktif")
	ErrNamaTidakBolehKosong = errors.New("nama lengkap tidak boleh kosong")
)

// ── Status & Tipe constants ───────────────────────────────────────────────────

type TingkatSantri string

const (
	TingkatMTS     TingkatSantri = "MTS"
	TingkatMA      TingkatSantri = "MA"
	TingkatS1      TingkatSantri = "S1"
	TingkatTahfidz TingkatSantri = "TAHFIDZ"
)

// ── Santri ────────────────────────────────────────────────────────────────────

// Santri merepresentasikan peserta didik aktif di pondok pesantren.
// Relasi nasabah ↔ santri adalah 1:1 (nasabah_id opsional jika belum punya rekening BMT).
type Santri struct {
	ID                uuid.UUID     `json:"id"`
	BMTID             uuid.UUID     `json:"bmt_id"`
	CabangID          uuid.UUID     `json:"cabang_id"`
	NomorIndukSantri  string        `json:"nomor_induk_santri"`
	NamaLengkap       string        `json:"nama_lengkap"`
	NasabahID         *uuid.UUID    `json:"nasabah_id,omitempty"`
	Tingkat           TingkatSantri `json:"tingkat"`
	KelasID           *uuid.UUID    `json:"kelas_id,omitempty"`
	Asrama            string        `json:"asrama,omitempty"`
	Kamar             string        `json:"kamar,omitempty"`
	Angkatan          *int16        `json:"angkatan,omitempty"`
	StatusAktif       bool          `json:"status_aktif"`
	TanggalMasuk      *time.Time    `json:"tanggal_masuk,omitempty"`
	TanggalKeluar     *time.Time    `json:"tanggal_keluar,omitempty"`
	FotoURL           string        `json:"foto_url,omitempty"`
	NamaWali          string        `json:"nama_wali,omitempty"`
	TeleponWali       string        `json:"telepon_wali,omitempty"`
	NasabahWaliID     *uuid.UUID    `json:"nasabah_wali_id,omitempty"`
	// FingerprintTemplate disimpan dalam bentuk terenkripsi
	FingerprintTemplate string    `json:"-"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// NewSantri membuat entitas Santri baru dengan validasi dasar.
func NewSantri(bmtID, cabangID uuid.UUID, nis, namaLengkap string, tingkat TingkatSantri) (*Santri, error) {
	if nis == "" {
		return nil, errors.New("nomor induk santri wajib diisi")
	}
	if namaLengkap == "" {
		return nil, ErrNamaTidakBolehKosong
	}
	now := time.Now()
	return &Santri{
		ID:               uuid.New(),
		BMTID:            bmtID,
		CabangID:         cabangID,
		NomorIndukSantri: nis,
		NamaLengkap:      namaLengkap,
		Tingkat:          tingkat,
		StatusAktif:      true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Nonaktifkan menonaktifkan santri dan mencatat tanggal keluar.
func (s *Santri) Nonaktifkan(tanggalKeluar time.Time) {
	s.StatusAktif = false
	s.TanggalKeluar = &tanggalKeluar
	s.UpdatedAt = time.Now()
}

// ── Kelas ─────────────────────────────────────────────────────────────────────

// Kelas merepresentasikan rombongan belajar dalam satu tahun ajaran.
type Kelas struct {
	ID           uuid.UUID  `json:"id"`
	BMTID        uuid.UUID  `json:"bmt_id"`
	CabangID     uuid.UUID  `json:"cabang_id"`
	Nama         string     `json:"nama"`
	Tingkat      string     `json:"tingkat"`
	TahunAjaran  string     `json:"tahun_ajaran"`
	WaliKelasID  *uuid.UUID `json:"wali_kelas_id,omitempty"`
	Kapasitas    *int16     `json:"kapasitas,omitempty"`
	JumlahSantri int        `json:"jumlah_santri"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// NewKelas membuat entitas Kelas baru.
func NewKelas(bmtID, cabangID uuid.UUID, nama, tingkat, tahunAjaran string) (*Kelas, error) {
	if nama == "" {
		return nil, errors.New("nama kelas wajib diisi")
	}
	if tahunAjaran == "" {
		return nil, errors.New("tahun ajaran wajib diisi")
	}
	now := time.Now()
	return &Kelas{
		ID:          uuid.New(),
		BMTID:       bmtID,
		CabangID:    cabangID,
		Nama:        nama,
		Tingkat:     tingkat,
		TahunAjaran: tahunAjaran,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ValidasiTambahSantri memastikan kelas belum penuh sebelum menambah santri.
func (k *Kelas) ValidasiTambahSantri() error {
	if k.Kapasitas != nil && k.JumlahSantri >= int(*k.Kapasitas) {
		return fmt.Errorf("%w: kapasitas %d sudah terpenuhi", ErrKelasFull, *k.Kapasitas)
	}
	return nil
}

// ── Pengajar ──────────────────────────────────────────────────────────────────

// Pengajar merepresentasikan guru / ustadz di pondok pesantren.
type Pengajar struct {
	ID                  uuid.UUID `json:"id"`
	BMTID               uuid.UUID `json:"bmt_id"`
	CabangID            uuid.UUID `json:"cabang_id"`
	NIP                 string    `json:"nip,omitempty"`
	NamaLengkap         string    `json:"nama_lengkap"`
	Jabatan             string    `json:"jabatan,omitempty"`
	Spesialisasi        string    `json:"spesialisasi,omitempty"`
	NasabahID           *uuid.UUID `json:"nasabah_id,omitempty"`
	FingerprintTemplate string    `json:"-"`
	StatusAktif         bool      `json:"status_aktif"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// NewPengajar membuat entitas Pengajar baru.
func NewPengajar(bmtID, cabangID uuid.UUID, nip, namaLengkap string) (*Pengajar, error) {
	if namaLengkap == "" {
		return nil, ErrNamaTidakBolehKosong
	}
	now := time.Now()
	return &Pengajar{
		ID:          uuid.New(),
		BMTID:       bmtID,
		CabangID:    cabangID,
		NIP:         nip,
		NamaLengkap: namaLengkap,
		StatusAktif: true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ── Karyawan ──────────────────────────────────────────────────────────────────

// Karyawan merepresentasikan staf non-pengajar di pondok pesantren.
type Karyawan struct {
	ID                  uuid.UUID  `json:"id"`
	BMTID               uuid.UUID  `json:"bmt_id"`
	CabangID            uuid.UUID  `json:"cabang_id"`
	NIKKaryawan         string     `json:"nik_karyawan,omitempty"`
	NamaLengkap         string     `json:"nama_lengkap"`
	Jabatan             string     `json:"jabatan,omitempty"`
	Departemen          string     `json:"departemen,omitempty"`
	NasabahID           *uuid.UUID `json:"nasabah_id,omitempty"`
	FingerprintTemplate string     `json:"-"`
	StatusAktif         bool       `json:"status_aktif"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// NewKaryawan membuat entitas Karyawan baru.
func NewKaryawan(bmtID, cabangID uuid.UUID, nik, namaLengkap string) (*Karyawan, error) {
	if namaLengkap == "" {
		return nil, ErrNamaTidakBolehKosong
	}
	now := time.Now()
	return &Karyawan{
		ID:          uuid.New(),
		BMTID:       bmtID,
		CabangID:    cabangID,
		NIKKaryawan: nik,
		NamaLengkap: namaLengkap,
		StatusAktif: true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ── Filter & Input types ──────────────────────────────────────────────────────

type ListSantriFilter struct {
	BMTID       uuid.UUID
	CabangID    uuid.UUID
	KelasID     *uuid.UUID
	Tingkat     TingkatSantri
	StatusAktif *bool
	Keyword     string
	Page        int
	PerPage     int
}

type ListPengajarFilter struct {
	BMTID       uuid.UUID
	CabangID    uuid.UUID
	StatusAktif *bool
	Page        int
	PerPage     int
}

type ListKaryawanFilter struct {
	BMTID       uuid.UUID
	CabangID    uuid.UUID
	Departemen  string
	StatusAktif *bool
	Page        int
	PerPage     int
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// SantriRepository mendefinisikan kontrak akses data untuk entitas Santri.
type SantriRepository interface {
	Create(ctx context.Context, s *Santri) error
	GetByID(ctx context.Context, id uuid.UUID) (*Santri, error)
	GetByNIS(ctx context.Context, bmtID uuid.UUID, nis string) (*Santri, error)
	List(ctx context.Context, filter ListSantriFilter) ([]*Santri, int64, error)
	Update(ctx context.Context, s *Santri) error
	UpdateKelas(ctx context.Context, santriID, kelasID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// KelasRepository mendefinisikan kontrak akses data untuk entitas Kelas.
type KelasRepository interface {
	Create(ctx context.Context, k *Kelas) error
	GetByID(ctx context.Context, id uuid.UUID) (*Kelas, error)
	List(ctx context.Context, bmtID, cabangID uuid.UUID, tahunAjaran string) ([]*Kelas, error)
	Update(ctx context.Context, k *Kelas) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountSantri(ctx context.Context, kelasID uuid.UUID) (int, error)
}

// PengajarRepository mendefinisikan kontrak akses data untuk entitas Pengajar.
type PengajarRepository interface {
	Create(ctx context.Context, p *Pengajar) error
	GetByID(ctx context.Context, id uuid.UUID) (*Pengajar, error)
	GetByNIP(ctx context.Context, bmtID uuid.UUID, nip string) (*Pengajar, error)
	List(ctx context.Context, filter ListPengajarFilter) ([]*Pengajar, int64, error)
	Update(ctx context.Context, p *Pengajar) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// KaryawanRepository mendefinisikan kontrak akses data untuk entitas Karyawan.
type KaryawanRepository interface {
	Create(ctx context.Context, k *Karyawan) error
	GetByID(ctx context.Context, id uuid.UUID) (*Karyawan, error)
	GetByNIK(ctx context.Context, bmtID uuid.UUID, nik string) (*Karyawan, error)
	List(ctx context.Context, filter ListKaryawanFilter) ([]*Karyawan, int64, error)
	Update(ctx context.Context, k *Karyawan) error
	Delete(ctx context.Context, id uuid.UUID) error
}
