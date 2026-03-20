package auth

import (
	"encoding/json"
	"net/http"

	"github.com/bmt-saas/api/internal/domain/keamanan"
	"github.com/bmt-saas/api/internal/domain/nasabah"
	"github.com/bmt-saas/api/internal/domain/pengguna"
	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/jwt"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	authSvc    *service.AuthService
	sessionSvc *service.SessionService
	jwtManager *jwt.Manager
}

func NewHandler(authSvc *service.AuthService, sessionSvc *service.SessionService, jwtManager *jwt.Manager) *Handler {
	return &Handler{authSvc: authSvc, sessionSvc: sessionSvc, jwtManager: jwtManager}
}

// RegisterRoutes kompatibilitas lama — gunakan NewHandler untuk wiring lengkap.
func RegisterRoutes(r chi.Router, jwtManager *jwt.Manager) {
	h := &Handler{jwtManager: jwtManager}
	h.registerRoutes(r)
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	h.registerRoutes(r)
}

func (h *Handler) registerRoutes(r chi.Router) {
	r.Post("/staf/login", h.handleStafLogin)
	r.Post("/staf/login/otp", h.handleStafVerifikasiOTP)
	r.Post("/staf/refresh", h.handleRefresh)
	r.Post("/staf/logout", h.handleLogout)
	r.Post("/nasabah/login", h.handleNasabahLogin)
	r.Post("/nasabah/refresh", h.handleRefresh)
	r.Post("/pondok/login", h.handlePondokLogin)
	r.Post("/pondok/refresh", h.handleRefresh)
	r.Post("/merchant/login", h.handleMerchantLogin)
	r.Post("/merchant/refresh", h.handleRefresh)
}

// POST /auth/staf/login
func (h *Handler) handleStafLogin(w http.ResponseWriter, r *http.Request) {
	if h.authSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "auth service belum dikonfigurasi")
		return
	}

	var req struct {
		BMTID    uuid.UUID `json:"bmt_id"`
		Username string    `json:"username"`
		Password string    `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}
	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "username dan password wajib diisi")
		return
	}

	result, err := h.authSvc.LoginStaf(r.Context(), service.LoginStafInput{
		BMTID:     req.BMTID,
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: r.RemoteAddr,
		UserAgent: r.Header.Get("User-Agent"),
	})
	if err != nil {
		switch err {
		case pengguna.ErrPasswordSalah:
			response.Error(w, http.StatusUnauthorized, err.Error())
		case pengguna.ErrPenggunaBlokir, pengguna.ErrPenggunaNonAktif:
			response.Error(w, http.StatusForbidden, err.Error())
		case keamanan.ErrOTPBlokir:
			response.Error(w, http.StatusTooManyRequests, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Success(w, result)
}

// POST /auth/staf/login/otp  — verifikasi OTP 2FA
func (h *Handler) handleStafVerifikasiOTP(w http.ResponseWriter, r *http.Request) {
	if h.authSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "auth service belum dikonfigurasi")
		return
	}

	var req struct {
		BMTID    uuid.UUID `json:"bmt_id"`
		Username string    `json:"username"`
		Tujuan   string    `json:"tujuan"`
		KodeOTP  string    `json:"kode_otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}

	tokens, err := h.authSvc.VerifikasiOTPStaf(r.Context(),
		req.BMTID, req.Username, req.Tujuan, req.KodeOTP,
		r.RemoteAddr, r.Header.Get("User-Agent"),
	)
	if err != nil {
		switch err {
		case keamanan.ErrOTPSalah:
			response.Error(w, http.StatusUnauthorized, err.Error())
		case keamanan.ErrOTPExpired:
			response.Error(w, http.StatusUnauthorized, err.Error())
		case keamanan.ErrOTPBlokir:
			response.Error(w, http.StatusTooManyRequests, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Success(w, tokens)
}

// POST /auth/nasabah/login
func (h *Handler) handleNasabahLogin(w http.ResponseWriter, r *http.Request) {
	if h.authSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "auth service belum dikonfigurasi")
		return
	}

	var req struct {
		BMTID   uuid.UUID `json:"bmt_id"`
		Telepon string    `json:"telepon"`
		PIN     string    `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}

	tokens, err := h.authSvc.LoginNasabah(r.Context(), service.LoginNasabahInput{
		BMTID:     req.BMTID,
		Telepon:   req.Telepon,
		PIN:       req.PIN,
		IPAddress: r.RemoteAddr,
		UserAgent: r.Header.Get("User-Agent"),
	})
	if err != nil {
		switch err {
		case nasabah.ErrPINSalah:
			response.Error(w, http.StatusUnauthorized, err.Error())
		case nasabah.ErrNasabahNonAktif:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Success(w, tokens)
}

// POST /auth/staf/refresh
func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if h.sessionSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "session service belum dikonfigurasi")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		response.BadRequest(w, "refresh_token wajib diisi")
		return
	}

	tokens, err := h.sessionSvc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "sesi tidak valid atau sudah berakhir")
		return
	}
	response.Success(w, tokens)
}

// POST /auth/staf/logout
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if h.sessionSvc == nil {
		response.Success(w, map[string]string{"message": "logout berhasil"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	userID := middleware.GetUserID(r.Context())
	if req.RefreshToken != "" {
		_ = h.sessionSvc.CabutSesi(r.Context(), userID, req.RefreshToken)
	}

	response.Success(w, map[string]string{"message": "logout berhasil"})
}

func (h *Handler) handlePondokLogin(w http.ResponseWriter, r *http.Request) {
	// Pondok login gunakan mekanisme sama dengan staf — role ADMIN_PONDOK/OPERATOR_PONDOK
	h.handleStafLogin(w, r)
}

func (h *Handler) handleMerchantLogin(w http.ResponseWriter, r *http.Request) {
	h.handleStafLogin(w, r)
}
