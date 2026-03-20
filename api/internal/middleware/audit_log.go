package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuditLogger menulis log audit ke persistent store.
type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
}

// AuditEntry adalah satu record audit.
type AuditEntry struct {
	ID          uuid.UUID
	BMTID       uuid.UUID
	CabangID    uuid.UUID
	SubjekTipe  string
	SubjekID    uuid.UUID
	Aksi        string
	Entitas     string
	EntitasID   string
	DataRequest *string
	IP          string
	UserAgent   string
	CreatedAt   time.Time
}

// AuditLog middleware mencatat semua mutasi (POST/PUT/DELETE/PATCH) ke audit_log.
func AuditLog(logger AuditLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Hanya catat mutasi
			if r.Method == http.MethodGet || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			userID := GetUserID(ctx)
			bmtID := GetBMTID(ctx)
			cabangID := GetCabangID(ctx)
			role := GetRole(ctx)

			// Baca body untuk disimpan (max 8KB)
			var dataReq *string
			if r.Body != nil && r.ContentLength > 0 && r.ContentLength < 8192 {
				body, err := io.ReadAll(io.LimitReader(r.Body, 8192))
				if err == nil && len(body) > 0 {
					// Restore body untuk handler berikutnya
					r.Body = io.NopCloser(bytes.NewBuffer(body))
					// Mask field sensitif sebelum simpan
					masked := maskSensitiveFields(body)
					s := string(masked)
					dataReq = &s
				}
			}

			// Tentukan aksi dari method + path
			aksi := httpMethodToAksi(r.Method)
			entitas, entitasID := pathToEntitas(r.URL.Path)

			// Tentukan tipe subjek dari role
			subjekTipe := roleToSubjekTipe(role)

			entry := AuditEntry{
				ID:          uuid.New(),
				BMTID:       bmtID,
				CabangID:    cabangID,
				SubjekTipe:  subjekTipe,
				SubjekID:    userID,
				Aksi:        aksi,
				Entitas:     entitas,
				EntitasID:   entitasID,
				DataRequest: dataReq,
				IP:          r.RemoteAddr,
				UserAgent:   r.Header.Get("User-Agent"),
				CreatedAt:   time.Now(),
			}

			next.ServeHTTP(w, r)

			// Log setelah handler selesai (fire-and-forget, jangan block response)
			go func() {
				_ = logger.Log(context.Background(), entry)
			}()
		})
	}
}

func httpMethodToAksi(method string) string {
	switch method {
	case http.MethodPost:
		return "CREATE"
	case http.MethodPut, http.MethodPatch:
		return "UPDATE"
	case http.MethodDelete:
		return "DELETE"
	default:
		return "OTHER"
	}
}

func pathToEntitas(path string) (entitas, entitasID string) {
	// Contoh: /teller/rekening/uuid-xxx/setor → "rekening", "uuid-xxx"
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return "unknown", ""
	}
	for i, p := range parts {
		if _, err := uuid.Parse(p); err == nil {
			if i > 0 {
				entitas = parts[i-1]
			}
			entitasID = p
			return
		}
	}
	entitas = parts[len(parts)-1]
	return
}

func roleToSubjekTipe(role string) string {
	switch role {
	case "NASABAH", "ALUMNI":
		return "NASABAH"
	case "DEVELOPER":
		return "DEVELOPER"
	case "":
		return "SISTEM"
	default:
		return "STAF"
	}
}

// maskSensitiveFields menyembunyikan field sensitif dari JSON body.
func maskSensitiveFields(body []byte) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return body
	}
	sensitif := []string{"password", "pin", "password_hash", "pin_hash", "token", "secret"}
	for _, k := range sensitif {
		if _, ok := m[k]; ok {
			m[k] = "***"
		}
	}
	out, err := json.Marshal(m)
	if err != nil {
		return body
	}
	return out
}
