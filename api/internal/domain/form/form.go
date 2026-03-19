package form

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFormNotFound           = errors.New("form tidak ditemukan")
	ErrFormTidakBisaDiubah    = errors.New("status form tidak mengizinkan perubahan")
	ErrApproverTidakBerwenang = errors.New("role tidak berwenang menyetujui form ini")
	ErrFormSudahDiproses      = errors.New("form sudah diproses")
)

type JenisForm string

const (
	FormDaftarNasabah   JenisForm = "FORM_DAFTAR_NASABAH"
	FormBukaRekening    JenisForm = "FORM_BUKA_REKENING"
	FormTutupRekening   JenisForm = "FORM_TUTUP_REKENING"
	FormBlokirRekening  JenisForm = "FORM_BLOKIR_REKENING"
	FormBukaPembiayaan  JenisForm = "FORM_BUKA_PEMBIAYAAN"
	FormUbahDataNasabah JenisForm = "FORM_UBAH_DATA_NASABAH"
)

type StatusForm string

const (
	StatusDraft     StatusForm = "DRAFT"
	StatusDiajukan  StatusForm = "DIAJUKAN"
	StatusDisetujui StatusForm = "DISETUJUI"
	StatusDitolak   StatusForm = "DITOLAK"
	StatusBatal     StatusForm = "BATAL"
)

type FormPengajuan struct {
	ID           uuid.UUID              `json:"id"`
	BMTID        uuid.UUID              `json:"bmt_id"`
	CabangID     uuid.UUID              `json:"cabang_id"`
	JenisForm    JenisForm              `json:"jenis_form"`
	NomorForm    string                 `json:"nomor_form"`
	Status       StatusForm             `json:"status"`
	DataForm     map[string]interface{} `json:"data_form"`
	Catatan      string                 `json:"catatan"`
	DiajukanOleh uuid.UUID              `json:"diajukan_oleh"`
	DiajukanAt   *time.Time             `json:"diajukan_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type FormApproval struct {
	ID           uuid.UUID  `json:"id"`
	FormID       uuid.UUID  `json:"form_id"`
	Urutan       int16      `json:"urutan"`
	RoleApprover string     `json:"role_approver"`
	ApproverID   *uuid.UUID `json:"approver_id,omitempty"`
	Status       string     `json:"status"`
	Catatan      string     `json:"catatan,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type AjukanInput struct {
	FormID       uuid.UUID `json:"form_id" validate:"required"`
	DiajukanOleh uuid.UUID `json:"diajukan_oleh" validate:"required"`
}

type ApprovalInput struct {
	FormID       uuid.UUID `json:"form_id" validate:"required"`
	ApproverID   uuid.UUID `json:"approver_id" validate:"required"`
	RoleApprover string    `json:"role_approver" validate:"required"`
	Catatan      string    `json:"catatan"`
	Setujui      bool      `json:"setujui"`
}

type Repository interface {
	Create(ctx context.Context, f *FormPengajuan) error
	GetByID(ctx context.Context, id uuid.UUID) (*FormPengajuan, error)
	ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, status StatusForm, page, perPage int) ([]*FormPengajuan, int64, error)
	Update(ctx context.Context, f *FormPengajuan) error

	CreateApproval(ctx context.Context, a *FormApproval) error
	GetApprovalsByForm(ctx context.Context, formID uuid.UUID) ([]*FormApproval, error)
	UpdateApproval(ctx context.Context, a *FormApproval) error

	GenerateNomorForm(ctx context.Context, bmtID uuid.UUID, jenis JenisForm) (string, error)
}

func (f *FormPengajuan) BisaDiubah() bool {
	return f.Status == StatusDraft
}

func (f *FormPengajuan) Ajukan(diajukanOleh uuid.UUID) error {
	if f.Status != StatusDraft {
		return ErrFormTidakBisaDiubah
	}
	now := time.Now()
	f.Status = StatusDiajukan
	f.DiajukanOleh = diajukanOleh
	f.DiajukanAt = &now
	f.UpdatedAt = now
	return nil
}
