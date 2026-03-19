package finance

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJurnalNotFound          = errors.New("jurnal tidak ditemukan")
	ErrJurnalTidakBalance      = errors.New("jurnal tidak balance: debit ≠ kredit")
	ErrVendorNotFound          = errors.New("vendor tidak ditemukan")
	ErrTransaksiOperasionalNotFound = errors.New("transaksi operasional tidak ditemukan")
)

type PosisiJurnal string

const (
	PosisiDebit  PosisiJurnal = "DEBIT"
	PosisiKredit PosisiJurnal = "KREDIT"
)

type StatusJurnal string

const (
	StatusJurnalDraft  StatusJurnal = "DRAFT"
	StatusJurnalPosted StatusJurnal = "POSTED"
	StatusJurnalVoid   StatusJurnal = "VOID"
)

// EntriJurnal is a single debit or credit line in a journal entry
type EntriJurnal struct {
	ID        uuid.UUID    `json:"id"`
	JurnalID  uuid.UUID    `json:"jurnal_id"`
	KodeAkun  string       `json:"kode_akun"`
	NamaAkun  string       `json:"nama_akun"`
	Posisi    PosisiJurnal `json:"posisi"`
	Nominal   int64        `json:"nominal"`
}

// JurnalManual is a manually created journal entry by finance staff
type JurnalManual struct {
	ID          uuid.UUID       `json:"id"`
	BMTID       uuid.UUID       `json:"bmt_id"`
	CabangID    uuid.UUID       `json:"cabang_id"`
	Tanggal     time.Time       `json:"tanggal"`
	Keterangan  string          `json:"keterangan"`
	Referensi   string          `json:"referensi,omitempty"`
	Entries     []*EntriJurnal  `json:"entries"`
	Status      StatusJurnal    `json:"status"`
	DibuatOleh  uuid.UUID       `json:"dibuat_oleh"`
	DisetujuiOleh *uuid.UUID    `json:"disetujui_oleh,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// Vendor represents a supplier or payee for operational expenses
type Vendor struct {
	ID         uuid.UUID `json:"id"`
	BMTID      uuid.UUID `json:"bmt_id"`
	Nama       string    `json:"nama"`
	NPWP       string    `json:"npwp,omitempty"`
	Alamat     string    `json:"alamat,omitempty"`
	Telepon    string    `json:"telepon,omitempty"`
	Email      string    `json:"email,omitempty"`
	RekeningBank string  `json:"rekening_bank,omitempty"`
	NamaBank   string    `json:"nama_bank,omitempty"`
	IsAktif    bool      `json:"is_aktif"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TransaksiOperasional records an operational expense or income
type TransaksiOperasional struct {
	ID           uuid.UUID       `json:"id"`
	BMTID        uuid.UUID       `json:"bmt_id"`
	CabangID     uuid.UUID       `json:"cabang_id"`
	VendorID     *uuid.UUID      `json:"vendor_id,omitempty"`
	Tanggal      time.Time       `json:"tanggal"`
	Jenis        string          `json:"jenis"`  // PENGELUARAN | PENERIMAAN
	Kategori     string          `json:"kategori"`
	Keterangan   string          `json:"keterangan"`
	Nominal      int64           `json:"nominal"`
	KodeAkunDebit  string        `json:"kode_akun_debit"`
	KodeAkunKredit string        `json:"kode_akun_kredit"`
	Lampiran     json.RawMessage `json:"lampiran,omitempty"` // ["url1", "url2"]
	JurnalID     *uuid.UUID      `json:"jurnal_id,omitempty"`
	DibuatOleh   uuid.UUID       `json:"dibuat_oleh"`
	CreatedAt    time.Time       `json:"created_at"`
}

type CreateJurnalInput struct {
	BMTID      uuid.UUID
	CabangID   uuid.UUID
	Tanggal    time.Time
	Keterangan string
	Referensi  string
	Entries    []EntriInput
	DibuatOleh uuid.UUID
}

type EntriInput struct {
	KodeAkun string
	NamaAkun string
	Posisi   PosisiJurnal
	Nominal  int64
}

type ListJurnalFilter struct {
	BMTID      *uuid.UUID
	CabangID   *uuid.UUID
	TanggalDari *time.Time
	TanggalSampai *time.Time
	Status     *StatusJurnal
	Page       int
	PerPage    int
}

type Repository interface {
	// Jurnal Manual
	CreateJurnal(ctx context.Context, j *JurnalManual) error
	GetJurnalByID(ctx context.Context, id uuid.UUID) (*JurnalManual, error)
	ListJurnal(ctx context.Context, filter ListJurnalFilter) ([]*JurnalManual, int64, error)
	PostJurnal(ctx context.Context, id uuid.UUID, disetujuiOleh uuid.UUID) error
	VoidJurnal(ctx context.Context, id uuid.UUID) error

	// Vendor
	CreateVendor(ctx context.Context, v *Vendor) error
	GetVendorByID(ctx context.Context, id uuid.UUID) (*Vendor, error)
	ListVendor(ctx context.Context, bmtID uuid.UUID, page, perPage int) ([]*Vendor, int64, error)
	UpdateVendor(ctx context.Context, v *Vendor) error

	// Transaksi Operasional
	CreateTransaksiOperasional(ctx context.Context, t *TransaksiOperasional) error
	GetTransaksiOperasionalByID(ctx context.Context, id uuid.UUID) (*TransaksiOperasional, error)
	ListTransaksiOperasional(ctx context.Context, bmtID, cabangID uuid.UUID, dari, sampai time.Time, page, perPage int) ([]*TransaksiOperasional, int64, error)
}

func NewJurnalManual(input CreateJurnalInput) (*JurnalManual, error) {
	if input.Keterangan == "" {
		return nil, errors.New("keterangan jurnal wajib diisi")
	}
	if len(input.Entries) < 2 {
		return nil, errors.New("jurnal minimal memiliki 2 entri")
	}

	// Validate balance: total debit must equal total kredit
	var totalDebit, totalKredit int64
	entries := make([]*EntriJurnal, 0, len(input.Entries))
	jurnalID := uuid.New()

	for _, e := range input.Entries {
		if e.Nominal <= 0 {
			return nil, errors.New("nominal entri jurnal harus lebih dari 0")
		}
		if e.KodeAkun == "" {
			return nil, errors.New("kode akun wajib diisi")
		}
		switch e.Posisi {
		case PosisiDebit:
			totalDebit += e.Nominal
		case PosisiKredit:
			totalKredit += e.Nominal
		default:
			return nil, errors.New("posisi jurnal harus DEBIT atau KREDIT")
		}
		entries = append(entries, &EntriJurnal{
			ID:       uuid.New(),
			JurnalID: jurnalID,
			KodeAkun: e.KodeAkun,
			NamaAkun: e.NamaAkun,
			Posisi:   e.Posisi,
			Nominal:  e.Nominal,
		})
	}

	if totalDebit != totalKredit {
		return nil, ErrJurnalTidakBalance
	}

	now := time.Now()
	return &JurnalManual{
		ID:         jurnalID,
		BMTID:      input.BMTID,
		CabangID:   input.CabangID,
		Tanggal:    input.Tanggal,
		Keterangan: input.Keterangan,
		Referensi:  input.Referensi,
		Entries:    entries,
		Status:     StatusJurnalDraft,
		DibuatOleh: input.DibuatOleh,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func NewVendor(bmtID uuid.UUID, nama string) (*Vendor, error) {
	if nama == "" {
		return nil, errors.New("nama vendor wajib diisi")
	}
	now := time.Now()
	return &Vendor{
		ID:        uuid.New(),
		BMTID:     bmtID,
		Nama:      nama,
		IsAktif:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
