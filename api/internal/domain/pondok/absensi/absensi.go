package absensi

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrAbsensiNotFound         = errors.New("data absensi tidak ditemukan")
	ErrAbsensiSudahAda         = errors.New("absensi untuk sesi ini sudah tercatat")
	ErrMetodeTidakDiizinkan    = errors.New("metode absensi tidak diizinkan oleh settings BMT")
	ErrSubjekTidakValid        = errors.New("tipe subjek absensi tidak valid")
	ErrStatusTidakValid        = errors.New("status absensi tidak valid")
	ErrTanggalWajibDiisi       = errors.New("tanggal absensi wajib diisi")
)

// ── Tipe Subjek Absensi ───────────────────────────────────────────────────────

// TipeSubjek mendefinisikan siapa yang diabsen.
// Nilai dari settings BMT, bukan hardcode.
type TipeSubjek string

const (
	TipeSubjekSantri    TipeSubjek = "SANTRI"
	TipeSubjekPengajar  TipeSubjek = "PENGAJAR"
	TipeSubjekKaryawan  TipeSubjek = "KARYAWAN"
)

// ── Status Absensi ────────────────────────────────────────────────────────────

type StatusAbsensi string

const (
	StatusHadir     StatusAbsensi = "HADIR"
	StatusSakit     StatusAbsensi = "SAKIT"
	StatusIzin      StatusAbsensi = "IZIN"
	StatusAlfa      StatusAbsensi = "ALFA"
	StatusTerlambat StatusAbsensi = "TERLAMBAT"
)

// ── Metode Absensi ────────────────────────────────────────────────────────────

// MetodeAbsensi mendefinisikan cara pencatatan absensi.
// Metode yang diizinkan diambil dari settings BMT ("pondok.absensi_metode"),
// bukan dari konstanta ini — konstanta hanya digunakan untuk type-safety.
type MetodeAbsensi string

const (
	MetodeManual    MetodeAbsensi = "MANUAL"
	MetodeNFC       MetodeAbsensi = "NFC"
	MetodeBiometrik MetodeAbsensi = "BIOMETRIK"
)

// ── Absensi ───────────────────────────────────────────────────────────────────

// Absensi merepresentasikan satu catatan kehadiran.
//
// Penting: metode yang valid diambil dari settings DB via settings.ResolveJSON,
// tidak pernah divalidasi hardcode terhadap konstanta MetodeAbsensi di atas.
// Contoh: metodeDiizinkan := settings.ResolveJSON(ctx, bmtID, cabangID, "pondok.absensi_metode")
type Absensi struct {
	ID          uuid.UUID     `json:"id"`
	BMTID       uuid.UUID     `json:"bmt_id"`
	CabangID    uuid.UUID     `json:"cabang_id"`
	SubjekID    uuid.UUID     `json:"subjek_id"`
	SubjekTipe  TipeSubjek    `json:"subjek_tipe"`
	Tanggal     time.Time     `json:"tanggal"`
	Sesi        string        `json:"sesi,omitempty"`
	JadwalID    *uuid.UUID    `json:"jadwal_id,omitempty"`
	Status      StatusAbsensi `json:"status"`
	Keterangan  string        `json:"keterangan,omitempty"`
	Metode      MetodeAbsensi `json:"metode"`
	WaktuScan   *time.Time    `json:"waktu_scan,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	// CreatedBy NULL jika via NFC/biometrik otomatis
	CreatedBy   *uuid.UUID    `json:"created_by,omitempty"`
}

// NewAbsensiManual membuat catatan absensi via input manual oleh pengguna pondok.
func NewAbsensiManual(bmtID, cabangID, subjekID uuid.UUID, subjekTipe TipeSubjek, tanggal time.Time, sesi string, status StatusAbsensi, keterangan string, createdBy uuid.UUID) (*Absensi, error) {
	if err := validateAbsensiInput(subjekTipe, status); err != nil {
		return nil, err
	}
	return &Absensi{
		ID:         uuid.New(),
		BMTID:      bmtID,
		CabangID:   cabangID,
		SubjekID:   subjekID,
		SubjekTipe: subjekTipe,
		Tanggal:    tanggal,
		Sesi:       sesi,
		Status:     status,
		Keterangan: keterangan,
		Metode:     MetodeManual,
		CreatedBy:  &createdBy,
		CreatedAt:  time.Now(),
	}, nil
}

// NewAbsensiNFC membuat catatan absensi via scan kartu NFC.
// CreatedBy dikosongkan karena proses otomatis.
func NewAbsensiNFC(bmtID, cabangID, subjekID uuid.UUID, subjekTipe TipeSubjek, tanggal time.Time, sesi string, jadwalID *uuid.UUID) (*Absensi, error) {
	if err := validateSubjekTipe(subjekTipe); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Absensi{
		ID:         uuid.New(),
		BMTID:      bmtID,
		CabangID:   cabangID,
		SubjekID:   subjekID,
		SubjekTipe: subjekTipe,
		Tanggal:    tanggal,
		Sesi:       sesi,
		JadwalID:   jadwalID,
		Status:     StatusHadir,
		Metode:     MetodeNFC,
		WaktuScan:  &now,
		CreatedAt:  now,
	}, nil
}

// NewAbsensiBiometrik membuat catatan absensi via sidik jari atau wajah.
// CreatedBy dikosongkan karena proses otomatis.
func NewAbsensiBiometrik(bmtID, cabangID, subjekID uuid.UUID, subjekTipe TipeSubjek, tanggal time.Time, sesi string, jadwalID *uuid.UUID) (*Absensi, error) {
	if err := validateSubjekTipe(subjekTipe); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Absensi{
		ID:         uuid.New(),
		BMTID:      bmtID,
		CabangID:   cabangID,
		SubjekID:   subjekID,
		SubjekTipe: subjekTipe,
		Tanggal:    tanggal,
		Sesi:       sesi,
		JadwalID:   jadwalID,
		Status:     StatusHadir,
		Metode:     MetodeBiometrik,
		WaktuScan:  &now,
		CreatedAt:  now,
	}, nil
}

// ── Rekap Absensi ─────────────────────────────────────────────────────────────

// RekapAbsensi merepresentasikan rangkuman kehadiran subjek untuk satu periode.
type RekapAbsensi struct {
	SubjekID    uuid.UUID  `json:"subjek_id"`
	SubjekTipe  TipeSubjek `json:"subjek_tipe"`
	Periode     string     `json:"periode"`
	TotalHadir  int        `json:"total_hadir"`
	TotalSakit  int        `json:"total_sakit"`
	TotalIzin   int        `json:"total_izin"`
	TotalAlfa   int        `json:"total_alfa"`
	TotalTerlambat int     `json:"total_terlambat"`
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListAbsensiFilter struct {
	BMTID      uuid.UUID
	CabangID   uuid.UUID
	SubjekID   *uuid.UUID
	SubjekTipe TipeSubjek
	KelasID    *uuid.UUID
	Tanggal    *time.Time
	DariTgl    *time.Time
	SampaiTgl  *time.Time
	Sesi       string
	Status     StatusAbsensi
	Metode     MetodeAbsensi
	Page       int
	PerPage    int
}

// ── Repository interface ──────────────────────────────────────────────────────

// Repository mendefinisikan kontrak akses data untuk entitas Absensi.
type Repository interface {
	Create(ctx context.Context, a *Absensi) error
	GetByID(ctx context.Context, id uuid.UUID) (*Absensi, error)
	// GetBySubjekTanggalSesi mencari absensi unik untuk kombinasi subjek+tanggal+sesi.
	GetBySubjekTanggalSesi(ctx context.Context, bmtID, subjekID uuid.UUID, tanggal time.Time, sesi string) (*Absensi, error)
	List(ctx context.Context, filter ListAbsensiFilter) ([]*Absensi, int64, error)
	Update(ctx context.Context, a *Absensi) error
	Delete(ctx context.Context, id uuid.UUID) error
	// Rekap menghitung total hadir/sakit/izin/alfa per subjek dalam rentang tanggal.
	Rekap(ctx context.Context, bmtID, cabangID uuid.UUID, subjekTipe TipeSubjek, dariTgl, sampaiTgl time.Time) ([]*RekapAbsensi, error)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func validateSubjekTipe(t TipeSubjek) error {
	switch t {
	case TipeSubjekSantri, TipeSubjekPengajar, TipeSubjekKaryawan:
		return nil
	}
	return ErrSubjekTidakValid
}

func validateAbsensiInput(subjekTipe TipeSubjek, status StatusAbsensi) error {
	if err := validateSubjekTipe(subjekTipe); err != nil {
		return err
	}
	switch status {
	case StatusHadir, StatusSakit, StatusIzin, StatusAlfa, StatusTerlambat:
		return nil
	}
	return ErrStatusTidakValid
}
