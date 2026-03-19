package pembayaran

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPembayaranNotFound      = errors.New("pembayaran tidak ditemukan")
	ErrPembayaranSudahSettled  = errors.New("pembayaran sudah selesai")
	ErrPembayaranExpired       = errors.New("pembayaran sudah expired")
	ErrIdempotencyKeyDuplikat  = errors.New("transaksi sudah pernah diproses (idempotency)")
)

type MetodePembayaran string

const (
	MetodeMidtrans    MetodePembayaran = "MIDTRANS"
	MetodeRekeningBMT MetodePembayaran = "REKENING_BMT"
	MetodeNFC         MetodePembayaran = "NFC"
)

type StatusPembayaran string

const (
	StatusPending    StatusPembayaran = "PENDING"
	StatusSettlement StatusPembayaran = "SETTLEMENT"
	StatusExpire     StatusPembayaran = "EXPIRE"
	StatusCancel     StatusPembayaran = "CANCEL"
)

type PembayaranPesanan struct {
	ID               uuid.UUID        `json:"id"`
	PesananID        uuid.UUID        `json:"pesanan_id"`
	Metode           MetodePembayaran `json:"metode"`
	Nominal          int64            `json:"nominal"`
	Status           StatusPembayaran `json:"status"`
	MidtransOrderID  *string          `json:"midtrans_order_id,omitempty"`
	RekeningID       *uuid.UUID       `json:"rekening_id,omitempty"`
	KartuNFCID       *uuid.UUID       `json:"kartu_nfc_id,omitempty"`
	IdempotencyKey   *uuid.UUID       `json:"idempotency_key,omitempty"`
	SettledAt        *time.Time       `json:"settled_at,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
}

type CreatePembayaranInput struct {
	PesananID       uuid.UUID
	Metode          MetodePembayaran
	Nominal         int64
	MidtransOrderID *string
	RekeningID      *uuid.UUID
	KartuNFCID      *uuid.UUID
	IdempotencyKey  *uuid.UUID
}

type Repository interface {
	Create(ctx context.Context, p *PembayaranPesanan) error
	GetByID(ctx context.Context, id uuid.UUID) (*PembayaranPesanan, error)
	GetByPesananID(ctx context.Context, pesananID uuid.UUID) (*PembayaranPesanan, error)
	GetByMidtransOrderID(ctx context.Context, orderID string) (*PembayaranPesanan, error)
	GetByIdempotencyKey(ctx context.Context, key uuid.UUID) (*PembayaranPesanan, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusPembayaran, settledAt *time.Time) error
}

func NewPembayaran(input CreatePembayaranInput) (*PembayaranPesanan, error) {
	if input.Nominal <= 0 {
		return nil, errors.New("nominal pembayaran harus lebih dari 0")
	}
	if input.Metode == MetodeRekeningBMT && input.RekeningID == nil {
		return nil, errors.New("rekening_id wajib untuk metode rekening BMT")
	}
	if input.Metode == MetodeNFC && input.KartuNFCID == nil {
		return nil, errors.New("kartu_nfc_id wajib untuk metode NFC")
	}
	return &PembayaranPesanan{
		ID:              uuid.New(),
		PesananID:       input.PesananID,
		Metode:          input.Metode,
		Nominal:         input.Nominal,
		Status:          StatusPending,
		MidtransOrderID: input.MidtransOrderID,
		RekeningID:      input.RekeningID,
		KartuNFCID:      input.KartuNFCID,
		IdempotencyKey:  input.IdempotencyKey,
		CreatedAt:       time.Now(),
	}, nil
}

func (p *PembayaranPesanan) Settle() error {
	if p.Status == StatusSettlement {
		return ErrPembayaranSudahSettled
	}
	if p.Status == StatusExpire || p.Status == StatusCancel {
		return errors.New("pembayaran sudah tidak aktif")
	}
	now := time.Now()
	p.Status = StatusSettlement
	p.SettledAt = &now
	return nil
}
