package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

// DistribusiService menghitung dan mendistribusikan bagi hasil deposito setiap akhir bulan.
// Bagi hasil dihitung berdasarkan:
//   saldo_deposito × (nisbah_nasabah / 100) × rate_bulanan_dari_settings
// Rate bulanan diambil dari setting "DEPOSITO_RATE_BULANAN_PERSEN" (misal: 0.5 = 0.5% per bulan).
type DistribusiService struct {
	rekeningRepo     rekening.Repository
	rekeningService  *RekeningService
	settingsResolver *settings.Resolver
}

func NewDistribusiService(
	rekeningRepo rekening.Repository,
	rekeningService *RekeningService,
	settingsResolver *settings.Resolver,
) *DistribusiService {
	return &DistribusiService{
		rekeningRepo:    rekeningRepo,
		rekeningService: rekeningService,
		settingsResolver: settingsResolver,
	}
}

// HasilDistribusi menyimpan ringkasan distribusi bagi hasil untuk satu rekening.
type HasilDistribusi struct {
	RekeningID uuid.UUID
	Nominal    money.Money
}

// DistribusiBagiHasil mendistribusikan bagi hasil ke semua rekening deposito aktif milik BMT.
// Parameter bulan adalah bulan yang sedang diselesaikan (biasanya time.Now()).
// Mengembalikan slice HasilDistribusi dan error jika query gagal.
func (s *DistribusiService) DistribusiBagiHasil(ctx context.Context, bmtID uuid.UUID, bulan time.Time) ([]HasilDistribusi, error) {
	// Ambil rate bulanan dari settings (dalam persen, misal 0.5 = 0,5%/bulan)
	rateStr := s.settingsResolver.ResolveWithDefault(ctx, bmtID, uuid.Nil, "DEPOSITO_RATE_BULANAN_PERSEN", "0.5")
	var ratePersen float64
	fmt.Sscanf(rateStr, "%f", &ratePersen)

	depositoList, err := s.rekeningRepo.ListDepositoAktif(ctx, bmtID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil rekening deposito: %w", err)
	}

	var hasil []HasilDistribusi

	for _, rek := range depositoList {
		nisbah := float64(50) // default 50% jika nisbah_nasabah tidak di-set
		if rek.NisbahNasabah != nil {
			nisbah = float64(*rek.NisbahNasabah)
		}

		// Bagi hasil = saldo × rate_bulanan × (nisbah_nasabah / 100)
		saldo := float64(rek.Saldo.Int64())
		nominalFloat := saldo * (ratePersen / 100) * (nisbah / 100)
		nominal := money.New(int64(nominalFloat))

		if nominal <= 0 {
			continue
		}

		// Setor bagi hasil ke rekening nasabah
		_, err := s.rekeningService.Setor(ctx, rekening.SetoranInput{
			RekeningID: rek.ID,
			Nominal:    nominal.Int64(),
			Keterangan: fmt.Sprintf("Bagi hasil deposito %s %d/%d", rek.NomorRekening, bulan.Month(), bulan.Year()),
		})
		if err != nil {
			fmt.Printf("gagal setor bagi hasil rekening %s: %v\n", rek.NomorRekening, err)
			continue
		}

		hasil = append(hasil, HasilDistribusi{
			RekeningID: rek.ID,
			Nominal:    nominal,
		})
	}

	return hasil, nil
}
