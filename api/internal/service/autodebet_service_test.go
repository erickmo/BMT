package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainAutodebet "github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/bmt-saas/api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── Tests: HitungTanggalJatuhTempo ───────────────────────────────────────────

// TestAutodebet_TanggalDariRekeningConfig_Benar adalah test wajib per CLAUDE.md.
// Tanggal autodebet harus dari konfigurasi DB, bukan hardcode.
func TestAutodebet_TanggalDariRekeningConfig_Benar(t *testing.T) {
	skenario := []struct {
		nama     string
		bulan    time.Time
		tanggal  int
		expected time.Time
	}{
		{
			nama:     "tanggal 5 Januari — normal",
			bulan:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  5,
			expected: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 28 Februari — dalam batas",
			bulan:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  28,
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 31 Februari — disesuaikan ke 28",
			bulan:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  31,
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 29 Februari — disesuaikan ke 28 (bukan tahun kabisat)",
			bulan:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  29,
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 29 Februari — valid di tahun kabisat 2024",
			bulan:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  29,
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 31 Maret — dalam batas",
			bulan:    time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  31,
			expected: time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 31 April — disesuaikan ke 30",
			bulan:    time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  31,
			expected: time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			nama:     "tanggal 1 Desember — batas bawah",
			bulan:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			tanggal:  1,
			expected: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, sk := range skenario {
		t.Run(sk.nama, func(t *testing.T) {
			result := service.HitungTanggalJatuhTempo(sk.bulan, sk.tanggal)
			assert.Equal(t, sk.expected, result,
				"tanggal jatuh tempo harus dihitung dari konfigurasi DB, bukan hardcode")
		})
	}
}

// ── Tests: GenerateJadwalBulanan ─────────────────────────────────────────────

func TestAutodebetService_GenerateJadwalBulanan_SatuKonfigurasi(t *testing.T) {
	autodebetRepo := new(MockAutodebetRepo)
	rekeningRepo := new(MockRekeningRepo)
	jurnalSvc := new(MockJurnalService)
	rekeningService := service.NewRekeningService(rekeningRepo, autodebetRepo, nil, jurnalSvc)
	svc := service.NewAutodebetService(autodebetRepo, rekeningService)

	bmtID := uuid.New()
	rekeningID := uuid.New()
	configID := uuid.New()
	bulan := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	configs := []*domainAutodebet.Config{
		{
			ID:           configID,
			BMTID:        bmtID,
			RekeningID:   rekeningID,
			Jenis:        domainAutodebet.JenisSimpananWajib,
			TanggalDebet: 10,
			IsAktif:      true,
		},
	}

	autodebetRepo.On("ListConfigByRekening", mock.Anything, rekeningID).Return(configs, nil)
	autodebetRepo.On("CreateJadwal", mock.Anything, mock.MatchedBy(func(j *domainAutodebet.Jadwal) bool {
		expected := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
		return j.RekeningID == rekeningID &&
			j.BMTID == bmtID &&
			j.Jenis == domainAutodebet.JenisSimpananWajib &&
			j.TanggalJatuhTempo.Equal(expected) &&
			j.Status == domainAutodebet.StatusMenunggu
	})).Return(nil)

	err := svc.GenerateJadwalBulanan(context.Background(), []uuid.UUID{rekeningID}, bmtID, bulan)

	assert.NoError(t, err)
	autodebetRepo.AssertExpectations(t)
}

func TestAutodebetService_GenerateJadwalBulanan_KonfigurasiTidakAktif(t *testing.T) {
	autodebetRepo := new(MockAutodebetRepo)
	rekeningRepo := new(MockRekeningRepo)
	jurnalSvc := new(MockJurnalService)
	rekeningService := service.NewRekeningService(rekeningRepo, autodebetRepo, nil, jurnalSvc)
	svc := service.NewAutodebetService(autodebetRepo, rekeningService)

	bmtID := uuid.New()
	rekeningID := uuid.New()
	bulan := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	configs := []*domainAutodebet.Config{
		{
			ID:           uuid.New(),
			BMTID:        bmtID,
			RekeningID:   rekeningID,
			Jenis:        domainAutodebet.JenisSimpananWajib,
			TanggalDebet: 10,
			IsAktif:      false, // tidak aktif, harus diskip
		},
	}

	autodebetRepo.On("ListConfigByRekening", mock.Anything, rekeningID).Return(configs, nil)

	err := svc.GenerateJadwalBulanan(context.Background(), []uuid.UUID{rekeningID}, bmtID, bulan)

	assert.NoError(t, err)
	// CreateJadwal tidak boleh dipanggil untuk config tidak aktif
	autodebetRepo.AssertNotCalled(t, "CreateJadwal", mock.Anything, mock.Anything)
}

func TestAutodebetService_GenerateJadwalBulanan_ErrorDB_LanjutKeRekeningBerikutnya(t *testing.T) {
	autodebetRepo := new(MockAutodebetRepo)
	rekeningRepo := new(MockRekeningRepo)
	jurnalSvc := new(MockJurnalService)
	rekeningService := service.NewRekeningService(rekeningRepo, autodebetRepo, nil, jurnalSvc)
	svc := service.NewAutodebetService(autodebetRepo, rekeningService)

	bmtID := uuid.New()
	rek1ID := uuid.New()
	rek2ID := uuid.New()
	bulan := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	// Rekening 1: error DB
	autodebetRepo.On("ListConfigByRekening", mock.Anything, rek1ID).
		Return([]*domainAutodebet.Config{}, errors.New("connection error"))

	// Rekening 2: sukses
	configs2 := []*domainAutodebet.Config{
		{
			ID:           uuid.New(),
			BMTID:        bmtID,
			RekeningID:   rek2ID,
			Jenis:        domainAutodebet.JenisBiayaAdmin,
			TanggalDebet: 1,
			IsAktif:      true,
		},
	}
	autodebetRepo.On("ListConfigByRekening", mock.Anything, rek2ID).Return(configs2, nil)
	autodebetRepo.On("CreateJadwal", mock.Anything, mock.Anything).Return(nil)

	// Harus mengembalikan nil meski rekening pertama error
	err := svc.GenerateJadwalBulanan(context.Background(), []uuid.UUID{rek1ID, rek2ID}, bmtID, bulan)

	assert.NoError(t, err, "error satu rekening tidak boleh menghentikan rekening lain")
	autodebetRepo.AssertExpectations(t)
}

func TestAutodebetService_GenerateJadwalBulanan_SliceKosong_TidakPanicTidakError(t *testing.T) {
	autodebetRepo := new(MockAutodebetRepo)
	rekeningRepo := new(MockRekeningRepo)
	jurnalSvc := new(MockJurnalService)
	rekeningService := service.NewRekeningService(rekeningRepo, autodebetRepo, nil, jurnalSvc)
	svc := service.NewAutodebetService(autodebetRepo, rekeningService)

	bmtID := uuid.New()
	bulan := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)

	// Slice kosong — tidak boleh panic
	err := svc.GenerateJadwalBulanan(context.Background(), []uuid.UUID{}, bmtID, bulan)

	assert.NoError(t, err)
	autodebetRepo.AssertNotCalled(t, "ListConfigByRekening", mock.Anything, mock.Anything)
	autodebetRepo.AssertNotCalled(t, "CreateJadwal", mock.Anything, mock.Anything)
}
