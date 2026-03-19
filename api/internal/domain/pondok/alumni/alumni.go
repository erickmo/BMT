package alumni

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrAlumniNotFound        = errors.New("data alumni tidak ditemukan")
	ErrAlumniSudahTerdaftar  = errors.New("santri ini sudah terdaftar sebagai alumni")
	ErrNamaTidakBolehKosong  = errors.New("nama lengkap alumni wajib diisi")
	ErrAngkatanTidakValid    = errors.New("angkatan tidak valid")
	ErrTahunLulusTidakValid  = errors.New("tahun lulus tidak valid")
)

// ── Alumni ────────────────────────────────────────────────────────────────────

// Alumni merepresentasikan santri yang telah menyelesaikan pendidikan di pondok.
// Alumni bisa terhubung ke nasabah BMT jika masih aktif sebagai anggota.
type Alumni struct {
	ID               uuid.UUID  `json:"id"`
	BMTID            uuid.UUID  `json:"bmt_id"`
	// SantriID terhubung ke pondok_santri jika data santri masih ada
	SantriID         *uuid.UUID `json:"santri_id,omitempty"`
	NamaLengkap      string     `json:"nama_lengkap"`
	Angkatan         int16      `json:"angkatan"`
	TahunLulus       int16      `json:"tahun_lulus"`
	Pekerjaan        string     `json:"pekerjaan,omitempty"`
	Instansi         string     `json:"instansi,omitempty"`
	KotaDomisili     string     `json:"kota_domisili,omitempty"`
	Telepon          string     `json:"telepon,omitempty"`
	Email            string     `json:"email,omitempty"`
	LinkedInURL      string     `json:"linkedin_url,omitempty"`
	FotoURL          string     `json:"foto_url,omitempty"`
	// IsVerified menandakan profil alumni sudah diverifikasi admin pondok
	IsVerified       bool       `json:"is_verified"`
	// NasabahID diisi jika alumni masih aktif sebagai nasabah BMT
	NasabahID        *uuid.UUID `json:"nasabah_id,omitempty"`
	IsAktifJaringan  bool       `json:"is_aktif_jaringan"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// NewAlumni membuat entitas Alumni baru.
func NewAlumni(bmtID uuid.UUID, namaLengkap string, angkatan, tahunLulus int16) (*Alumni, error) {
	if namaLengkap == "" {
		return nil, ErrNamaTidakBolehKosong
	}
	if angkatan <= 0 {
		return nil, ErrAngkatanTidakValid
	}
	if tahunLulus < int16(angkatan) {
		return nil, ErrTahunLulusTidakValid
	}
	now := time.Now()
	return &Alumni{
		ID:              uuid.New(),
		BMTID:           bmtID,
		NamaLengkap:     namaLengkap,
		Angkatan:        angkatan,
		TahunLulus:      tahunLulus,
		IsVerified:      false,
		IsAktifJaringan: true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Verifikasi menandai profil alumni sebagai terverifikasi oleh admin pondok.
func (a *Alumni) Verifikasi() {
	a.IsVerified = true
	a.UpdatedAt = time.Now()
}

// HubungkanNasabah menautkan alumni ke akun nasabah BMT-nya.
func (a *Alumni) HubungkanNasabah(nasabahID uuid.UUID) {
	a.NasabahID = &nasabahID
	a.UpdatedAt = time.Now()
}

// HubungkanSantri menautkan alumni ke data santri asalnya.
func (a *Alumni) HubungkanSantri(santriID uuid.UUID) {
	a.SantriID = &santriID
	a.UpdatedAt = time.Now()
}

// UpdateProfil memperbarui informasi profesional alumni.
func (a *Alumni) UpdateProfil(pekerjaan, instansi, kotaDomisili, telepon, email string) {
	a.Pekerjaan = pekerjaan
	a.Instansi = instansi
	a.KotaDomisili = kotaDomisili
	a.Telepon = telepon
	a.Email = email
	a.UpdatedAt = time.Now()
}

// NonaktifkanDariJaringan mengeluarkan alumni dari jaringan aktif pondok.
func (a *Alumni) NonaktifkanDariJaringan() {
	a.IsAktifJaringan = false
	a.UpdatedAt = time.Now()
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListAlumniFilter struct {
	BMTID           uuid.UUID
	Angkatan        *int16
	TahunLulus      *int16
	Pekerjaan       string
	KotaDomisili    string
	IsVerified      *bool
	IsAktifJaringan *bool
	Keyword         string
	Page            int
	PerPage         int
}

// ── Repository interface ──────────────────────────────────────────────────────

// Repository mendefinisikan kontrak akses data untuk entitas Alumni.
type Repository interface {
	Create(ctx context.Context, a *Alumni) error
	GetByID(ctx context.Context, id uuid.UUID) (*Alumni, error)
	GetBySantriID(ctx context.Context, santriID uuid.UUID) (*Alumni, error)
	List(ctx context.Context, filter ListAlumniFilter) ([]*Alumni, int64, error)
	Update(ctx context.Context, a *Alumni) error
	Delete(ctx context.Context, id uuid.UUID) error
}
