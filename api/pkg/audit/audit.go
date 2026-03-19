package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ActionType string

const (
	ActionCreate  ActionType = "CREATE"
	ActionRead    ActionType = "READ"
	ActionUpdate  ActionType = "UPDATE"
	ActionDelete  ActionType = "DELETE"
	ActionLogin   ActionType = "LOGIN"
	ActionLogout  ActionType = "LOGOUT"
	ActionApprove ActionType = "APPROVE"
	ActionReject  ActionType = "REJECT"
)

type AktorType string

const (
	AktorStaf           AktorType = "STAF"
	AktorNasabah        AktorType = "NASABAH"
	AktorSistem         AktorType = "SISTEM"
	AktorDeveloper      AktorType = "DEVELOPER"
	AktorPenggunaPondok AktorType = "PENGGUNA_PONDOK"
)

type LogEntry struct {
	ID           uuid.UUID   `json:"id"`
	BMTID        *uuid.UUID  `json:"bmt_id,omitempty"`
	SubjekID     uuid.UUID   `json:"subjek_id"`
	SubjekTipe   AktorType   `json:"subjek_tipe"`
	Aksi         ActionType  `json:"aksi"`
	ResourceTipe string      `json:"resource_tipe"`
	ResourceID   *uuid.UUID  `json:"resource_id,omitempty"`
	DataSebelum  interface{} `json:"data_sebelum,omitempty"`
	DataSesudah  interface{} `json:"data_sesudah,omitempty"`
	IPAddress    string      `json:"ip_address"`
	UserAgent    string      `json:"user_agent"`
	CreatedAt    time.Time   `json:"created_at"`
}

type Logger interface {
	Log(ctx context.Context, entry LogEntry) error
}

type contextKey string

const auditKey contextKey = "audit_context"

type AuditContext struct {
	SubjekID   uuid.UUID
	SubjekTipe AktorType
	BMTID      *uuid.UUID
	IPAddress  string
	UserAgent  string
}

func WithContext(ctx context.Context, ac AuditContext) context.Context {
	return context.WithValue(ctx, auditKey, ac)
}

func FromContext(ctx context.Context) (AuditContext, bool) {
	ac, ok := ctx.Value(auditKey).(AuditContext)
	return ac, ok
}
