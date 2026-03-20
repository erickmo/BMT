package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

// ZakatService mengelola perhitungan zakat mal akhir tahun.
// Berjalan via worker HandleHitungZakat setiap 1 Muharram / akhir tahun hijriyah.
//
// Konfigurasi dari settings:
//   - ZAKAT_NISAB_RUPIAH  — batas minimal saldo wajib zakat (0 = fitur nonaktif)
//   - ZAKAT_RATE_PERSEN   — persentase zakat (default "2.5")
type ZakatService struct {
	rekeningRepo     rekening.Repository
	settingsResolver *settings.Resolver
}

func NewZakatService(rekeningRepo rekening.Repository, settingsResolver *settings.Resolver) *ZakatService {
	return &ZakatService{
		rekeningRepo:     rekeningRepo,
		settingsResolver: settingsResolver,
	}
}

// HitungZakatBMT menghitung kewajiban zakat mal untuk deposito aktif BMT.
// Mengembalikan jumlah rekening yang memenuhi syarat nisab.
// Penyimpanan ke tabel zakat dan notifikasi akan diimplementasikan pada sprint berikutnya.
func (s *ZakatService) HitungZakatBMT(ctx context.Context, bmtID uuid.UUID) (int, error) {
	nisabStr := s.settingsResolver.ResolveWithDefault(ctx, bmtID, uuid.Nil, "ZAKAT_NISAB_RUPIAH", "0")
	nisab, err := strconv.ParseInt(nisabStr, 10, 64)
	if err != nil || nisab <= 0 {
		return 0, nil // fitur zakat tidak dikonfigurasi
	}

	rateStr := s.settingsResolver.ResolveWithDefault(ctx, bmtID, uuid.Nil, "ZAKAT_RATE_PERSEN", "2.5")
	rate, _ := strconv.ParseFloat(rateStr, 64)
	if rate <= 0 {
		rate = 2.5
	}

	// Gunakan deposito aktif sebagai basis perhitungan zakat mal
	rekeningen, err := s.rekeningRepo.ListDepositoAktif(ctx, bmtID)
	if err != nil {
		return 0, fmt.Errorf("gagal ambil rekening deposito: %w", err)
	}

	count := 0
	for _, rek := range rekeningen {
		if int64(rek.Saldo) >= nisab {
			nominal := int64(float64(int64(rek.Saldo)) * rate / 100)
			// TODO Sprint 6: simpan ke tabel kewajiban_zakat dan kirim notifikasi
			_ = nominal
			count++
		}
	}

	fmt.Printf("[ZakatService] BMT %s: %d rekening deposito memenuhi nisab (Rp %d)\n",
		bmtID, count, nisab)
	return count, nil
}
