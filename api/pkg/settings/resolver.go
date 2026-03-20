package settings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Store interface {
	GetPlatform(ctx context.Context, kunci string) (string, error)
	GetBMT(ctx context.Context, bmtID uuid.UUID, kunci string) (string, error)
	GetCabang(ctx context.Context, cabangID uuid.UUID, kunci string) (string, error)
}

type Resolver struct {
	store Store
}

func NewResolver(store Store) *Resolver {
	return &Resolver{store: store}
}

// Resolve resolves setting with priority: cabang > bmt > platform
func (r *Resolver) Resolve(ctx context.Context, bmtID, cabangID uuid.UUID, kunci string) string {
	// Try cabang first
	if cabangID != uuid.Nil {
		if val, err := r.store.GetCabang(ctx, cabangID, kunci); err == nil && val != "" {
			return val
		}
	}

	// Try BMT
	if bmtID != uuid.Nil {
		if val, err := r.store.GetBMT(ctx, bmtID, kunci); err == nil && val != "" {
			return val
		}
	}

	// Fallback to platform
	if val, err := r.store.GetPlatform(ctx, kunci); err == nil {
		return val
	}

	return ""
}

// ResolveWithDefault resolves setting with a default fallback value
func (r *Resolver) ResolveWithDefault(ctx context.Context, bmtID, cabangID uuid.UUID, kunci, defaultVal string) string {
	val := r.Resolve(ctx, bmtID, cabangID, kunci)
	if val == "" {
		return defaultVal
	}
	return val
}

// ResolveInt resolves a setting as int
func (r *Resolver) ResolveInt(ctx context.Context, bmtID, cabangID uuid.UUID, kunci string, defaultVal int) int {
	val := r.Resolve(ctx, bmtID, cabangID, kunci)
	if val == "" {
		return defaultVal
	}
	var n int
	if _, err := fmt.Sscanf(val, "%d", &n); err != nil {
		return defaultVal
	}
	return n
}

// ResolveBool resolves a setting as bool
func (r *Resolver) ResolveBool(ctx context.Context, bmtID, cabangID uuid.UUID, kunci string, defaultVal bool) bool {
	val := r.Resolve(ctx, bmtID, cabangID, kunci)
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1" || val == "yes"
}

// GetApprovers mengembalikan daftar role approver untuk jenis form.
func (r *Resolver) GetApprovers(ctx context.Context, bmtID, cabangID uuid.UUID, jenisForm string) []string {
	var approvers []string
	key := "approval." + jenisForm
	if err := r.ResolveJSON(ctx, bmtID, cabangID, key, &approvers); err != nil {
		return []string{"MANAJER_CABANG"}
	}
	return approvers
}

// ResolveJSON resolves a setting as JSON into target
func (r *Resolver) ResolveJSON(ctx context.Context, bmtID, cabangID uuid.UUID, kunci string, target interface{}) error {
	val := r.Resolve(ctx, bmtID, cabangID, kunci)
	if val == "" {
		return fmt.Errorf("setting %s not found", kunci)
	}
	return json.Unmarshal([]byte(val), target)
}
