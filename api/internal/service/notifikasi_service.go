package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bmt-saas/api/internal/domain/keamanan"
	"github.com/bmt-saas/api/internal/domain/notifikasi"
	"github.com/bmt-saas/api/pkg/notif"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

// NotifikasiService mengelola pengiriman notifikasi.
// Kirim(): render template dari DB, buat antrian.
// Providers dipakai saat worker HandleKirimNotifikasi memproses antrian.
type NotifikasiService struct {
	repo             notifikasi.Repository
	settingsResolver *settings.Resolver
	providers        map[notifikasi.Channel]notif.Provider
}

func NewNotifikasiService(
	repo notifikasi.Repository,
	settingsResolver *settings.Resolver,
) *NotifikasiService {
	return &NotifikasiService{
		repo:             repo,
		settingsResolver: settingsResolver,
		providers:        make(map[notifikasi.Channel]notif.Provider),
	}
}

// SetProvider mendaftarkan provider untuk channel tertentu.
func (s *NotifikasiService) SetProvider(channel notifikasi.Channel, p notif.Provider) {
	s.providers[channel] = p
}

type KirimInput struct {
	BMTID        uuid.UUID
	Channel      notifikasi.Channel
	TemplateKode string
	Tujuan       string
	Variables    map[string]string // untuk render template: {{nama}} → "Ali"
	DataEkstra   map[string]string // dikirim ke FCM data payload
}

// Kirim merender template dan membuat antrian notifikasi.
// Jika template tidak ditemukan, gunakan pesan fallback dari Variables["pesan"].
func (s *NotifikasiService) Kirim(ctx context.Context, input KirimInput) error {
	var pesan, subjek string

	tmpl, err := s.repo.GetTemplate(ctx, &input.BMTID, input.TemplateKode, input.Channel)
	if err != nil {
		// fallback ke pesan langsung dari Variables["pesan"]
		fallback, ok := input.Variables["pesan"]
		if !ok || fallback == "" {
			return fmt.Errorf("template %s tidak ditemukan dan tidak ada pesan fallback: %w", input.TemplateKode, err)
		}
		pesan = fallback
	} else {
		pesan = renderTemplate(tmpl.Isi, input.Variables)
		subjek = renderTemplate(tmpl.Judul, input.Variables)
	}

	var dataEkstraRaw json.RawMessage
	if len(input.DataEkstra) > 0 {
		raw, err := json.Marshal(input.DataEkstra)
		if err != nil {
			return fmt.Errorf("gagal marshal data ekstra: %w", err)
		}
		dataEkstraRaw = raw
	}

	antrian, err := notifikasi.NewAntrian(notifikasi.CreateAntrianInput{
		BMTID:      input.BMTID,
		Channel:    input.Channel,
		Tujuan:     input.Tujuan,
		Subjek:     subjek,
		Pesan:      pesan,
		DataEkstra: dataEkstraRaw,
	})
	if err != nil {
		return fmt.Errorf("gagal membuat antrian notifikasi: %w", err)
	}

	if err := s.repo.CreateAntrian(ctx, antrian); err != nil {
		return fmt.Errorf("gagal simpan antrian notifikasi: %w", err)
	}

	return nil
}

// renderTemplate menggantikan {{key}} dengan value dari variables map.
func renderTemplate(tpl string, variables map[string]string) string {
	result := tpl
	for key, val := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, val)
	}
	return result
}

// TriggerTransaksiInput berisi data untuk trigger notifikasi transaksi.
type TriggerTransaksiInput struct {
	BMTID          uuid.UUID
	JenisTransaksi string // SETOR | TARIK | TRANSFER
	NomorRekening  string
	Nominal        int64
	Saldo          int64
	FCMToken       string // token FCM nasabah (bisa kosong)
	Telepon        string // nomor WA nasabah (bisa kosong)
}

// TriggerTransaksi dipanggil setelah setor/tarik/transfer berhasil.
// Trigger FCM + WA dengan template TRANSAKSI_SETOR / TRANSAKSI_TARIK / TRANSAKSI_TRANSFER.
func (s *NotifikasiService) TriggerTransaksi(ctx context.Context, input TriggerTransaksiInput) error {
	templateKode := fmt.Sprintf("TRANSAKSI_%s", input.JenisTransaksi)

	variables := map[string]string{
		"nomor_rekening": input.NomorRekening,
		"nominal":        fmt.Sprintf("Rp %s", formatRupiah(input.Nominal)),
		"saldo":          fmt.Sprintf("Rp %s", formatRupiah(input.Saldo)),
		"jenis":          input.JenisTransaksi,
	}

	if input.FCMToken != "" {
		if err := s.Kirim(ctx, KirimInput{
			BMTID:        input.BMTID,
			Channel:      notifikasi.ChannelFCM,
			TemplateKode: templateKode,
			Tujuan:       input.FCMToken,
			Variables:    variables,
		}); err != nil {
			return fmt.Errorf("gagal kirim notifikasi FCM transaksi: %w", err)
		}
	}

	if input.Telepon != "" {
		if err := s.Kirim(ctx, KirimInput{
			BMTID:        input.BMTID,
			Channel:      notifikasi.ChannelWhatsApp,
			TemplateKode: templateKode,
			Tujuan:       input.Telepon,
			Variables:    variables,
		}); err != nil {
			return fmt.Errorf("gagal kirim notifikasi WhatsApp transaksi: %w", err)
		}
	}

	return nil
}

// KirimOTP dipanggil oleh OTPService untuk mengirim kode OTP.
// Channel menggunakan keamanan.ChannelOTP agar compatible dengan NotifSender interface.
func (s *NotifikasiService) KirimOTP(ctx context.Context, tujuan, kode string, channel keamanan.ChannelOTP) error {
	var notifChannel notifikasi.Channel
	switch channel {
	case keamanan.ChannelSMS:
		notifChannel = notifikasi.ChannelSMS
	case keamanan.ChannelEmail:
		notifChannel = notifikasi.ChannelEmail
	default:
		return fmt.Errorf("channel OTP tidak dikenal: %s", channel)
	}

	pesan := fmt.Sprintf("Kode OTP Anda: %s", kode)

	antrian, err := notifikasi.NewAntrian(notifikasi.CreateAntrianInput{
		BMTID:   uuid.Nil,
		Channel: notifChannel,
		Tujuan:  tujuan,
		Pesan:   pesan,
	})
	if err != nil {
		return fmt.Errorf("gagal membuat antrian OTP: %w", err)
	}

	if err := s.repo.CreateAntrian(ctx, antrian); err != nil {
		return fmt.Errorf("gagal simpan antrian OTP: %w", err)
	}

	return nil
}

// DeliverAntrian memproses satu antrian: kirim via provider, update status.
// Dipanggil oleh worker HandleKirimNotifikasi.
// Max percobaan: 3 (setelah itu status GAGAL).
func (s *NotifikasiService) DeliverAntrian(ctx context.Context, a *notifikasi.NotifikasiAntrian) error {
	if a.Percobaan >= 3 {
		if err := s.repo.UpdateStatusAntrian(ctx, a.ID, notifikasi.StatusGagal, "melebihi batas percobaan"); err != nil {
			return fmt.Errorf("gagal update status antrian GAGAL: %w", err)
		}
		return nil
	}

	provider, ok := s.providers[a.Channel]
	if !ok {
		if err := s.repo.IncrementPercobaan(ctx, a.ID); err != nil {
			return fmt.Errorf("gagal increment percobaan (provider tidak ditemukan): %w", err)
		}
		return fmt.Errorf("provider untuk channel %s tidak terdaftar", a.Channel)
	}

	err := provider.Kirim(ctx, a.Tujuan, a.Subjek, a.Pesan)
	if err != nil {
		incrErr := s.repo.IncrementPercobaan(ctx, a.ID)
		if incrErr != nil {
			return fmt.Errorf("gagal kirim dan gagal increment percobaan: kirim=%w, incr=%v", err, incrErr)
		}
		return fmt.Errorf("gagal kirim notifikasi via %s: %w", a.Channel, err)
	}

	now := time.Now()
	a.DikirimAt = &now

	if err := s.repo.UpdateStatusAntrian(ctx, a.ID, notifikasi.StatusTerkirim, ""); err != nil {
		return fmt.Errorf("gagal update status antrian TERKIRIM: %w", err)
	}

	return nil
}

// GetPendingAntrian mengembalikan antrian MENUNGGU untuk diproses oleh worker.
func (s *NotifikasiService) GetPendingAntrian(ctx context.Context, limit int) ([]*notifikasi.NotifikasiAntrian, error) {
	return s.repo.GetAntrianMenunggu(ctx, limit)
}

// formatRupiah memformat angka int64 dengan separator titik.
// Contoh: 1000000 → "1.000.000"
func formatRupiah(n int64) string {
	s := fmt.Sprintf("%d", n)
	result := make([]byte, 0, len(s)+(len(s)-1)/3)
	remainder := len(s) % 3
	if remainder == 0 {
		remainder = 3
	}
	result = append(result, s[:remainder]...)
	for i := remainder; i < len(s); i += 3 {
		result = append(result, '.')
		result = append(result, s[i:i+3]...)
	}
	return string(result)
}
