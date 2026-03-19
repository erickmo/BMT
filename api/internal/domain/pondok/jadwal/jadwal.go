package jadwal

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Error sentinels ──────────────────────────────────────────────────────────

var (
	ErrJadwalPelajaranNotFound = errors.New("jadwal pelajaran tidak ditemukan")
	ErrJadwalKegiatanNotFound  = errors.New("jadwal kegiatan tidak ditemukan")
	ErrJadwalPiketNotFound     = errors.New("jadwal piket tidak ditemukan")
	ErrJadwalShiftNotFound     = errors.New("jadwal shift tidak ditemukan")
	ErrKalenderNotFound        = errors.New("entri kalender tidak ditemukan")
	ErrJadwalBentrok           = errors.New("jadwal bentrok dengan slot yang sudah ada")
	ErrTanggalTidakValid       = errors.New("tanggal tidak valid")
	ErrHariTidakValid          = errors.New("hari harus antara 1 (Senin) dan 7 (Minggu)")
	ErrJamTidakValid           = errors.New("jam mulai harus sebelum jam selesai")
)

// ── Jenis Kalender ────────────────────────────────────────────────────────────

type JenisKalender string

const (
	JenisLibur         JenisKalender = "LIBUR"
	JenisUjian         JenisKalender = "UJIAN"
	JenisAcara         JenisKalender = "ACARA"
	JenisHariEfektif   JenisKalender = "HARI_EFEKTIF"
	JenisLiburNasional JenisKalender = "LIBUR_NASIONAL"
)

// ── Kategori Kegiatan ─────────────────────────────────────────────────────────

type KategoriKegiatan string

const (
	KategoriPengajian KategoriKegiatan = "PENGAJIAN"
	KategoriOlahraga  KategoriKegiatan = "OLAHRAGA"
	KategoriEkstra    KategoriKegiatan = "EKSTRA"
	KategoriRapat     KategoriKegiatan = "RAPAT"
	KategoriAcara     KategoriKegiatan = "ACARA"
	KategoriLainnya   KategoriKegiatan = "LAINNYA"
)

// ── Target Peserta Kegiatan ───────────────────────────────────────────────────

type TargetPeserta string

const (
	TargetSemua     TargetPeserta = "SEMUA"
	TargetSantri    TargetPeserta = "SANTRI"
	TargetPengajar  TargetPeserta = "PENGAJAR"
	TargetKaryawan  TargetPeserta = "KARYAWAN"
	TargetTertentu  TargetPeserta = "TERTENTU"
)

// ── Jenis Pengguna Shift ──────────────────────────────────────────────────────

type JenisPenggunaShift string

const (
	JenisShiftPengajar  JenisPenggunaShift = "PENGAJAR"
	JenisShiftKaryawan  JenisPenggunaShift = "KARYAWAN"
)

// ── Kalender Akademik ─────────────────────────────────────────────────────────

// KalenderAkademik merepresentasikan entri kalender pondok per hari.
type KalenderAkademik struct {
	ID          uuid.UUID     `json:"id"`
	BMTID       uuid.UUID     `json:"bmt_id"`
	CabangID    uuid.UUID     `json:"cabang_id"`
	TahunAjaran string        `json:"tahun_ajaran"`
	Tanggal     time.Time     `json:"tanggal"`
	Jenis       JenisKalender `json:"jenis"`
	Keterangan  string        `json:"keterangan,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

// NewKalenderAkademik membuat entri kalender akademik baru.
func NewKalenderAkademik(bmtID, cabangID uuid.UUID, tahunAjaran string, tanggal time.Time, jenis JenisKalender) (*KalenderAkademik, error) {
	if tahunAjaran == "" {
		return nil, errors.New("tahun ajaran wajib diisi")
	}
	return &KalenderAkademik{
		ID:          uuid.New(),
		BMTID:       bmtID,
		CabangID:    cabangID,
		TahunAjaran: tahunAjaran,
		Tanggal:     tanggal,
		Jenis:       jenis,
		CreatedAt:   time.Now(),
	}, nil
}

// ── Jadwal Pelajaran ──────────────────────────────────────────────────────────

// JadwalPelajaran merepresentasikan slot pelajaran mingguan untuk satu kelas.
// Hari menggunakan konvensi: 1=Senin ... 7=Minggu.
type JadwalPelajaran struct {
	ID          uuid.UUID `json:"id"`
	BMTID       uuid.UUID `json:"bmt_id"`
	KelasID     uuid.UUID `json:"kelas_id"`
	MapelID     uuid.UUID `json:"mapel_id"`
	PengajarID  uuid.UUID `json:"pengajar_id"`
	Hari        int16     `json:"hari"`
	JamMulai    string    `json:"jam_mulai"`
	JamSelesai  string    `json:"jam_selesai"`
	Ruangan     string    `json:"ruangan,omitempty"`
	TahunAjaran string    `json:"tahun_ajaran"`
	Semester    int16     `json:"semester"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewJadwalPelajaran membuat jadwal pelajaran baru dengan validasi.
func NewJadwalPelajaran(bmtID, kelasID, mapelID, pengajarID uuid.UUID, hari int16, jamMulai, jamSelesai, tahunAjaran string, semester int16) (*JadwalPelajaran, error) {
	if hari < 1 || hari > 7 {
		return nil, ErrHariTidakValid
	}
	if jamMulai == "" || jamSelesai == "" {
		return nil, ErrJamTidakValid
	}
	if jamMulai >= jamSelesai {
		return nil, ErrJamTidakValid
	}
	now := time.Now()
	return &JadwalPelajaran{
		ID:          uuid.New(),
		BMTID:       bmtID,
		KelasID:     kelasID,
		MapelID:     mapelID,
		PengajarID:  pengajarID,
		Hari:        hari,
		JamMulai:    jamMulai,
		JamSelesai:  jamSelesai,
		TahunAjaran: tahunAjaran,
		Semester:    semester,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ── Jadwal Kegiatan ───────────────────────────────────────────────────────────

// JadwalKegiatan merepresentasikan kegiatan pondok (pengajian, olahraga, rapat, dll.).
type JadwalKegiatan struct {
	ID             uuid.UUID        `json:"id"`
	BMTID          uuid.UUID        `json:"bmt_id"`
	CabangID       uuid.UUID        `json:"cabang_id"`
	Nama           string           `json:"nama"`
	Kategori       KategoriKegiatan `json:"kategori"`
	TanggalMulai   time.Time        `json:"tanggal_mulai"`
	TanggalSelesai *time.Time       `json:"tanggal_selesai,omitempty"`
	Lokasi         string           `json:"lokasi,omitempty"`
	Peserta        TargetPeserta    `json:"peserta"`
	Deskripsi      string           `json:"deskripsi,omitempty"`
	CreatedBy      uuid.UUID        `json:"created_by"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// NewJadwalKegiatan membuat jadwal kegiatan baru.
func NewJadwalKegiatan(bmtID, cabangID uuid.UUID, nama string, kategori KategoriKegiatan, tanggalMulai time.Time, peserta TargetPeserta, createdBy uuid.UUID) (*JadwalKegiatan, error) {
	if nama == "" {
		return nil, errors.New("nama kegiatan wajib diisi")
	}
	now := time.Now()
	return &JadwalKegiatan{
		ID:           uuid.New(),
		BMTID:        bmtID,
		CabangID:     cabangID,
		Nama:         nama,
		Kategori:     kategori,
		TanggalMulai: tanggalMulai,
		Peserta:      peserta,
		CreatedBy:    createdBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ── Jadwal Piket ──────────────────────────────────────────────────────────────

// JadwalPiket merepresentasikan jadwal piket santri (kebersihan, keamanan, dll.).
type JadwalPiket struct {
	ID             uuid.UUID `json:"id"`
	BMTID          uuid.UUID `json:"bmt_id"`
	SantriID       uuid.UUID `json:"santri_id"`
	JenisPiket     string    `json:"jenis_piket"`
	Hari           int16     `json:"hari"`
	Lokasi         string    `json:"lokasi,omitempty"`
	PeriodeMulai   time.Time `json:"periode_mulai"`
	PeriodeSelesai time.Time `json:"periode_selesai"`
	CreatedAt      time.Time `json:"created_at"`
}

// NewJadwalPiket membuat jadwal piket baru.
func NewJadwalPiket(bmtID, santriID uuid.UUID, jenisPiket string, hari int16, periodeMulai, periodeSelesai time.Time) (*JadwalPiket, error) {
	if jenisPiket == "" {
		return nil, errors.New("jenis piket wajib diisi")
	}
	if hari < 1 || hari > 7 {
		return nil, ErrHariTidakValid
	}
	if !periodeMulai.Before(periodeSelesai) {
		return nil, ErrTanggalTidakValid
	}
	return &JadwalPiket{
		ID:             uuid.New(),
		BMTID:          bmtID,
		SantriID:       santriID,
		JenisPiket:     jenisPiket,
		Hari:           hari,
		PeriodeMulai:   periodeMulai,
		PeriodeSelesai: periodeSelesai,
		CreatedAt:      time.Now(),
	}, nil
}

// ── Jadwal Shift ──────────────────────────────────────────────────────────────

// JadwalShift merepresentasikan jadwal kerja harian pengajar atau karyawan.
type JadwalShift struct {
	ID             uuid.UUID          `json:"id"`
	BMTID          uuid.UUID          `json:"bmt_id"`
	PenggunaID     *uuid.UUID         `json:"pengguna_id,omitempty"`
	JenisPengguna  JenisPenggunaShift `json:"jenis_pengguna"`
	Tanggal        time.Time          `json:"tanggal"`
	JamMasuk       string             `json:"jam_masuk"`
	JamKeluar      string             `json:"jam_keluar"`
	Keterangan     string             `json:"keterangan,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// NewJadwalShift membuat jadwal shift baru.
func NewJadwalShift(bmtID uuid.UUID, penggunaID *uuid.UUID, jenisPengguna JenisPenggunaShift, tanggal time.Time, jamMasuk, jamKeluar string) (*JadwalShift, error) {
	if jamMasuk == "" || jamKeluar == "" {
		return nil, ErrJamTidakValid
	}
	if jamMasuk >= jamKeluar {
		return nil, ErrJamTidakValid
	}
	now := time.Now()
	return &JadwalShift{
		ID:            uuid.New(),
		BMTID:         bmtID,
		PenggunaID:    penggunaID,
		JenisPengguna: jenisPengguna,
		Tanggal:       tanggal,
		JamMasuk:      jamMasuk,
		JamKeluar:     jamKeluar,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// ── Filter types ──────────────────────────────────────────────────────────────

type ListJadwalPelajaranFilter struct {
	BMTID       uuid.UUID
	KelasID     *uuid.UUID
	PengajarID  *uuid.UUID
	TahunAjaran string
	Semester    *int16
	Hari        *int16
}

type ListJadwalKegiatanFilter struct {
	BMTID     uuid.UUID
	CabangID  uuid.UUID
	Kategori  KategoriKegiatan
	Peserta   TargetPeserta
	DariTgl   *time.Time
	SampaiTgl *time.Time
}

type ListJadwalShiftFilter struct {
	BMTID         uuid.UUID
	PenggunaID    *uuid.UUID
	JenisPengguna JenisPenggunaShift
	DariTgl       *time.Time
	SampaiTgl     *time.Time
}

// ── Repository interfaces ─────────────────────────────────────────────────────

// KalenderRepository mendefinisikan kontrak akses data untuk KalenderAkademik.
type KalenderRepository interface {
	Create(ctx context.Context, k *KalenderAkademik) error
	GetByID(ctx context.Context, id uuid.UUID) (*KalenderAkademik, error)
	GetByTanggal(ctx context.Context, bmtID, cabangID uuid.UUID, tanggal time.Time) (*KalenderAkademik, error)
	List(ctx context.Context, bmtID, cabangID uuid.UUID, tahunAjaran string, dariTgl, sampaiTgl *time.Time) ([]*KalenderAkademik, error)
	Update(ctx context.Context, k *KalenderAkademik) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// JadwalPelajaranRepository mendefinisikan kontrak akses data untuk JadwalPelajaran.
type JadwalPelajaranRepository interface {
	Create(ctx context.Context, j *JadwalPelajaran) error
	GetByID(ctx context.Context, id uuid.UUID) (*JadwalPelajaran, error)
	List(ctx context.Context, filter ListJadwalPelajaranFilter) ([]*JadwalPelajaran, error)
	Update(ctx context.Context, j *JadwalPelajaran) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// JadwalKegiatanRepository mendefinisikan kontrak akses data untuk JadwalKegiatan.
type JadwalKegiatanRepository interface {
	Create(ctx context.Context, j *JadwalKegiatan) error
	GetByID(ctx context.Context, id uuid.UUID) (*JadwalKegiatan, error)
	List(ctx context.Context, filter ListJadwalKegiatanFilter) ([]*JadwalKegiatan, int64, error)
	Update(ctx context.Context, j *JadwalKegiatan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// JadwalPiketRepository mendefinisikan kontrak akses data untuk JadwalPiket.
type JadwalPiketRepository interface {
	Create(ctx context.Context, j *JadwalPiket) error
	GetByID(ctx context.Context, id uuid.UUID) (*JadwalPiket, error)
	ListBySantri(ctx context.Context, bmtID, santriID uuid.UUID) ([]*JadwalPiket, error)
	ListByHari(ctx context.Context, bmtID uuid.UUID, hari int16) ([]*JadwalPiket, error)
	Update(ctx context.Context, j *JadwalPiket) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// JadwalShiftRepository mendefinisikan kontrak akses data untuk JadwalShift.
type JadwalShiftRepository interface {
	Create(ctx context.Context, j *JadwalShift) error
	GetByID(ctx context.Context, id uuid.UUID) (*JadwalShift, error)
	List(ctx context.Context, filter ListJadwalShiftFilter) ([]*JadwalShift, error)
	Update(ctx context.Context, j *JadwalShift) error
	Delete(ctx context.Context, id uuid.UUID) error
}
