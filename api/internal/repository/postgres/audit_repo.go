package postgres

import (
	"context"

	"github.com/bmt-saas/api/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Log(ctx context.Context, e middleware.AuditEntry) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO audit_log (id, bmt_id, cabang_id, subjek_tipe, subjek_id, aksi,
		entitas, entitas_id, data_request, ip, user_agent, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, e.ID, e.BMTID, e.CabangID, e.SubjekTipe, e.SubjekID, e.Aksi,
		e.Entitas, e.EntitasID, e.DataRequest, e.IP, e.UserAgent, e.CreatedAt)
	return err
}
