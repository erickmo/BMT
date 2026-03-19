package autodebet

import (
	"context"
	"errors"
	"time"

	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
)

var (
	ErrConfigNotFound = errors.New("konfigurasi autodebet tidak ditemukan")
	ErrJadwalNotFound = errors.New("jadwal autodebet tidak ditemukan")
)

type JenisAutodebet string

const (
	JenisSimpananWajib      JenisAutodebet = "SIMPANAN_WAJIB"
	JenisBiayaAdmin         JenisAutodebet = "BIAYA_ADMIN_REKENING"
	JenisAngsuranPembiayaan JenisAutodebet = "ANGSURAN_PEMBIAYAAN"
	JenisSPPPondok          JenisAutodebet = "SPP_PONDOK"
)

type StatusJadwal string

const (
	StatusMenunggu StatusJadwal = "MENUNGGU"
	StatusSukses   StatusJadwal = "SUKSES"
	StatusGagal    StatusJadwal = "GAGAL"
	StatusPartial  StatusJadwal = "PARTIAL"
)

type Config struct {
	ID           uuid.UUID      `json:"id"`
	BMTID        uuid.UUID      `json:"bmt_id"`
	RekeningID   uuid.UUID      `json:"rekening_id"`
	Jenis        JenisAutodebet `json:"jenis"`
	TanggalDebet int16          `json:"tanggal_debet"`
	IsAktif      bool           `json:"is_aktif"`
	ReferensiID  *uuid.UUID     `json:"referensi_id,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	UpdatedBy    uuid.UUID      `json:"updated_by"`
}

type Jadwal struct {
	ID                uuid.UUID      `json:"id"`
	BMTID             uuid.UUID      `json:"bmt_id"`
	RekeningID        uuid.UUID      `json:"rekening_id"`
	ConfigID          uuid.UUID      `json:"config_id"`
	Jenis             JenisAutodebet `json:"jenis"`
	NominalTarget     money.Money    `json:"nominal_target"`
	TanggalJatuhTempo time.Time      `json:"tanggal_jatuh_tempo"`
	Status            StatusJadwal   `json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
}

type Tunggakan struct {
	ID              uuid.UUID      `json:"id"`
	BMTID           uuid.UUID      `json:"bmt_id"`
	RekeningID      uuid.UUID      `json:"rekening_id"`
	JadwalID        uuid.UUID      `json:"jadwal_id"`
	Jenis           JenisAutodebet `json:"jenis"`
	NominalTarget   money.Money    `json:"nominal_target"`
	NominalTerbayar money.Money    `json:"nominal_terbayar"`
	NominalSisa     money.Money    `json:"nominal_sisa"`
	Status          string         `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type HasilAutodebet struct {
	JadwalID         uuid.UUID   `json:"jadwal_id"`
	RekeningID       uuid.UUID   `json:"rekening_id"`
	NominalTarget    money.Money `json:"nominal_target"`
	NominalDidebit   money.Money `json:"nominal_didebit"`
	NominalTunggakan money.Money `json:"nominal_tunggakan"`
	IsPartial        bool        `json:"is_partial"`
	TunggakanID      *uuid.UUID  `json:"tunggakan_id,omitempty"`
}

type Repository interface {
	CreateConfig(ctx context.Context, c *Config) error
	GetConfig(ctx context.Context, id uuid.UUID) (*Config, error)
	ListConfigByRekening(ctx context.Context, rekeningID uuid.UUID) ([]*Config, error)
	UpdateConfig(ctx context.Context, c *Config) error

	CreateJadwal(ctx context.Context, j *Jadwal) error
	ListJadwalByTanggal(ctx context.Context, bmtID uuid.UUID, tanggal time.Time) ([]*Jadwal, error)
	UpdateJadwalStatus(ctx context.Context, id uuid.UUID, status StatusJadwal) error

	CreateTunggakan(ctx context.Context, t *Tunggakan) error
	ListTunggakanByRekening(ctx context.Context, rekeningID uuid.UUID) ([]*Tunggakan, error)
	UpdateTunggakan(ctx context.Context, t *Tunggakan) error
}

// EksekusiAutodebet menghitung partial debit.
// Jika saldo < target: debit semampu saldo, sisanya jadi tunggakan.
func EksekusiAutodebet(saldoRekening, nominalTarget money.Money) HasilAutodebet {
	berhasil := money.Min(saldoRekening, nominalTarget)
	sisa := nominalTarget.Sub(berhasil)

	return HasilAutodebet{
		NominalTarget:    nominalTarget,
		NominalDidebit:   berhasil,
		NominalTunggakan: sisa,
		IsPartial:        sisa > 0,
	}
}
