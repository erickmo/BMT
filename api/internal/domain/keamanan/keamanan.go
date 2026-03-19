package keamanan

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// SubjekTipe identifies the type of the authenticated subject
type SubjekTipe string

const (
	SubjekNasabah       SubjekTipe = "NASABAH"
	SubjekPengguna      SubjekTipe = "PENGGUNA"
	SubjekPenggunaPondok SubjekTipe = "PENGGUNA_PONDOK"
	SubjekDeveloper     SubjekTipe = "DEVELOPER"
)

// TipeOTP identifies the purpose of an OTP code
type TipeOTP string

const (
	TipeOTPLogin              TipeOTP = "LOGIN"
	TipeOTPResetPIN           TipeOTP = "RESET_PIN"
	TipeOTPKonfirmasiTransaksi TipeOTP = "KONFIRMASI_TRANSAKSI"
	TipeOTPUbahEmail          TipeOTP = "UBAH_EMAIL"
)

// ChannelOTP identifies the delivery channel for OTP
type ChannelOTP string

const (
	ChannelSMS   ChannelOTP = "SMS"
	ChannelEmail ChannelOTP = "EMAIL"
)

// StatusFraud identifies the review status of a fraud alert
type StatusFraud string

const (
	StatusFraudOpen          StatusFraud = "OPEN"
	StatusFraudReviewed      StatusFraud = "REVIEWED"
	StatusFraudFalsePositive StatusFraud = "FALSE_POSITIVE"
	StatusFraudConfirmed     StatusFraud = "CONFIRMED"
)

// TipeFraudRule identifies the detection mechanism of a fraud rule
type TipeFraudRule string

const (
	TipeFrekuensi TipeFraudRule = "FREKUENSI"
	TipeNominal   TipeFraudRule = "NOMINAL"
	TipeLokasi    TipeFraudRule = "LOKASI"
	TipeWaktu     TipeFraudRule = "WAKTU"
	TipeVelocity  TipeFraudRule = "VELOCITY"
)

// AksiFraud identifies the action taken when a fraud rule is triggered
type AksiFraud string

const (
	AksiLog           AksiFraud = "LOG"
	AksiNotifikasi    AksiFraud = "NOTIFIKASI"
	AksiBlokir        AksiFraud = "BLOKIR_SEMENTARA"
	AksiRequireOTP    AksiFraud = "REQUIRE_OTP"
)

// OTPLog stores a hashed OTP and its metadata — never store plaintext
type OTPLog struct {
	ID          uuid.UUID  `json:"id"`
	Tujuan      string     `json:"tujuan"`       // phone or email
	Channel     ChannelOTP `json:"channel"`
	KodeHash    string     `json:"-"`            // bcrypt hash — never expose
	Tipe        TipeOTP    `json:"tipe"`
	ReferensiID *uuid.UUID `json:"referensi_id,omitempty"`
	IsDigunakan bool       `json:"is_digunakan"`
	ExpiredAt   time.Time  `json:"expired_at"`
	IPAddress   string     `json:"ip_address,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// SesiAktif tracks active authentication sessions
type SesiAktif struct {
	ID                uuid.UUID  `json:"id"`
	SubjekID          uuid.UUID  `json:"subjek_id"`
	SubjekTipe        SubjekTipe `json:"subjek_tipe"`
	RefreshTokenHash  string     `json:"-"` // bcrypt hash — never expose
	DeviceInfo        json.RawMessage `json:"device_info,omitempty"`
	IPAddress         string     `json:"ip_address,omitempty"`
	LastActiveAt      time.Time  `json:"last_active_at"`
	ExpiredAt         time.Time  `json:"expired_at"`
	IsAktif           bool       `json:"is_aktif"`
	CreatedAt         time.Time  `json:"created_at"`
}

// FraudRule defines a rule for detecting suspicious transactions
type FraudRule struct {
	ID        uuid.UUID       `json:"id"`
	BMTID     *uuid.UUID      `json:"bmt_id,omitempty"` // NULL = applies to all BMTs
	Nama      string          `json:"nama"`
	Tipe      TipeFraudRule   `json:"tipe"`
	Kondisi   json.RawMessage `json:"kondisi"`
	Aksi      AksiFraud       `json:"aksi"`
	IsAktif   bool            `json:"is_aktif"`
	CreatedAt time.Time       `json:"created_at"`
}

// FraudAlert records a triggered fraud rule instance
type FraudAlert struct {
	ID            uuid.UUID   `json:"id"`
	BMTID         uuid.UUID   `json:"bmt_id"`
	RuleID        uuid.UUID   `json:"rule_id"`
	NasabahID     *uuid.UUID  `json:"nasabah_id,omitempty"`
	TransaksiID   *uuid.UUID  `json:"transaksi_id,omitempty"`
	Deskripsi     string      `json:"deskripsi"`
	Status        StatusFraud `json:"status"`
	DireviewOleh  *uuid.UUID  `json:"direview_oleh,omitempty"`
	DireviewAt    *time.Time  `json:"direview_at,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}

// AuditLog records every significant action in the system
type AuditLog struct {
	ID           uuid.UUID  `json:"id"`
	BMTID        *uuid.UUID `json:"bmt_id,omitempty"`
	SubjekID     uuid.UUID  `json:"subjek_id"`
	SubjekTipe   SubjekTipe `json:"subjek_tipe"`
	Aksi         string     `json:"aksi"`
	ResourceTipe string     `json:"resource_tipe,omitempty"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty"`
	DataSebelum  json.RawMessage `json:"data_sebelum,omitempty"`
	DataSesudah  json.RawMessage `json:"data_sesudah,omitempty"`
	IPAddress    string     `json:"ip_address,omitempty"`
	UserAgent    string     `json:"user_agent,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Repository interface {
	// OTP
	CreateOTP(ctx context.Context, o *OTPLog) error
	GetOTPByTujuanAndTipe(ctx context.Context, tujuan string, tipe TipeOTP) (*OTPLog, error)
	MarkOTPDigunakan(ctx context.Context, id uuid.UUID) error

	// Sesi
	CreateSesi(ctx context.Context, s *SesiAktif) error
	GetSesiByRefreshHash(ctx context.Context, hash string) (*SesiAktif, error)
	ListSesiBySubjek(ctx context.Context, subjekID uuid.UUID) ([]*SesiAktif, error)
	NonaktifkanSesi(ctx context.Context, id uuid.UUID) error
	NonaktifkanSemuaSesi(ctx context.Context, subjekID uuid.UUID) error
	UpdateLastActive(ctx context.Context, id uuid.UUID) error
	DeleteExpiredSesi(ctx context.Context) (int64, error)

	// Fraud
	CreateFraudRule(ctx context.Context, r *FraudRule) error
	GetFraudRuleByID(ctx context.Context, id uuid.UUID) (*FraudRule, error)
	ListFraudRuleAktif(ctx context.Context, bmtID *uuid.UUID) ([]*FraudRule, error)
	CreateFraudAlert(ctx context.Context, a *FraudAlert) error
	GetFraudAlertByID(ctx context.Context, id uuid.UUID) (*FraudAlert, error)
	ListFraudAlert(ctx context.Context, bmtID uuid.UUID, status *StatusFraud, page, perPage int) ([]*FraudAlert, int64, error)
	UpdateStatusFraudAlert(ctx context.Context, id uuid.UUID, status StatusFraud, reviewOleh uuid.UUID) error

	// Audit
	CreateAuditLog(ctx context.Context, l *AuditLog) error
	ListAuditLog(ctx context.Context, bmtID *uuid.UUID, subjekID *uuid.UUID, resourceTipe string, page, perPage int) ([]*AuditLog, int64, error)
	DeleteAuditLogOlderThan(ctx context.Context, cutoff time.Time) (int64, error)
}

func NewOTPLog(tujuan string, channel ChannelOTP, tipe TipeOTP, kodeHash string, expiredAt time.Time, ip string) (*OTPLog, error) {
	if tujuan == "" {
		return nil, errors.New("tujuan OTP wajib diisi")
	}
	if kodeHash == "" {
		return nil, errors.New("kode OTP wajib di-hash sebelum disimpan")
	}
	return &OTPLog{
		ID:          uuid.New(),
		Tujuan:      tujuan,
		Channel:     channel,
		KodeHash:    kodeHash,
		Tipe:        tipe,
		IsDigunakan: false,
		ExpiredAt:   expiredAt,
		IPAddress:   ip,
		CreatedAt:   time.Now(),
	}, nil
}

func (o *OTPLog) IsExpired() bool {
	return time.Now().After(o.ExpiredAt)
}

func NewSesiAktif(subjekID uuid.UUID, subjekTipe SubjekTipe, refreshHash string, deviceInfo json.RawMessage, ip string, expiredAt time.Time) (*SesiAktif, error) {
	if refreshHash == "" {
		return nil, errors.New("refresh token hash wajib diisi")
	}
	now := time.Now()
	return &SesiAktif{
		ID:               uuid.New(),
		SubjekID:         subjekID,
		SubjekTipe:       subjekTipe,
		RefreshTokenHash: refreshHash,
		DeviceInfo:       deviceInfo,
		IPAddress:        ip,
		LastActiveAt:     now,
		ExpiredAt:        expiredAt,
		IsAktif:          true,
		CreatedAt:        now,
	}, nil
}

func (s *SesiAktif) IsExpired() bool {
	return time.Now().After(s.ExpiredAt)
}

func NewAuditLog(
	bmtID *uuid.UUID,
	subjekID uuid.UUID,
	subjekTipe SubjekTipe,
	aksi string,
	resourceTipe string,
	resourceID *uuid.UUID,
	sebelum, sesudah json.RawMessage,
	ip, userAgent string,
) *AuditLog {
	return &AuditLog{
		ID:           uuid.New(),
		BMTID:        bmtID,
		SubjekID:     subjekID,
		SubjekTipe:   subjekTipe,
		Aksi:         aksi,
		ResourceTipe: resourceTipe,
		ResourceID:   resourceID,
		DataSebelum:  sebelum,
		DataSesudah:  sesudah,
		IPAddress:    ip,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}
}

func NewFraudAlert(bmtID, ruleID uuid.UUID, nasabahID, transaksiID *uuid.UUID, deskripsi string) (*FraudAlert, error) {
	if deskripsi == "" {
		return nil, errors.New("deskripsi fraud alert wajib diisi")
	}
	return &FraudAlert{
		ID:          uuid.New(),
		BMTID:       bmtID,
		RuleID:      ruleID,
		NasabahID:   nasabahID,
		TransaksiID: transaksiID,
		Deskripsi:   deskripsi,
		Status:      StatusFraudOpen,
		CreatedAt:   time.Now(),
	}, nil
}
