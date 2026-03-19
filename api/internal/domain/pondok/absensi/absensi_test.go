package absensi_test

import (
	"context"
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/pondok/absensi"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSettingsStore implementasi settings.Store untuk testing
type mockSettingsStore struct {
	platform map[string]string
	bmt      map[string]string
	cabang   map[string]string
}

func (m *mockSettingsStore) GetPlatform(_ context.Context, kunci string) (string, error) {
	if v, ok := m.platform[kunci]; ok {
		return v, nil
	}
	return "", assert.AnError
}

func (m *mockSettingsStore) GetBMT(_ context.Context, _ uuid.UUID, kunci string) (string, error) {
	if v, ok := m.bmt[kunci]; ok {
		return v, nil
	}
	return "", assert.AnError
}

func (m *mockSettingsStore) GetCabang(_ context.Context, _ uuid.UUID, kunci string) (string, error) {
	if v, ok := m.cabang[kunci]; ok {
		return v, nil
	}
	return "", assert.AnError
}

// TestAbsensi_MetodeDariSettings_BukanKonstanta memastikan validasi metode absensi
// membaca dari settings DB, bukan dari konstanta hardcode.
//
// Prinsip: "pondok.absensi_metode" dari settings DB menentukan metode yang diizinkan.
// Domain layer hanya menyediakan konstanta type-safe; validasi aktual menggunakan resolver.
func TestAbsensi_MetodeDariSettings_BukanKonstanta(t *testing.T) {
	bmtID := uuid.New()
	cabangID := uuid.New()

	t.Run("metode dari settings DB - MANUAL dan NFC aktif", func(t *testing.T) {
		store := &mockSettingsStore{
			bmt: map[string]string{
				"pondok.absensi_metode": `["MANUAL","NFC"]`,
			},
			platform: map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)

		// Ambil metode dari settings — bukan dari konstanta
		var metodeDiizinkan []string
		err := resolver.ResolveJSON(context.Background(), bmtID, cabangID, "pondok.absensi_metode", &metodeDiizinkan)
		require.NoError(t, err)

		// Verifikasi nilai berasal dari DB (settings), bukan hardcode
		assert.Contains(t, metodeDiizinkan, "MANUAL")
		assert.Contains(t, metodeDiizinkan, "NFC")
		assert.NotContains(t, metodeDiizinkan, "BIOMETRIK") // tidak dikonfigurasi
	})

	t.Run("metode dari settings DB - semua metode aktif", func(t *testing.T) {
		store := &mockSettingsStore{
			bmt: map[string]string{
				"pondok.absensi_metode": `["MANUAL","NFC","BIOMETRIK"]`,
			},
			platform: map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)

		var metodeDiizinkan []string
		err := resolver.ResolveJSON(context.Background(), bmtID, cabangID, "pondok.absensi_metode", &metodeDiizinkan)
		require.NoError(t, err)

		assert.Len(t, metodeDiizinkan, 3)
		assert.Contains(t, metodeDiizinkan, string(absensi.MetodeManual))
		assert.Contains(t, metodeDiizinkan, string(absensi.MetodeNFC))
		assert.Contains(t, metodeDiizinkan, string(absensi.MetodeBiometrik))
	})

	t.Run("settings tidak ada - error bukan default hardcode", func(t *testing.T) {
		store := &mockSettingsStore{
			bmt:      map[string]string{},
			platform: map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)

		var metodeDiizinkan []string
		err := resolver.ResolveJSON(context.Background(), bmtID, cabangID, "pondok.absensi_metode", &metodeDiizinkan)

		// Jika tidak ada settings, harus error — bukan fallback ke daftar hardcode
		assert.Error(t, err, "jika settings tidak ada, harus error — bukan pakai hardcode default")
	})

	t.Run("cabang override settings BMT", func(t *testing.T) {
		store := &mockSettingsStore{
			bmt: map[string]string{
				"pondok.absensi_metode": `["MANUAL","NFC","BIOMETRIK"]`,
			},
			cabang: map[string]string{
				"pondok.absensi_metode": `["MANUAL"]`, // cabang ini hanya izinkan manual
			},
			platform: map[string]string{},
		}
		resolver := settings.NewResolver(store)

		var metodeDiizinkan []string
		err := resolver.ResolveJSON(context.Background(), bmtID, cabangID, "pondok.absensi_metode", &metodeDiizinkan)
		require.NoError(t, err)

		// Cabang override harus menang
		assert.Len(t, metodeDiizinkan, 1)
		assert.Contains(t, metodeDiizinkan, "MANUAL")
		assert.NotContains(t, metodeDiizinkan, "NFC")
	})
}

// TestAbsensi_NewAbsensiManual_Valid memastikan entitas absensi manual terbentuk benar.
func TestAbsensi_NewAbsensiManual_Valid(t *testing.T) {
	bmtID := uuid.New()
	cabangID := uuid.New()
	santriID := uuid.New()
	createdBy := uuid.New()

	a, err := absensi.NewAbsensiManual(
		bmtID, cabangID, santriID,
		absensi.TipeSubjekSantri,
		time.Now(),
		"PAGI",
		absensi.StatusHadir,
		"",
		createdBy,
	)

	require.NoError(t, err)
	assert.Equal(t, absensi.MetodeManual, a.Metode)
	assert.Equal(t, absensi.StatusHadir, a.Status)
	assert.NotNil(t, a.CreatedBy)
	assert.Equal(t, createdBy, *a.CreatedBy)
	// Absensi manual bukan NFC/biometrik, tidak ada waktu scan
	assert.Nil(t, a.WaktuScan)
}

// TestAbsensi_NewAbsensiNFC_AutoHadir memastikan absensi NFC otomatis set status HADIR.
func TestAbsensi_NewAbsensiNFC_AutoHadir(t *testing.T) {
	bmtID := uuid.New()
	cabangID := uuid.New()
	santriID := uuid.New()

	a, err := absensi.NewAbsensiNFC(
		bmtID, cabangID, santriID,
		absensi.TipeSubjekSantri,
		time.Now(),
		"PAGI",
		nil,
	)

	require.NoError(t, err)
	assert.Equal(t, absensi.MetodeNFC, a.Metode)
	assert.Equal(t, absensi.StatusHadir, a.Status)
	// NFC otomatis, tidak ada created_by
	assert.Nil(t, a.CreatedBy)
	// NFC memiliki waktu scan
	assert.NotNil(t, a.WaktuScan)
}

// TestAbsensi_StatusTidakValid memastikan status absensi yang tidak valid ditolak.
func TestAbsensi_StatusTidakValid(t *testing.T) {
	bmtID := uuid.New()
	cabangID := uuid.New()
	santriID := uuid.New()

	_, err := absensi.NewAbsensiManual(
		bmtID, cabangID, santriID,
		absensi.TipeSubjekSantri,
		time.Now(),
		"PAGI",
		absensi.StatusAbsensi("STATUS_TIDAK_VALID"),
		"",
		uuid.New(),
	)

	assert.ErrorIs(t, err, absensi.ErrStatusTidakValid)
}

// TestAbsensi_SubjekTipeTidakValid memastikan tipe subjek yang tidak valid ditolak.
func TestAbsensi_SubjekTipeTidakValid(t *testing.T) {
	_, err := absensi.NewAbsensiNFC(
		uuid.New(), uuid.New(), uuid.New(),
		absensi.TipeSubjek("TIPE_TIDAK_VALID"),
		time.Now(),
		"PAGI",
		nil,
	)

	assert.ErrorIs(t, err, absensi.ErrSubjekTidakValid)
}
