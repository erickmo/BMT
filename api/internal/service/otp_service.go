package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/keamanan"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// NotifSender mengirim OTP ke tujuan (SMS / Email).
// Implementasi nyata ada di Sprint 4 (NotifikasiService).
type NotifSender interface {
	KirimOTP(ctx context.Context, tujuan, kode string, channel keamanan.ChannelOTP) error
}

// LogNotifSender adalah stub yang hanya log — dipakai sebelum Sprint 4.
type LogNotifSender struct{}

func (l *LogNotifSender) KirimOTP(_ context.Context, tujuan, kode string, channel keamanan.ChannelOTP) error {
	// Pada production Sprint 4 ini akan diganti integrasi nyata
	fmt.Printf("[OTP STUB] kirim ke %s via %s: %s\n", tujuan, channel, kode)
	return nil
}

type OTPService struct {
	repo     keamanan.Repository
	redis    *redis.Client
	settings *settings.Resolver
	sender   NotifSender
}

func NewOTPService(
	repo keamanan.Repository,
	redisClient *redis.Client,
	settingsResolver *settings.Resolver,
	sender NotifSender,
) *OTPService {
	if sender == nil {
		sender = &LogNotifSender{}
	}
	return &OTPService{repo: repo, redis: redisClient, settings: settingsResolver, sender: sender}
}

const otpBruteKey = "otp:brute:%s:%s" // tujuan:tipe

// Generate membuat OTP baru dan mengirimnya ke tujuan.
func (s *OTPService) Generate(ctx context.Context, tujuan string, channel keamanan.ChannelOTP, tipe keamanan.TipeOTP, ip string) error {
	// Cek brute-force
	bruteKey := fmt.Sprintf(otpBruteKey, tujuan, tipe)
	attempts, _ := s.redis.Get(ctx, bruteKey).Int()
	maxAttempts := 3
	if attempts >= maxAttempts {
		return keamanan.ErrOTPBlokir
	}

	// Generate 6-digit OTP
	kode, err := generateOTPKode()
	if err != nil {
		return fmt.Errorf("gagal generate OTP: %w", err)
	}

	// Hash sebelum simpan
	hash, err := bcrypt.GenerateFromPassword([]byte(kode), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash OTP: %w", err)
	}

	// TTL dari settings (default 5 menit)
	ttlMenit := s.settings.ResolveInt(ctx, uuid.Nil, uuid.Nil, "otp.ttl_menit", 5)
	expiredAt := time.Now().Add(time.Duration(ttlMenit) * time.Minute)

	otpLog, err := keamanan.NewOTPLog(tujuan, channel, tipe, string(hash), expiredAt, ip)
	if err != nil {
		return err
	}

	if err := s.repo.CreateOTP(ctx, otpLog); err != nil {
		return fmt.Errorf("gagal simpan OTP: %w", err)
	}

	// Kirim ke tujuan
	return s.sender.KirimOTP(ctx, tujuan, kode, channel)
}

// Validasi memverifikasi OTP yang diinput user.
func (s *OTPService) Validasi(ctx context.Context, tujuan, kodeInput string, tipe keamanan.TipeOTP) error {
	bruteKey := fmt.Sprintf(otpBruteKey, tujuan, tipe)

	otpLog, err := s.repo.GetOTPByTujuanAndTipe(ctx, tujuan, tipe)
	if err != nil {
		s.incrBrute(ctx, bruteKey)
		return keamanan.ErrOTPNotFound
	}

	if otpLog.IsExpired() {
		return keamanan.ErrOTPExpired
	}

	if err := bcrypt.CompareHashAndPassword([]byte(otpLog.KodeHash), []byte(kodeInput)); err != nil {
		s.incrBrute(ctx, bruteKey)
		return keamanan.ErrOTPSalah
	}

	// Tandai sudah digunakan
	if err := s.repo.MarkOTPDigunakan(ctx, otpLog.ID); err != nil {
		return err
	}

	// Reset brute-force counter
	s.redis.Del(ctx, bruteKey)
	return nil
}

func (s *OTPService) incrBrute(ctx context.Context, key string) {
	pipe := s.redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 15*time.Minute)
	pipe.Exec(ctx)
}

func generateOTPKode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Konversi ke 6 digit (000000–999999)
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", n%1000000), nil
}
