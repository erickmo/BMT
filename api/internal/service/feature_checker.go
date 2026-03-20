package service

import (
	"context"

	"github.com/bmt-saas/api/internal/domain/platform"
	"github.com/google/uuid"
)

// PlatformFeatureChecker mengimplementasikan middleware.FeatureChecker
// dengan mengecek kontrak BMT aktif.
type PlatformFeatureChecker struct {
	repo platform.Repository
}

func NewPlatformFeatureChecker(repo platform.Repository) *PlatformFeatureChecker {
	return &PlatformFeatureChecker{repo: repo}
}

// FiturAktif memeriksa apakah kode fitur ada dan bernilai true di kontrak BMT aktif.
func (c *PlatformFeatureChecker) FiturAktif(ctx context.Context, bmtID uuid.UUID, kode string) error {
	kontrak, err := c.repo.GetKontrakAktif(ctx, bmtID)
	if err != nil {
		return platform.ErrFiturTidakAktif
	}
	val, ok := kontrak.Fitur[kode]
	if !ok {
		return platform.ErrFiturTidakAktif
	}
	if aktif, ok := val.(bool); !ok || !aktif {
		return platform.ErrFiturTidakAktif
	}
	return nil
}
