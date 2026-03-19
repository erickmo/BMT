package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const idempotencyHeader = "X-Idempotency-Key"

type idempotencyKey struct{}

func Idempotency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(idempotencyHeader)
		if key != "" {
			parsed, err := uuid.Parse(key)
			if err == nil {
				ctx := context.WithValue(r.Context(), idempotencyKey{}, parsed)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func GetIdempotencyKey(ctx context.Context) (uuid.UUID, bool) {
	key, ok := ctx.Value(idempotencyKey{}).(uuid.UUID)
	return key, ok
}
