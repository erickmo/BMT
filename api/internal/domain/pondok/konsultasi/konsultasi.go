package konsultasi

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrSesiNotFound      = errors.New("sesi konsultasi tidak ditemukan")
	ErrPesanNotFound     = errors.New("pesan konsultasi tidak ditemukan")
	ErrSesiSudahTutup    = errors.New("sesi konsultasi sudah ditutup")
	ErrSesiSudahDijawab  = errors.New("sesi konsultasi sudah dijawab")
	ErrTopikWajibDiisi   = errors.New("topik konsultasi wajib diisi")
	ErrJudulWajibDiisi   = errors.New("judul sesi konsultasi wajib diisi")
	ErrPesanKosong       = errors.New("isi pesan tidak boleh kosong")
	ErrTidakBerhak       = errors.New("tidak berhak mengakses sesi konsultasi ini")
)

// ── Topik Konsultasi ──────────────────────────────────────────────────────────

type TopikKonsultasi string

const (
	TopikAkademik   TopikKonsultasi = "AKADEMIK"
	TopikBK         TopikKonsultasi = "BK"
	TopikKesehatan  TopikKonsultasi = "KESEHATAN"
	TopikKeuangan   TopikKonsultasi = "KEUANGAN"
	TopikUmum       TopikKonsultasi = "UMUM"
)

// ── Status Sesi ───────────────────────────────────────────────────────────────

type StatusSesi string

const (
	StatusOpen     StatusSesi = "OPEN"
	StatusDijawab  StatusSesi = "DIJAWAB"
	StatusDitutup  StatusSesi = "DITUTUP"
)

// ── Tipe Penanya ──────────────────────────────────────────────────────────────

type TipePenanya string

const (
	TipeSantri TipePenanya = "SANTRI"
	TipeWali   TipePenanya = "WALI"
)

// ── Tipe Pengirim Pesan ───────────────────────────────────────────────────────

type TipePengirim string

const (
	TipePengirimSantri         TipePengirim = "SANTRI"
	TipePengirimWali           TipePengirim = "WALI"
	TipePengirimPenggunaPondok TipePengirim = "PENGGUNA_PONDOK"
)

// ── KonsultasiSesi ────────────────────────────────────────────────────────────

// KonsultasiSesi merepresentasikan sesi konsultasi antara santri/wali dan konselor.
// CatatanPrivat hanya bisa dibaca oleh konselor yang ditugaskan.
type KonsultasiSesi struct {
	ID            uuid.UUID       `json:"id"`
	BMTID         uuid.UUID       `json:"bmt_id"`
	PenanyaID     uuid.UUID       `json:"penanya_id"`
	PenanyaTipe   TipePenanya     `json:"penanya_tipe"`
	// PenjawabID NULL jika belum ditugaskan ke konselor
	PenjawabID    *uuid.UUID      `json:"penjawab_id,omitempty"`
	Topik         TopikKonsultasi `json:"topik"`
	Judul         string          `json:"judul"`
	Status        StatusSesi      `json:"status"`
	// CatatanPrivat hanya bisa diakses konselor — tidak pernah dikirim ke client selain konselor
	CatatanPrivat string          `json:"-"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// NewKonsultasiSesi membuat sesi konsultasi baru.
func NewKonsultasiSesi(bmtID, penanyaID uuid.UUID, penanyaTipe TipePenanya, topik TopikKonsultasi, judul string) (*KonsultasiSesi, error) {
	if judul == "" {
		return nil, ErrJudulWajibDiisi
	}
	now := time.Now()
	return &KonsultasiSesi{
		ID:          uuid.New(),
		BMTID:       bmtID,
		PenanyaID:   penanyaID,
		PenanyaTipe: penanyaTipe,
		Topik:       topik,
		Judul:       judul,
		Status:      StatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// TugaskanKonselor menugaskan konselor ke sesi ini.
func (s *KonsultasiSesi) TugaskanKonselor(penjawabID uuid.UUID) error {
	if s.Status == StatusDitutup {
		return ErrSesiSudahTutup
	}
	s.PenjawabID = &penjawabID
	s.UpdatedAt = time.Now()
	return nil
}

// Jawab mengubah status sesi menjadi DIJAWAB setelah konselor membalas.
func (s *KonsultasiSesi) Jawab() error {
	if s.Status == StatusDitutup {
		return ErrSesiSudahTutup
	}
	s.Status = StatusDijawab
	s.UpdatedAt = time.Now()
	return nil
}

// Tutup menutup sesi konsultasi.
func (s *KonsultasiSesi) Tutup() {
	s.Status = StatusDitutup
	s.UpdatedAt = time.Now()
}

// ── KonsultasiPesan ───────────────────────────────────────────────────────────

// KonsultasiPesan merepresentasikan satu pesan dalam thread konsultasi.
// LampiranURL mengacu ke file di MinIO.
type KonsultasiPesan struct {
	ID           uuid.UUID    `json:"id"`
	SesiID       uuid.UUID    `json:"sesi_id"`
	PengirimID   uuid.UUID    `json:"pengirim_id"`
	PengirimTipe TipePengirim `json:"pengirim_tipe"`
	Pesan        string       `json:"pesan"`
	LampiranURL  string       `json:"lampiran_url,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

// NewKonsultasiPesan membuat pesan baru dalam sesi konsultasi.
func NewKonsultasiPesan(sesiID, pengirimID uuid.UUID, pengirimTipe TipePengirim, pesan string) (*KonsultasiPesan, error) {
	if pesan == "" {
		return nil, ErrPesanKosong
	}
	return &KonsultasiPesan{
		ID:           uuid.New(),
		SesiID:       sesiID,
		PengirimID:   pengirimID,
		PengirimTipe: pengirimTipe,
		Pesan:        pesan,
		CreatedAt:    time.Now(),
	}, nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListSesiFilter struct {
	BMTID       uuid.UUID
	PenanyaID   *uuid.UUID
	PenjawabID  *uuid.UUID
	Topik       TopikKonsultasi
	Status      StatusSesi
	Page        int
	PerPage     int
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// SesiRepository mendefinisikan kontrak akses data untuk KonsultasiSesi.
type SesiRepository interface {
	Create(ctx context.Context, s *KonsultasiSesi) error
	GetByID(ctx context.Context, id uuid.UUID) (*KonsultasiSesi, error)
	List(ctx context.Context, filter ListSesiFilter) ([]*KonsultasiSesi, int64, error)
	Update(ctx context.Context, s *KonsultasiSesi) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PesanRepository mendefinisikan kontrak akses data untuk KonsultasiPesan.
type PesanRepository interface {
	Create(ctx context.Context, p *KonsultasiPesan) error
	GetByID(ctx context.Context, id uuid.UUID) (*KonsultasiPesan, error)
	ListBySesi(ctx context.Context, sesiID uuid.UUID, limit, offset int) ([]*KonsultasiPesan, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
