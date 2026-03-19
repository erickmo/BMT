package platform

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/settings", func(r chi.Router) {
		r.Get("/", handleListSettings)
		r.Put("/{kunci}", handleUpdateSettings)
	})
	r.Route("/jenis-rekening", func(r chi.Router) {
		r.Get("/", handleListJenisRekening)
		r.Post("/", handleCreateJenisRekening)
		r.Put("/{id}", handleUpdateJenisRekening)
	})
	r.Route("/pengguna", func(r chi.Router) {
		r.Get("/", handleListPengguna)
		r.Post("/", handleCreatePengguna)
		r.Put("/{id}", handleUpdatePengguna)
	})
	r.Route("/cabang", func(r chi.Router) {
		r.Get("/", handleListCabang)
		r.Post("/", handleCreateCabang)
		r.Put("/{id}", handleUpdateCabang)
	})
	r.Route("/merchant", func(r chi.Router) {
		r.Get("/", handleListMerchant)
		r.Post("/", handleCreateMerchant)
		r.Put("/{id}", handleUpdateMerchant)
	})
	r.Route("/terminal-kiosk", func(r chi.Router) {
		r.Get("/", handleListTerminal)
		r.Post("/", handleCreateTerminal)
		r.Put("/{id}", handleUpdateTerminal)
	})
	r.Get("/laporan/konsolidasi", handleLaporanKonsolidasi)
	r.Get("/usage", handleUsage)
}

func handleListSettings(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "settings berhasil diupdate"})
}

func handleListJenisRekening(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJenisRekening(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "jenis rekening berhasil dibuat"})
}

func handleUpdateJenisRekening(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "jenis rekening berhasil diupdate"})
}

func handleListPengguna(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreatePengguna(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "pengguna berhasil dibuat"})
}

func handleUpdatePengguna(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "pengguna berhasil diupdate"})
}

func handleListCabang(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateCabang(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "cabang berhasil dibuat"})
}

func handleUpdateCabang(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "cabang berhasil diupdate"})
}

func handleListMerchant(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateMerchant(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "merchant berhasil dibuat"})
}

func handleUpdateMerchant(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "merchant berhasil diupdate"})
}

func handleListTerminal(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateTerminal(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "terminal kiosk berhasil dibuat"})
}

func handleUpdateTerminal(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "terminal kiosk berhasil diupdate"})
}

func handleLaporanKonsolidasi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleUsage(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}
