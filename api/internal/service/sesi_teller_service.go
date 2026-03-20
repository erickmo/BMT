package service

import (
	"context"
	"time"

	"github.com/bmt-saas/api/internal/domain/sesi_teller"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

type SesiTellerService struct {
	repo     sesi_teller.Repository
	settings *settings.Resolver
}

func NewSesiTellerService(repo sesi_teller.Repository, settingsResolver *settings.Resolver) *SesiTellerService {
	return &SesiTellerService{repo: repo, settings: settingsResolver}
}

type BukaSesiInput struct {
	BMTID        uuid.UUID
	CabangID     uuid.UUID
	TellerID     uuid.UUID
	Redenominasi []sesi_teller.ItemPecahan
}

func (s *SesiTellerService) BukaSesi(ctx context.Context, input BukaSesiInput) (*sesi_teller.SesiTeller, error) {
	existing, err := s.repo.GetAktifByTeller(ctx, input.TellerID)
	if err == nil && existing != nil {
		return nil, sesi_teller.ErrSesiSudahAktif
	}

	toleransiNominal := s.settings.ResolveInt(ctx, input.BMTID, input.CabangID, "sesi_teller.toleransi_selisih", 0)

	sesi := &sesi_teller.SesiTeller{
		ID:               uuid.New(),
		BMTID:            input.BMTID,
		CabangID:         input.CabangID,
		TellerID:         input.TellerID,
		Tanggal:          time.Now(),
		Redenominasi:     input.Redenominasi,
		Status:           sesi_teller.StatusAktif,
		ToleransiSelisih: money.New(int64(toleransiNominal)),
		DibukaPada:       time.Now(),
	}
	sesi.SaldoAwal = sesi.HitungSaldoAwal()

	if err := s.repo.Create(ctx, sesi); err != nil {
		return nil, err
	}
	return sesi, nil
}

func (s *SesiTellerService) GetSesiAktif(ctx context.Context, tellerID uuid.UUID) (*sesi_teller.SesiTeller, error) {
	return s.repo.GetAktifByTeller(ctx, tellerID)
}

func (s *SesiTellerService) TutupSesi(ctx context.Context, sesiID, tellerID uuid.UUID, redenominasiAkhir []sesi_teller.ItemPecahan) (*sesi_teller.SesiTeller, error) {
	sesi, err := s.repo.GetByID(ctx, sesiID)
	if err != nil {
		return nil, err
	}
	if sesi.TellerID != tellerID {
		return nil, sesi_teller.ErrSesiTidakAktif
	}

	tutupErr := sesi.TutupSesi(redenominasiAkhir, sesi.ToleransiSelisih)
	if updateErr := s.repo.Update(ctx, sesi); updateErr != nil {
		return nil, updateErr
	}
	return sesi, tutupErr
}
