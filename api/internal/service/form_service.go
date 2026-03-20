package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/form"
	"github.com/bmt-saas/api/internal/domain/nasabah"
	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

type FormService struct {
	repo        form.Repository
	nasabahRepo nasabah.Repository
	rekeningRepo rekening.Repository
	settings    *settings.Resolver
}

func NewFormService(
	repo form.Repository,
	nasabahRepo nasabah.Repository,
	rekeningRepo rekening.Repository,
	settingsResolver *settings.Resolver,
) *FormService {
	return &FormService{
		repo:         repo,
		nasabahRepo:  nasabahRepo,
		rekeningRepo: rekeningRepo,
		settings:     settingsResolver,
	}
}

type BuatFormInput struct {
	BMTID     uuid.UUID
	CabangID  uuid.UUID
	JenisForm form.JenisForm
	DataForm  map[string]interface{}
	CreatedBy uuid.UUID
}

// BuatForm membuat form baru dengan status DRAFT.
func (s *FormService) BuatForm(ctx context.Context, input BuatFormInput) (*form.FormPengajuan, error) {
	nomor, err := s.repo.GenerateNomorForm(ctx, input.BMTID, input.JenisForm)
	if err != nil {
		return nil, fmt.Errorf("gagal generate nomor form: %w", err)
	}

	f := &form.FormPengajuan{
		ID:           uuid.New(),
		BMTID:        input.BMTID,
		CabangID:     input.CabangID,
		JenisForm:    input.JenisForm,
		NomorForm:    nomor,
		Status:       form.StatusDraft,
		DataForm:     input.DataForm,
		DiajukanOleh: input.CreatedBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Buat record approval berdasarkan konfigurasi approver dari settings
	approvers := s.settings.GetApprovers(ctx, input.BMTID, input.CabangID, string(input.JenisForm))
	for i, role := range approvers {
		a := &form.FormApproval{
			ID:           uuid.New(),
			FormID:       f.ID,
			Urutan:       int16(i + 1),
			RoleApprover: role,
			Status:       "MENUNGGU",
			CreatedAt:    time.Now(),
		}
		if err := s.repo.CreateApproval(ctx, a); err != nil {
			return nil, fmt.Errorf("gagal buat approval: %w", err)
		}
	}

	if err := s.repo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

// AjukanForm transisi DRAFT → DIAJUKAN.
func (s *FormService) AjukanForm(ctx context.Context, formID, userID uuid.UUID) (*form.FormPengajuan, error) {
	f, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	if err := f.Ajukan(userID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

// ProsesApproval menyetujui atau menolak form.
func (s *FormService) ProsesApproval(ctx context.Context, input form.ApprovalInput) (*form.FormPengajuan, error) {
	f, err := s.repo.GetByID(ctx, input.FormID)
	if err != nil {
		return nil, err
	}
	if f.Status != form.StatusDiajukan {
		return nil, form.ErrFormTidakBisaDiubah
	}

	approvals, err := s.repo.GetApprovalsByForm(ctx, f.ID)
	if err != nil {
		return nil, err
	}

	// Cari approval yang perlu diproses (MENUNGGU + role sesuai)
	var targetApproval *form.FormApproval
	for _, a := range approvals {
		if a.Status == "MENUNGGU" && a.RoleApprover == input.RoleApprover {
			targetApproval = a
			break
		}
	}
	if targetApproval == nil {
		return nil, form.ErrApproverTidakBerwenang
	}

	now := time.Now()
	targetApproval.ApproverID = &input.ApproverID
	targetApproval.Catatan = input.Catatan
	targetApproval.ApprovedAt = &now

	if input.Setujui {
		targetApproval.Status = "DISETUJUI"
	} else {
		targetApproval.Status = "DITOLAK"
		f.Status = form.StatusDitolak
		f.Catatan = input.Catatan
		_ = s.repo.UpdateApproval(ctx, targetApproval)
		_ = s.repo.Update(ctx, f)
		return f, nil
	}

	_ = s.repo.UpdateApproval(ctx, targetApproval)

	// Cek apakah semua approval sudah DISETUJUI
	semuaDisetujui := true
	for _, a := range approvals {
		if a.ID == targetApproval.ID {
			continue
		}
		if a.Status != "DISETUJUI" {
			semuaDisetujui = false
			break
		}
	}

	if semuaDisetujui {
		f.Status = form.StatusDisetujui
		_ = s.repo.Update(ctx, f)
		// Auto-execute
		if err := s.eksekusiForm(ctx, f); err != nil {
			// Log error tapi jangan rollback status form
			fmt.Printf("[FORM] gagal eksekusi form %s: %v\n", f.NomorForm, err)
		}
	}

	return f, nil
}

// eksekusiForm menjalankan aksi yang sesuai dengan jenis form setelah DISETUJUI.
func (s *FormService) eksekusiForm(ctx context.Context, f *form.FormPengajuan) error {
	switch f.JenisForm {
	case form.FormDaftarNasabah:
		return s.eksekusiDaftarNasabah(ctx, f)
	case form.FormBukaRekening:
		return s.eksekusiBukaRekening(ctx, f)
	case form.FormBlokirRekening:
		return s.eksekusiBlokirRekening(ctx, f)
	case form.FormTutupRekening:
		return s.eksekusiTutupRekening(ctx, f)
	case form.FormUbahDataNasabah:
		return s.eksekusiUbahDataNasabah(ctx, f)
	}
	return nil
}

func (s *FormService) eksekusiDaftarNasabah(ctx context.Context, f *form.FormPengajuan) error {
	data, _ := json.Marshal(f.DataForm)
	var input nasabah.CreateNasabahInput
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("data form tidak valid: %w", err)
	}
	input.BMTID = f.BMTID
	input.CabangID = f.CabangID

	nomor, err := s.nasabahRepo.GenerateNomorNasabah(ctx, f.BMTID)
	if err != nil {
		return err
	}
	n, err := nasabah.New(input, nomor)
	if err != nil {
		return err
	}
	return s.nasabahRepo.Create(ctx, n)
}

func (s *FormService) eksekusiBukaRekening(ctx context.Context, f *form.FormPengajuan) error {
	nasabahIDStr, _ := f.DataForm["nasabah_id"].(string)
	jenisIDStr, _ := f.DataForm["jenis_rekening_id"].(string)

	nasabahID, err := uuid.Parse(nasabahIDStr)
	if err != nil {
		return fmt.Errorf("nasabah_id tidak valid: %w", err)
	}
	jenisID, err := uuid.Parse(jenisIDStr)
	if err != nil {
		return fmt.Errorf("jenis_rekening_id tidak valid: %w", err)
	}

	jenis, err := s.rekeningRepo.GetJenisByID(ctx, jenisID)
	if err != nil {
		return err
	}

	nomor, err := s.rekeningRepo.GenerateNomorRekening(ctx, f.BMTID, f.CabangID, jenis.Kode)
	if err != nil {
		return err
	}

	rek, err := rekening.NewRekening(f.BMTID, f.CabangID, nasabahID, jenisID, nomor, jenis)
	if err != nil {
		return err
	}
	formID := f.ID
	rek.CreatedByFormID = &formID
	return s.rekeningRepo.Create(ctx, rek)
}

func (s *FormService) eksekusiBlokirRekening(ctx context.Context, f *form.FormPengajuan) error {
	rekeningIDStr, _ := f.DataForm["rekening_id"].(string)
	alasan, _ := f.DataForm["alasan"].(string)

	rekeningID, err := uuid.Parse(rekeningIDStr)
	if err != nil {
		return fmt.Errorf("rekening_id tidak valid: %w", err)
	}
	return s.rekeningRepo.UpdateStatus(ctx, rekeningID, rekening.StatusBlokir, alasan)
}

func (s *FormService) eksekusiTutupRekening(ctx context.Context, f *form.FormPengajuan) error {
	rekeningIDStr, _ := f.DataForm["rekening_id"].(string)
	rekeningID, err := uuid.Parse(rekeningIDStr)
	if err != nil {
		return fmt.Errorf("rekening_id tidak valid: %w", err)
	}
	return s.rekeningRepo.UpdateStatus(ctx, rekeningID, rekening.StatusTutup, "form tutup rekening disetujui")
}

func (s *FormService) eksekusiUbahDataNasabah(ctx context.Context, f *form.FormPengajuan) error {
	nasabahIDStr, _ := f.DataForm["nasabah_id"].(string)
	nasabahID, err := uuid.Parse(nasabahIDStr)
	if err != nil {
		return fmt.Errorf("nasabah_id tidak valid: %w", err)
	}

	n, err := s.nasabahRepo.GetByID(ctx, nasabahID)
	if err != nil {
		return err
	}

	if v, ok := f.DataForm["nama_lengkap"].(string); ok && v != "" {
		n.NamaLengkap = v
	}
	if v, ok := f.DataForm["alamat"].(string); ok {
		n.Alamat = v
	}
	if v, ok := f.DataForm["telepon"].(string); ok && v != "" {
		n.Telepon = v
	}
	if v, ok := f.DataForm["email"].(string); ok {
		n.Email = v
	}
	if v, ok := f.DataForm["pekerjaan"].(string); ok {
		n.Pekerjaan = v
	}

	return s.nasabahRepo.Update(ctx, n)
}

// GetForm mengembalikan form dengan cek tenant.
func (s *FormService) GetForm(ctx context.Context, formID, bmtID uuid.UUID) (*form.FormPengajuan, error) {
	f, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	if f.BMTID != bmtID {
		return nil, form.ErrFormNotFound
	}
	return f, nil
}

// ListForm mengembalikan daftar form dengan paginasi.
func (s *FormService) ListForm(ctx context.Context, bmtID, cabangID uuid.UUID, status form.StatusForm, page, perPage int) ([]*form.FormPengajuan, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return s.repo.ListByBMT(ctx, bmtID, cabangID, status, page, perPage)
}

// GetApprovals mengembalikan daftar approval untuk form.
func (s *FormService) GetApprovals(ctx context.Context, formID, bmtID uuid.UUID) ([]*form.FormApproval, error) {
	f, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	if f.BMTID != bmtID {
		return nil, form.ErrFormNotFound
	}
	return s.repo.GetApprovalsByForm(ctx, formID)
}

// GetApprovers adalah helper agar SettingsService mudah diakses di handler.
func (s *FormService) GetApprovers(ctx context.Context, bmtID, cabangID uuid.UUID, jenisForm string) []string {
	return s.settings.GetApprovers(ctx, bmtID, cabangID, jenisForm)
}
