package nfc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrKartuNFCNotFound      = errors.New("kartu NFC tidak ditemukan")
	ErrKartuNFCTidakAktif    = errors.New("kartu NFC tidak aktif atau expired")
	ErrKartuNFCPINSalah      = errors.New("PIN kartu tidak sesuai")
	ErrLimitHarianTerlampaui = errors.New("limit transaksi harian kartu NFC terlampaui")
	ErrLimitPerTransaksi     = errors.New("nominal melebihi limit per transaksi NFC")
	ErrIPKioskTidakDiizinkan = errors.New("IP tidak terdaftar sebagai terminal kiosk")
	ErrTransaksiNFCNotFound  = errors.New("transaksi NFC tidak ditemukan")
	ErrIdempotencyDuplikat   = errors.New("transaksi sudah pernah diproses (idempotency)")
)

type StatusKartuNFC string

const (
	StatusKartuAktif   StatusKartuNFC = "AKTIF"
	StatusKartuBlokir  StatusKartuNFC = "BLOKIR"
	StatusKartuExpired StatusKartuNFC = "EXPIRED"
)

type StatusTransaksiNFC string

const (
	StatusTxPending  StatusTransaksiNFC = "PENDING"
	StatusTxBerhasil StatusTransaksiNFC = "BERHASIL"
	StatusTxGagal    StatusTransaksiNFC = "GAGAL"
)

// KartuNFC represents an NFC card issued to a santri
type KartuNFC struct {
	ID              uuid.UUID      `json:"id"`
	BMTID           uuid.UUID      `json:"bmt_id"`
	NasabahID       uuid.UUID      `json:"nasabah_id"`
	RekeningID      uuid.UUID      `json:"rekening_id"`
	UID             string         `json:"uid"` // physical NFC card UID
	PINHash         string         `json:"-"`   // bcrypt hash — never expose
	Status          StatusKartuNFC `json:"status"`
	LimitPerTransaksi int64        `json:"limit_per_transaksi"`
	LimitHarian     int64          `json:"limit_harian"`
	TotalHarian     int64          `json:"total_harian"`
	TanggalReset    time.Time      `json:"tanggal_reset"`
	TanggalExpired  *time.Time     `json:"tanggal_expired,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// TransaksiNFC records a single NFC card transaction
type TransaksiNFC struct {
	ID             uuid.UUID          `json:"id"`
	BMTID          uuid.UUID          `json:"bmt_id"`
	KartuNFCID     uuid.UUID          `json:"kartu_nfc_id"`
	MerchantID     uuid.UUID          `json:"merchant_id"`
	RekeningID     uuid.UUID          `json:"rekening_id"`
	Nominal        int64              `json:"nominal"`
	SaldoSebelum   int64              `json:"saldo_sebelum"`
	SaldoSesudah   int64              `json:"saldo_sesudah"`
	Status         StatusTransaksiNFC `json:"status"`
	KeteranganGagal string            `json:"keterangan_gagal,omitempty"`
	IdempotencyKey uuid.UUID          `json:"idempotency_key"`
	CreatedAt      time.Time          `json:"created_at"`
}

// TerminalKiosk represents an NFC kiosk terminal with IP whitelist
type TerminalKiosk struct {
	ID        uuid.UUID `json:"id"`
	BMTID     uuid.UUID `json:"bmt_id"`
	CabangID  uuid.UUID `json:"cabang_id"`
	Nama      string    `json:"nama"`
	IPAddress string    `json:"ip_address"`
	Lokasi    string    `json:"lokasi,omitempty"`
	IsAktif   bool      `json:"is_aktif"`
	CreatedAt time.Time `json:"created_at"`
}

type TransaksiNFCInput struct {
	KartuNFCID     uuid.UUID
	MerchantID     uuid.UUID
	Nominal        int64
	IdempotencyKey uuid.UUID
}

type Repository interface {
	// KartuNFC
	CreateKartu(ctx context.Context, k *KartuNFC) error
	GetKartuByID(ctx context.Context, id uuid.UUID) (*KartuNFC, error)
	GetKartuByUID(ctx context.Context, uid string) (*KartuNFC, error)
	GetKartuByNasabah(ctx context.Context, nasabahID uuid.UUID) ([]*KartuNFC, error)
	UpdateStatusKartu(ctx context.Context, id uuid.UUID, status StatusKartuNFC) error
	UpdateLimitKartu(ctx context.Context, id uuid.UUID, limitPerTx, limitHarian int64) error
	UpdateTotalHarian(ctx context.Context, id uuid.UUID, tambah int64) error
	ResetTotalHarian(ctx context.Context, bmtID uuid.UUID) error
	LockKartuForUpdate(ctx context.Context, uid string) (*KartuNFC, error)

	// TransaksiNFC
	CreateTransaksi(ctx context.Context, t *TransaksiNFC) error
	GetTransaksiByID(ctx context.Context, id uuid.UUID) (*TransaksiNFC, error)
	GetTransaksiByIdempotency(ctx context.Context, key uuid.UUID) (*TransaksiNFC, error)
	ListTransaksiByKartu(ctx context.Context, kartuID uuid.UUID, limit, offset int) ([]*TransaksiNFC, int64, error)

	// Kiosk
	CreateTerminalKiosk(ctx context.Context, t *TerminalKiosk) error
	GetTerminalByIP(ctx context.Context, ip string) (*TerminalKiosk, error)
	ListTerminalByBMT(ctx context.Context, bmtID uuid.UUID) ([]*TerminalKiosk, error)
	UpdateTerminalStatus(ctx context.Context, id uuid.UUID, isAktif bool) error
}

func NewKartuNFC(bmtID, nasabahID, rekeningID uuid.UUID, uid, pinHash string, limitPerTx, limitHarian int64) (*KartuNFC, error) {
	if uid == "" {
		return nil, errors.New("UID kartu NFC wajib diisi")
	}
	if pinHash == "" {
		return nil, errors.New("PIN harus di-hash sebelum disimpan")
	}
	now := time.Now()
	return &KartuNFC{
		ID:                uuid.New(),
		BMTID:             bmtID,
		NasabahID:         nasabahID,
		RekeningID:        rekeningID,
		UID:               uid,
		PINHash:           pinHash,
		Status:            StatusKartuAktif,
		LimitPerTransaksi: limitPerTx,
		LimitHarian:       limitHarian,
		TotalHarian:       0,
		TanggalReset:      now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (k *KartuNFC) ValidasiTransaksi(nominal int64) error {
	if k.Status != StatusKartuAktif {
		return ErrKartuNFCTidakAktif
	}
	if k.TanggalExpired != nil && time.Now().After(*k.TanggalExpired) {
		return ErrKartuNFCTidakAktif
	}
	if nominal > k.LimitPerTransaksi {
		return ErrLimitPerTransaksi
	}
	if k.TotalHarian+nominal > k.LimitHarian {
		return ErrLimitHarianTerlampaui
	}
	return nil
}

func NewTransaksiNFC(kartu *KartuNFC, merchantID uuid.UUID, nominal int64, saldoRekening int64, idempotencyKey uuid.UUID) (*TransaksiNFC, error) {
	if err := kartu.ValidasiTransaksi(nominal); err != nil {
		return nil, err
	}
	return &TransaksiNFC{
		ID:             uuid.New(),
		BMTID:          kartu.BMTID,
		KartuNFCID:     kartu.ID,
		MerchantID:     merchantID,
		RekeningID:     kartu.RekeningID,
		Nominal:        nominal,
		SaldoSebelum:   saldoRekening,
		SaldoSesudah:   saldoRekening - nominal,
		Status:         StatusTxPending,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now(),
	}, nil
}
