package integrasi

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrIntegrasiLogNotFound = errors.New("log integrasi tidak ditemukan")
	ErrProviderTidakValid   = errors.New("provider integrasi tidak valid")
	ErrIntegrasiBelumAktif  = errors.New("integrasi eksternal belum dikonfigurasi atau tidak aktif")
)

type Provider string

const (
	ProviderDAPODIK Provider = "DAPODIK"
	ProviderEMIS    Provider = "EMIS"
	ProviderPPDB    Provider = "PPDB"
)

type ArahSinkron string

const (
	ArahPull ArahSinkron = "PULL"
	ArahPush ArahSinkron = "PUSH"
)

type StatusSinkron string

const (
	StatusSukses  StatusSinkron = "SUKSES"
	StatusGagal   StatusSinkron = "GAGAL"
	StatusPartial StatusSinkron = "PARTIAL"
)

type JenisSinkron string

const (
	SinkronSantri    JenisSinkron = "SANTRI"
	SinkronLembaga   JenisSinkron = "LEMBAGA"
	SinkronGuru      JenisSinkron = "GURU"
	SinkronKurikulum JenisSinkron = "KURIKULUM"
)

// IntegrasiLog records a synchronization event with an external system
type IntegrasiLog struct {
	ID            uuid.UUID       `json:"id"`
	BMTID         uuid.UUID       `json:"bmt_id"`
	Provider      Provider        `json:"provider"`
	Jenis         JenisSinkron    `json:"jenis"`
	Arah          ArahSinkron     `json:"arah"`
	Status        StatusSinkron   `json:"status"`
	JumlahRecord  int             `json:"jumlah_record"`
	Berhasil      int             `json:"berhasil"`
	Gagal         int             `json:"gagal"`
	ErrorDetail   json.RawMessage `json:"error_detail,omitempty"`
	DijalankanOleh *uuid.UUID     `json:"dijalankan_oleh,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type CreateLogInput struct {
	BMTID          uuid.UUID
	Provider       Provider
	Jenis          JenisSinkron
	Arah           ArahSinkron
	Status         StatusSinkron
	JumlahRecord   int
	Berhasil       int
	Gagal          int
	ErrorDetail    json.RawMessage
	DijalankanOleh *uuid.UUID
}

type ListLogFilter struct {
	BMTID    *uuid.UUID
	Provider *Provider
	Status   *StatusSinkron
	DariTgl  *time.Time
	Page     int
	PerPage  int
}

type Repository interface {
	Create(ctx context.Context, l *IntegrasiLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*IntegrasiLog, error)
	List(ctx context.Context, filter ListLogFilter) ([]*IntegrasiLog, int64, error)
	GetLatestByBMTAndProvider(ctx context.Context, bmtID uuid.UUID, provider Provider) (*IntegrasiLog, error)
}

func NewIntegrasiLog(input CreateLogInput) (*IntegrasiLog, error) {
	switch input.Provider {
	case ProviderDAPODIK, ProviderEMIS, ProviderPPDB:
	default:
		return nil, ErrProviderTidakValid
	}
	return &IntegrasiLog{
		ID:             uuid.New(),
		BMTID:          input.BMTID,
		Provider:       input.Provider,
		Jenis:          input.Jenis,
		Arah:           input.Arah,
		Status:         input.Status,
		JumlahRecord:   input.JumlahRecord,
		Berhasil:       input.Berhasil,
		Gagal:          input.Gagal,
		ErrorDetail:    input.ErrorDetail,
		DijalankanOleh: input.DijalankanOleh,
		CreatedAt:      time.Now(),
	}, nil
}

func (l *IntegrasiLog) RingkasanStatus() string {
	if l.Status == StatusSukses {
		return "Sinkronisasi berhasil: " + string(l.Provider)
	}
	if l.Status == StatusPartial {
		return "Sinkronisasi sebagian: " + string(l.Provider)
	}
	return "Sinkronisasi gagal: " + string(l.Provider)
}

// ValidasiProvider checks that a provider string is one of the supported values
func ValidasiProvider(p string) error {
	switch Provider(p) {
	case ProviderDAPODIK, ProviderEMIS, ProviderPPDB:
		return nil
	}
	return errors.New("provider tidak dikenal: " + p)
}
