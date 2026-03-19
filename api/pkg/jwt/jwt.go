package jwt

import (
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token tidak valid")
	ErrExpiredToken = errors.New("token sudah expired")
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	BMTID     uuid.UUID `json:"bmt_id"`
	CabangID  uuid.UUID `json:"cabang_id"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"` // access | refresh
	gojwt.RegisteredClaims
}

type Manager struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewManager(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *Manager {
	return &Manager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (m *Manager) GenerateAccessToken(userID, bmtID, cabangID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID:    userID,
		BMTID:     bmtID,
		CabangID:  cabangID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(m.accessExpiry)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.accessSecret))
}

func (m *Manager) GenerateRefreshToken(userID, bmtID, cabangID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID:    userID,
		BMTID:     bmtID,
		CabangID:  cabangID,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(m.refreshExpiry)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.refreshSecret))
}

func (m *Manager) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return m.verify(tokenStr, m.accessSecret)
}

func (m *Manager) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return m.verify(tokenStr, m.refreshSecret)
}

func (m *Manager) verify(tokenStr, secret string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenStr, &Claims{}, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, gojwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
