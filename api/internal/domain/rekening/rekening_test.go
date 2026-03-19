package rekening_test

import (
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buatRekening(saldo int64, status rekening.StatusRekening) *rekening.Rekening {
	return &rekening.Rekening{
		ID:              uuid.New(),
		BMTID:           uuid.New(),
		CabangID:        uuid.New(),
		NasabahID:       uuid.New(),
		JenisRekeningID: uuid.New(),
		NomorRekening:   "ANNUR-KDR-SU-00000001",
		Saldo:           money.New(saldo),
		Status:          status,
		TanggalBuka:     time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func TestRekening_Setor_Valid(t *testing.T) {
	rek := buatRekening(100000, rekening.StatusAktif)
	err := rek.ValidasiSetor(50000, 10000)
	assert.NoError(t, err)
}

func TestRekening_Setor_DibawahMinimum(t *testing.T) {
	rek := buatRekening(100000, rekening.StatusAktif)
	err := rek.ValidasiSetor(5000, 10000)
	assert.ErrorIs(t, err, rekening.ErrSetoranDibawahMin)
}

func TestRekening_Setor_Blokir(t *testing.T) {
	rek := buatRekening(100000, rekening.StatusBlokir)
	err := rek.ValidasiSetor(50000, 0)
	assert.ErrorIs(t, err, rekening.ErrRekeningBeku)
}

func TestRekening_Tarik_Valid(t *testing.T) {
	rek := buatRekening(500000, rekening.StatusAktif)
	err := rek.ValidasiTarik(100000, true)
	assert.NoError(t, err)
}

func TestRekening_Tarik_SaldoKurang(t *testing.T) {
	rek := buatRekening(50000, rekening.StatusAktif)
	err := rek.ValidasiTarik(100000, true)
	assert.ErrorIs(t, err, rekening.ErrSaldoTidakCukup)
}

func TestRekening_Tarik_TidakBisaDitarik(t *testing.T) {
	rek := buatRekening(500000, rekening.StatusAktif)
	err := rek.ValidasiTarik(100000, false)
	assert.ErrorIs(t, err, rekening.ErrPenarikanTidakBisa)
}

func TestRekening_NewTransaksiSetor_SaldoBenar(t *testing.T) {
	rek := buatRekening(100000, rekening.StatusAktif)
	createdBy := uuid.New()
	tr := rek.NewTransaksiSetor(50000, "setoran tunai", &createdBy, nil)

	require.NotNil(t, tr)
	assert.Equal(t, int64(100000), tr.SaldoSebelum)
	assert.Equal(t, int64(150000), tr.SaldoSesudah)
	assert.Equal(t, int64(50000), tr.Nominal)
	assert.Equal(t, "KREDIT", tr.Posisi)
}

func TestRekening_NewTransaksiTarik_SaldoBenar(t *testing.T) {
	rek := buatRekening(500000, rekening.StatusAktif)
	createdBy := uuid.New()
	tr := rek.NewTransaksiTarik(100000, "penarikan tunai", &createdBy, nil)

	require.NotNil(t, tr)
	assert.Equal(t, int64(500000), tr.SaldoSebelum)
	assert.Equal(t, int64(400000), tr.SaldoSesudah)
	assert.Equal(t, int64(100000), tr.Nominal)
	assert.Equal(t, "DEBIT", tr.Posisi)
}
