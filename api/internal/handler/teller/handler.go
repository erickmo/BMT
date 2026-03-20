package teller

import (
	"encoding/json"
	"net/http"

	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/internal/domain/sesi_teller"
	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	sesiSvc    *service.SesiTellerService
	rekeningService *service.RekeningService
	nasabahSvc *service.NasabahService
}

func NewHandler(
	sesiSvc *service.SesiTellerService,
	rekeningService *service.RekeningService,
	nasabahSvc *service.NasabahService,
) *Handler {
	return &Handler{
		sesiSvc:         sesiSvc,
		rekeningService: rekeningService,
		nasabahSvc:      nasabahSvc,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/sesi/buka", h.handleBukaSesi)
	r.Get("/sesi/aktif", h.handleGetSesiAktif)
	r.Post("/sesi/tutup", h.handleTutupSesi)
	r.Get("/nasabah/cari", h.handleCariNasabah)
	r.Post("/rekening/{id}/setor", h.handleSetor)
	r.Post("/rekening/{id}/tarik", h.handleTarik)
	r.Post("/pembiayaan/{id}/angsuran", h.handleBayarAngsuran)
	r.Post("/spp/{id}/bayar", h.handleBayarSPP)
}

// RegisterRoutes adalah fungsi kompatibilitas lama (diperlukan di main.go lama).
// Gunakan NewHandler(...).RegisterRoutes(r) untuk wiring lengkap.
func RegisterRoutes(r chi.Router) {
	h := &Handler{}
	h.RegisterRoutes(r)
}

// POST /teller/sesi/buka
func (h *Handler) handleBukaSesi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tellerID := middleware.GetUserID(ctx)
	bmtID := middleware.GetBMTID(ctx)
	cabangID := middleware.GetCabangID(ctx)

	var req struct {
		Redenominasi []sesi_teller.ItemPecahan `json:"redenominasi"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}

	sesi, err := h.sesiSvc.BukaSesi(ctx, service.BukaSesiInput{
		BMTID:        bmtID,
		CabangID:     cabangID,
		TellerID:     tellerID,
		Redenominasi: req.Redenominasi,
	})
	if err != nil {
		switch err {
		case sesi_teller.ErrSesiSudahAktif:
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Created(w, sesi)
}

// GET /teller/sesi/aktif
func (h *Handler) handleGetSesiAktif(w http.ResponseWriter, r *http.Request) {
	tellerID := middleware.GetUserID(r.Context())
	sesi, err := h.sesiSvc.GetSesiAktif(r.Context(), tellerID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, sesi)
}

// POST /teller/sesi/tutup
func (h *Handler) handleTutupSesi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tellerID := middleware.GetUserID(ctx)

	var req struct {
		SesiID           uuid.UUID                 `json:"sesi_id"`
		RedenominasiAkhir []sesi_teller.ItemPecahan `json:"redenominasi_akhir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}

	sesi, err := h.sesiSvc.TutupSesi(ctx, req.SesiID, tellerID, req.RedenominasiAkhir)
	if err != nil {
		switch err {
		case sesi_teller.ErrSesiSelisih:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		case sesi_teller.ErrSesiTidakAktif:
			response.Error(w, http.StatusNotFound, err.Error())
		case sesi_teller.ErrSesiSudahTutup:
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Success(w, sesi)
}

// GET /teller/nasabah/cari?q=...&page=1&per_page=20
func (h *Handler) handleCariNasabah(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bmtID := middleware.GetBMTID(ctx)
	q := r.URL.Query().Get("q")

	page := 1
	perPage := 20

	nasabahList, total, err := h.nasabahSvc.Search(ctx, bmtID, q, page, perPage)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.WithMeta(w, nasabahList, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
}

// POST /teller/rekening/{id}/setor
func (h *Handler) handleSetor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tellerID := middleware.GetUserID(ctx)

	rekeningID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "rekening_id tidak valid")
		return
	}

	// Pastikan sesi teller aktif
	if _, err := h.sesiSvc.GetSesiAktif(ctx, tellerID); err != nil {
		response.Error(w, http.StatusForbidden, "tidak ada sesi teller aktif")
		return
	}

	var req struct {
		Nominal        int64      `json:"nominal"`
		Keterangan     string     `json:"keterangan"`
		IdempotencyKey *uuid.UUID `json:"idempotency_key,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}
	if req.Nominal <= 0 {
		response.BadRequest(w, "nominal harus lebih dari 0")
		return
	}

	// Gunakan idempotency key dari header jika ada
	idempKey, hasKey := middleware.GetIdempotencyKey(ctx)
	if hasKey && req.IdempotencyKey == nil {
		req.IdempotencyKey = &idempKey
	}

	tr, err := h.rekeningService.Setor(ctx, rekening.SetoranInput{
		RekeningID:     rekeningID,
		Nominal:        req.Nominal,
		Keterangan:     req.Keterangan,
		IdempotencyKey: req.IdempotencyKey,
		CreatedBy:      tellerID,
	})
	if err != nil {
		switch err {
		case rekening.ErrRekeningBeku:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		case rekening.ErrRekeningTutup:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		case rekening.ErrSetoranDibawahMin:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Created(w, tr)
}

// POST /teller/rekening/{id}/tarik
func (h *Handler) handleTarik(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tellerID := middleware.GetUserID(ctx)

	rekeningID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "rekening_id tidak valid")
		return
	}

	// Pastikan sesi teller aktif
	if _, err := h.sesiSvc.GetSesiAktif(ctx, tellerID); err != nil {
		response.Error(w, http.StatusForbidden, "tidak ada sesi teller aktif")
		return
	}

	var req struct {
		Nominal        int64      `json:"nominal"`
		Keterangan     string     `json:"keterangan"`
		IdempotencyKey *uuid.UUID `json:"idempotency_key,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}
	if req.Nominal <= 0 {
		response.BadRequest(w, "nominal harus lebih dari 0")
		return
	}

	idempKey, hasKey := middleware.GetIdempotencyKey(ctx)
	if hasKey && req.IdempotencyKey == nil {
		req.IdempotencyKey = &idempKey
	}

	tr, err := h.rekeningService.Tarik(ctx, rekening.PenarikanInput{
		RekeningID:     rekeningID,
		Nominal:        req.Nominal,
		Keterangan:     req.Keterangan,
		IdempotencyKey: req.IdempotencyKey,
		CreatedBy:      tellerID,
	})
	if err != nil {
		switch err {
		case rekening.ErrSaldoTidakCukup:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		case rekening.ErrRekeningBeku:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		case rekening.ErrRekeningTutup:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		case rekening.ErrPenarikanTidakBisa:
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Created(w, tr)
}

func (h *Handler) handleBayarAngsuran(w http.ResponseWriter, r *http.Request) {
	// Sprint 5: pembiayaan service
	response.Error(w, http.StatusNotImplemented, "coming in sprint 5")
}

func (h *Handler) handleBayarSPP(w http.ResponseWriter, r *http.Request) {
	// Sprint 6: pondok service
	response.Error(w, http.StatusNotImplemented, "coming in sprint 6")
}
