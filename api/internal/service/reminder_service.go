package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/pembiayaan"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

// ReminderService mengirim reminder angsuran H-N sebelum jatuh tempo.
// N diambil dari setting "REMINDER_ANGSURAN_HARI" (default: 3).
// Notifikasi aktual dikirim via NotifikasiService (Sprint 4).
// Sprint 3: log saja, wiring notifikasi menyusul.
type ReminderService struct {
	pembiayaanRepo   pembiayaan.Repository
	settingsResolver *settings.Resolver
}

func NewReminderService(pembiayaanRepo pembiayaan.Repository, settingsResolver *settings.Resolver) *ReminderService {
	return &ReminderService{
		pembiayaanRepo:   pembiayaanRepo,
		settingsResolver: settingsResolver,
	}
}

// ReminderAngsuranBMT mengirim reminder untuk semua angsuran yang jatuh tempo H+N hari.
func (s *ReminderService) ReminderAngsuranBMT(ctx context.Context, bmtID uuid.UUID) error {
	// Ambil N dari settings
	nStr := s.settingsResolver.ResolveWithDefault(ctx, bmtID, uuid.Nil, "REMINDER_ANGSURAN_HARI", "3")
	var n int
	fmt.Sscanf(nStr, "%d", &n)
	if n <= 0 {
		n = 3
	}

	targetTanggal := time.Now().AddDate(0, 0, n).Truncate(24 * time.Hour)

	angsurans, err := s.pembiayaanRepo.GetAngsuranJatuhTempo(ctx, bmtID, targetTanggal)
	if err != nil {
		return fmt.Errorf("gagal ambil angsuran jatuh tempo: %w", err)
	}

	for _, a := range angsurans {
		// Sprint 4: panggil NotifikasiService.Kirim() di sini
		// Untuk sekarang, log saja sebagai placeholder
		fmt.Printf("[REMINDER] angsuran %s pembiayaan %s jatuh tempo %s nominal %d\n",
			a.ID, a.PembiayaanID, a.TanggalJatuhTempo.Format("2006-01-02"), a.TotalAngsuran)
	}

	return nil
}
