package teller

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Post("/sesi/buka", handleBukaSesi)
	r.Get("/sesi/aktif", handleGetSesiAktif)
	r.Post("/sesi/tutup", handleTutupSesi)
	r.Get("/nasabah/cari", handleCariNasabah)
	r.Post("/rekening/{id}/setor", handleSetor)
	r.Post("/rekening/{id}/tarik", handleTarik)
	r.Post("/pembiayaan/{id}/angsuran", handleBayarAngsuran)
	r.Post("/spp/{id}/bayar", handleBayarSPP)
}

func handleBukaSesi(w http.ResponseWriter, r *http.Request) {
	// TODO: validasi redenominasi (pecahan dari DB), bukan konstanta
	response.Created(w, map[string]string{"message": "sesi teller dibuka"})
}

func handleGetSesiAktif(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleTutupSesi(w http.ResponseWriter, r *http.Request) {
	// TODO: validasi selisih saldo kas sesuai settings sesi_teller.toleransi_selisih
	response.Success(w, map[string]string{"message": "sesi teller ditutup"})
}

func handleCariNasabah(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleSetor(w http.ResponseWriter, r *http.Request) {
	// TODO: delegasi ke RekeningService.Setor, cek sesi teller aktif
	response.Created(w, map[string]string{"message": "setoran berhasil"})
}

func handleTarik(w http.ResponseWriter, r *http.Request) {
	// TODO: delegasi ke RekeningService.Tarik, cek sesi teller aktif
	response.Created(w, map[string]string{"message": "penarikan berhasil"})
}

func handleBayarAngsuran(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "angsuran berhasil dibayar"})
}

func handleBayarSPP(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "SPP berhasil dibayar"})
}
