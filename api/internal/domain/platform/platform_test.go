package platform_test

import (
	"testing"

	"github.com/bmt-saas/api/internal/domain/platform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBMT_Valid(t *testing.T) {
	input := platform.CreateBMTInput{
		Kode:    "ANNUR",
		Nama:    "BMT An-Nur Kudus",
		Alamat:  "Jl. Kudus Raya No. 1",
		Telepon: "0291-123456",
		Email:   "bmt@annur.co.id",
	}

	bmt, err := platform.NewBMT(input)
	require.NoError(t, err)
	require.NotNil(t, bmt)

	assert.Equal(t, "ANNUR", bmt.Kode)
	assert.Equal(t, "BMT An-Nur Kudus", bmt.Nama)
	assert.Equal(t, platform.StatusBMTAktif, bmt.Status)
}

func TestNewBMT_KodeKosong(t *testing.T) {
	input := platform.CreateBMTInput{
		Kode: "",
		Nama: "BMT Test",
	}

	_, err := platform.NewBMT(input)
	assert.Error(t, err)
}

func TestNewBMT_NamaKosong(t *testing.T) {
	input := platform.CreateBMTInput{
		Kode: "TEST",
		Nama: "",
	}

	_, err := platform.NewBMT(input)
	assert.Error(t, err)
}

func TestNewCabang_Valid(t *testing.T) {
	input := platform.CreateCabangInput{
		BMTID: uuid.New(),
		Kode:  "KDR",
		Nama:  "Cabang Kudus",
	}

	cabang, err := platform.NewCabang(input)
	require.NoError(t, err)
	require.NotNil(t, cabang)

	assert.Equal(t, "KDR", cabang.Kode)
	assert.Equal(t, "AKTIF", cabang.Status)
}
