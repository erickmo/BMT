package webhook

import (
	"io"
	"log"
	"net/http"

	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

// Handler menangani webhook dari provider eksternal.
type Handler struct {
	midtransSvc *service.MidtransService
}

// NewHandler membuat instance Handler baru.
func NewHandler(midtransSvc *service.MidtransService) *Handler {
	return &Handler{midtransSvc: midtransSvc}
}

// RegisterRoutes mendaftarkan semua route webhook ke router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/midtrans", h.handleMidtransWebhook)
}

// POST /webhook/midtrans
// Midtrans akan retry jika menerima non-200, sehingga kita selalu return 200.
func (h *Handler) handleMidtransWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[webhook] gagal baca body Midtrans: %v", err)
		response.Success(w, map[string]string{"status": "ok"})
		return
	}

	if err := h.midtransSvc.HandleWebhook(r.Context(), body); err != nil {
		log.Printf("[webhook] gagal proses webhook Midtrans: %v", err)
	}

	response.Success(w, map[string]string{"status": "ok"})
}
