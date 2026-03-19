package developer

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/health/detail", handleHealthDetail)
	r.Route("/bmt", func(r chi.Router) {
		r.Get("/", handleListBMT)
		r.Post("/", handleCreateBMT)
		r.Get("/{id}", handleGetBMT)
		r.Put("/{id}", handleUpdateBMT)
		r.Put("/{id}/status", handleUpdateBMTStatus)
		r.Post("/{id}/kontrak", handleCreateKontrak)
		r.Get("/{id}/cabang", handleListCabang)
		r.Post("/{id}/cabang", handleCreateCabang)
		r.Post("/{id}/pengguna/seed", handleSeedPengguna)
	})
	r.Route("/pecahan-uang", func(r chi.Router) {
		r.Get("/", handleListPecahan)
		r.Post("/", handleCreatePecahan)
		r.Put("/{id}", handleUpdatePecahan)
		r.Delete("/{id}", handleDeletePecahan)
	})
	r.Route("/platform-settings", func(r chi.Router) {
		r.Get("/", handleListPlatformSettings)
		r.Put("/{kunci}", handleUpdatePlatformSettings)
	})
	r.Get("/usage-log", handleListUsageLog)
	r.Get("/metrics", handleMetrics)
	r.Post("/maintenance", handleMaintenance)
}

func handleHealthDetail(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"status": "ok", "db": "ok", "redis": "ok"})
}

func handleListBMT(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateBMT(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "BMT berhasil dibuat"})
}

func handleGetBMT(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateBMT(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "BMT berhasil diupdate"})
}

func handleUpdateBMTStatus(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "Status BMT berhasil diupdate"})
}

func handleCreateKontrak(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "Kontrak berhasil dibuat"})
}

func handleListCabang(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateCabang(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "Cabang berhasil dibuat"})
}

func handleSeedPengguna(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "Pengguna berhasil di-seed"})
}

func handleListPecahan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreatePecahan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "Pecahan uang berhasil dibuat"})
}

func handleUpdatePecahan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "Pecahan uang berhasil diupdate"})
}

func handleDeletePecahan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "Pecahan uang berhasil dihapus"})
}

func handleListPlatformSettings(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleUpdatePlatformSettings(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "Settings berhasil diupdate"})
}

func handleListUsageLog(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{"uptime": "ok"})
}

func handleMaintenance(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "Maintenance mode updated"})
}
