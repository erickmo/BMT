package service

import (
	"context"

	"github.com/bmt-saas/api/internal/domain/nasabah"
	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/google/uuid"
)

type NasabahService struct {
	repo    nasabah.Repository
	rekRepo rekening.Repository
}

func NewNasabahService(repo nasabah.Repository, rekRepo rekening.Repository) *NasabahService {
	return &NasabahService{repo: repo, rekRepo: rekRepo}
}

func (s *NasabahService) GetByID(ctx context.Context, id, bmtID uuid.UUID) (*nasabah.Nasabah, error) {
	n, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if n.BMTID != bmtID {
		return nil, nasabah.ErrNasabahNotFound
	}
	return n, nil
}

func (s *NasabahService) GetByNomor(ctx context.Context, bmtID uuid.UUID, nomor string) (*nasabah.Nasabah, error) {
	return s.repo.GetByNomor(ctx, bmtID, nomor)
}

func (s *NasabahService) Search(ctx context.Context, bmtID uuid.UUID, query string, page, perPage int) ([]*nasabah.Nasabah, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.repo.Search(ctx, bmtID, query, perPage, offset)
}

func (s *NasabahService) ListRekening(ctx context.Context, nasabahID, bmtID uuid.UUID) ([]*rekening.Rekening, error) {
	n, err := s.repo.GetByID(ctx, nasabahID)
	if err != nil {
		return nil, err
	}
	if n.BMTID != bmtID {
		return nil, nasabah.ErrNasabahNotFound
	}
	return s.rekRepo.ListByNasabah(ctx, nasabahID)
}

func (s *NasabahService) GetMutasi(ctx context.Context, rekeningID, bmtID uuid.UUID, limit, offset int) ([]*rekening.TransaksiRekening, int64, error) {
	rek, err := s.rekRepo.GetByID(ctx, rekeningID)
	if err != nil {
		return nil, 0, err
	}
	if rek.BMTID != bmtID {
		return nil, 0, rekening.ErrRekeningNotFound
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.rekRepo.ListTransaksi(ctx, rekeningID, limit, offset)
}
