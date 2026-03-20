package keuangan_test

import (
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/pondok/keuangan"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ── Helper ────────────────────────────────────────────────────────────────────

func buatTagihanSPP(nominal int64) *keuangan.TagihanSPP {
	tagihan, _ := keuangan.NewTagihanSPP(
		uuid.New(), uuid.New(), uuid.New(), uuid.New(),
		"2025-03",
		nominal,
		time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
	)
	return tagihan
}

// ── Tests: TagihanSPP — NewTagihanSPP ─────────────────────────────────────────

func TestNewTagihanSPP_Valid(t *testing.T) {
	tagihan, err := keuangan.NewTagihanSPP(
		uuid.New(), uuid.New(), uuid.New(), uuid.New(),
		"2025-03",
		500000,
		time.Now().Add(30*24*time.Hour),
	)

	assert.NoError(t, err)
	assert.NotNil(t, tagihan)
	assert.Equal(t, int64(500000), tagihan.Nominal)
	assert.Equal(t, int64(500000), tagihan.NominalEfektif)
	assert.Equal(t, int64(500000), tagihan.NominalSisa)
	assert.Equal(t, int64(0), tagihan.NominalTerbayar)
	assert.Equal(t, keuangan.StatusBelumBayar, tagihan.Status)
	assert.Equal(t, float64(0), tagihan.BeasiswaPersen)
}

func TestNewTagihanSPP_PeriodeKosong_Error(t *testing.T) {
	_, err := keuangan.NewTagihanSPP(
		uuid.New(), uuid.New(), uuid.New(), uuid.New(),
		"", // periode kosong
		500000,
		time.Now(),
	)
	assert.Error(t, err)
}

func TestNewTagihanSPP_NominalNol_Error(t *testing.T) {
	_, err := keuangan.NewTagihanSPP(
		uuid.New(), uuid.New(), uuid.New(), uuid.New(),
		"2025-03",
		0, // nominal nol
		time.Now(),
	)
	assert.ErrorIs(t, err, keuangan.ErrNominalHarusPositif)
}

// ── Tests: TagihanSPP — TerapkanBeasiswa ──────────────────────────────────────

// TestBeasiswa_TagihanSPP_NominalEfektifBenar adalah test wajib per CLAUDE.md.
// NominalEfektif = Nominal - BeasiswaNominal setelah beasiswa diterapkan.
func TestBeasiswa_TagihanSPP_NominalEfektifBenar(t *testing.T) {
	skenario := []struct {
		nama            string
		nominal         int64
		persenBeasiswa  float64
		expectedNominal int64 // BeasiswaNominal
		expectedEfektif int64 // NominalEfektif
	}{
		{
			nama:            "beasiswa 50% dari SPP 500rb",
			nominal:         500000,
			persenBeasiswa:  50,
			expectedNominal: 250000,
			expectedEfektif: 250000,
		},
		{
			nama:            "beasiswa 100% (beasiswa penuh)",
			nominal:         500000,
			persenBeasiswa:  100,
			expectedNominal: 500000,
			expectedEfektif: 0,
		},
		{
			nama:            "beasiswa 0% (tidak ada beasiswa)",
			nominal:         500000,
			persenBeasiswa:  0,
			expectedNominal: 0,
			expectedEfektif: 500000,
		},
		{
			nama:            "beasiswa 25% dari SPP 1000rb",
			nominal:         1000000,
			persenBeasiswa:  25,
			expectedNominal: 250000,
			expectedEfektif: 750000,
		},
		{
			nama:            "beasiswa 33.33% dari SPP 300rb",
			nominal:         300000,
			persenBeasiswa:  33.33,
			expectedNominal: 99990,  // int64(300000 * 33.33 / 100) = 99990
			expectedEfektif: 200010, // 300000 - 99990
		},
	}

	for _, sk := range skenario {
		t.Run(sk.nama, func(t *testing.T) {
			tagihan := buatTagihanSPP(sk.nominal)

			err := tagihan.TerapkanBeasiswa(sk.persenBeasiswa)

			assert.NoError(t, err)
			assert.Equal(t, sk.persenBeasiswa, tagihan.BeasiswaPersen,
				"persentase beasiswa harus tersimpan")
			assert.Equal(t, sk.expectedNominal, tagihan.BeasiswaNominal,
				"nominal beasiswa harus dihitung dengan benar")
			assert.Equal(t, sk.expectedEfektif, tagihan.NominalEfektif,
				"NominalEfektif = Nominal - BeasiswaNominal")
			assert.Equal(t, sk.expectedEfektif, tagihan.NominalSisa,
				"NominalSisa harus diperbarui berdasarkan NominalEfektif")
		})
	}
}

func TestBeasiswa_PersenTidakValid_Error(t *testing.T) {
	tagihan := buatTagihanSPP(500000)

	t.Run("persen negatif", func(t *testing.T) {
		err := tagihan.TerapkanBeasiswa(-1)
		assert.ErrorIs(t, err, keuangan.ErrBeasiswaPersenTidakValid)
	})

	t.Run("persen lebih dari 100", func(t *testing.T) {
		err := tagihan.TerapkanBeasiswa(101)
		assert.ErrorIs(t, err, keuangan.ErrBeasiswaPersenTidakValid)
	})
}

// ── Tests: TagihanSPP — Bayar ─────────────────────────────────────────────────

func TestTagihanSPP_Bayar_Lunas(t *testing.T) {
	tagihan := buatTagihanSPP(500000)

	dibayar, err := tagihan.Bayar(500000)

	assert.NoError(t, err)
	assert.Equal(t, int64(500000), dibayar)
	assert.Equal(t, keuangan.StatusLunas, tagihan.Status)
	assert.Equal(t, int64(0), tagihan.NominalSisa)
	assert.NotNil(t, tagihan.TanggalLunas)
}

func TestTagihanSPP_Bayar_Parsial(t *testing.T) {
	tagihan := buatTagihanSPP(500000)

	dibayar, err := tagihan.Bayar(200000)

	assert.NoError(t, err)
	assert.Equal(t, int64(200000), dibayar)
	assert.Equal(t, keuangan.StatusSebagian, tagihan.Status)
	assert.Equal(t, int64(300000), tagihan.NominalSisa)
	assert.Equal(t, int64(200000), tagihan.NominalTerbayar)
	assert.Nil(t, tagihan.TanggalLunas)
}

func TestTagihanSPP_Bayar_MelebihiSisa_Error(t *testing.T) {
	tagihan := buatTagihanSPP(500000)

	_, err := tagihan.Bayar(600000)

	assert.ErrorIs(t, err, keuangan.ErrPembayaranMelebihi)
}

func TestTagihanSPP_Bayar_SudahLunas_Error(t *testing.T) {
	tagihan := buatTagihanSPP(500000)

	tagihan.Bayar(500000) // lunas
	_, err := tagihan.Bayar(1000)

	assert.ErrorIs(t, err, keuangan.ErrTagihanSudahLunas)
}

func TestTagihanSPP_Bayar_SetelaBeasiswa(t *testing.T) {
	tagihan := buatTagihanSPP(500000)

	// Terapkan beasiswa 50%: nominal efektif jadi 250rb
	err := tagihan.TerapkanBeasiswa(50)
	assert.NoError(t, err)

	// Bayar 250rb (nominal efektif penuh) → harus lunas
	dibayar, err := tagihan.Bayar(250000)

	assert.NoError(t, err)
	assert.Equal(t, int64(250000), dibayar)
	assert.Equal(t, keuangan.StatusLunas, tagihan.Status,
		"setelah beasiswa 50%, bayar nominal efektif harus langsung lunas")
}

// ── Tests: JenisTagihan ───────────────────────────────────────────────────────

func TestNewJenisTagihan_Valid(t *testing.T) {
	jenis, err := keuangan.NewJenisTagihan(
		uuid.New(), "SPP-REG", "SPP Reguler", 500000, keuangan.FrekuensiBulanan,
	)

	assert.NoError(t, err)
	assert.NotNil(t, jenis)
	assert.Equal(t, "SPP-REG", jenis.Kode)
	assert.Equal(t, int64(500000), jenis.Nominal)
	assert.True(t, jenis.IsAktif)
}

func TestNewJenisTagihan_KodeKosong_Error(t *testing.T) {
	_, err := keuangan.NewJenisTagihan(uuid.New(), "", "SPP Reguler", 500000, keuangan.FrekuensiBulanan)
	assert.Error(t, err)
}

func TestNewJenisTagihan_NamaKosong_Error(t *testing.T) {
	_, err := keuangan.NewJenisTagihan(uuid.New(), "SPP-REG", "", 500000, keuangan.FrekuensiBulanan)
	assert.Error(t, err)
}

func TestNewJenisTagihan_NominalNol_Error(t *testing.T) {
	_, err := keuangan.NewJenisTagihan(uuid.New(), "SPP-REG", "SPP Reguler", 0, keuangan.FrekuensiBulanan)
	assert.ErrorIs(t, err, keuangan.ErrNominalHarusPositif)
}
