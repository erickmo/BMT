package perpustakaan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrBukuNotFound          = errors.New("buku tidak ditemukan")
	ErrPeminjamanNotFound     = errors.New("data peminjaman tidak ditemukan")
	ErrBukuTidakTersedia     = errors.New("stok buku tidak tersedia untuk dipinjam")
	ErrBukuSudahDikembalikan = errors.New("buku sudah dikembalikan")
	ErrBukuTidakAktif        = errors.New("buku tidak aktif dalam katalog")
	ErrJudulWajibDiisi       = errors.New("judul buku wajib diisi")
	ErrStokTidakValid        = errors.New("stok total harus minimal 1")
)

// ── Status Peminjaman ─────────────────────────────────────────────────────────

type StatusPeminjaman string

const (
	StatusDipinjam      StatusPeminjaman = "DIPINJAM"
	StatusDikembalikan  StatusPeminjaman = "DIKEMBALIKAN"
	StatusTerlambat     StatusPeminjaman = "TERLAMBAT"
	StatusHilang        StatusPeminjaman = "HILANG"
)

// ── Tipe Peminjam ─────────────────────────────────────────────────────────────

type TipePeminjam string

const (
	TipeSantri   TipePeminjam = "SANTRI"
	TipePengajar TipePeminjam = "PENGAJAR"
)

// ── Buku ──────────────────────────────────────────────────────────────────────

// Buku merepresentasikan koleksi perpustakaan pondok.
type Buku struct {
	ID               uuid.UUID `json:"id"`
	BMTID            uuid.UUID `json:"bmt_id"`
	CabangID         uuid.UUID `json:"cabang_id"`
	Judul            string    `json:"judul"`
	Pengarang        string    `json:"pengarang,omitempty"`
	Penerbit         string    `json:"penerbit,omitempty"`
	TahunTerbit      *int16    `json:"tahun_terbit,omitempty"`
	ISBN             string    `json:"isbn,omitempty"`
	Kategori         string    `json:"kategori,omitempty"`
	JumlahTotal      int16     `json:"jumlah_total"`
	JumlahTersedia   int16     `json:"jumlah_tersedia"`
	FotoURL          string    `json:"foto_url,omitempty"`
	IsAktif          bool      `json:"is_aktif"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewBuku membuat entitas Buku baru.
func NewBuku(bmtID, cabangID uuid.UUID, judul string, stokTotal int16) (*Buku, error) {
	if judul == "" {
		return nil, ErrJudulWajibDiisi
	}
	if stokTotal < 1 {
		return nil, ErrStokTidakValid
	}
	now := time.Now()
	return &Buku{
		ID:             uuid.New(),
		BMTID:          bmtID,
		CabangID:       cabangID,
		Judul:          judul,
		JumlahTotal:    stokTotal,
		JumlahTersedia: stokTotal,
		IsAktif:        true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// ValidasiPinjam memastikan buku dapat dipinjam.
func (b *Buku) ValidasiPinjam() error {
	if !b.IsAktif {
		return ErrBukuTidakAktif
	}
	if b.JumlahTersedia <= 0 {
		return ErrBukuTidakTersedia
	}
	return nil
}

// KurangiStok mengurangi stok tersedia saat buku dipinjam.
func (b *Buku) KurangiStok() error {
	if err := b.ValidasiPinjam(); err != nil {
		return err
	}
	b.JumlahTersedia--
	b.UpdatedAt = time.Now()
	return nil
}

// TambahStok menambah stok tersedia saat buku dikembalikan.
func (b *Buku) TambahStok() {
	if b.JumlahTersedia < b.JumlahTotal {
		b.JumlahTersedia++
	}
	b.UpdatedAt = time.Now()
}

// ── Peminjaman ────────────────────────────────────────────────────────────────

// Peminjaman merepresentasikan transaksi peminjaman buku perpustakaan.
// Denda dihitung di layer service berdasarkan hari keterlambatan.
// Denda masuk ke dana sosial (akun 611), bukan pendapatan BMT (prinsip syariah).
type Peminjaman struct {
	ID                    uuid.UUID        `json:"id"`
	BMTID                 uuid.UUID        `json:"bmt_id"`
	BukuID                uuid.UUID        `json:"buku_id"`
	PeminjamID            uuid.UUID        `json:"peminjam_id"`
	PeminjamTipe          TipePeminjam     `json:"peminjam_tipe"`
	TanggalPinjam         time.Time        `json:"tanggal_pinjam"`
	TanggalKembaliRencana time.Time        `json:"tanggal_kembali_rencana"`
	TanggalKembaliAktual  *time.Time       `json:"tanggal_kembali_aktual,omitempty"`
	Status                StatusPeminjaman `json:"status"`
	// Denda dalam satuan int64 (Rupiah) — masuk dana sosial, bukan pendapatan
	Denda     int64     `json:"denda"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPeminjaman membuat transaksi peminjaman baru.
func NewPeminjaman(bmtID, bukuID, peminjamID uuid.UUID, peminjamTipe TipePeminjam, tanggalPinjam, tanggalKembaliRencana time.Time) (*Peminjaman, error) {
	if !tanggalPinjam.Before(tanggalKembaliRencana) {
		return nil, errors.New("tanggal kembali harus setelah tanggal pinjam")
	}
	now := time.Now()
	return &Peminjaman{
		ID:                    uuid.New(),
		BMTID:                 bmtID,
		BukuID:                bukuID,
		PeminjamID:            peminjamID,
		PeminjamTipe:          peminjamTipe,
		TanggalPinjam:         tanggalPinjam,
		TanggalKembaliRencana: tanggalKembaliRencana,
		Status:                StatusDipinjam,
		Denda:                 0,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

// Kembalikan mencatat pengembalian buku dan menghitung denda jika terlambat.
// dendaPerHari diambil dari settings BMT — tidak hardcode.
func (p *Peminjaman) Kembalikan(tanggalKembali time.Time, dendaPerHari int64) error {
	if p.Status == StatusDikembalikan {
		return ErrBukuSudahDikembalikan
	}
	p.TanggalKembaliAktual = &tanggalKembali
	p.Status = StatusDikembalikan
	// Hitung denda jika terlambat
	if tanggalKembali.After(p.TanggalKembaliRencana) {
		hariTerlambat := int64(tanggalKembali.Sub(p.TanggalKembaliRencana).Hours() / 24)
		if hariTerlambat < 1 {
			hariTerlambat = 1
		}
		p.Denda = hariTerlambat * dendaPerHari
	}
	p.UpdatedAt = time.Now()
	return nil
}

// IsKadaluarsa memeriksa apakah peminjaman sudah melewati tanggal kembali rencana.
func (p *Peminjaman) IsKadaluarsa() bool {
	return p.Status == StatusDipinjam && time.Now().After(p.TanggalKembaliRencana)
}

// TandaiTerlambat mengubah status menjadi TERLAMBAT.
func (p *Peminjaman) TandaiTerlambat() error {
	if p.Status != StatusDipinjam {
		return fmt.Errorf("status peminjaman %s tidak bisa ditandai terlambat", p.Status)
	}
	p.Status = StatusTerlambat
	p.UpdatedAt = time.Now()
	return nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListBukuFilter struct {
	BMTID    uuid.UUID
	CabangID uuid.UUID
	Kategori string
	Keyword  string
	IsAktif  *bool
	Page     int
	PerPage  int
}

type ListPeminjamanFilter struct {
	BMTID       uuid.UUID
	BukuID      *uuid.UUID
	PeminjamID  *uuid.UUID
	PeminjamTipe TipePeminjam
	Status      StatusPeminjaman
	DariTgl     *time.Time
	SampaiTgl   *time.Time
	Page        int
	PerPage     int
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// BukuRepository mendefinisikan kontrak akses data untuk entitas Buku.
type BukuRepository interface {
	Create(ctx context.Context, b *Buku) error
	GetByID(ctx context.Context, id uuid.UUID) (*Buku, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Buku, error)
	List(ctx context.Context, filter ListBukuFilter) ([]*Buku, int64, error)
	Update(ctx context.Context, b *Buku) error
	UpdateStok(ctx context.Context, id uuid.UUID, jumlahTersedia int16) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PeminjamanRepository mendefinisikan kontrak akses data untuk entitas Peminjaman.
type PeminjamanRepository interface {
	Create(ctx context.Context, p *Peminjaman) error
	GetByID(ctx context.Context, id uuid.UUID) (*Peminjaman, error)
	List(ctx context.Context, filter ListPeminjamanFilter) ([]*Peminjaman, int64, error)
	// ListHampirJatuhTempo digunakan worker reminder perpustakaan.
	ListHampirJatuhTempo(ctx context.Context, bmtID uuid.UUID, hariSebelum int) ([]*Peminjaman, error)
	Update(ctx context.Context, p *Peminjaman) error
	Delete(ctx context.Context, id uuid.UUID) error
}
