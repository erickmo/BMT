package service_test

import (
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHitungTanggalJatuhTempo_NormalDate(t *testing.T) {
	bulan := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tanggal := service.HitungTanggalJatuhTempo(bulan, 15)

	assert.Equal(t, 15, tanggal.Day())
	assert.Equal(t, time.January, tanggal.Month())
	assert.Equal(t, 2025, tanggal.Year())
}

func TestHitungTanggalJatuhTempo_TanggalMelebihiBulan(t *testing.T) {
	// Februari hanya 28 hari di 2025
	bulan := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	tanggal := service.HitungTanggalJatuhTempo(bulan, 30)

	assert.Equal(t, 28, tanggal.Day()) // Disesuaikan ke akhir bulan
	assert.Equal(t, time.February, tanggal.Month())
}

func TestHitungTanggalJatuhTempo_Tanggal28(t *testing.T) {
	bulan := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	tanggal := service.HitungTanggalJatuhTempo(bulan, 28)

	assert.Equal(t, 28, tanggal.Day())
	assert.Equal(t, time.March, tanggal.Month())
}
