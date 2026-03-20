package finance

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domainfinance "github.com/bmt-saas/api/internal/domain/finance"
	"github.com/bmt-saas/api/internal/domain/pembiayaan"
	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler menangani request HTTP untuk domain finance.
type Handler struct {
	financeSvc    *service.FinanceService
	pembiayaanSvc *service.PembiayaanService
}

// NewHandler membuat instance baru Handler.
func NewHandler(financeSvc *service.FinanceService, pembiayaanSvc *service.PembiayaanService) *Handler {
	return &Handler{financeSvc: financeSvc, pembiayaanSvc: pembiayaanSvc}
}

// RegisterRoutes mendaftarkan semua route finance.
func (h *Handler) RegisterRoutes(r chi.Router) {
	// Jurnal manual (double-entry via module-vernon-accounting)
	r.Get("/jurnal", h.listJurnal)
	r.Post("/jurnal", h.createJurnal)
	r.Get("/jurnal/{id}", h.getJurnal)
	r.Post("/jurnal/{id}/post", h.postJurnal)

	// Transaksi operasional
	r.Get("/transaksi", h.listTransaksi)
	r.Post("/transaksi", h.createTransaksi)

	// Pembiayaan
	r.Get("/pembiayaan", h.listPembiayaan)
	r.Post("/pembiayaan", h.ajukanPembiayaan)
	r.Get("/pembiayaan/{id}", h.getPembiayaan)
	r.Post("/pembiayaan/{id}/status", h.updateStatusPembiayaan)
	r.Post("/pembiayaan/{id}/angsuran", h.bayarAngsuran)
	r.Get("/pembiayaan/{id}/jadwal", h.getJadwalAngsuran)

	// Laporan keuangan
	r.Get("/laporan/neraca", h.laporanNeraca)
	r.Get("/laporan/shu", h.laporanSHU)
	r.Get("/laporan/arus-kas", h.laporanArusKas)
	r.Get("/laporan/kolektibilitas", h.laporanKolektibilitas)
	r.Get("/laporan/bagi-hasil-deposito", h.laporanBagiHasilDeposito)
}

// parseDateRange mengambil query param "dari" dan "sampai" (format YYYY-MM-DD).
// Default: dari = 1 bulan lalu, sampai = hari ini.
func parseDateRange(r *http.Request) (dari, sampai time.Time) {
	dari = time.Now().AddDate(0, -1, 0)
	sampai = time.Now()
	if d, err := time.Parse("2006-01-02", r.URL.Query().Get("dari")); err == nil {
		dari = d
	}
	if s, err := time.Parse("2006-01-02", r.URL.Query().Get("sampai")); err == nil {
		sampai = s
	}
	return
}

// ─── Jurnal ───────────────────────────────────────────────────────────────────

func (h *Handler) listJurnal(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())

	filter := domainfinance.ListJurnalFilter{
		BMTID:   &bmtID,
		Page:    1,
		PerPage: 50,
	}
	if s := r.URL.Query().Get("status"); s != "" {
		st := domainfinance.StatusJurnal(s)
		filter.Status = &st
	}
	if dari, err := time.Parse("2006-01-02", r.URL.Query().Get("dari")); err == nil {
		filter.TanggalDari = &dari
	}
	if sampai, err := time.Parse("2006-01-02", r.URL.Query().Get("sampai")); err == nil {
		filter.TanggalSampai = &sampai
	}

	list, total, err := h.financeSvc.ListJurnal(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.WithMeta(w, list, &response.Meta{Total: total, Page: filter.Page, PerPage: filter.PerPage})
}

type createJurnalReq struct {
	Tanggal    string `json:"tanggal"`
	Keterangan string `json:"keterangan"`
	Referensi  string `json:"referensi,omitempty"`
	Entries    []struct {
		KodeAkun string `json:"kode_akun"`
		NamaAkun string `json:"nama_akun"`
		Posisi   string `json:"posisi"`
		Nominal  int64  `json:"nominal"`
	} `json:"entries"`
}

func (h *Handler) createJurnal(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	cabangID := middleware.GetCabangID(r.Context())
	userID := middleware.GetUserID(r.Context())

	var req createJurnalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "request tidak valid")
		return
	}

	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		response.BadRequest(w, "format tanggal tidak valid, gunakan YYYY-MM-DD")
		return
	}

	entries := make([]domainfinance.EntriInput, 0, len(req.Entries))
	for _, e := range req.Entries {
		entries = append(entries, domainfinance.EntriInput{
			KodeAkun: e.KodeAkun,
			NamaAkun: e.NamaAkun,
			Posisi:   domainfinance.PosisiJurnal(e.Posisi),
			Nominal:  e.Nominal,
		})
	}

	j, err := h.financeSvc.CreateJurnal(r.Context(), domainfinance.CreateJurnalInput{
		BMTID:      bmtID,
		CabangID:   cabangID,
		Tanggal:    tanggal,
		Keterangan: req.Keterangan,
		Referensi:  req.Referensi,
		Entries:    entries,
		DibuatOleh: userID,
	})
	if err != nil {
		if errors.Is(err, domainfinance.ErrJurnalTidakBalance) {
			response.BadRequest(w, "jurnal tidak balance: total debit ≠ total kredit")
			return
		}
		response.InternalError(w)
		return
	}
	response.Created(w, j)
}

func (h *Handler) getJurnal(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id tidak valid")
		return
	}

	j, err := h.financeSvc.GetJurnal(r.Context(), id, bmtID)
	if err != nil {
		if errors.Is(err, domainfinance.ErrJurnalNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
	response.Success(w, j)
}

func (h *Handler) postJurnal(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	userID := middleware.GetUserID(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id tidak valid")
		return
	}

	j, err := h.financeSvc.PostJurnal(r.Context(), id, bmtID, userID)
	if err != nil {
		if errors.Is(err, domainfinance.ErrJurnalNotFound) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, j)
}

// ─── Transaksi Operasional ───────────────────────────────────────────────────

func (h *Handler) listTransaksi(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	cabangID := middleware.GetCabangID(r.Context())
	dari, sampai := parseDateRange(r)

	list, total, err := h.financeSvc.ListTransaksiOperasional(r.Context(), bmtID, cabangID, dari, sampai, 1, 100)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.WithMeta(w, list, &response.Meta{Total: total, Page: 1, PerPage: 100})
}

type createTransaksiReq struct {
	VendorID       *uuid.UUID `json:"vendor_id,omitempty"`
	Tanggal        string     `json:"tanggal"`
	Jenis          string     `json:"jenis"`
	Kategori       string     `json:"kategori"`
	Keterangan     string     `json:"keterangan"`
	Nominal        int64      `json:"nominal"`
	KodeAkunDebit  string     `json:"kode_akun_debit"`
	KodeAkunKredit string     `json:"kode_akun_kredit"`
}

func (h *Handler) createTransaksi(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	cabangID := middleware.GetCabangID(r.Context())
	userID := middleware.GetUserID(r.Context())

	var req createTransaksiReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "request tidak valid")
		return
	}

	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		response.BadRequest(w, "format tanggal tidak valid, gunakan YYYY-MM-DD")
		return
	}

	t, err := h.financeSvc.CreateTransaksiOperasional(r.Context(), &domainfinance.TransaksiOperasional{
		BMTID:          bmtID,
		CabangID:       cabangID,
		VendorID:       req.VendorID,
		Tanggal:        tanggal,
		Jenis:          req.Jenis,
		Kategori:       req.Kategori,
		Keterangan:     req.Keterangan,
		Nominal:        req.Nominal,
		KodeAkunDebit:  req.KodeAkunDebit,
		KodeAkunKredit: req.KodeAkunKredit,
		DibuatOleh:     userID,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Created(w, t)
}

// ─── Pembiayaan ───────────────────────────────────────────────────────────────

func (h *Handler) listPembiayaan(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	cabangID := middleware.GetCabangID(r.Context())

	list, total, err := h.pembiayaanSvc.ListPembiayaan(r.Context(), pembiayaan.ListPembiayaanFilter{
		BMTID:    &bmtID,
		CabangID: &cabangID,
		Page:     1,
		PerPage:  50,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.WithMeta(w, list, &response.Meta{Total: total, Page: 1, PerPage: 50})
}

func (h *Handler) ajukanPembiayaan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "gunakan endpoint /api/form/pembiayaan untuk pengajuan"})
}

func (h *Handler) getPembiayaan(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id tidak valid")
		return
	}

	p, err := h.pembiayaanSvc.GetByID(r.Context(), id, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, p)
}

func (h *Handler) updateStatusPembiayaan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id tidak valid")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "request tidak valid")
		return
	}

	p, err := h.pembiayaanSvc.MajukanStatus(r.Context(), id, pembiayaan.StatusPembiayaan(req.Status), userID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, p)
}

func (h *Handler) bayarAngsuran(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	pembiayaanID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id tidak valid")
		return
	}

	var req struct {
		RekeningID uuid.UUID `json:"rekening_id"`
		Nominal    int64     `json:"nominal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "request tidak valid")
		return
	}

	angsuran, err := h.pembiayaanSvc.BayarAngsuran(r.Context(), pembiayaanID, req.RekeningID, req.Nominal, userID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Success(w, angsuran)
}

func (h *Handler) getJadwalAngsuran(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "id tidak valid")
		return
	}

	jadwal, err := h.pembiayaanSvc.GetJadwalAngsuran(r.Context(), id, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, jadwal)
}

// ─── Laporan ─────────────────────────────────────────────────────────────────

func (h *Handler) laporanNeraca(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	dari, sampai := parseDateRange(r)

	laporan, err := h.financeSvc.GetLaporanNeraca(r.Context(), bmtID, dari, sampai)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Success(w, laporan)
}

func (h *Handler) laporanSHU(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	dari, sampai := parseDateRange(r)

	laporan, err := h.financeSvc.GetLaporanSHU(r.Context(), bmtID, dari, sampai)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Success(w, laporan)
}

func (h *Handler) laporanArusKas(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	dari, sampai := parseDateRange(r)

	laporan, err := h.financeSvc.GetLaporanArusKas(r.Context(), bmtID, dari, sampai)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Success(w, laporan)
}

func (h *Handler) laporanKolektibilitas(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())

	laporan, err := h.financeSvc.GetLaporanKolektibilitas(r.Context(), bmtID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Success(w, laporan)
}

func (h *Handler) laporanBagiHasilDeposito(w http.ResponseWriter, r *http.Request) {
	bmtID := middleware.GetBMTID(r.Context())
	bulan := time.Now()
	if b, err := time.Parse("2006-01", r.URL.Query().Get("bulan")); err == nil {
		bulan = b
	}

	laporan, err := h.financeSvc.GetLaporanBagiHasilDeposito(r.Context(), bmtID, bulan)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Success(w, laporan)
}
