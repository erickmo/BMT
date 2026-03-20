package service

import (
	"context"
	"errors"

	"github.com/bmt-saas/api/internal/domain/keamanan"
	"github.com/bmt-saas/api/internal/domain/nasabah"
	"github.com/bmt-saas/api/internal/domain/pengguna"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	penggunaRepo pengguna.Repository
	nasabahRepo  nasabah.Repository
	sessionSvc   *SessionService
	otpSvc       *OTPService
	settings     *settings.Resolver
}

func NewAuthService(
	penggunaRepo pengguna.Repository,
	nasabahRepo nasabah.Repository,
	sessionSvc *SessionService,
	otpSvc *OTPService,
	settingsResolver *settings.Resolver,
) *AuthService {
	return &AuthService{
		penggunaRepo: penggunaRepo,
		nasabahRepo:  nasabahRepo,
		sessionSvc:   sessionSvc,
		otpSvc:       otpSvc,
		settings:     settingsResolver,
	}
}

type LoginStafInput struct {
	BMTID     uuid.UUID
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

type LoginNasabahInput struct {
	BMTID     uuid.UUID
	Telepon   string
	PIN       string
	IPAddress string
	UserAgent string
}

type LoginResult struct {
	Tokens       *TokenPair `json:"tokens,omitempty"`
	PendingOTP   bool       `json:"pending_otp,omitempty"`
	OTPChannel   string     `json:"otp_channel,omitempty"`
	OTPTujuan    string     `json:"otp_tujuan,omitempty"`
}

// LoginStaf memvalidasi kredensial staf dan membuat sesi.
func (s *AuthService) LoginStaf(ctx context.Context, input LoginStafInput) (*LoginResult, error) {
	p, err := s.penggunaRepo.GetByUsername(ctx, input.BMTID, input.Username)
	if err != nil {
		if errors.Is(err, pengguna.ErrPenggunaNotFound) {
			return nil, pengguna.ErrPasswordSalah
		}
		return nil, err
	}

	if p.Status == pengguna.StatusBlokir {
		return nil, pengguna.ErrPenggunaBlokir
	}
	if p.Status == pengguna.StatusNonAktif {
		return nil, pengguna.ErrPenggunaNonAktif
	}

	if err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(input.Password)); err != nil {
		return nil, pengguna.ErrPasswordSalah
	}

	// Cek apakah 2FA wajib
	is2FA := s.settings.ResolveBool(ctx, input.BMTID, uuid.Nil, "keamanan.2fa_wajib_staf", false)
	if is2FA && p.Email != "" {
		channel := keamanan.ChannelEmail
		if p.Telepon != "" {
			channel = keamanan.ChannelSMS
		}
		tujuan := p.Email
		if channel == keamanan.ChannelSMS {
			tujuan = p.Telepon
		}
		if err := s.otpSvc.Generate(ctx, tujuan, channel, keamanan.TipeOTPLogin, input.IPAddress); err != nil {
			return nil, err
		}
		return &LoginResult{
			PendingOTP: true,
			OTPChannel: string(channel),
			OTPTujuan:  tujuan,
		}, nil
	}

	// Update last login
	_ = s.penggunaRepo.UpdateLastLogin(ctx, p.ID)

	cabangID := uuid.Nil
	if p.CabangID != nil {
		cabangID = *p.CabangID
	}

	tokens, err := s.sessionSvc.BuatSesi(ctx, CreateSessionInput{
		SubjekID:   p.ID,
		SubjekTipe: keamanan.SubjekPengguna,
		BMTID:      p.BMTID,
		CabangID:   cabangID,
		Role:       p.Role,
		DeviceInfo: map[string]string{"user_agent": input.UserAgent},
		IPAddress:  input.IPAddress,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{Tokens: tokens}, nil
}

// VerifikasiOTPStaf memvalidasi OTP 2FA lalu issue token.
func (s *AuthService) VerifikasiOTPStaf(ctx context.Context, bmtID uuid.UUID, username, tujuan, kodeOTP, ip, userAgent string) (*TokenPair, error) {
	if err := s.otpSvc.Validasi(ctx, tujuan, kodeOTP, keamanan.TipeOTPLogin); err != nil {
		return nil, err
	}

	p, err := s.penggunaRepo.GetByUsername(ctx, bmtID, username)
	if err != nil {
		return nil, err
	}

	_ = s.penggunaRepo.UpdateLastLogin(ctx, p.ID)

	cabangID := uuid.Nil
	if p.CabangID != nil {
		cabangID = *p.CabangID
	}

	return s.sessionSvc.BuatSesi(ctx, CreateSessionInput{
		SubjekID:   p.ID,
		SubjekTipe: keamanan.SubjekPengguna,
		BMTID:      p.BMTID,
		CabangID:   cabangID,
		Role:       p.Role,
		DeviceInfo: map[string]string{"user_agent": userAgent},
		IPAddress:  ip,
	})
}

// LoginNasabah memvalidasi telepon + PIN nasabah dan membuat sesi.
func (s *AuthService) LoginNasabah(ctx context.Context, input LoginNasabahInput) (*TokenPair, error) {
	// Cari nasabah by telepon (search by exact match)
	list, _, err := s.nasabahRepo.Search(ctx, input.BMTID, input.Telepon, 1, 0)
	if err != nil || len(list) == 0 {
		return nil, nasabah.ErrPINSalah
	}
	n := list[0]

	if n.Status == nasabah.StatusBlokir {
		return nil, nasabah.ErrNasabahNonAktif
	}
	if n.Status == nasabah.StatusNonAktif {
		return nil, nasabah.ErrNasabahNonAktif
	}

	if err := bcrypt.CompareHashAndPassword([]byte(n.PINHash), []byte(input.PIN)); err != nil {
		return nil, nasabah.ErrPINSalah
	}

	return s.sessionSvc.BuatSesi(ctx, CreateSessionInput{
		SubjekID:   n.ID,
		SubjekTipe: keamanan.SubjekNasabah,
		BMTID:      n.BMTID,
		CabangID:   n.CabangID,
		Role:       "NASABAH",
		DeviceInfo: map[string]string{"user_agent": input.UserAgent},
		IPAddress:  input.IPAddress,
	})
}
