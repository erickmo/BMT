package health_record

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrKesehatanNotFound   = errors.New("data kesehatan santri tidak ditemukan")
	ErrKunjunganNotFound   = errors.New("data kunjungan UKS tidak ditemukan")
	ErrKeluhanWajibDiisi   = errors.New("keluhan wajib diisi")
	ErrSantriSudahTerdaftar = errors.New("data kesehatan dasar santri sudah terdaftar")
)

// ── Jenis Kunjungan ───────────────────────────────────────────────────────────

type JenisKunjungan string

const (
	JenisKunjunganSakit             JenisKunjungan = "SAKIT"
	JenisKunjunganPemeriksaanRutin  JenisKunjungan = "PEMERIKSAAN_RUTIN"
	JenisKunjunganKecelakaan        JenisKunjungan = "KECELAKAAN"
	JenisKunjunganRujukan           JenisKunjungan = "RUJUKAN"
)

// ── KesehatanSantri ───────────────────────────────────────────────────────────

// KesehatanSantri merepresentasikan data kesehatan dasar santri (one-per-santri).
// Alergi disimpan sebagai slice string yang dikonversi ke TEXT[] di PostgreSQL.
type KesehatanSantri struct {
	ID              uuid.UUID `json:"id"`
	BMTID           uuid.UUID `json:"bmt_id"`
	SantriID        uuid.UUID `json:"santri_id"`
	GolonganDarah   string    `json:"golongan_darah,omitempty"`
	Alergi          []string  `json:"alergi,omitempty"`
	RiwayatPenyakit string    `json:"riwayat_penyakit,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NewKesehatanSantri membuat data kesehatan dasar santri baru.
func NewKesehatanSantri(bmtID, santriID uuid.UUID) (*KesehatanSantri, error) {
	now := time.Now()
	return &KesehatanSantri{
		ID:        uuid.New(),
		BMTID:     bmtID,
		SantriID:  santriID,
		Alergi:    []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// TambahAlergi menambahkan entri alergi baru ke daftar.
func (k *KesehatanSantri) TambahAlergi(alergi string) {
	if alergi == "" {
		return
	}
	// Cegah duplikat
	for _, a := range k.Alergi {
		if a == alergi {
			return
		}
	}
	k.Alergi = append(k.Alergi, alergi)
	k.UpdatedAt = time.Now()
}

// ── KesehatanKunjungan ────────────────────────────────────────────────────────

// KesehatanKunjungan merepresentasikan satu kunjungan santri ke UKS.
// DicatatOleh adalah pengguna_pondok.id (petugas UKS).
type KesehatanKunjungan struct {
	ID               uuid.UUID      `json:"id"`
	BMTID            uuid.UUID      `json:"bmt_id"`
	SantriID         uuid.UUID      `json:"santri_id"`
	JenisKunjungan   JenisKunjungan `json:"jenis_kunjungan"`
	Tanggal          time.Time      `json:"tanggal"`
	Keluhan          string         `json:"keluhan"`
	Diagnosa         string         `json:"diagnosa,omitempty"`
	Tindakan         string         `json:"tindakan,omitempty"`
	ObatDiberikan    string         `json:"obat_diberikan,omitempty"`
	PerluRujukan     bool           `json:"perlu_rujukan"`
	FasilitasRujukan string         `json:"fasilitas_rujukan,omitempty"`
	DicatatOleh      uuid.UUID      `json:"dicatat_oleh"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// NewKesehatanKunjungan membuat catatan kunjungan UKS baru.
func NewKesehatanKunjungan(bmtID, santriID uuid.UUID, jenisKunjungan JenisKunjungan, tanggal time.Time, keluhan string, dicatatOleh uuid.UUID) (*KesehatanKunjungan, error) {
	if keluhan == "" {
		return nil, ErrKeluhanWajibDiisi
	}
	now := time.Now()
	return &KesehatanKunjungan{
		ID:             uuid.New(),
		BMTID:          bmtID,
		SantriID:       santriID,
		JenisKunjungan: jenisKunjungan,
		Tanggal:        tanggal,
		Keluhan:        keluhan,
		PerluRujukan:   false,
		DicatatOleh:    dicatatOleh,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// SetRujukan menandai kunjungan memerlukan rujukan ke fasilitas kesehatan luar.
func (k *KesehatanKunjungan) SetRujukan(fasilitas string) error {
	if fasilitas == "" {
		return errors.New("nama fasilitas rujukan wajib diisi")
	}
	k.PerluRujukan = true
	k.FasilitasRujukan = fasilitas
	k.UpdatedAt = time.Now()
	return nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListKunjunganFilter struct {
	BMTID          uuid.UUID
	SantriID       *uuid.UUID
	JenisKunjungan JenisKunjungan
	DariTgl        *time.Time
	SampaiTgl      *time.Time
	PerluRujukan   *bool
	Page           int
	PerPage        int
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// KesehatanSantriRepository mendefinisikan kontrak akses data untuk KesehatanSantri.
type KesehatanSantriRepository interface {
	Create(ctx context.Context, k *KesehatanSantri) error
	GetBySantriID(ctx context.Context, santriID uuid.UUID) (*KesehatanSantri, error)
	Update(ctx context.Context, k *KesehatanSantri) error
}

// KesehatanKunjunganRepository mendefinisikan kontrak akses data untuk KesehatanKunjungan.
type KesehatanKunjunganRepository interface {
	Create(ctx context.Context, k *KesehatanKunjungan) error
	GetByID(ctx context.Context, id uuid.UUID) (*KesehatanKunjungan, error)
	List(ctx context.Context, filter ListKunjunganFilter) ([]*KesehatanKunjungan, int64, error)
	Update(ctx context.Context, k *KesehatanKunjungan) error
	Delete(ctx context.Context, id uuid.UUID) error
}
