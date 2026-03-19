package nfc

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// X-Idempotency-Key wajib untuk /nfc/transaksi
	r.Post("/transaksi", handleTransaksiNFC)

	// Kiosk — cek saldo tanpa PIN, IP whitelist diverifikasi middleware
	r.Get("/ceksaldo/{uid}", handleCekSaldoKiosk)
}

func handleTransaksiNFC(w http.ResponseWriter, r *http.Request) {
	// Alur: tap NFC → PIN 6 digit → debit rekening
	// Idempotency key wajib (X-Idempotency-Key header)
	// Validasi: kartu aktif, PIN benar, limit per transaksi & harian dari settings BMT
	response.Created(w, map[string]string{"message": "transaksi NFC berhasil"})
}

func handleCekSaldoKiosk(w http.ResponseWriter, r *http.Request) {
	// Kiosk: tampil nama + saldo + 5 transaksi terakhir
	// Tidak ada PIN. IP whitelist terminal kiosk.
	uid := chi.URLParam(r, "uid")
	response.Success(w, map[string]interface{}{
		"uid":        uid,
		"nama":       "",
		"saldo":      0,
		"transaksi":  []interface{}{},
	})
}
