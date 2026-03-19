package nasabah

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/profil", handleGetProfil)
	r.Get("/rekening", handleListRekening)
	r.Get("/rekening/{id}/transaksi", handleListTransaksi)
	r.Post("/rekening/{id}/setor-online", handleSetorOnline)
	r.Get("/pembiayaan", handleListPembiayaan)
	r.Get("/santri", handleGetSantri)
	r.Get("/nfc/saldo", handleGetNFCSaldo)
	r.Get("/nfc/transaksi", handleListNFCTransaksi)
	r.Post("/nfc/topup", handleTopupNFC)
	r.Get("/spp/tagihan", handleListTagihanSPP)
	r.Post("/spp/{id}/bayar", handleBayarSPP)

	// E-commerce
	r.Get("/shop/toko", handleListToko)
	r.Get("/shop/toko/{slug}/produk", handleListProdukToko)
	r.Post("/shop/keranjang", handleAddKeranjang)
	r.Post("/shop/pesanan", handleCreatePesanan)
	r.Get("/shop/pesanan", handleListPesanan)
	r.Get("/shop/pesanan/{id}", handleGetPesanan)
	r.Post("/shop/ulasan", handleCreateUlasan)
}

func handleGetProfil(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleListRekening(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleListTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleSetorOnline(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "setoran online berhasil"})
}

func handleListPembiayaan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleGetSantri(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleGetNFCSaldo(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleListNFCTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleTopupNFC(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "top-up NFC berhasil"})
}

func handleListTagihanSPP(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleBayarSPP(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "SPP berhasil dibayar"})
}

func handleListToko(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleListProdukToko(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleAddKeranjang(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "produk ditambahkan ke keranjang"})
}

func handleCreatePesanan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "pesanan berhasil dibuat"})
}

func handleListPesanan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleGetPesanan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleCreateUlasan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "ulasan berhasil ditambahkan"})
}
