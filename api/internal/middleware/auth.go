package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bmt-saas/api/pkg/jwt"
	"github.com/bmt-saas/api/pkg/response"
	"github.com/google/uuid"
)

type contextKey string

const (
	CtxUserID   contextKey = "user_id"
	CtxBMTID    contextKey = "bmt_id"
	CtxCabangID contextKey = "cabang_id"
	CtxRole     contextKey = "role"
)

func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Unauthorized(w)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w)
				return
			}

			claims, err := jwtManager.VerifyAccessToken(parts[1])
			if err != nil {
				response.Unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxBMTID, claims.BMTID)
			ctx = context.WithValue(ctx, CtxCabangID, claims.CabangID)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(CtxRole).(string)
			if !ok || !allowed[role] {
				response.Forbidden(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(CtxUserID).(uuid.UUID)
	return id
}

func GetBMTID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(CtxBMTID).(uuid.UUID)
	return id
}

func GetCabangID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(CtxCabangID).(uuid.UUID)
	return id
}

func GetRole(ctx context.Context) string {
	role, _ := ctx.Value(CtxRole).(string)
	return role
}
