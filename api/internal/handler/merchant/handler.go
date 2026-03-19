package merchant

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// Kasir mode: tap NFC → PIN → konfirmasi → struk
	r.Post("/transaksi", handleTransaksiKasir)
	r.Get("/transaksi/{id}", handleGetTransaksi)
	r.Get("/transaksi", handleListTransaksi)

	// Owner mode: laporan penjualan
	r.Get("/laporan/penjualan", handleLaporanPenjualan)
	r.Get("/laporan/bulanan", handleLaporanBulanan)
}

func handleTransaksiKasir(w http.ResponseWriter, r *http.Request) {
	// Mode Kasir: input nominal → tap NFC nasabah → input PIN → konfirmasi → struk
	response.Created(w, map[string]string{"message": "transaksi kasir berhasil"})
}

func handleGetTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleListTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleLaporanPenjualan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleLaporanBulanan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}
