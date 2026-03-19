package pembiayaan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPembiayaanNotFound     = errors.New("pembiayaan tidak ditemukan")
	ErrPembiayaanTidakAktif   = errors.New("pembiayaan tidak dalam status aktif")
	ErrAngsuranMelebihiSaldo  = errors.New("nominal angsuran melebihi saldo kewajiban")
	ErrBeasiswaTidakValid     = errors.New("persen beasiswa harus antara 0 dan 100")
	ErrAkadTidakValid         = errors.New("akad pembiayaan tidak valid")
)

type StatusPembiayaan string

const (
	StatusPengajuan StatusPembiayaan = "PENGAJUAN"
	StatusAktif     StatusPembiayaan = "AKTIF"
	StatusLunas     StatusPembiayaan = "LUNAS"
	StatusMacet     StatusPembiayaan = "MACET"
	StatusWriteOff  StatusPembiayaan = "WRITEOFF"
)

type AkadPembiayaan string

const (
	AkadMurabahah AkadPembiayaan = "MURABAHAH"
	AkadMudharabah AkadPembiayaan = "MUDHARABAH"
	AkadMusyarakah AkadPembiayaan = "MUSYARAKAH"
	AkadIjarah    AkadPembiayaan = "IJARAH"
	AkadQardh     AkadPembiayaan = "QARDH"
)

type Pembiayaan struct {
	ID                    uuid.UUID        `json:"id"`
	BMTID                 uuid.UUID        `json:"bmt_id"`
	CabangID              uuid.UUID        `json:"cabang_id"`
	NasabahID             uuid.UUID        `json:"nasabah_id"`
	ProdukPembiayaanID    uuid.UUID        `json:"produk_pembiayaan_id"`
	NomorPembiayaan       string           `json:"nomor_pembiayaan"`
	Akad                  AkadPembiayaan   `json:"akad"`
	Pokok                 int64            `json:"pokok"`
	MarginPersen          *float64         `json:"margin_persen,omitempty"`
	NisbahNasabah         *int16           `json:"nisbah_nasabah,omitempty"`
	JangkaBulan           int16            `json:"jangka_bulan"`
	AngsuranPerBulan      int64            `json:"angsuran_per_bulan"`
	TotalKewajiban        int64            `json:"total_kewajiban"`
	// Beasiswa
	AdaBeasiswa           bool             `json:"ada_beasiswa"`
	BeasiswaPersen        *float64         `json:"beasiswa_persen,omitempty"`
	BeasiswaNominal       int64            `json:"beasiswa_nominal"`
	BeasiswaSumber        string           `json:"beasiswa_sumber,omitempty"`
	BeasiswaDitetapkanOleh *uuid.UUID      `json:"beasiswa_ditetapkan_oleh,omitempty"`
	BeasiswaDitetapkanAt  *time.Time       `json:"beasiswa_ditetapkan_at,omitempty"`
	// Status
	Status                StatusPembiayaan `json:"status"`
	Kolektibilitas        int16            `json:"kolektibilitas"`
	HariTunggak           int              `json:"hari_tunggak"`
	SaldoPokok            int64            `json:"saldo_pokok"`
	SaldoMargin           int64            `json:"saldo_margin"`
	// Audit
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`
	CreatedBy             uuid.UUID        `json:"created_by"`
	UpdatedBy             uuid.UUID        `json:"updated_by"`
	IsVoided              bool             `json:"is_voided"`
}

type AngsuranPembiayaan struct {
	ID              uuid.UUID  `json:"id"`
	BMTID           uuid.UUID  `json:"bmt_id"`
	PembiayaanID    uuid.UUID  `json:"pembiayaan_id"`
	PeriodeBulan    int16      `json:"periode_bulan"`
	NominalPokok    int64      `json:"nominal_pokok"`
	NominalMargin   int64      `json:"nominal_margin"`
	TotalAngsuran   int64      `json:"total_angsuran"`
	TanggalJatuhTempo time.Time `json:"tanggal_jatuh_tempo"`
	TanggalBayar    *time.Time `json:"tanggal_bayar,omitempty"`
	NominalTerbayar int64      `json:"nominal_terbayar"`
	Status          string     `json:"status"` // MENUNGGU | TERBAYAR | SEBAGIAN | LEWAT
	TransaksiID     *uuid.UUID `json:"transaksi_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type BeasiswaRiwayat struct {
	ID             uuid.UUID  `json:"id"`
	PembiayaanID   uuid.UUID  `json:"pembiayaan_id"`
	PersenSebelum  *float64   `json:"persen_sebelum,omitempty"`
	PersenSesudah  float64    `json:"persen_sesudah"`
	NominalSebelum *int64     `json:"nominal_sebelum,omitempty"`
	NominalSesudah int64      `json:"nominal_sesudah"`
	Alasan         string     `json:"alasan,omitempty"`
	DitetapkanOleh uuid.UUID  `json:"ditetapkan_oleh"`
	CreatedAt      time.Time  `json:"created_at"`
}

type CreatePembiayaanInput struct {
	BMTID              uuid.UUID
	CabangID           uuid.UUID
	NasabahID          uuid.UUID
	ProdukPembiayaanID uuid.UUID
	Akad               AkadPembiayaan
	Pokok              int64
	MarginPersen       *float64
	NisbahNasabah      *int16
	JangkaBulan        int16
	CreatedBy          uuid.UUID
}

type SetBeasiswaInput struct {
	PembiayaanID   uuid.UUID
	Persen         float64
	Sumber         string
	DitetapkanOleh uuid.UUID
}

type ListPembiayaanFilter struct {
	BMTID         *uuid.UUID
	CabangID      *uuid.UUID
	NasabahID     *uuid.UUID
	Status        *StatusPembiayaan
	Kolektibilitas *int16
	Page          int
	PerPage       int
}

type Repository interface {
	Create(ctx context.Context, p *Pembiayaan) error
	GetByID(ctx context.Context, id uuid.UUID) (*Pembiayaan, error)
	GetByNomor(ctx context.Context, nomor string) (*Pembiayaan, error)
	List(ctx context.Context, filter ListPembiayaanFilter) ([]*Pembiayaan, int64, error)
	Update(ctx context.Context, p *Pembiayaan) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusPembiayaan) error
	UpdateKolektibilitas(ctx context.Context, id uuid.UUID, kolektibilitas int16, hariTunggak int) error
	UpdateSaldo(ctx context.Context, id uuid.UUID, saldoPokok, saldoMargin int64) error
	SetBeasiswa(ctx context.Context, id uuid.UUID, persen float64, nominal int64, sumber string, oleh uuid.UUID) error
	LockForUpdate(ctx context.Context, id uuid.UUID) (*Pembiayaan, error)
	GenerateNomor(ctx context.Context, bmtID, cabangID uuid.UUID) (string, error)

	// Angsuran
	CreateAngsuran(ctx context.Context, a *AngsuranPembiayaan) error
	GetAngsuranByID(ctx context.Context, id uuid.UUID) (*AngsuranPembiayaan, error)
	ListAngsuran(ctx context.Context, pembiayaanID uuid.UUID) ([]*AngsuranPembiayaan, error)
	GetAngsuranJatuhTempo(ctx context.Context, bmtID uuid.UUID, tanggal time.Time) ([]*AngsuranPembiayaan, error)
	UpdateAngsuranTerbayar(ctx context.Context, id uuid.UUID, nominal int64, tanggalBayar time.Time, transaksiID uuid.UUID) error

	// Beasiswa riwayat
	CreateBeasiswaRiwayat(ctx context.Context, r *BeasiswaRiwayat) error
	ListBeasiswaRiwayat(ctx context.Context, pembiayaanID uuid.UUID) ([]*BeasiswaRiwayat, error)
}

func NewPembiayaan(input CreatePembiayaanInput, nomor string) (*Pembiayaan, error) {
	if input.Pokok <= 0 {
		return nil, errors.New("pokok pembiayaan harus lebih dari 0")
	}
	if input.JangkaBulan <= 0 {
		return nil, errors.New("jangka bulan harus lebih dari 0")
	}
	switch input.Akad {
	case AkadMurabahah, AkadMudharabah, AkadMusyarakah, AkadIjarah, AkadQardh:
	default:
		return nil, ErrAkadTidakValid
	}

	// Calculate angsuran per bulan based on akad
	totalKewajiban := hitungTotalKewajiban(input.Pokok, input.MarginPersen, input.JangkaBulan, input.Akad)
	angsuranPerBulan := totalKewajiban / int64(input.JangkaBulan)

	now := time.Now()
	return &Pembiayaan{
		ID:                 uuid.New(),
		BMTID:              input.BMTID,
		CabangID:           input.CabangID,
		NasabahID:          input.NasabahID,
		ProdukPembiayaanID: input.ProdukPembiayaanID,
		NomorPembiayaan:    nomor,
		Akad:               input.Akad,
		Pokok:              input.Pokok,
		MarginPersen:       input.MarginPersen,
		NisbahNasabah:      input.NisbahNasabah,
		JangkaBulan:        input.JangkaBulan,
		AngsuranPerBulan:   angsuranPerBulan,
		TotalKewajiban:     totalKewajiban,
		Status:             StatusPengajuan,
		Kolektibilitas:     1,
		SaldoPokok:         input.Pokok,
		SaldoMargin:        totalKewajiban - input.Pokok,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          input.CreatedBy,
		UpdatedBy:          input.CreatedBy,
	}, nil
}

func hitungTotalKewajiban(pokok int64, marginPersen *float64, jangka int16, akad AkadPembiayaan) int64 {
	if akad == AkadQardh || marginPersen == nil || *marginPersen == 0 {
		return pokok
	}
	// Flat margin for murabahah: total = pokok + (pokok * margin% * jangka / 12)
	if akad == AkadMurabahah {
		margin := int64(float64(pokok) * (*marginPersen) / 100 * float64(jangka) / 12)
		return pokok + margin
	}
	return pokok
}

func (p *Pembiayaan) ValidasiAktif() error {
	if p.Status != StatusAktif {
		return fmt.Errorf("%w: status saat ini %s", ErrPembiayaanTidakAktif, p.Status)
	}
	return nil
}

func (p *Pembiayaan) SetBeasiswa(persen float64, sumber string, oleh uuid.UUID) error {
	if persen < 0 || persen > 100 {
		return ErrBeasiswaTidakValid
	}
	nominal := int64(float64(p.TotalKewajiban) * persen / 100)
	now := time.Now()
	p.AdaBeasiswa = persen > 0
	p.BeasiswaPersen = &persen
	p.BeasiswaNominal = nominal
	p.BeasiswaSumber = sumber
	p.BeasiswaDitetapkanOleh = &oleh
	p.BeasiswaDitetapkanAt = &now
	p.UpdatedAt = now
	return nil
}
