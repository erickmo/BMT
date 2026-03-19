package penilaian

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrNilaiNotFound          = errors.New("nilai tidak ditemukan")
	ErrNilaiTahfidzNotFound   = errors.New("nilai tahfidz tidak ditemukan")
	ErrNilaiAkhlakNotFound    = errors.New("nilai akhlak tidak ditemukan")
	ErrRaportNotFound         = errors.New("raport tidak ditemukan")
	ErrRaportSudahFinal       = errors.New("raport sudah final, tidak bisa diubah")
	ErrRaportSudahDiterbitkan = errors.New("raport sudah diterbitkan")
	ErrNilaiDiluarRentang     = errors.New("nilai harus antara 0 dan 100")
	ErrKomponenSudahDinilai   = errors.New("komponen ini sudah dinilai untuk santri tersebut")
	ErrPoinHarusPositif       = errors.New("poin prestasi harus positif")
	ErrPoinHarusNegatif       = errors.New("poin pelanggaran harus negatif")
)

// ── Status Tahfidz ────────────────────────────────────────────────────────────

type StatusTahfidz string

const (
	StatusTahfidzLulus      StatusTahfidz = "LULUS"
	StatusTahfidzMengulang  StatusTahfidz = "MENGULANG"
	StatusTahfidzBelumDiuji StatusTahfidz = "BELUM_DIUJI"
)

// ── Jenis Akhlak ──────────────────────────────────────────────────────────────

type JenisAkhlak string

const (
	JenisPelanggaran JenisAkhlak = "PELANGGARAN"
	JenisPrestasi    JenisAkhlak = "PRESTASI"
)

// ── Status Raport ─────────────────────────────────────────────────────────────

type StatusRaport string

const (
	StatusRaportDraft       StatusRaport = "DRAFT"
	StatusRaportFinal       StatusRaport = "FINAL"
	StatusRaportDiterbitkan StatusRaport = "DITERBITKAN"
)

// ── Nilai ──────────────────────────────────────────────────────────────────────

// Nilai merepresentasikan nilai santri untuk satu komponen penilaian.
// Kombinasi (santri_id, komponen_id) bersifat unik.
type Nilai struct {
	ID           uuid.UUID `json:"id"`
	BMTID        uuid.UUID `json:"bmt_id"`
	SantriID     uuid.UUID `json:"santri_id"`
	KomponenID   uuid.UUID `json:"komponen_id"`
	Nilai        float64   `json:"nilai"`
	Catatan      string    `json:"catatan,omitempty"`
	DiinputOleh  uuid.UUID `json:"diinput_oleh"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewNilai membuat entitas Nilai baru dengan validasi rentang 0–100.
func NewNilai(bmtID, santriID, komponenID uuid.UUID, nilaiAngka float64, catatan string, diinputOleh uuid.UUID) (*Nilai, error) {
	if nilaiAngka < 0 || nilaiAngka > 100 {
		return nil, ErrNilaiDiluarRentang
	}
	now := time.Now()
	return &Nilai{
		ID:          uuid.New(),
		BMTID:       bmtID,
		SantriID:    santriID,
		KomponenID:  komponenID,
		Nilai:       nilaiAngka,
		Catatan:     catatan,
		DiinputOleh: diinputOleh,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ── NilaiTahfidz ──────────────────────────────────────────────────────────────

// NilaiTahfidz merepresentasikan hasil ujian hafalan Al-Quran santri.
type NilaiTahfidz struct {
	ID           uuid.UUID     `json:"id"`
	BMTID        uuid.UUID     `json:"bmt_id"`
	SantriID     uuid.UUID     `json:"santri_id"`
	Surah        string        `json:"surah"`
	AyatMulai    int16         `json:"ayat_mulai"`
	AyatSelesai  int16         `json:"ayat_selesai"`
	Nilai        *float64      `json:"nilai,omitempty"`
	Status       StatusTahfidz `json:"status"`
	TanggalUjian time.Time     `json:"tanggal_ujian"`
	PengujiID    *uuid.UUID    `json:"penguji_id,omitempty"`
	Catatan      string        `json:"catatan,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
}

// NewNilaiTahfidz membuat entitas NilaiTahfidz baru.
func NewNilaiTahfidz(bmtID, santriID uuid.UUID, surah string, ayatMulai, ayatSelesai int16, tanggalUjian time.Time, pengujiID *uuid.UUID) (*NilaiTahfidz, error) {
	if surah == "" {
		return nil, errors.New("nama surah wajib diisi")
	}
	if ayatMulai < 1 || ayatSelesai < ayatMulai {
		return nil, errors.New("rentang ayat tidak valid")
	}
	return &NilaiTahfidz{
		ID:           uuid.New(),
		BMTID:        bmtID,
		SantriID:     santriID,
		Surah:        surah,
		AyatMulai:    ayatMulai,
		AyatSelesai:  ayatSelesai,
		Status:       StatusTahfidzBelumDiuji,
		TanggalUjian: tanggalUjian,
		PengujiID:    pengujiID,
		CreatedAt:    time.Now(),
	}, nil
}

// Lulus mencatat hasil kelulusan tahfidz beserta nilai.
func (n *NilaiTahfidz) Lulus(nilaiAngka float64) error {
	if nilaiAngka < 0 || nilaiAngka > 100 {
		return ErrNilaiDiluarRentang
	}
	n.Status = StatusTahfidzLulus
	n.Nilai = &nilaiAngka
	return nil
}

// Mengulang menandai tahfidz harus diulang.
func (n *NilaiTahfidz) Mengulang() {
	n.Status = StatusTahfidzMengulang
}

// ── NilaiAkhlak ───────────────────────────────────────────────────────────────

// NilaiAkhlak merepresentasikan catatan poin akhlak/kedisiplinan santri.
// Poin positif = prestasi, poin negatif = pelanggaran.
type NilaiAkhlak struct {
	ID          uuid.UUID   `json:"id"`
	BMTID       uuid.UUID   `json:"bmt_id"`
	SantriID    uuid.UUID   `json:"santri_id"`
	Tanggal     time.Time   `json:"tanggal"`
	Jenis       JenisAkhlak `json:"jenis"`
	Kategori    string      `json:"kategori"`
	Deskripsi   string      `json:"deskripsi"`
	Poin        int16       `json:"poin"`
	DicatatOleh uuid.UUID   `json:"dicatat_oleh"`
	CreatedAt   time.Time   `json:"created_at"`
}

// NewNilaiAkhlak membuat catatan akhlak baru dengan validasi poin sesuai jenis.
func NewNilaiAkhlak(bmtID, santriID uuid.UUID, tanggal time.Time, jenis JenisAkhlak, kategori, deskripsi string, poin int16, dicatatOleh uuid.UUID) (*NilaiAkhlak, error) {
	if deskripsi == "" {
		return nil, errors.New("deskripsi catatan akhlak wajib diisi")
	}
	if jenis == JenisPrestasi && poin <= 0 {
		return nil, ErrPoinHarusPositif
	}
	if jenis == JenisPelanggaran && poin >= 0 {
		return nil, ErrPoinHarusNegatif
	}
	return &NilaiAkhlak{
		ID:          uuid.New(),
		BMTID:       bmtID,
		SantriID:    santriID,
		Tanggal:     tanggal,
		Jenis:       jenis,
		Kategori:    kategori,
		Deskripsi:   deskripsi,
		Poin:        poin,
		DicatatOleh: dicatatOleh,
		CreatedAt:   time.Now(),
	}, nil
}

// ── NilaiMapelSnapshot (embedded dalam Raport) ────────────────────────────────

// NilaiMapelSnapshot adalah snapshot nilai akhir per mapel yang disimpan dalam raport.
type NilaiMapelSnapshot struct {
	MapelID    uuid.UUID `json:"mapel_id"`
	NamaMapel  string    `json:"nama_mapel"`
	NilaiAkhir float64   `json:"nilai_akhir"`
	Predikat   string    `json:"predikat"`
}

// ── Raport ────────────────────────────────────────────────────────────────────

// Raport merepresentasikan laporan hasil belajar digital per santri per semester.
// NilaiMapel, NilaiTahfidz, dan NilaiAkhlak disimpan sebagai JSONB (snapshot saat generate).
// FileURL mengacu ke PDF raport di MinIO.
type Raport struct {
	ID              uuid.UUID    `json:"id"`
	BMTID           uuid.UUID    `json:"bmt_id"`
	SantriID        uuid.UUID    `json:"santri_id"`
	KelasID         uuid.UUID    `json:"kelas_id"`
	TahunAjaran     string       `json:"tahun_ajaran"`
	Semester        int16        `json:"semester"`
	NilaiMapel      json.RawMessage `json:"nilai_mapel"`
	NilaiTahfidz    json.RawMessage `json:"nilai_tahfidz,omitempty"`
	NilaiAkhlak     json.RawMessage `json:"nilai_akhlak,omitempty"`
	TotalHadir      *int16       `json:"total_hadir,omitempty"`
	TotalSakit      *int16       `json:"total_sakit,omitempty"`
	TotalIzin       *int16       `json:"total_izin,omitempty"`
	TotalAlfa       *int16       `json:"total_alfa,omitempty"`
	Peringkat       *int16       `json:"peringkat,omitempty"`
	CatatanWaliKelas string      `json:"catatan_wali_kelas,omitempty"`
	FileURL         string       `json:"file_url,omitempty"`
	Status          StatusRaport `json:"status"`
	DiterbitkanAt   *time.Time   `json:"diterbitkan_at,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
}

// NewRaport membuat entitas Raport baru dalam status DRAFT.
func NewRaport(bmtID, santriID, kelasID uuid.UUID, tahunAjaran string, semester int16) (*Raport, error) {
	if tahunAjaran == "" {
		return nil, errors.New("tahun ajaran wajib diisi")
	}
	if semester != 1 && semester != 2 {
		return nil, errors.New("semester harus 1 atau 2")
	}
	return &Raport{
		ID:          uuid.New(),
		BMTID:       bmtID,
		SantriID:    santriID,
		KelasID:     kelasID,
		TahunAjaran: tahunAjaran,
		Semester:    semester,
		Status:      StatusRaportDraft,
		NilaiMapel:  json.RawMessage("[]"),
		CreatedAt:   time.Now(),
	}, nil
}

// Finalisasi mengunci raport agar tidak bisa diubah lagi.
func (r *Raport) Finalisasi() error {
	if r.Status == StatusRaportDiterbitkan {
		return ErrRaportSudahDiterbitkan
	}
	r.Status = StatusRaportFinal
	return nil
}

// Terbitkan mempublikasikan raport ke wali santri.
func (r *Raport) Terbitkan() error {
	if r.Status != StatusRaportFinal {
		return errors.New("raport harus difinalisasi dahulu sebelum diterbitkan")
	}
	now := time.Now()
	r.Status = StatusRaportDiterbitkan
	r.DiterbitkanAt = &now
	return nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListNilaiFilter struct {
	BMTID       uuid.UUID
	SantriID    *uuid.UUID
	KomponenID  *uuid.UUID
	TahunAjaran string
	Semester    *int16
}

type ListNilaiTahfidzFilter struct {
	BMTID    uuid.UUID
	SantriID *uuid.UUID
	Status   StatusTahfidz
}

type ListNilaiAkhlakFilter struct {
	BMTID     uuid.UUID
	SantriID  *uuid.UUID
	Jenis     JenisAkhlak
	DariTgl   *time.Time
	SampaiTgl *time.Time
}

type ListRaportFilter struct {
	BMTID       uuid.UUID
	SantriID    *uuid.UUID
	KelasID     *uuid.UUID
	TahunAjaran string
	Semester    *int16
	Status      StatusRaport
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// NilaiRepository mendefinisikan kontrak akses data untuk entitas Nilai.
type NilaiRepository interface {
	Create(ctx context.Context, n *Nilai) error
	GetByID(ctx context.Context, id uuid.UUID) (*Nilai, error)
	GetBySantriKomponen(ctx context.Context, santriID, komponenID uuid.UUID) (*Nilai, error)
	List(ctx context.Context, filter ListNilaiFilter) ([]*Nilai, error)
	Update(ctx context.Context, n *Nilai) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NilaiTahfidzRepository mendefinisikan kontrak akses data untuk NilaiTahfidz.
type NilaiTahfidzRepository interface {
	Create(ctx context.Context, n *NilaiTahfidz) error
	GetByID(ctx context.Context, id uuid.UUID) (*NilaiTahfidz, error)
	List(ctx context.Context, filter ListNilaiTahfidzFilter) ([]*NilaiTahfidz, error)
	Update(ctx context.Context, n *NilaiTahfidz) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NilaiAkhlakRepository mendefinisikan kontrak akses data untuk NilaiAkhlak.
type NilaiAkhlakRepository interface {
	Create(ctx context.Context, n *NilaiAkhlak) error
	GetByID(ctx context.Context, id uuid.UUID) (*NilaiAkhlak, error)
	List(ctx context.Context, filter ListNilaiAkhlakFilter) ([]*NilaiAkhlak, error)
	SumPoin(ctx context.Context, bmtID, santriID uuid.UUID, dariTgl, sampaiTgl time.Time) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// RaportRepository mendefinisikan kontrak akses data untuk entitas Raport.
type RaportRepository interface {
	Create(ctx context.Context, r *Raport) error
	GetByID(ctx context.Context, id uuid.UUID) (*Raport, error)
	GetBySantriSemester(ctx context.Context, santriID uuid.UUID, tahunAjaran string, semester int16) (*Raport, error)
	List(ctx context.Context, filter ListRaportFilter) ([]*Raport, error)
	Update(ctx context.Context, r *Raport) error
	Delete(ctx context.Context, id uuid.UUID) error
}
