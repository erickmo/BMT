package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/bmt-saas/api/pkg/response"
)

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup goroutine
	go func() {
		ticker := time.NewTicker(window)
		for range ticker.C {
			rl.mu.Lock()
			now := time.Now()
			for ip, times := range rl.requests {
				var valid []time.Time
				for _, t := range times {
					if now.Sub(t) <= window {
						valid = append(valid, t)
					}
				}
				if len(valid) == 0 {
					delete(rl.requests, ip)
				} else {
					rl.requests[ip] = valid
				}
			}
			rl.mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			rl.mu.Lock()
			now := time.Now()
			times := rl.requests[ip]
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) <= window {
					valid = append(valid, t)
				}
			}

			if len(valid) >= limit {
				rl.mu.Unlock()
				response.Error(w, http.StatusTooManyRequests, "too many requests")
				return
			}

			rl.requests[ip] = append(valid, now)
			rl.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
