package form

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bmt-saas/api/internal/domain/form"
	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	formSvc *service.FormService
}

func NewHandler(formSvc *service.FormService) *Handler {
	return &Handler{formSvc: formSvc}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.handleBuatForm)
	r.Get("/", h.handleListForm)
	r.Get("/{id}", h.handleGetForm)
	r.Post("/{id}/ajukan", h.handleAjukanForm)
	r.Post("/{id}/approval", h.handleProsesApproval)
	r.Get("/{id}/approval", h.handleGetApprovals)
}

// POST /api/form
func (h *Handler) handleBuatForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	bmtID := middleware.GetBMTID(ctx)
	cabangID := middleware.GetCabangID(ctx)

	var req struct {
		JenisForm form.JenisForm         `json:"jenis_form"`
		DataForm  map[string]interface{} `json:"data_form"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}
	if req.JenisForm == "" {
		response.BadRequest(w, "jenis_form wajib diisi")
		return
	}

	f, err := h.formSvc.BuatForm(ctx, service.BuatFormInput{
		BMTID:     bmtID,
		CabangID:  cabangID,
		JenisForm: req.JenisForm,
		DataForm:  req.DataForm,
		CreatedBy: userID,
	})
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Created(w, f)
}

// GET /api/form?status=DIAJUKAN&page=1&per_page=20
func (h *Handler) handleListForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bmtID := middleware.GetBMTID(ctx)
	cabangID := middleware.GetCabangID(ctx)

	status := form.StatusForm(r.URL.Query().Get("status"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	list, total, err := h.formSvc.ListForm(ctx, bmtID, cabangID, status, page, perPage)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.WithMeta(w, list, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
}

// GET /api/form/{id}
func (h *Handler) handleGetForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bmtID := middleware.GetBMTID(ctx)

	formID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "form_id tidak valid")
		return
	}

	f, err := h.formSvc.GetForm(ctx, formID, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, f)
}

// POST /api/form/{id}/ajukan
func (h *Handler) handleAjukanForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	bmtID := middleware.GetBMTID(ctx)

	formID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "form_id tidak valid")
		return
	}

	// Pastikan form milik BMT ini
	f, err := h.formSvc.GetForm(ctx, formID, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}

	f, err = h.formSvc.AjukanForm(ctx, f.ID, userID)
	if err != nil {
		switch err {
		case form.ErrFormTidakBisaDiubah:
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Success(w, f)
}

// POST /api/form/{id}/approval
func (h *Handler) handleProsesApproval(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	approverID := middleware.GetUserID(ctx)
	role := middleware.GetRole(ctx)
	bmtID := middleware.GetBMTID(ctx)

	formID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "form_id tidak valid")
		return
	}

	var req struct {
		Setujui bool   `json:"setujui"`
		Catatan string `json:"catatan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid")
		return
	}

	// Pastikan form milik BMT ini
	if _, err := h.formSvc.GetForm(ctx, formID, bmtID); err != nil {
		response.NotFound(w)
		return
	}

	f, err := h.formSvc.ProsesApproval(ctx, form.ApprovalInput{
		FormID:       formID,
		ApproverID:   approverID,
		RoleApprover: role,
		Catatan:      req.Catatan,
		Setujui:      req.Setujui,
	})
	if err != nil {
		switch err {
		case form.ErrApproverTidakBerwenang:
			response.Forbidden(w)
		case form.ErrFormTidakBisaDiubah:
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.InternalError(w)
		}
		return
	}
	response.Success(w, f)
}

// GET /api/form/{id}/approval
func (h *Handler) handleGetApprovals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bmtID := middleware.GetBMTID(ctx)

	formID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "form_id tidak valid")
		return
	}

	approvals, err := h.formSvc.GetApprovals(ctx, formID, bmtID)
	if err != nil {
		response.NotFound(w)
		return
	}
	response.Success(w, approvals)
}
