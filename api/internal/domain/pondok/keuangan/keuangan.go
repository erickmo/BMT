package keuangan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrJenisTagihanNotFound  = errors.New("jenis tagihan tidak ditemukan")
	ErrTagihanSPPNotFound    = errors.New("tagihan SPP tidak ditemukan")
	ErrTagihanSudahLunas     = errors.New("tagihan sudah lunas")
	ErrTagihanSudahAda       = errors.New("tagihan untuk periode ini sudah dibuat")
	ErrNominalBeasiswaMelebihi = errors.New("nominal beasiswa melebihi nominal tagihan")
	ErrBeasiswaPersenTidakValid = errors.New("persentase beasiswa harus antara 0 dan 100")
	ErrPembayaranMelebihi    = errors.New("nominal pembayaran melebihi sisa tagihan")
	ErrKodeTagihanSudahAda   = errors.New("kode jenis tagihan sudah terdaftar")
	ErrNominalHarusPositif   = errors.New("nominal harus lebih dari 0")
)

// ── Frekuensi Tagihan ─────────────────────────────────────────────────────────

type FrekuensiTagihan string

const (
	FrekuensiBulanan FrekuensiTagihan = "BULANAN"
	FrekuensiTahunan FrekuensiTagihan = "TAHUNAN"
	FrekuensiSekali  FrekuensiTagihan = "SEKALI"
	FrekuensiCustom  FrekuensiTagihan = "CUSTOM"
)

// ── Status Tagihan ────────────────────────────────────────────────────────────

type StatusTagihan string

const (
	StatusBelumBayar StatusTagihan = "BELUM_BAYAR"
	StatusSebagian   StatusTagihan = "SEBAGIAN"
	StatusLunas      StatusTagihan = "LUNAS"
)

// ── JenisTagihan ──────────────────────────────────────────────────────────────

// JenisTagihan mendefinisikan tipe tagihan yang berlaku di pondok.
// Dibuat oleh management BMT via CRUD — tidak ada hardcode jenis tagihan.
type JenisTagihan struct {
	ID         uuid.UUID        `json:"id"`
	BMTID      uuid.UUID        `json:"bmt_id"`
	Kode       string           `json:"kode"`
	Nama       string           `json:"nama"`
	Nominal    int64            `json:"nominal"`
	Frekuensi  FrekuensiTagihan `json:"frekuensi"`
	IsAktif    bool             `json:"is_aktif"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// NewJenisTagihan membuat entitas JenisTagihan baru.
func NewJenisTagihan(bmtID uuid.UUID, kode, nama string, nominal int64, frekuensi FrekuensiTagihan) (*JenisTagihan, error) {
	if kode == "" {
		return nil, errors.New("kode jenis tagihan wajib diisi")
	}
	if nama == "" {
		return nil, errors.New("nama jenis tagihan wajib diisi")
	}
	if nominal <= 0 {
		return nil, ErrNominalHarusPositif
	}
	now := time.Now()
	return &JenisTagihan{
		ID:        uuid.New(),
		BMTID:     bmtID,
		Kode:      kode,
		Nama:      nama,
		Nominal:   nominal,
		Frekuensi: frekuensi,
		IsAktif:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ── TagihanSPP ────────────────────────────────────────────────────────────────

// TagihanSPP merepresentasikan tagihan SPP santri untuk satu periode.
//
// Beasiswa ditetapkan oleh admin pondok dan mengurangi nominal efektif yang harus dibayar.
// NominalEfektif = Nominal - BeasiswaNominal
// Autodebet hanya men-debit NominalEfektif, bukan Nominal penuh.
type TagihanSPP struct {
	ID               uuid.UUID     `json:"id"`
	BMTID            uuid.UUID     `json:"bmt_id"`
	CabangID         uuid.UUID     `json:"cabang_id"`
	SantriID         uuid.UUID     `json:"santri_id"`
	JenisTagihanID   uuid.UUID     `json:"jenis_tagihan_id"`
	Periode          string        `json:"periode"`
	Nominal          int64         `json:"nominal"`
	NominalTerbayar  int64         `json:"nominal_terbayar"`
	NominalSisa      int64         `json:"nominal_sisa"`
	// Beasiswa — ditetapkan admin pondok, bisa 0–100%
	BeasiswaPersen   float64       `json:"beasiswa_persen"`
	BeasiswaNominal  int64         `json:"beasiswa_nominal"`
	// NominalEfektif = Nominal - BeasiswaNominal
	NominalEfektif   int64         `json:"nominal_efektif"`
	Status           StatusTagihan `json:"status"`
	TanggalJatuhTempo time.Time    `json:"tanggal_jatuh_tempo"`
	TanggalLunas     *time.Time    `json:"tanggal_lunas,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// NewTagihanSPP membuat tagihan SPP baru untuk santri.
func NewTagihanSPP(bmtID, cabangID, santriID, jenisTagihanID uuid.UUID, periode string, nominal int64, tanggalJatuhTempo time.Time) (*TagihanSPP, error) {
	if periode == "" {
		return nil, errors.New("periode tagihan wajib diisi")
	}
	if nominal <= 0 {
		return nil, ErrNominalHarusPositif
	}
	now := time.Now()
	return &TagihanSPP{
		ID:                uuid.New(),
		BMTID:             bmtID,
		CabangID:          cabangID,
		SantriID:          santriID,
		JenisTagihanID:    jenisTagihanID,
		Periode:           periode,
		Nominal:           nominal,
		NominalTerbayar:   0,
		NominalSisa:       nominal,
		BeasiswaPersen:    0,
		BeasiswaNominal:   0,
		NominalEfektif:    nominal,
		Status:            StatusBelumBayar,
		TanggalJatuhTempo: tanggalJatuhTempo,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// TerapkanBeasiswa menetapkan beasiswa pada tagihan.
// NominalEfektif dihitung ulang setelah beasiswa diterapkan.
func (t *TagihanSPP) TerapkanBeasiswa(persenBeasiswa float64) error {
	if persenBeasiswa < 0 || persenBeasiswa > 100 {
		return ErrBeasiswaPersenTidakValid
	}
	nominalBeasiswa := int64(float64(t.Nominal) * persenBeasiswa / 100)
	nominalEfektif := t.Nominal - nominalBeasiswa
	if nominalEfektif < 0 {
		return ErrNominalBeasiswaMelebihi
	}
	t.BeasiswaPersen = persenBeasiswa
	t.BeasiswaNominal = nominalBeasiswa
	t.NominalEfektif = nominalEfektif
	// Recalculate sisa berdasarkan nominal efektif
	t.NominalSisa = t.NominalEfektif - t.NominalTerbayar
	if t.NominalSisa < 0 {
		t.NominalSisa = 0
	}
	t.UpdatedAt = time.Now()
	return nil
}

// Bayar mencatat pembayaran (bisa parsial).
// Mengembalikan jumlah yang berhasil dibayar (partial jika sisa kurang dari nominal).
func (t *TagihanSPP) Bayar(nominal int64) (int64, error) {
	if t.Status == StatusLunas {
		return 0, ErrTagihanSudahLunas
	}
	if nominal <= 0 {
		return 0, ErrNominalHarusPositif
	}
	// Partial payment: hanya bayar sebesar sisa
	bayar := nominal
	if bayar > t.NominalSisa {
		return 0, fmt.Errorf("%w: nominal pembayaran %d, sisa tagihan %d", ErrPembayaranMelebihi, nominal, t.NominalSisa)
	}
	t.NominalTerbayar += bayar
	t.NominalSisa -= bayar
	now := time.Now()
	if t.NominalSisa == 0 {
		t.Status = StatusLunas
		t.TanggalLunas = &now
	} else {
		t.Status = StatusSebagian
	}
	t.UpdatedAt = now
	return bayar, nil
}

// SisaSetelahBeasiswa menghitung sisa tagihan yang masih harus dibayar.
func (t *TagihanSPP) SisaSetelahBeasiswa() int64 {
	return t.NominalSisa
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListTagihanFilter struct {
	BMTID     uuid.UUID
	CabangID  uuid.UUID
	SantriID  *uuid.UUID
	Periode   string
	Status    StatusTagihan
	DariTgl   *time.Time
	SampaiTgl *time.Time
	Page      int
	PerPage   int
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// JenisTagihanRepository mendefinisikan kontrak akses data untuk JenisTagihan.
type JenisTagihanRepository interface {
	Create(ctx context.Context, j *JenisTagihan) error
	GetByID(ctx context.Context, id uuid.UUID) (*JenisTagihan, error)
	GetByKode(ctx context.Context, bmtID uuid.UUID, kode string) (*JenisTagihan, error)
	List(ctx context.Context, bmtID uuid.UUID, aktifSaja bool) ([]*JenisTagihan, error)
	Update(ctx context.Context, j *JenisTagihan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// TagihanSPPRepository mendefinisikan kontrak akses data untuk TagihanSPP.
type TagihanSPPRepository interface {
	Create(ctx context.Context, t *TagihanSPP) error
	GetByID(ctx context.Context, id uuid.UUID) (*TagihanSPP, error)
	GetBySantriPeriode(ctx context.Context, santriID uuid.UUID, periode string) (*TagihanSPP, error)
	List(ctx context.Context, filter ListTagihanFilter) ([]*TagihanSPP, int64, error)
	// ListBelumLunas digunakan worker autodebet untuk mencari tagihan yang perlu di-debit.
	ListBelumLunas(ctx context.Context, bmtID uuid.UUID, tanggalJatuhTempo time.Time) ([]*TagihanSPP, error)
	Update(ctx context.Context, t *TagihanSPP) error
	UpdateBeasiswa(ctx context.Context, id uuid.UUID, persen float64, nominal, efektif, sisa int64) error
	Delete(ctx context.Context, id uuid.UUID) error
}
