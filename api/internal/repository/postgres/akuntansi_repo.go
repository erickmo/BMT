package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AkuntansiRepository struct {
	db *pgxpool.Pool
}

func NewAkuntansiRepository(db *pgxpool.Pool) *AkuntansiRepository {
	return &AkuntansiRepository{db: db}
}

func (r *AkuntansiRepository) CreateJurnal(ctx context.Context, j *service.Jurnal) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO jurnal_akuntansi (id, bmt_id, cabang_id, nomor, keterangan, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, j.ID, j.BMTID, j.CabangID, j.Nomor, j.Keterangan, time.Now())
	return err
}

func (r *AkuntansiRepository) CreateEntries(ctx context.Context, entries []*service.EntryJurnal) error {
	for _, e := range entries {
		_, err := r.db.Exec(ctx, `
			INSERT INTO jurnal_entry (id, jurnal_id, kode_akun, posisi, nominal)
			VALUES ($1,$2,$3,$4,$5)
		`, uuid.New(), e.JurnalID, e.KodeAkun, e.Posisi, e.Nominal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AkuntansiRepository) GenerateNomorJurnal(ctx context.Context, bmtID uuid.UUID) (string, error) {
	var seq int
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(COUNT(*), 0) + 1 FROM jurnal_akuntansi WHERE bmt_id = $1
	`, bmtID).Scan(&seq)
	if err != nil {
		return "", fmt.Errorf("gagal generate nomor jurnal: %w", err)
	}
	now := time.Now()
	return fmt.Sprintf("JNL-%d%02d-%06d", now.Year(), now.Month(), seq), nil
}
