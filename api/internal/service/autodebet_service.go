package service

import (
	"context"
	"fmt"
	"time"

	domainAutodebet "github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/google/uuid"
)

type AutodebetService struct {
	autodebetRepo   domainAutodebet.Repository
	rekeningService *RekeningService
}

func NewAutodebetService(autodebetRepo domainAutodebet.Repository, rekeningService *RekeningService) *AutodebetService {
	return &AutodebetService{
		autodebetRepo:   autodebetRepo,
		rekeningService: rekeningService,
	}
}

// GenerateJadwalBulanan membuat jadwal autodebet untuk bulan target
// berdasarkan konfigurasi autodebet masing-masing rekening.
func (s *AutodebetService) GenerateJadwalBulanan(ctx context.Context, rekeningIDs []uuid.UUID, bmtID uuid.UUID, bulan time.Time) error {
	for _, rekeningID := range rekeningIDs {
		configs, err := s.autodebetRepo.ListConfigByRekening(ctx, rekeningID)
		if err != nil {
			// Log error tapi lanjutkan ke rekening berikutnya
			fmt.Printf("gagal ambil config autodebet rekening %s: %v\n", rekeningID, err)
			continue
		}

		for _, cfg := range configs {
			if !cfg.IsAktif {
				continue
			}

			// Hitung tanggal jatuh tempo bulan target.
			// Tanggal 29/30/31 otomatis disesuaikan ke hari terakhir bulan.
			tanggal := HitungTanggalJatuhTempo(bulan, int(cfg.TanggalDebet))

			// Nominal target ditentukan dari jenis autodebet.
			// Untuk ANGSURAN_PEMBIAYAAN dan SPP_PONDOK, nominal di-resolve
			// saat eksekusi dari referensi (pembiayaan_id / tagihan_spp_id).
			// Jadwal dengan nominal 0 tetap dibuat sebagai placeholder.
			jadwal := &domainAutodebet.Jadwal{
				ID:                uuid.New(),
				BMTID:             bmtID,
				RekeningID:        cfg.RekeningID,
				ConfigID:          cfg.ID,
				Jenis:             cfg.Jenis,
				NominalTarget:     0,
				TanggalJatuhTempo: tanggal,
				Status:            domainAutodebet.StatusMenunggu,
				CreatedAt:         time.Now(),
			}

			if err := s.autodebetRepo.CreateJadwal(ctx, jadwal); err != nil {
				fmt.Printf("gagal buat jadwal autodebet rekening %s: %v\n", rekeningID, err)
			}
		}
	}

	return nil
}

// HitungTanggalJatuhTempo menghitung tanggal jatuh tempo autodebet.
// Jika tanggal melebihi hari terakhir bulan, disesuaikan ke akhir bulan.
// Parameter bulan adalah awal bulan target (misal time.Date(2025, 2, 1, ...)).
// Parameter tanggal adalah hari dalam bulan (1-28).
func HitungTanggalJatuhTempo(bulan time.Time, tanggal int) time.Time {
	// Cari hari terakhir bulan
	bulanBerikutnya := time.Date(bulan.Year(), bulan.Month()+1, 1, 0, 0, 0, 0, bulan.Location())
	hariTerakhir := bulanBerikutnya.AddDate(0, 0, -1).Day()

	hari := tanggal
	if hari > hariTerakhir {
		hari = hariTerakhir
	}

	return time.Date(bulan.Year(), bulan.Month(), hari, 0, 0, 0, 0, bulan.Location())
}

// EksekusiBulanan menjalankan autodebet untuk BMT tertentu pada tanggal hari ini.
// Merupakan alias EksekusiHarian untuk digunakan oleh worker bulanan.
func (s *AutodebetService) EksekusiBulanan(ctx context.Context, bmtID uuid.UUID) error {
	return s.EksekusiHarian(ctx, bmtID, time.Now())
}

// GenerateJadwalBulanDepan membuat jadwal autodebet bulan depan untuk semua rekening aktif BMT.
// Berbeda dengan GenerateJadwalBulanan yang memerlukan rekeningIDs eksplisit.
func (s *AutodebetService) GenerateJadwalBulanDepan(ctx context.Context, bmtID uuid.UUID) error {
	configs, err := s.autodebetRepo.ListConfigAktifByBMT(ctx, bmtID)
	if err != nil {
		return fmt.Errorf("gagal ambil config autodebet BMT %s: %w", bmtID, err)
	}

	now := time.Now()
	bulanDepan := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	for _, cfg := range configs {
		tanggal := HitungTanggalJatuhTempo(bulanDepan, int(cfg.TanggalDebet))
		jadwal := &domainAutodebet.Jadwal{
			ID:                uuid.New(),
			BMTID:             bmtID,
			RekeningID:        cfg.RekeningID,
			ConfigID:          cfg.ID,
			Jenis:             cfg.Jenis,
			NominalTarget:     0,
			TanggalJatuhTempo: tanggal,
			Status:            domainAutodebet.StatusMenunggu,
			CreatedAt:         time.Now(),
		}
		if err := s.autodebetRepo.CreateJadwal(ctx, jadwal); err != nil {
			fmt.Printf("gagal buat jadwal autodebet rekening %s: %v\n", cfg.RekeningID, err)
		}
	}

	return nil
}

// EksekusiHarian menjalankan semua jadwal autodebet untuk hari ini.
// Menggunakan partial debit: jika saldo tidak cukup, debit semampu saldo,
// sisanya dicatat sebagai tunggakan.
func (s *AutodebetService) EksekusiHarian(ctx context.Context, bmtID uuid.UUID, tanggal time.Time) error {
	jadwals, err := s.autodebetRepo.ListJadwalByTanggal(ctx, bmtID, tanggal)
	if err != nil {
		return fmt.Errorf("gagal ambil jadwal: %w", err)
	}

	for _, jadwal := range jadwals {
		if jadwal.Status != domainAutodebet.StatusMenunggu {
			continue
		}
		if err := s.rekeningService.EksekusiAutodebetJadwal(ctx, jadwal); err != nil {
			// Log error tapi lanjutkan ke jadwal berikutnya
			fmt.Printf("gagal eksekusi autodebet jadwal %s: %v\n", jadwal.ID, err)
		}
	}

	return nil
}
