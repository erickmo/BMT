package ecommerce

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// Toko (owner/pondok)
	r.Get("/toko", handleListToko)
	r.Post("/toko", handleCreateToko)
	r.Get("/toko/{id}", handleGetToko)
	r.Put("/toko/{id}", handleUpdateToko)

	// Produk
	r.Get("/produk", handleListProduk)
	r.Post("/produk", handleCreateProduk)
	r.Get("/produk/{id}", handleGetProduk)
	r.Put("/produk/{id}", handleUpdateProduk)
	r.Put("/produk/{id}/stok", handleUpdateStokProduk)

	// Pesanan (dari sisi seller/pondok)
	r.Get("/pesanan", handleListPesanan)
	r.Get("/pesanan/{id}", handleGetPesanan)
	r.Put("/pesanan/{id}/status", handleUpdateStatusPesanan)

	// Laporan
	r.Get("/laporan/penjualan", handleLaporanPenjualan)
}

func RegisterOPOPRoutes(r chi.Router) {
	// OPOP marketplace lintas pondok
	r.Get("/toko", handleListTokoOPOP)
	r.Get("/toko/{slug}", handleGetTokoOPOP)
	r.Get("/produk", handleListProdukOPOP)
	r.Post("/pesanan", handleCreatePesananB2B)
}

func handleListToko(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateToko(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "toko berhasil dibuat"})
}

func handleGetToko(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateToko(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "toko berhasil diupdate"})
}

func handleListProduk(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateProduk(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "produk berhasil ditambahkan"})
}

func handleGetProduk(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateProduk(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "produk berhasil diupdate"})
}

func handleUpdateStokProduk(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "stok produk berhasil diupdate"})
}

func handleListPesanan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleGetPesanan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateStatusPesanan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "status pesanan berhasil diupdate"})
}

func handleLaporanPenjualan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleListTokoOPOP(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleGetTokoOPOP(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"slug": chi.URLParam(r, "slug")})
}

func handleListProdukOPOP(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreatePesananB2B(w http.ResponseWriter, r *http.Request) {
	// B2B: pondok A pesan dari pondok B
	response.Created(w, map[string]string{"message": "pesanan B2B berhasil dibuat"})
}
