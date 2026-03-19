package service

import (
	"context"

	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

type SettingsService struct {
	resolver *settings.Resolver
}

func NewSettingsService(resolver *settings.Resolver) *SettingsService {
	return &SettingsService{resolver: resolver}
}

func (s *SettingsService) GetJamBuka(ctx context.Context, bmtID, cabangID uuid.UUID) string {
	return s.resolver.ResolveWithDefault(ctx, bmtID, cabangID, "operasional.jam_buka", "08:00")
}

func (s *SettingsService) GetJamTutup(ctx context.Context, bmtID, cabangID uuid.UUID) string {
	return s.resolver.ResolveWithDefault(ctx, bmtID, cabangID, "operasional.jam_tutup", "16:00")
}

func (s *SettingsService) GetZonaWaktu(ctx context.Context, bmtID, cabangID uuid.UUID) string {
	return s.resolver.ResolveWithDefault(ctx, bmtID, cabangID, "operasional.zona_waktu", "Asia/Jakarta")
}

func (s *SettingsService) GetTanggalAutodebetSimpananWajib(ctx context.Context, bmtID, cabangID uuid.UUID) int {
	return s.resolver.ResolveInt(ctx, bmtID, cabangID, "autodebet.tanggal_simpanan_wajib", 1)
}

func (s *SettingsService) GetRetryAutodebet(ctx context.Context, bmtID, cabangID uuid.UUID) int {
	return s.resolver.ResolveInt(ctx, bmtID, cabangID, "autodebet.retry_hari", 3)
}

func (s *SettingsService) GetToleransiSelisihSesiTeller(ctx context.Context, bmtID, cabangID uuid.UUID) int64 {
	return int64(s.resolver.ResolveInt(ctx, bmtID, cabangID, "sesi_teller.toleransi_selisih", 0))
}

func (s *SettingsService) GetApprovers(ctx context.Context, bmtID, cabangID uuid.UUID, jenisForm string) []string {
	var approvers []string
	key := "approval." + jenisForm
	err := s.resolver.ResolveJSON(ctx, bmtID, cabangID, key, &approvers)
	if err != nil {
		return []string{"TELLER", "MANAJER_CABANG"}
	}
	return approvers
}

func (s *SettingsService) GetMetodeAbsensi(ctx context.Context, bmtID, cabangID uuid.UUID) []string {
	var metode []string
	err := s.resolver.ResolveJSON(ctx, bmtID, cabangID, "pondok.absensi_metode", &metode)
	if err != nil {
		return []string{"MANUAL"}
	}
	return metode
}

func (s *SettingsService) GetNFCLimitDefault(ctx context.Context, bmtID, cabangID uuid.UUID) int64 {
	return int64(s.resolver.ResolveInt(ctx, bmtID, cabangID, "nfc.limit_default_per_transaksi", 500000))
}
