package akademik

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrMapelNotFound          = errors.New("mata pelajaran tidak ditemukan")
	ErrSilabusNotFound        = errors.New("silabus tidak ditemukan")
	ErrRPPNotFound            = errors.New("RPP tidak ditemukan")
	ErrKomponenNilaiNotFound  = errors.New("komponen nilai tidak ditemukan")
	ErrKodeMapelSudahAda      = errors.New("kode mata pelajaran sudah terdaftar")
	ErrBobotMelebihiSeratus   = errors.New("total bobot komponen nilai tidak boleh melebihi 100%")
	ErrMapelTidakAktif        = errors.New("mata pelajaran tidak aktif")
	ErrNamaWajibDiisi         = errors.New("nama wajib diisi")
	ErrKodeWajibDiisi         = errors.New("kode wajib diisi")
)

// ── Mapel (Mata Pelajaran) ────────────────────────────────────────────────────

// Mapel merepresentasikan mata pelajaran yang diajarkan di pondok pesantren.
type Mapel struct {
	ID        uuid.UUID `json:"id"`
	BMTID     uuid.UUID `json:"bmt_id"`
	CabangID  uuid.UUID `json:"cabang_id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	Tingkat   string    `json:"tingkat"`
	IsAktif   bool      `json:"is_aktif"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMapel membuat entitas Mapel baru dengan validasi.
func NewMapel(bmtID, cabangID uuid.UUID, kode, nama, tingkat string) (*Mapel, error) {
	if kode == "" {
		return nil, ErrKodeWajibDiisi
	}
	if nama == "" {
		return nil, ErrNamaWajibDiisi
	}
	now := time.Now()
	return &Mapel{
		ID:        uuid.New(),
		BMTID:     bmtID,
		CabangID:  cabangID,
		Kode:      kode,
		Nama:      nama,
		Tingkat:   tingkat,
		IsAktif:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ── Silabus ───────────────────────────────────────────────────────────────────

// Silabus merepresentasikan silabus per mata pelajaran per tahun ajaran dan semester.
// FileURL mengacu ke objek di MinIO.
type Silabus struct {
	ID          uuid.UUID `json:"id"`
	BMTID       uuid.UUID `json:"bmt_id"`
	MapelID     uuid.UUID `json:"mapel_id"`
	TahunAjaran string    `json:"tahun_ajaran"`
	Semester    int16     `json:"semester"`
	Deskripsi   string    `json:"deskripsi,omitempty"`
	FileURL     string    `json:"file_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

// NewSilabus membuat entitas Silabus baru.
func NewSilabus(bmtID, mapelID uuid.UUID, tahunAjaran string, semester int16, createdBy uuid.UUID) (*Silabus, error) {
	if tahunAjaran == "" {
		return nil, errors.New("tahun ajaran wajib diisi")
	}
	if semester != 1 && semester != 2 {
		return nil, errors.New("semester harus 1 atau 2")
	}
	return &Silabus{
		ID:          uuid.New(),
		BMTID:       bmtID,
		MapelID:     mapelID,
		TahunAjaran: tahunAjaran,
		Semester:    semester,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
	}, nil
}

// ── RPP (Rencana Pelaksanaan Pembelajaran) ────────────────────────────────────

// RPP merepresentasikan rencana pelaksanaan pembelajaran per pertemuan.
// MateriURL mengacu ke file materi ajar di MinIO.
type RPP struct {
	ID           uuid.UUID `json:"id"`
	BMTID        uuid.UUID `json:"bmt_id"`
	SilabusID    uuid.UUID `json:"silabus_id"`
	PengajarID   uuid.UUID `json:"pengajar_id"`
	PertemuanKe  int16     `json:"pertemuan_ke"`
	Topik        string    `json:"topik"`
	Tujuan       string    `json:"tujuan,omitempty"`
	MateriURL    string    `json:"materi_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewRPP membuat entitas RPP baru.
func NewRPP(bmtID, silabusID, pengajarID uuid.UUID, pertemuanKe int16, topik string) (*RPP, error) {
	if topik == "" {
		return nil, errors.New("topik RPP wajib diisi")
	}
	if pertemuanKe < 1 {
		return nil, errors.New("nomor pertemuan harus >= 1")
	}
	return &RPP{
		ID:          uuid.New(),
		BMTID:       bmtID,
		SilabusID:   silabusID,
		PengajarID:  pengajarID,
		PertemuanKe: pertemuanKe,
		Topik:       topik,
		CreatedAt:   time.Now(),
	}, nil
}

// ── KomponenNilai ─────────────────────────────────────────────────────────────

// KomponenNilai merepresentasikan komponen penilaian (UH, UTS, UAS, Tugas) beserta bobot persentasenya.
// Total bobot semua komponen untuk satu mapel/semester tidak boleh melebihi 100.
type KomponenNilai struct {
	ID          uuid.UUID `json:"id"`
	BMTID       uuid.UUID `json:"bmt_id"`
	MapelID     uuid.UUID `json:"mapel_id"`
	TahunAjaran string    `json:"tahun_ajaran"`
	Semester    int16     `json:"semester"`
	Nama        string    `json:"nama"`
	BobotPersen int16     `json:"bobot_persen"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewKomponenNilai membuat entitas KomponenNilai baru.
func NewKomponenNilai(bmtID, mapelID uuid.UUID, tahunAjaran string, semester int16, nama string, bobotPersen int16) (*KomponenNilai, error) {
	if nama == "" {
		return nil, ErrNamaWajibDiisi
	}
	if bobotPersen <= 0 || bobotPersen > 100 {
		return nil, errors.New("bobot persen harus antara 1 dan 100")
	}
	now := time.Now()
	return &KomponenNilai{
		ID:          uuid.New(),
		BMTID:       bmtID,
		MapelID:     mapelID,
		TahunAjaran: tahunAjaran,
		Semester:    semester,
		Nama:        nama,
		BobotPersen: bobotPersen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListMapelFilter struct {
	BMTID    uuid.UUID
	CabangID uuid.UUID
	Tingkat  string
	IsAktif  *bool
}

type ListSilabusFilter struct {
	BMTID       uuid.UUID
	MapelID     *uuid.UUID
	TahunAjaran string
	Semester    *int16
}

type ListRPPFilter struct {
	BMTID      uuid.UUID
	SilabusID  *uuid.UUID
	PengajarID *uuid.UUID
}

type ListKomponenFilter struct {
	BMTID       uuid.UUID
	MapelID     uuid.UUID
	TahunAjaran string
	Semester    *int16
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// MapelRepository mendefinisikan kontrak akses data untuk entitas Mapel.
type MapelRepository interface {
	Create(ctx context.Context, m *Mapel) error
	GetByID(ctx context.Context, id uuid.UUID) (*Mapel, error)
	GetByKode(ctx context.Context, bmtID, cabangID uuid.UUID, kode string) (*Mapel, error)
	List(ctx context.Context, filter ListMapelFilter) ([]*Mapel, error)
	Update(ctx context.Context, m *Mapel) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SilabusRepository mendefinisikan kontrak akses data untuk entitas Silabus.
type SilabusRepository interface {
	Create(ctx context.Context, s *Silabus) error
	GetByID(ctx context.Context, id uuid.UUID) (*Silabus, error)
	List(ctx context.Context, filter ListSilabusFilter) ([]*Silabus, error)
	Update(ctx context.Context, s *Silabus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// RPPRepository mendefinisikan kontrak akses data untuk entitas RPP.
type RPPRepository interface {
	Create(ctx context.Context, r *RPP) error
	GetByID(ctx context.Context, id uuid.UUID) (*RPP, error)
	List(ctx context.Context, filter ListRPPFilter) ([]*RPP, error)
	Update(ctx context.Context, r *RPP) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// KomponenNilaiRepository mendefinisikan kontrak akses data untuk KomponenNilai.
type KomponenNilaiRepository interface {
	Create(ctx context.Context, k *KomponenNilai) error
	GetByID(ctx context.Context, id uuid.UUID) (*KomponenNilai, error)
	List(ctx context.Context, filter ListKomponenFilter) ([]*KomponenNilai, error)
	SumBobot(ctx context.Context, mapelID uuid.UUID, tahunAjaran string, semester int16) (int16, error)
	Update(ctx context.Context, k *KomponenNilai) error
	Delete(ctx context.Context, id uuid.UUID) error
}
