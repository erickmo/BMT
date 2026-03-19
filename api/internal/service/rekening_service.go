package service

import (
	"context"
	"fmt"

	"github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

type RekeningService struct {
	repo             rekening.Repository
	autodebetRepo    autodebet.Repository
	settingsResolver *settings.Resolver
	jurnalService    JurnalService
}

type JurnalService interface {
	PostJurnal(ctx context.Context, j PostJurnalInput) error
}

// PostJurnalInput dan JurnalEntry didefinisikan di akuntansi_service.go

func NewRekeningService(
	repo rekening.Repository,
	autodebetRepo autodebet.Repository,
	settingsResolver *settings.Resolver,
	jurnalService JurnalService,
) *RekeningService {
	return &RekeningService{
		repo:             repo,
		autodebetRepo:    autodebetRepo,
		settingsResolver: settingsResolver,
		jurnalService:    jurnalService,
	}
}

func (s *RekeningService) Setor(ctx context.Context, input rekening.SetoranInput) (*rekening.TransaksiRekening, error) {
	// Cek idempotency
	if input.IdempotencyKey != nil {
		existing, err := s.repo.GetTransaksiByIdempotency(ctx, *input.IdempotencyKey)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	// Lock rekening untuk update
	rek, err := s.repo.LockForUpdate(ctx, input.RekeningID)
	if err != nil {
		return nil, fmt.Errorf("gagal lock rekening: %w", err)
	}

	// Ambil jenis rekening untuk validasi
	jenis, err := s.repo.GetJenisByID(ctx, rek.JenisRekeningID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil jenis rekening: %w", err)
	}

	if err := rek.ValidasiSetor(input.Nominal, jenis.SetoranMin); err != nil {
		return nil, err
	}

	// Buat transaksi
	tr := rek.NewTransaksiSetor(input.Nominal, input.Keterangan, &input.CreatedBy, input.IdempotencyKey)

	// Update saldo
	if err := s.repo.UpdateSaldo(ctx, rek.ID, tr.SaldoSesudah); err != nil {
		return nil, fmt.Errorf("gagal update saldo: %w", err)
	}

	// Simpan transaksi
	if err := s.repo.CreateTransaksi(ctx, tr); err != nil {
		return nil, fmt.Errorf("gagal simpan transaksi: %w", err)
	}

	// Post jurnal: Debit Kas / Kredit Rekening
	_ = s.jurnalService.PostJurnal(ctx, PostJurnalInput{
		BMTID:      rek.BMTID,
		CabangID:   rek.CabangID,
		Keterangan: fmt.Sprintf("Setoran rekening %s", rek.NomorRekening),
		Referensi:  tr.ID.String(),
		Entries: []JurnalEntry{
			{KodeAkun: "101", Posisi: "DEBIT", Nominal: input.Nominal},
			{KodeAkun: "202", Posisi: "KREDIT", Nominal: input.Nominal},
		},
	})

	return tr, nil
}

func (s *RekeningService) Tarik(ctx context.Context, input rekening.PenarikanInput) (*rekening.TransaksiRekening, error) {
	// Cek idempotency
	if input.IdempotencyKey != nil {
		existing, err := s.repo.GetTransaksiByIdempotency(ctx, *input.IdempotencyKey)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	rek, err := s.repo.LockForUpdate(ctx, input.RekeningID)
	if err != nil {
		return nil, fmt.Errorf("gagal lock rekening: %w", err)
	}

	jenis, err := s.repo.GetJenisByID(ctx, rek.JenisRekeningID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil jenis rekening: %w", err)
	}

	if err := rek.ValidasiTarik(input.Nominal, jenis.BisaDitarik); err != nil {
		return nil, err
	}

	tr := rek.NewTransaksiTarik(input.Nominal, input.Keterangan, &input.CreatedBy, input.IdempotencyKey)

	if err := s.repo.UpdateSaldo(ctx, rek.ID, tr.SaldoSesudah); err != nil {
		return nil, fmt.Errorf("gagal update saldo: %w", err)
	}

	if err := s.repo.CreateTransaksi(ctx, tr); err != nil {
		return nil, fmt.Errorf("gagal simpan transaksi: %w", err)
	}

	// Post jurnal: Debit Rekening / Kredit Kas
	_ = s.jurnalService.PostJurnal(ctx, PostJurnalInput{
		BMTID:      rek.BMTID,
		CabangID:   rek.CabangID,
		Keterangan: fmt.Sprintf("Penarikan rekening %s", rek.NomorRekening),
		Referensi:  tr.ID.String(),
		Entries: []JurnalEntry{
			{KodeAkun: "202", Posisi: "DEBIT", Nominal: input.Nominal},
			{KodeAkun: "101", Posisi: "KREDIT", Nominal: input.Nominal},
		},
	})

	return tr, nil
}

// EksekusiAutodebetJadwal menjalankan autodebet untuk satu jadwal.
// Mengimplementasikan partial debit: jika saldo kurang, debit semampu saldo,
// sisanya jadi tunggakan.
func (s *RekeningService) EksekusiAutodebetJadwal(ctx context.Context, jadwal *autodebet.Jadwal) error {
	rek, err := s.repo.LockForUpdate(ctx, jadwal.RekeningID)
	if err != nil {
		return fmt.Errorf("gagal lock rekening: %w", err)
	}

	hasil := autodebet.EksekusiAutodebet(rek.Saldo, jadwal.NominalTarget)

	if hasil.NominalDidebit > 0 {
		tr := rek.NewTransaksiTarik(
			hasil.NominalDidebit.Int64(),
			fmt.Sprintf("Autodebet %s", jadwal.Jenis),
			nil,
			nil,
		)
		if err := s.repo.UpdateSaldo(ctx, rek.ID, tr.SaldoSesudah); err != nil {
			return err
		}
		if err := s.repo.CreateTransaksi(ctx, tr); err != nil {
			return err
		}
	}

	// Jika partial: buat tunggakan
	if hasil.IsPartial {
		tunggakanID := uuid.New()
		t := &autodebet.Tunggakan{
			ID:              tunggakanID,
			BMTID:           rek.BMTID,
			RekeningID:      rek.ID,
			JadwalID:        jadwal.ID,
			Jenis:           jadwal.Jenis,
			NominalTarget:   jadwal.NominalTarget,
			NominalTerbayar: hasil.NominalDidebit,
			NominalSisa:     hasil.NominalTunggakan,
			Status:          "OUTSTANDING",
		}
		hasil.TunggakanID = &tunggakanID
		if err := s.autodebetRepo.CreateTunggakan(ctx, t); err != nil {
			return err
		}
		if err := s.autodebetRepo.UpdateJadwalStatus(ctx, jadwal.ID, autodebet.StatusPartial); err != nil {
			return err
		}
	} else {
		if err := s.autodebetRepo.UpdateJadwalStatus(ctx, jadwal.ID, autodebet.StatusSukses); err != nil {
			return err
		}
	}

	return nil
}

// GetSaldo returns saldo rekening dengan tenant isolation
func (s *RekeningService) GetSaldo(ctx context.Context, rekeningID, bmtID uuid.UUID) (money.Money, error) {
	rek, err := s.repo.GetByID(ctx, rekeningID)
	if err != nil {
		return money.Zero, err
	}
	// Tenant isolation check
	if rek.BMTID != bmtID {
		return money.Zero, rekening.ErrRekeningNotFound
	}
	return rek.Saldo, nil
}
