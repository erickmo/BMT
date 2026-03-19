package finance

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// Jurnal manual (double-entry via module-vernon-accounting)
	r.Get("/jurnal", handleListJurnal)
	r.Post("/jurnal", handleCreateJurnal)
	r.Get("/jurnal/{id}", handleGetJurnal)

	// Transaksi operasional
	r.Get("/transaksi", handleListTransaksi)
	r.Post("/transaksi", handleCreateTransaksi)

	// Laporan keuangan
	r.Get("/laporan/neraca", handleLaporanNeraca)
	r.Get("/laporan/shu", handleLaporanSHU)
	r.Get("/laporan/arus-kas", handleLaporanArusKas)
	r.Get("/laporan/kolektibilitas", handleLaporanKolektibilitas)
	r.Get("/laporan/bagi-hasil-deposito", handleLaporanBagiHasilDeposito)
}

func handleListJurnal(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJurnal(w http.ResponseWriter, r *http.Request) {
	// Jurnal wajib balance: total debit == total kredit
	response.Created(w, map[string]string{"message": "jurnal berhasil diposting"})
}

func handleGetJurnal(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleListTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "transaksi berhasil dibuat"})
}

func handleLaporanNeraca(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleLaporanSHU(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleLaporanArusKas(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleLaporanKolektibilitas(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleLaporanBagiHasilDeposito(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}
