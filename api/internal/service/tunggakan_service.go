package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainAutodebet "github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
)

var (
	ErrTunggakanNotFound   = errors.New("tunggakan tidak ditemukan")
	ErrNominalMelebihiSisa = errors.New("nominal bayar melebihi sisa tunggakan")
)

type TunggakanService struct {
	autodebetRepo   domainAutodebet.Repository
	rekeningService *RekeningService
}

func NewTunggakanService(autodebetRepo domainAutodebet.Repository, rekeningService *RekeningService) *TunggakanService {
	return &TunggakanService{
		autodebetRepo:   autodebetRepo,
		rekeningService: rekeningService,
	}
}

// ListByRekening mengembalikan semua tunggakan outstanding milik sebuah rekening.
func (s *TunggakanService) ListByRekening(ctx context.Context, rekeningID, bmtID uuid.UUID) ([]*domainAutodebet.Tunggakan, error) {
	tunggakans, err := s.autodebetRepo.ListTunggakanByRekening(ctx, rekeningID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil tunggakan: %w", err)
	}
	return tunggakans, nil
}

// BayarTunggakan membayar sebagian atau seluruh tunggakan dengan mendebit rekening nasabah.
// Nominal boleh parsial — jika nominal < sisa, status tetap OUTSTANDING.
// Jika nominal == sisa, status berubah LUNAS.
func (s *TunggakanService) BayarTunggakan(
	ctx context.Context,
	jadwalID uuid.UUID,
	rekeningID uuid.UUID,
	nominal money.Money,
	operatorID uuid.UUID,
) (*domainAutodebet.Tunggakan, error) {
	tunggakans, err := s.autodebetRepo.ListTunggakanByRekening(ctx, rekeningID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil tunggakan: %w", err)
	}

	var target *domainAutodebet.Tunggakan
	for _, t := range tunggakans {
		if t.JadwalID == jadwalID {
			target = t
			break
		}
	}
	if target == nil {
		return nil, ErrTunggakanNotFound
	}
	if nominal > target.NominalSisa {
		return nil, ErrNominalMelebihiSisa
	}

	// Debit rekening
	_, err = s.rekeningService.Tarik(ctx, rekening.PenarikanInput{
		RekeningID: rekeningID,
		Nominal:    nominal.Int64(),
		Keterangan: fmt.Sprintf("Bayar tunggakan autodebet %s", target.Jenis),
		CreatedBy:  operatorID,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal debit rekening: %w", err)
	}

	// Update tunggakan
	target.NominalTerbayar = target.NominalTerbayar.Add(nominal)
	target.NominalSisa = target.NominalSisa.Sub(nominal)
	target.UpdatedAt = time.Now()
	if target.NominalSisa == 0 {
		target.Status = "LUNAS"
	}

	if err := s.autodebetRepo.UpdateTunggakan(ctx, target); err != nil {
		return nil, fmt.Errorf("gagal update tunggakan: %w", err)
	}

	return target, nil
}
