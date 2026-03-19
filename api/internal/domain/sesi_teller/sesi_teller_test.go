package sesi_teller_test

import (
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/sesi_teller"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func buatSesiAktif(saldoAwal int64) *sesi_teller.SesiTeller {
	return &sesi_teller.SesiTeller{
		ID:       uuid.New(),
		BMTID:    uuid.New(),
		CabangID: uuid.New(),
		TellerID: uuid.New(),
		Tanggal:  time.Now(),
		SaldoAwal: money.New(saldoAwal),
		Redenominasi: []sesi_teller.ItemPecahan{
			{Nominal: 100000, Jumlah: 3, Subtotal: 300000},
			{Nominal: 50000, Jumlah: 2, Subtotal: 100000},
		},
		Status:           sesi_teller.StatusAktif,
		ToleransiSelisih: money.Zero,
		DibukaPada:       time.Now(),
	}
}

func TestSesiTeller_TutupSesi_Seimbang(t *testing.T) {
	sesi := buatSesiAktif(400000)

	redenominasiAkhir := []sesi_teller.ItemPecahan{
		{Nominal: 100000, Jumlah: 3, Subtotal: 300000},
		{Nominal: 50000, Jumlah: 2, Subtotal: 100000},
	}

	err := sesi.TutupSesi(redenominasiAkhir, money.Zero)
	assert.NoError(t, err)
	assert.Equal(t, sesi_teller.StatusTutup, sesi.Status)
}

func TestSesiTeller_TutupSesi_Selisih_MelebihiToleransi(t *testing.T) {
	sesi := buatSesiAktif(400000)

	// Hanya ada 350000 di kasir (selisih 50000)
	redenominasiAkhir := []sesi_teller.ItemPecahan{
		{Nominal: 100000, Jumlah: 3, Subtotal: 300000},
		{Nominal: 50000, Jumlah: 1, Subtotal: 50000},
	}

	err := sesi.TutupSesi(redenominasiAkhir, money.Zero)
	assert.ErrorIs(t, err, sesi_teller.ErrSesiSelisih)
}

func TestSesiTeller_TutupSesi_Selisih_DalamToleransi(t *testing.T) {
	sesi := buatSesiAktif(400000)

	// Selisih 1000 (dalam toleransi 5000)
	redenominasiAkhir := []sesi_teller.ItemPecahan{
		{Nominal: 100000, Jumlah: 3, Subtotal: 300000},
		{Nominal: 50000, Jumlah: 2, Subtotal: 100000},
		{Nominal: 1000, Jumlah: -1, Subtotal: -1000}, // simulasi kurang 1000
	}

	// Dengan toleransi 5000
	err := sesi.TutupSesi(redenominasiAkhir, money.New(5000))
	// Selisih = 400000 - 399000 = 1000, dalam toleransi 5000
	assert.NoError(t, err)
}

func TestSesiTeller_TutupSesi_SudahTutup(t *testing.T) {
	sesi := buatSesiAktif(400000)
	sesi.Status = sesi_teller.StatusTutup

	err := sesi.TutupSesi([]sesi_teller.ItemPecahan{}, money.Zero)
	assert.ErrorIs(t, err, sesi_teller.ErrSesiSudahTutup)
}
