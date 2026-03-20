package middleware

import (
	"context"
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/google/uuid"
)

// FeatureChecker memeriksa apakah fitur aktif untuk BMT tertentu.
type FeatureChecker interface {
	FiturAktif(ctx context.Context, bmtID uuid.UUID, kode string) error
}

// RequireFeature middleware memblokir akses jika fitur tidak aktif di kontrak BMT.
func RequireFeature(checker FeatureChecker, kode string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bmtID := GetBMTID(r.Context())
			if bmtID == uuid.Nil {
				response.Forbidden(w)
				return
			}
			if err := checker.FiturAktif(r.Context(), bmtID, kode); err != nil {
				response.Error(w, http.StatusPaymentRequired, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
