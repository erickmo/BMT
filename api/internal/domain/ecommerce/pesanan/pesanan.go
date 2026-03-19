package pesanan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPesananNotFound          = errors.New("pesanan tidak ditemukan")
	ErrPesananSudahDibayar      = errors.New("pesanan sudah dibayar")
	ErrPesananTidakBisaDibatal  = errors.New("pesanan tidak dapat dibatalkan pada status ini")
	ErrPesananKosong            = errors.New("pesanan tidak boleh kosong")
	ErrStatusTransisiTidakValid = errors.New("transisi status pesanan tidak valid")
)

type StatusPesanan string

const (
	StatusMenungguPembayaran StatusPesanan = "MENUNGGU_PEMBAYARAN"
	StatusDibayar            StatusPesanan = "DIBAYAR"
	StatusDiproses           StatusPesanan = "DIPROSES"
	StatusDikirim            StatusPesanan = "DIKIRIM"
	StatusSelesai            StatusPesanan = "SELESAI"
	StatusDibatalkan         StatusPesanan = "DIBATALKAN"
)

type BuyerTipe string

const (
	BuyerWaliSantri BuyerTipe = "WALI_SANTRI"
	BuyerPondok     BuyerTipe = "PONDOK"
)

type MetodeBayar string

const (
	MetodeMidtrans    MetodeBayar = "MIDTRANS"
	MetodeRekeningBMT MetodeBayar = "REKENING_BMT"
	MetodeNFC         MetodeBayar = "NFC"
)

type Pesanan struct {
	ID            uuid.UUID       `json:"id"`
	BuyerTipe     BuyerTipe       `json:"buyer_tipe"`
	NasabahID     *uuid.UUID      `json:"nasabah_id,omitempty"`
	BMTBuyerID    *uuid.UUID      `json:"bmt_buyer_id,omitempty"`
	TokoID        uuid.UUID       `json:"toko_id"`
	BMTSellerID   uuid.UUID       `json:"bmt_seller_id"`
	NomorPesanan  string          `json:"nomor_pesanan"`
	Status        StatusPesanan   `json:"status"`
	Subtotal      int64           `json:"subtotal"`
	Ongkir        int64           `json:"ongkir"`
	Total         int64           `json:"total"`
	AlamatKirim   json.RawMessage `json:"alamat_kirim"`
	Kurir         string          `json:"kurir,omitempty"`
	NomorResi     string          `json:"nomor_resi,omitempty"`
	MetodeBayar   MetodeBayar     `json:"metode_bayar,omitempty"`
	Catatan       string          `json:"catatan,omitempty"`
	Items         []*PesananItem  `json:"items,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type PesananItem struct {
	ID          uuid.UUID `json:"id"`
	PesananID   uuid.UUID `json:"pesanan_id"`
	ProdukID    uuid.UUID `json:"produk_id"`
	NamaProduk  string    `json:"nama_produk"`
	Harga       int64     `json:"harga"`
	Jumlah      int       `json:"jumlah"`
	Subtotal    int64     `json:"subtotal"`
}

type CreatePesananInput struct {
	BuyerTipe   BuyerTipe
	NasabahID   *uuid.UUID
	BMTBuyerID  *uuid.UUID
	TokoID      uuid.UUID
	BMTSellerID uuid.UUID
	Items       []ItemInput
	AlamatKirim json.RawMessage
	Ongkir      int64
	Catatan     string
}

type ItemInput struct {
	ProdukID   uuid.UUID
	NamaProduk string
	Harga      int64
	Jumlah     int
}

type ListPesananFilter struct {
	NasabahID   *uuid.UUID
	TokoID      *uuid.UUID
	BMTSellerID *uuid.UUID
	BMTBuyerID  *uuid.UUID
	Status      *StatusPesanan
	Periode     string
	Page        int
	PerPage     int
}

type Repository interface {
	Create(ctx context.Context, p *Pesanan) error
	CreateItems(ctx context.Context, items []*PesananItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*Pesanan, error)
	GetByNomor(ctx context.Context, nomor string) (*Pesanan, error)
	GetWithItems(ctx context.Context, id uuid.UUID) (*Pesanan, error)
	List(ctx context.Context, filter ListPesananFilter) ([]*Pesanan, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status StatusPesanan) error
	UpdatePengiriman(ctx context.Context, id uuid.UUID, kurir, nomorResi string) error
	UpdateMetodeBayar(ctx context.Context, id uuid.UUID, metode MetodeBayar) error
	GenerateNomor(ctx context.Context, bmtID uuid.UUID) (string, error)
}

func NewPesanan(input CreatePesananInput) (*Pesanan, error) {
	if len(input.Items) == 0 {
		return nil, ErrPesananKosong
	}
	if input.BuyerTipe == BuyerWaliSantri && input.NasabahID == nil {
		return nil, errors.New("nasabah_id wajib untuk buyer wali santri")
	}
	if input.BuyerTipe == BuyerPondok && input.BMTBuyerID == nil {
		return nil, errors.New("bmt_buyer_id wajib untuk buyer pondok")
	}

	var subtotal int64
	items := make([]*PesananItem, 0, len(input.Items))
	pesananID := uuid.New()

	for _, it := range input.Items {
		if it.Jumlah <= 0 {
			return nil, fmt.Errorf("jumlah item harus lebih dari 0")
		}
		itemSubtotal := it.Harga * int64(it.Jumlah)
		subtotal += itemSubtotal
		items = append(items, &PesananItem{
			ID:         uuid.New(),
			PesananID:  pesananID,
			ProdukID:   it.ProdukID,
			NamaProduk: it.NamaProduk,
			Harga:      it.Harga,
			Jumlah:     it.Jumlah,
			Subtotal:   itemSubtotal,
		})
	}

	now := time.Now()
	return &Pesanan{
		ID:          pesananID,
		BuyerTipe:   input.BuyerTipe,
		NasabahID:   input.NasabahID,
		BMTBuyerID:  input.BMTBuyerID,
		TokoID:      input.TokoID,
		BMTSellerID: input.BMTSellerID,
		Status:      StatusMenungguPembayaran,
		Subtotal:    subtotal,
		Ongkir:      input.Ongkir,
		Total:       subtotal + input.Ongkir,
		AlamatKirim: input.AlamatKirim,
		Catatan:     input.Catatan,
		Items:       items,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// validTransisi defines allowed status transitions
var validTransisi = map[StatusPesanan][]StatusPesanan{
	StatusMenungguPembayaran: {StatusDibayar, StatusDibatalkan},
	StatusDibayar:            {StatusDiproses, StatusDibatalkan},
	StatusDiproses:           {StatusDikirim},
	StatusDikirim:            {StatusSelesai},
	StatusSelesai:            {},
	StatusDibatalkan:         {},
}

func (p *Pesanan) TransisiStatus(statusBaru StatusPesanan) error {
	allowed, ok := validTransisi[p.Status]
	if !ok {
		return ErrStatusTransisiTidakValid
	}
	for _, s := range allowed {
		if s == statusBaru {
			return nil
		}
	}
	return fmt.Errorf("%w: dari %s ke %s", ErrStatusTransisiTidakValid, p.Status, statusBaru)
}

func (p *Pesanan) BisaDibatalkan() bool {
	return p.Status == StatusMenungguPembayaran || p.Status == StatusDibayar
}
