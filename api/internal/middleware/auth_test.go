package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testAccessSecret  = "test-access-secret-minimal-32-chars-ok"
	testRefreshSecret = "test-refresh-secret-minimal-32-chars"
)

func buatJWTManager() *jwt.Manager {
	return jwt.NewManager(testAccessSecret, testRefreshSecret, 15*time.Minute, 7*24*time.Hour)
}

func buatTokenValid(t *testing.T, role string) string {
	t.Helper()
	mgr := buatJWTManager()
	token, err := mgr.GenerateAccessToken(uuid.New(), uuid.New(), uuid.New(), role)
	if err != nil {
		t.Fatalf("gagal generate token test: %v", err)
	}
	return token
}

// ── Tests: Auth middleware ────────────────────────────────────────────────────

func TestAuth_TanpaHeader_401(t *testing.T) {
	mgr := buatJWTManager()
	handler := middleware.Auth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_FormatHeaderSalah_401(t *testing.T) {
	mgr := buatJWTManager()
	handler := middleware.Auth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_TokenTidakValid_401(t *testing.T) {
	mgr := buatJWTManager()
	handler := middleware.Auth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ini.token.palsu")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_TokenValid_200_DanContextTerisi(t *testing.T) {
	mgr := buatJWTManager()

	var capturedCtx context.Context
	handler := middleware.Auth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))

	userID := uuid.New()
	bmtID := uuid.New()
	cabangID := uuid.New()
	token, err := mgr.GenerateAccessToken(userID, bmtID, cabangID, "TELLER")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedCtx)
	assert.Equal(t, userID, middleware.GetUserID(capturedCtx))
	assert.Equal(t, bmtID, middleware.GetBMTID(capturedCtx))
	assert.Equal(t, cabangID, middleware.GetCabangID(capturedCtx))
	assert.Equal(t, "TELLER", middleware.GetRole(capturedCtx))
}

func TestAuth_RefreshTokenDitolak_401(t *testing.T) {
	mgr := buatJWTManager()
	handler := middleware.Auth(mgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Generate refresh token — tidak boleh diterima oleh Auth (yang hanya terima access token)
	refreshToken, err := mgr.GenerateRefreshToken(uuid.New(), uuid.New(), uuid.New(), "NASABAH")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Refresh token menggunakan secret berbeda dari access token, harus ditolak
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Tests: RequireRole middleware ─────────────────────────────────────────────

func TestRequireRole_RoleCocok_200(t *testing.T) {
	handler := middleware.RequireRole("TELLER", "ADMIN_BMT")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxRole, "TELLER")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_RoleTidakCocok_403(t *testing.T) {
	handler := middleware.RequireRole("ADMIN_BMT")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), middleware.CtxRole, "NASABAH")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_TanpaRole_403(t *testing.T) {
	handler := middleware.RequireRole("TELLER")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// Context tanpa role
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ── Tests: DeveloperAuth middleware ──────────────────────────────────────────

func TestDeveloperAuth_TokenBenar_200(t *testing.T) {
	devToken := "secret-developer-token-xyz"
	handler := middleware.DeveloperAuth(devToken)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/dev/test", nil)
	req.Header.Set("Developer-Token", devToken)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeveloperAuth_TokenSalah_401(t *testing.T) {
	devToken := "secret-developer-token-xyz"
	handler := middleware.DeveloperAuth(devToken)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/dev/test", nil)
	req.Header.Set("Developer-Token", "wrong-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeveloperAuth_TanpaHeader_401(t *testing.T) {
	devToken := "secret-developer-token-xyz"
	handler := middleware.DeveloperAuth(devToken)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/dev/test", nil)
	// Tidak ada header sama sekali
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ── Tests: Idempotency middleware ─────────────────────────────────────────────

func TestIdempotency_HeaderValid_KeyDiContext(t *testing.T) {
	keyUUID := uuid.New()
	var capturedCtx context.Context

	handler := middleware.Idempotency(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/transaksi", nil)
	req.Header.Set("X-Idempotency-Key", keyUUID.String())
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	key, ok := middleware.GetIdempotencyKey(capturedCtx)
	assert.True(t, ok, "idempotency key harus ada di context")
	assert.Equal(t, keyUUID, key)
}

func TestIdempotency_TanpaHeader_KeyTidakDiContext(t *testing.T) {
	var capturedCtx context.Context

	handler := middleware.Idempotency(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/transaksi", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	_, ok := middleware.GetIdempotencyKey(capturedCtx)
	assert.False(t, ok, "tanpa header, idempotency key tidak boleh ada di context")
}

func TestIdempotency_HeaderBukanUUID_KeyTidakDiContext(t *testing.T) {
	var capturedCtx context.Context

	handler := middleware.Idempotency(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/transaksi", nil)
	req.Header.Set("X-Idempotency-Key", "bukan-uuid-valid")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Middleware tidak boleh error — hanya skip parsing key yang tidak valid
	assert.Equal(t, http.StatusOK, w.Code)
	_, ok := middleware.GetIdempotencyKey(capturedCtx)
	assert.False(t, ok, "idempotency key tidak valid tidak boleh masuk context")
}
