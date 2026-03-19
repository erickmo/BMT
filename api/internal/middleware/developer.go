package middleware

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
)

const developerTokenHeader = "Developer-Token"

func DeveloperAuth(devToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(developerTokenHeader)
			if token == "" || token != devToken {
				response.Unauthorized(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
