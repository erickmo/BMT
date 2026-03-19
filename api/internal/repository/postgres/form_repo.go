package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/form"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FormRepository struct {
	db *pgxpool.Pool
}

func NewFormRepository(db *pgxpool.Pool) *FormRepository {
	return &FormRepository{db: db}
}

func (r *FormRepository) Create(ctx context.Context, f *form.FormPengajuan) error {
	data, _ := json.Marshal(f.DataForm)
	_, err := r.db.Exec(ctx, `
		INSERT INTO form_pengajuan (id, bmt_id, cabang_id, jenis_form, nomor_form, status, data_form, catatan, diajukan_oleh, diajukan_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, f.ID, f.BMTID, f.CabangID, f.JenisForm, f.NomorForm, f.Status, data,
		f.Catatan, f.DiajukanOleh, f.DiajukanAt, f.CreatedAt, f.UpdatedAt)
	return err
}

func (r *FormRepository) GetByID(ctx context.Context, id uuid.UUID) (*form.FormPengajuan, error) {
	f := &form.FormPengajuan{}
	var data []byte
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, jenis_form, nomor_form, status, data_form, catatan, diajukan_oleh, diajukan_at, created_at, updated_at
		FROM form_pengajuan WHERE id = $1
	`, id).Scan(&f.ID, &f.BMTID, &f.CabangID, &f.JenisForm, &f.NomorForm, &f.Status,
		&data, &f.Catatan, &f.DiajukanOleh, &f.DiajukanAt, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, form.ErrFormNotFound
		}
		return nil, err
	}
	f.DataForm = make(map[string]interface{})
	json.Unmarshal(data, &f.DataForm)
	return f, nil
}

func (r *FormRepository) ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, status form.StatusForm, page, perPage int) ([]*form.FormPengajuan, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, jenis_form, nomor_form, status, data_form, catatan, diajukan_oleh, diajukan_at, created_at, updated_at
		FROM form_pengajuan
		WHERE bmt_id = $1 AND cabang_id = $2 AND ($3 = '' OR status = $3)
		ORDER BY created_at DESC LIMIT $4 OFFSET $5
	`, bmtID, cabangID, status, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*form.FormPengajuan
	for rows.Next() {
		f := &form.FormPengajuan{}
		var data []byte
		err := rows.Scan(&f.ID, &f.BMTID, &f.CabangID, &f.JenisForm, &f.NomorForm, &f.Status,
			&data, &f.Catatan, &f.DiajukanOleh, &f.DiajukanAt, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		f.DataForm = make(map[string]interface{})
		json.Unmarshal(data, &f.DataForm)
		result = append(result, f)
	}

	var total int64
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM form_pengajuan WHERE bmt_id = $1 AND cabang_id = $2 AND ($3 = '' OR status = $3)`,
		bmtID, cabangID, status).Scan(&total)

	return result, total, nil
}

func (r *FormRepository) Update(ctx context.Context, f *form.FormPengajuan) error {
	data, _ := json.Marshal(f.DataForm)
	_, err := r.db.Exec(ctx, `
		UPDATE form_pengajuan SET status=$1, data_form=$2, catatan=$3, diajukan_oleh=$4, diajukan_at=$5, updated_at=NOW()
		WHERE id=$6
	`, f.Status, data, f.Catatan, f.DiajukanOleh, f.DiajukanAt, f.ID)
	return err
}

func (r *FormRepository) CreateApproval(ctx context.Context, a *form.FormApproval) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO form_approval (id, form_id, urutan, role_approver, approver_id, status, catatan, approved_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, a.ID, a.FormID, a.Urutan, a.RoleApprover, a.ApproverID, a.Status, a.Catatan, a.ApprovedAt, a.CreatedAt)
	return err
}

func (r *FormRepository) GetApprovalsByForm(ctx context.Context, formID uuid.UUID) ([]*form.FormApproval, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, form_id, urutan, role_approver, approver_id, status, catatan, approved_at, created_at
		FROM form_approval WHERE form_id = $1 ORDER BY urutan
	`, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*form.FormApproval
	for rows.Next() {
		a := &form.FormApproval{}
		err := rows.Scan(&a.ID, &a.FormID, &a.Urutan, &a.RoleApprover, &a.ApproverID,
			&a.Status, &a.Catatan, &a.ApprovedAt, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *FormRepository) UpdateApproval(ctx context.Context, a *form.FormApproval) error {
	_, err := r.db.Exec(ctx, `
		UPDATE form_approval SET approver_id=$1, status=$2, catatan=$3, approved_at=$4
		WHERE id=$5
	`, a.ApproverID, a.Status, a.Catatan, a.ApprovedAt, a.ID)
	return err
}

func (r *FormRepository) GenerateNomorForm(ctx context.Context, bmtID uuid.UUID, jenis form.JenisForm) (string, error) {
	var seq int
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(COUNT(*), 0) + 1 FROM form_pengajuan WHERE bmt_id = $1 AND jenis_form = $2
	`, bmtID, jenis).Scan(&seq)
	if err != nil {
		return "", fmt.Errorf("gagal generate nomor form: %w", err)
	}
	now := time.Now()
	return fmt.Sprintf("%s-%d%02d-%06d", jenis, now.Year(), now.Month(), seq), nil
}
