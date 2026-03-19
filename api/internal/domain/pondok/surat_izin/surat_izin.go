package surat_izin

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrSuratIzinNotFound      = errors.New("surat izin tidak ditemukan")
	ErrSuratIzinSudahDisetujui = errors.New("surat izin sudah disetujui, tidak bisa diubah")
	ErrSuratIzinSudahDitolak  = errors.New("surat izin sudah ditolak")
	ErrSuratIzinTidakMenunggu = errors.New("surat izin tidak dalam status menunggu")
	ErrKeperluanWajibDiisi    = errors.New("keperluan izin wajib diisi")
	ErrTanggalTidakValid      = errors.New("tanggal kembali harus setelah tanggal keluar")
	ErrAlasanTolakWajibDiisi  = errors.New("alasan penolakan wajib diisi")
)

// ── Jenis Izin ────────────────────────────────────────────────────────────────

type JenisIzin string

const (
	JenisIzinKeluar   JenisIzin = "KELUAR"
	JenisIzinPulang   JenisIzin = "PULANG"
	JenisIzinSakit    JenisIzin = "SAKIT"
	JenisIzinLainnya  JenisIzin = "LAINNYA"
)

// ── Status Surat Izin ─────────────────────────────────────────────────────────

type StatusSuratIzin string

const (
	StatusMenunggu   StatusSuratIzin = "MENUNGGU"
	StatusDisetujui  StatusSuratIzin = "DISETUJUI"
	StatusDitolak    StatusSuratIzin = "DITOLAK"
	StatusDibatalkan StatusSuratIzin = "DIBATALKAN"
	StatusSelesai    StatusSuratIzin = "SELESAI"
)

// ── Pengaju Izin ──────────────────────────────────────────────────────────────

type PengajuIzin string

const (
	PengajuSantri PengajuIzin = "SANTRI"
	PengajuWali   PengajuIzin = "WALI"
)

// ── SuratIzin ─────────────────────────────────────────────────────────────────

// SuratIzin merepresentasikan permohonan izin keluar pondok oleh santri atau wali.
// Persetujuan dilakukan oleh pengguna pondok (admin/BK).
type SuratIzin struct {
	ID             uuid.UUID       `json:"id"`
	BMTID          uuid.UUID       `json:"bmt_id"`
	CabangID       uuid.UUID       `json:"cabang_id"`
	SantriID       uuid.UUID       `json:"santri_id"`
	Jenis          JenisIzin       `json:"jenis"`
	Keperluan      string          `json:"keperluan"`
	Tujuan         string          `json:"tujuan,omitempty"`
	TanggalMulai   time.Time       `json:"tanggal_mulai"`
	TanggalKembali time.Time       `json:"tanggal_kembali"`
	DiajukanOleh   PengajuIzin     `json:"diajukan_oleh"`
	// NasabahWaliID diisi jika pengaju adalah wali
	NasabahWaliID  *uuid.UUID      `json:"nasabah_wali_id,omitempty"`
	Status         StatusSuratIzin `json:"status"`
	// DisetujuiOleh adalah pengguna_pondok.id yang menyetujui/menolak
	DisetujuiOleh  *uuid.UUID      `json:"disetujui_oleh,omitempty"`
	AlasanTolak    string          `json:"alasan_tolak,omitempty"`
	// Konfirmasi keberangkatan & kepulangan aktual
	WaktuKeluarAktual  *time.Time  `json:"waktu_keluar_aktual,omitempty"`
	WaktuKembaliAktual *time.Time  `json:"waktu_kembali_aktual,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// NewSuratIzin membuat permohonan surat izin baru.
func NewSuratIzin(bmtID, cabangID, santriID uuid.UUID, jenis JenisIzin, keperluan string, tanggalMulai, tanggalKembali time.Time, diajukanOleh PengajuIzin) (*SuratIzin, error) {
	if keperluan == "" {
		return nil, ErrKeperluanWajibDiisi
	}
	if !tanggalMulai.Before(tanggalKembali) {
		return nil, ErrTanggalTidakValid
	}
	now := time.Now()
	return &SuratIzin{
		ID:             uuid.New(),
		BMTID:          bmtID,
		CabangID:       cabangID,
		SantriID:       santriID,
		Jenis:          jenis,
		Keperluan:      keperluan,
		TanggalMulai:   tanggalMulai,
		TanggalKembali: tanggalKembali,
		DiajukanOleh:   diajukanOleh,
		Status:         StatusMenunggu,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Setujui menyetujui surat izin.
func (s *SuratIzin) Setujui(disetujuiOleh uuid.UUID) error {
	if s.Status != StatusMenunggu {
		return ErrSuratIzinTidakMenunggu
	}
	s.Status = StatusDisetujui
	s.DisetujuiOleh = &disetujuiOleh
	s.UpdatedAt = time.Now()
	return nil
}

// Tolak menolak surat izin dengan alasan.
func (s *SuratIzin) Tolak(disetujuiOleh uuid.UUID, alasanTolak string) error {
	if s.Status != StatusMenunggu {
		return ErrSuratIzinTidakMenunggu
	}
	if alasanTolak == "" {
		return ErrAlasanTolakWajibDiisi
	}
	s.Status = StatusDitolak
	s.DisetujuiOleh = &disetujuiOleh
	s.AlasanTolak = alasanTolak
	s.UpdatedAt = time.Now()
	return nil
}

// BatalkanOlehPengaju membatalkan permohonan sebelum diproses.
func (s *SuratIzin) BatalkanOlehPengaju() error {
	if s.Status != StatusMenunggu {
		return errors.New("hanya surat izin berstatus menunggu yang bisa dibatalkan")
	}
	s.Status = StatusDibatalkan
	s.UpdatedAt = time.Now()
	return nil
}

// CatatKeluar mencatat waktu keberangkatan aktual santri.
func (s *SuratIzin) CatatKeluar(waktu time.Time) error {
	if s.Status != StatusDisetujui {
		return ErrSuratIzinSudahDitolak
	}
	s.WaktuKeluarAktual = &waktu
	s.UpdatedAt = time.Now()
	return nil
}

// CatatKembali mencatat waktu kepulangan aktual santri dan menyelesaikan izin.
func (s *SuratIzin) CatatKembali(waktu time.Time) error {
	if s.Status != StatusDisetujui {
		return errors.New("surat izin belum disetujui atau sudah selesai")
	}
	s.WaktuKembaliAktual = &waktu
	s.Status = StatusSelesai
	s.UpdatedAt = time.Now()
	return nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListSuratIzinFilter struct {
	BMTID         uuid.UUID
	CabangID      uuid.UUID
	SantriID      *uuid.UUID
	Jenis         JenisIzin
	Status        StatusSuratIzin
	DiajukanOleh  PengajuIzin
	DariTgl       *time.Time
	SampaiTgl     *time.Time
	Page          int
	PerPage       int
}

// ── Repository interface ──────────────────────────────────────────────────────

// Repository mendefinisikan kontrak akses data untuk entitas SuratIzin.
type Repository interface {
	Create(ctx context.Context, s *SuratIzin) error
	GetByID(ctx context.Context, id uuid.UUID) (*SuratIzin, error)
	List(ctx context.Context, filter ListSuratIzinFilter) ([]*SuratIzin, int64, error)
	Update(ctx context.Context, s *SuratIzin) error
	Delete(ctx context.Context, id uuid.UUID) error
}
