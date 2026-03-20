package rekening

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
)

var (
	ErrRekeningNotFound    = errors.New("rekening tidak ditemukan")
	ErrRekeningBeku        = errors.New("rekening dalam status blokir")
	ErrRekeningTutup       = errors.New("rekening sudah ditutup")
	ErrSaldoTidakCukup     = errors.New("saldo tidak mencukupi")
	ErrSetoranDibawahMin   = errors.New("setoran di bawah minimum")
	ErrPenarikanTidakBisa  = errors.New("jenis rekening tidak bisa ditarik")
)

type StatusRekening string

const (
	StatusAktif  StatusRekening = "AKTIF"
	StatusBlokir StatusRekening = "BLOKIR"
	StatusTutup  StatusRekening = "TUTUP"
)

type JenisRekening struct {
	ID                 uuid.UUID `json:"id"`
	BMTID              uuid.UUID `json:"bmt_id"`
	Kode               string    `json:"kode"`
	Nama               string    `json:"nama"`
	TipeDasar          string    `json:"tipe_dasar"`
	Akad               string    `json:"akad"`
	Deskripsi          string    `json:"deskripsi"`
	SetoranAwalMin     int64     `json:"setoran_awal_min"`
	SetoranMin         int64     `json:"setoran_min"`
	BisaDitarik        bool      `json:"bisa_ditarik"`
	SyaratPenarikan    string    `json:"syarat_penarikan"`
	NisbahNasabah      *int16    `json:"nisbah_nasabah,omitempty"`
	JangkaHari         *int16    `json:"jangka_hari,omitempty"`
	BiayaAdminBulanan  int64     `json:"biaya_admin_bulanan"`
	BisaNFC            bool      `json:"bisa_nfc"`
	BisaAutodebet      bool      `json:"bisa_autodebet"`
	BiayaAdminBuka     int64     `json:"biaya_admin_buka"`
	IsAktif            bool      `json:"is_aktif"`
	UrutanTampil       int16     `json:"urutan_tampil"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedBy          uuid.UUID `json:"created_by"`
	UpdatedBy          uuid.UUID `json:"updated_by"`
}

type Rekening struct {
	ID                uuid.UUID      `json:"id"`
	BMTID             uuid.UUID      `json:"bmt_id"`
	CabangID          uuid.UUID      `json:"cabang_id"`
	NasabahID         uuid.UUID      `json:"nasabah_id"`
	JenisRekeningID   uuid.UUID      `json:"jenis_rekening_id"`
	NomorRekening     string         `json:"nomor_rekening"`
	Saldo             money.Money    `json:"saldo"`
	Status            StatusRekening `json:"status"`
	AlasanBlokir      string         `json:"alasan_blokir,omitempty"`
	BiayaAdminBulanan int64          `json:"biaya_admin_bulanan"`
	NominalDeposito   *int64         `json:"nominal_deposito,omitempty"`
	NisbahNasabah     *int16         `json:"nisbah_nasabah,omitempty"`
	TanggalBuka       time.Time      `json:"tanggal_buka"`
	TanggalJatuhTempo *time.Time     `json:"tanggal_jatuh_tempo,omitempty"`
	TanggalTutup      *time.Time     `json:"tanggal_tutup,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	CreatedByFormID   *uuid.UUID     `json:"created_by_form_id,omitempty"`
}

type TransaksiRekening struct {
	ID             uuid.UUID  `json:"id"`
	BMTID          uuid.UUID  `json:"bmt_id"`
	CabangID       uuid.UUID  `json:"cabang_id"`
	RekeningID     uuid.UUID  `json:"rekening_id"`
	Jenis          string     `json:"jenis"`
	Posisi         string     `json:"posisi"`
	Nominal        int64      `json:"nominal"`
	SaldoSebelum   int64      `json:"saldo_sebelum"`
	SaldoSesudah   int64      `json:"saldo_sesudah"`
	Keterangan     string     `json:"keterangan"`
	ReferensiID    *uuid.UUID `json:"referensi_id,omitempty"`
	ReferensiTipe  string     `json:"referensi_tipe,omitempty"`
	IdempotencyKey *uuid.UUID `json:"idempotency_key,omitempty"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type SetoranInput struct {
	RekeningID     uuid.UUID  `json:"rekening_id" validate:"required"`
	Nominal        int64      `json:"nominal" validate:"required,min=1"`
	Keterangan     string     `json:"keterangan"`
	IdempotencyKey *uuid.UUID `json:"idempotency_key"`
	CreatedBy      uuid.UUID  `json:"created_by"`
}

type PenarikanInput struct {
	RekeningID     uuid.UUID  `json:"rekening_id" validate:"required"`
	Nominal        int64      `json:"nominal" validate:"required,min=1"`
	Keterangan     string     `json:"keterangan"`
	IdempotencyKey *uuid.UUID `json:"idempotency_key"`
	CreatedBy      uuid.UUID  `json:"created_by"`
}

type Repository interface {
	// JenisRekening
	CreateJenis(ctx context.Context, jr *JenisRekening) error
	GetJenisByID(ctx context.Context, id uuid.UUID) (*JenisRekening, error)
	ListJenisByBMT(ctx context.Context, bmtID uuid.UUID) ([]*JenisRekening, error)
	UpdateJenis(ctx context.Context, jr *JenisRekening) error

	// Rekening
	Create(ctx context.Context, r *Rekening) error
	GetByID(ctx context.Context, id uuid.UUID) (*Rekening, error)
	GetByNomor(ctx context.Context, nomor string) (*Rekening, error)
	ListByNasabah(ctx context.Context, nasabahID uuid.UUID) ([]*Rekening, error)
	ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, page, perPage int) ([]*Rekening, int64, error)
	ListDepositoAktif(ctx context.Context, bmtID uuid.UUID) ([]*Rekening, error)
	UpdateSaldo(ctx context.Context, id uuid.UUID, saldoBaru int64) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusRekening, alasan string) error
	LockForUpdate(ctx context.Context, id uuid.UUID) (*Rekening, error)

	// Transaksi
	CreateTransaksi(ctx context.Context, t *TransaksiRekening) error
	ListTransaksi(ctx context.Context, rekeningID uuid.UUID, limit, offset int) ([]*TransaksiRekening, int64, error)
	GetTransaksiByIdempotency(ctx context.Context, key uuid.UUID) (*TransaksiRekening, error)

	GenerateNomorRekening(ctx context.Context, bmtID, cabangID uuid.UUID, kodeJenis string) (string, error)
}

func NewRekening(
	bmtID, cabangID, nasabahID, jenisID uuid.UUID,
	nomor string,
	jenis *JenisRekening,
) (*Rekening, error) {
	if nomor == "" {
		return nil, errors.New("nomor rekening wajib diisi")
	}
	return &Rekening{
		ID:                uuid.New(),
		BMTID:             bmtID,
		CabangID:          cabangID,
		NasabahID:         nasabahID,
		JenisRekeningID:   jenisID,
		NomorRekening:     nomor,
		Saldo:             money.Zero,
		Status:            StatusAktif,
		BiayaAdminBulanan: jenis.BiayaAdminBulanan,
		NisbahNasabah:     jenis.NisbahNasabah,
		TanggalBuka:       time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}

func (r *Rekening) ValidasiSetor(nominal int64, jenisMin int64) error {
	if r.Status == StatusBlokir {
		return ErrRekeningBeku
	}
	if r.Status == StatusTutup {
		return ErrRekeningTutup
	}
	if nominal < jenisMin {
		return fmt.Errorf("%w: minimum setoran %d", ErrSetoranDibawahMin, jenisMin)
	}
	return nil
}

func (r *Rekening) ValidasiTarik(nominal int64, bisaDitarik bool) error {
	if r.Status == StatusBlokir {
		return ErrRekeningBeku
	}
	if r.Status == StatusTutup {
		return ErrRekeningTutup
	}
	if !bisaDitarik {
		return ErrPenarikanTidakBisa
	}
	if money.New(nominal) > r.Saldo {
		return ErrSaldoTidakCukup
	}
	return nil
}

func (r *Rekening) NewTransaksiSetor(nominal int64, keterangan string, createdBy *uuid.UUID, idempotencyKey *uuid.UUID) *TransaksiRekening {
	saldoBaru := r.Saldo.Add(money.New(nominal))
	return &TransaksiRekening{
		ID:             uuid.New(),
		BMTID:          r.BMTID,
		CabangID:       r.CabangID,
		RekeningID:     r.ID,
		Jenis:          "SETOR",
		Posisi:         "KREDIT",
		Nominal:        nominal,
		SaldoSebelum:   r.Saldo.Int64(),
		SaldoSesudah:   saldoBaru.Int64(),
		Keterangan:     keterangan,
		IdempotencyKey: idempotencyKey,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}
}

func (r *Rekening) NewTransaksiTarik(nominal int64, keterangan string, createdBy *uuid.UUID, idempotencyKey *uuid.UUID) *TransaksiRekening {
	saldoBaru := r.Saldo.Sub(money.New(nominal))
	return &TransaksiRekening{
		ID:             uuid.New(),
		BMTID:          r.BMTID,
		CabangID:       r.CabangID,
		RekeningID:     r.ID,
		Jenis:          "TARIK",
		Posisi:         "DEBIT",
		Nominal:        nominal,
		SaldoSebelum:   r.Saldo.Int64(),
		SaldoSesudah:   saldoBaru.Int64(),
		Keterangan:     keterangan,
		IdempotencyKey: idempotencyKey,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}
}
