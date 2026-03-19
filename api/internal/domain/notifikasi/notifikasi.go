package notifikasi

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTemplateNotFound    = errors.New("template notifikasi tidak ditemukan")
	ErrPengumumanNotFound  = errors.New("pengumuman tidak ditemukan")
	ErrChannelTidakValid   = errors.New("channel notifikasi tidak valid")
)

type Channel string

const (
	ChannelFCM       Channel = "FCM"
	ChannelWhatsApp  Channel = "WHATSAPP"
	ChannelSMS       Channel = "SMS"
	ChannelEmail     Channel = "EMAIL"
)

type StatusAntrian string

const (
	StatusMenunggu  StatusAntrian = "MENUNGGU"
	StatusTerkirim  StatusAntrian = "TERKIRIM"
	StatusGagal     StatusAntrian = "GAGAL"
)

type TargetPengumuman string

const (
	TargetSemua    TargetPengumuman = "SEMUA"
	TargetSantri   TargetPengumuman = "SANTRI"
	TargetWali     TargetPengumuman = "WALI"
	TargetPengajar TargetPengumuman = "PENGAJAR"
	TargetKaryawan TargetPengumuman = "KARYAWAN"
	TargetKelas    TargetPengumuman = "KELAS"
	TargetAsrama   TargetPengumuman = "ASRAMA"
)

// NotifikasiTemplate holds reusable notification templates per BMT
type NotifikasiTemplate struct {
	ID        uuid.UUID  `json:"id"`
	BMTID     *uuid.UUID `json:"bmt_id,omitempty"` // NULL = global platform template
	Kode      string     `json:"kode"`
	Channel   Channel    `json:"channel"`
	Judul     string     `json:"judul,omitempty"`
	Isi       string     `json:"isi"`
	IsAktif   bool       `json:"is_aktif"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NotifikasiAntrian is a queued notification pending delivery
type NotifikasiAntrian struct {
	ID            uuid.UUID       `json:"id"`
	BMTID         uuid.UUID       `json:"bmt_id"`
	Channel       Channel         `json:"channel"`
	Tujuan        string          `json:"tujuan"`   // FCM token / phone / email
	Subjek        string          `json:"subjek,omitempty"`
	Pesan         string          `json:"pesan"`
	DataEkstra    json.RawMessage `json:"data_ekstra,omitempty"`
	Status        StatusAntrian   `json:"status"`
	Percobaan     int16           `json:"percobaan"`
	ErrorTerakhir string          `json:"error_terakhir,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	DikirimAt     *time.Time      `json:"dikirim_at,omitempty"`
}

// NotifikasiLog records the outcome of each notification delivery attempt
type NotifikasiLog struct {
	ID            uuid.UUID  `json:"id"`
	BMTID         uuid.UUID  `json:"bmt_id"`
	TemplateKode  string     `json:"template_kode"`
	Channel       Channel    `json:"channel"`
	Tujuan        string     `json:"tujuan"`
	IsiTerkirim   string     `json:"isi_terkirim"`
	Status        string     `json:"status"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	ReferensiID   *uuid.UUID `json:"referensi_id,omitempty"`
	ReferensiTipe string     `json:"referensi_tipe,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// Pengumuman is an in-app bulletin board post for a pondok
type Pengumuman struct {
	ID             uuid.UUID        `json:"id"`
	BMTID          uuid.UUID        `json:"bmt_id"`
	CabangID       uuid.UUID        `json:"cabang_id"`
	Judul          string           `json:"judul"`
	Isi            string           `json:"isi"`
	Tipe           TargetPengumuman `json:"tipe"`
	TargetID       *uuid.UUID       `json:"target_id,omitempty"`
	TargetAsrama   string           `json:"target_asrama,omitempty"`
	FileURL        string           `json:"file_url,omitempty"`
	IsPinned       bool             `json:"is_pinned"`
	TanggalMulai   time.Time        `json:"tanggal_mulai"`
	TanggalSelesai *time.Time       `json:"tanggal_selesai,omitempty"`
	DibuatOleh     uuid.UUID        `json:"dibuat_oleh"`
	IsAktif        bool             `json:"is_aktif"`
	CreatedAt      time.Time        `json:"created_at"`
}

// PengumumanBaca tracks which users have read a pengumuman
type PengumumanBaca struct {
	PengumumanID uuid.UUID  `json:"pengumuman_id"`
	NasabahID    *uuid.UUID `json:"nasabah_id,omitempty"`
	PenggunaID   *uuid.UUID `json:"pengguna_id,omitempty"`
	DibacaAt     time.Time  `json:"dibaca_at"`
}

type CreateAntrianInput struct {
	BMTID      uuid.UUID
	Channel    Channel
	Tujuan     string
	Subjek     string
	Pesan      string
	DataEkstra json.RawMessage
}

type ListAntrianFilter struct {
	BMTID   *uuid.UUID
	Channel *Channel
	Status  *StatusAntrian
	Limit   int
}

type Repository interface {
	// Template
	GetTemplate(ctx context.Context, bmtID *uuid.UUID, kode string, channel Channel) (*NotifikasiTemplate, error)
	ListTemplate(ctx context.Context, bmtID uuid.UUID) ([]*NotifikasiTemplate, error)
	UpsertTemplate(ctx context.Context, t *NotifikasiTemplate) error

	// Antrian
	CreateAntrian(ctx context.Context, a *NotifikasiAntrian) error
	GetAntrianMenunggu(ctx context.Context, limit int) ([]*NotifikasiAntrian, error)
	UpdateStatusAntrian(ctx context.Context, id uuid.UUID, status StatusAntrian, errorMsg string) error
	IncrementPercobaan(ctx context.Context, id uuid.UUID) error

	// Log
	CreateLog(ctx context.Context, l *NotifikasiLog) error

	// Pengumuman
	CreatePengumuman(ctx context.Context, p *Pengumuman) error
	GetPengumumanByID(ctx context.Context, id uuid.UUID) (*Pengumuman, error)
	ListPengumuman(ctx context.Context, bmtID, cabangID uuid.UUID, tipe TargetPengumuman, page, perPage int) ([]*Pengumuman, int64, error)
	UpdatePengumuman(ctx context.Context, p *Pengumuman) error
	MarkBaca(ctx context.Context, b *PengumumanBaca) error
}

func NewAntrian(input CreateAntrianInput) (*NotifikasiAntrian, error) {
	if input.Tujuan == "" {
		return nil, errors.New("tujuan notifikasi wajib diisi")
	}
	if input.Pesan == "" {
		return nil, errors.New("pesan notifikasi wajib diisi")
	}
	switch input.Channel {
	case ChannelFCM, ChannelWhatsApp, ChannelSMS, ChannelEmail:
	default:
		return nil, ErrChannelTidakValid
	}
	return &NotifikasiAntrian{
		ID:         uuid.New(),
		BMTID:      input.BMTID,
		Channel:    input.Channel,
		Tujuan:     input.Tujuan,
		Subjek:     input.Subjek,
		Pesan:      input.Pesan,
		DataEkstra: input.DataEkstra,
		Status:     StatusMenunggu,
		Percobaan:  0,
		CreatedAt:  time.Now(),
	}, nil
}

func NewPengumuman(
	bmtID, cabangID uuid.UUID,
	judul, isi string,
	tipe TargetPengumuman,
	dibuatOleh uuid.UUID,
) (*Pengumuman, error) {
	if judul == "" {
		return nil, errors.New("judul pengumuman wajib diisi")
	}
	if isi == "" {
		return nil, errors.New("isi pengumuman wajib diisi")
	}
	now := time.Now()
	return &Pengumuman{
		ID:           uuid.New(),
		BMTID:        bmtID,
		CabangID:     cabangID,
		Judul:        judul,
		Isi:          isi,
		Tipe:         tipe,
		IsPinned:     false,
		TanggalMulai: now,
		DibuatOleh:   dibuatOleh,
		IsAktif:      true,
		CreatedAt:    now,
	}, nil
}
