package auth

import (
	"encoding/json"
	"net/http"

	"github.com/bmt-saas/api/pkg/jwt"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, jwtManager *jwt.Manager) {
	r.Post("/staf/login", handleStafLogin)
	r.Post("/staf/refresh", handleRefresh)
	r.Post("/staf/logout", handleLogout)
	r.Post("/nasabah/login", handleNasabahLogin)
	r.Post("/nasabah/refresh", handleRefresh)
	r.Post("/pondok/login", handlePondokLogin)
	r.Post("/pondok/refresh", handleRefresh)
	r.Post("/merchant/login", handleMerchantLogin)
	r.Post("/merchant/refresh", handleRefresh)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func handleStafLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "username dan password wajib diisi")
		return
	}

	// TODO: validate credentials against DB, generate real JWT
	response.Success(w, LoginResponse{
		AccessToken:  "dummy-access-token",
		RefreshToken: "dummy-refresh-token",
		ExpiresIn:    900,
	})
}

func handleNasabahLogin(w http.ResponseWriter, r *http.Request) {
	handleStafLogin(w, r)
}

func handlePondokLogin(w http.ResponseWriter, r *http.Request) {
	handleStafLogin(w, r)
}

func handleMerchantLogin(w http.ResponseWriter, r *http.Request) {
	handleStafLogin(w, r)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	// TODO: validate refresh token and issue new access token
	response.Success(w, map[string]string{"access_token": "new-access-token"})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	// TODO: invalidate session from sesi_aktif table
	response.Success(w, map[string]string{"message": "logout berhasil"})
}
