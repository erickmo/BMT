package nasabah

import (
	"net/http"
	"strconv"

	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	nasabahSvc *service.NasabahService
	rekeningService *service.RekeningService
}

func NewHandler(nasabahSvc *service.NasabahService, rekeningService *service.RekeningService) *Handler {
	return &Handler{nasabahSvc: nasabahSvc, rekeningService: rekeningService}
}

// RegisterRoutes kompatibilitas lama.
func RegisterRoutes(r chi.Router) {
	h := &Handler{}
	h.RegisterRoutes(r)
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/profil", h.handleGetProfil)
	r.Get("/rekening", h.handleListRekening)
	r.Get("/rekening/{id}/transaksi", h.handleListTransaksi)
	r.Post("/rekening/{id}/setor-online", h.handleSetorOnline)
	r.Get("/pembiayaan", h.handleListPembiayaan)
	r.Get("/santri", h.handleGetSantri)
	r.Get("/nfc/saldo", h.handleGetNFCSaldo)
	r.Get("/nfc/transaksi", h.handleListNFCTransaksi)
	r.Post("/nfc/topup", h.handleTopupNFC)
	r.Get("/spp/tagihan", h.handleListTagihanSPP)
	r.Post("/spp/{id}/bayar", h.handleBayarSPP)

	// E-commerce
	r.Get("/shop/toko", h.handleListToko)
	r.Get("/shop/toko/{slug}/produk", h.handleListProdukToko)
	r.Post("/shop/keranjang", h.handleAddKeranjang)
	r.Post("/shop/pesanan", h.handleCreatePesanan)
	r.Get("/shop/pesanan", h.handleListPesanan)
	r.Get("/shop/pesanan/{id}", h.handleGetPesanan)
	r.Post("/shop/ulasan", h.handleCreateUlasan)
}

// GET /nasabah/profil
func (h *Handler) handleGetProfil(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nasabahID := middleware.GetUserID(ctx)
	bmtID := middleware.GetBMTID(ctx)

	n, err := h.nasabahSvc.GetByID(ctx, nasabahID, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, n)
}

// GET /nasabah/rekening
func (h *Handler) handleListRekening(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nasabahID := middleware.GetUserID(ctx)
	bmtID := middleware.GetBMTID(ctx)

	rekeningList, err := h.nasabahSvc.ListRekening(ctx, nasabahID, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, rekeningList)
}

// GET /nasabah/rekening/{id}/transaksi?limit=50&offset=0
func (h *Handler) handleListTransaksi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bmtID := middleware.GetBMTID(ctx)

	rekeningID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "rekening_id tidak valid")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 50
	}

	transaksi, total, err := h.nasabahSvc.GetMutasi(ctx, rekeningID, bmtID, limit, offset)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.WithMeta(w, transaksi, &response.Meta{
		Total:   total,
		PerPage: limit,
	})
}

func (h *Handler) handleSetorOnline(w http.ResponseWriter, r *http.Request) {
	// Sprint 4: Midtrans integration
	response.Error(w, http.StatusNotImplemented, "coming in sprint 4")
}

func (h *Handler) handleListPembiayaan(w http.ResponseWriter, r *http.Request) {
	// Sprint 5: pembiayaan service
	response.Error(w, http.StatusNotImplemented, "coming in sprint 5")
}

func (h *Handler) handleGetSantri(w http.ResponseWriter, r *http.Request) {
	// Sprint 6: pondok service
	response.Error(w, http.StatusNotImplemented, "coming in sprint 6")
}

func (h *Handler) handleGetNFCSaldo(w http.ResponseWriter, r *http.Request) {
	// Sprint 7: NFC service
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleListNFCTransaksi(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleTopupNFC(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleListTagihanSPP(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 6")
}

func (h *Handler) handleBayarSPP(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 6")
}

func (h *Handler) handleListToko(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleListProdukToko(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleAddKeranjang(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleCreatePesanan(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleListPesanan(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleGetPesanan(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}

func (h *Handler) handleCreateUlasan(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "coming in sprint 7")
}
