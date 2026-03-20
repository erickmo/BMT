package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bmt-saas/api/internal/domain/keamanan"
	"github.com/bmt-saas/api/pkg/jwt"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionService struct {
	repo       keamanan.Repository
	jwtManager *jwt.Manager
	settings   *settings.Resolver
}

func NewSessionService(repo keamanan.Repository, jwtManager *jwt.Manager, settingsResolver *settings.Resolver) *SessionService {
	return &SessionService{repo: repo, jwtManager: jwtManager, settings: settingsResolver}
}

type CreateSessionInput struct {
	SubjekID    uuid.UUID
	SubjekTipe  keamanan.SubjekTipe
	BMTID       uuid.UUID
	CabangID    uuid.UUID
	Role        string
	DeviceInfo  map[string]string
	IPAddress   string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

// BuatSesi membuat sesi baru, simpan refresh token hash, dan return token pair.
func (s *SessionService) BuatSesi(ctx context.Context, input CreateSessionInput) (*TokenPair, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(input.SubjekID, input.BMTID, input.CabangID, input.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(input.SubjekID, input.BMTID, input.CabangID, input.Role)
	if err != nil {
		return nil, err
	}

	// Hash refresh token sebelum simpan
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// TTL sesi dari settings (default 30 hari)
	ttlHari := s.settings.ResolveInt(ctx, input.BMTID, input.CabangID, "auth.refresh_ttl_hari", 30)

	deviceInfoJSON, _ := json.Marshal(input.DeviceInfo)

	sesi, err := keamanan.NewSesiAktif(
		input.SubjekID,
		input.SubjekTipe,
		string(refreshHash),
		json.RawMessage(deviceInfoJSON),
		input.IPAddress,
		time.Now().Add(time.Duration(ttlHari)*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateSesi(ctx, sesi); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 menit (access token)
	}, nil
}

// Refresh memvalidasi refresh token dan issue access token baru.
func (s *SessionService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, keamanan.ErrSesiNotFound
	}

	// Cari sesi berdasarkan prefix (bcrypt tidak bisa di-lookup langsung)
	// Strategi: cari semua sesi subjek, verify hash
	sesiList, err := s.repo.ListSesiBySubjek(ctx, claims.UserID)
	if err != nil {
		return nil, keamanan.ErrSesiNotFound
	}

	var sesiAktif *keamanan.SesiAktif
	for _, sesi := range sesiList {
		if !sesi.IsAktif {
			continue
		}
		if bcrypt.CompareHashAndPassword([]byte(sesi.RefreshTokenHash), []byte(refreshToken)) == nil {
			sesiAktif = sesi
			break
		}
	}
	if sesiAktif == nil {
		return nil, keamanan.ErrSesiNotFound
	}
	if sesiAktif.IsExpired() {
		return nil, keamanan.ErrSesiNotFound
	}

	// Update last active
	_ = s.repo.UpdateLastActive(ctx, sesiAktif.ID)

	// Issue access token baru
	var cabangID uuid.UUID
	if claims.CabangID != uuid.Nil {
		cabangID = claims.CabangID
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(claims.UserID, claims.BMTID, cabangID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

// CabutSesi menonaktifkan sesi berdasarkan refresh token.
func (s *SessionService) CabutSesi(ctx context.Context, subjekID uuid.UUID, refreshToken string) error {
	sesiList, err := s.repo.ListSesiBySubjek(ctx, subjekID)
	if err != nil {
		return err
	}
	for _, sesi := range sesiList {
		if bcrypt.CompareHashAndPassword([]byte(sesi.RefreshTokenHash), []byte(refreshToken)) == nil {
			return s.repo.NonaktifkanSesi(ctx, sesi.ID)
		}
	}
	return nil
}

// CabutSemuaSesi menonaktifkan semua sesi subjek (logout from all devices).
func (s *SessionService) CabutSemuaSesi(ctx context.Context, subjekID uuid.UUID) error {
	return s.repo.NonaktifkanSemuaSesi(ctx, subjekID)
}
