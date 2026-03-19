package middleware

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/google/uuid"
)

// TenantRequired ensures BMTID is present in context
func TenantRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bmtID := GetBMTID(r.Context())
		if bmtID == uuid.Nil {
			response.Forbidden(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
