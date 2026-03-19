package settings_test

import (
	"context"
	"testing"

	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockStore implements settings.Store
type MockStore struct {
	platform map[string]string
	bmt      map[string]string
	cabang   map[string]string
}

func (m *MockStore) GetPlatform(ctx context.Context, kunci string) (string, error) {
	val, ok := m.platform[kunci]
	if !ok {
		return "", assert.AnError
	}
	return val, nil
}

func (m *MockStore) GetBMT(ctx context.Context, bmtID uuid.UUID, kunci string) (string, error) {
	val, ok := m.bmt[kunci]
	if !ok {
		return "", assert.AnError
	}
	return val, nil
}

func (m *MockStore) GetCabang(ctx context.Context, cabangID uuid.UUID, kunci string) (string, error) {
	val, ok := m.cabang[kunci]
	if !ok {
		return "", assert.AnError
	}
	return val, nil
}

func TestResolver_CabangMengoverrideBMT(t *testing.T) {
	store := &MockStore{
		platform: map[string]string{"operasional.jam_buka": "08:00"},
		bmt:      map[string]string{"operasional.jam_buka": "09:00"},
		cabang:   map[string]string{"operasional.jam_buka": "07:30"},
	}

	resolver := settings.NewResolver(store)
	val := resolver.Resolve(context.Background(), uuid.New(), uuid.New(), "operasional.jam_buka")

	// Cabang harus menang
	assert.Equal(t, "07:30", val)
}

func TestResolver_BMTMengoverridePlatform(t *testing.T) {
	store := &MockStore{
		platform: map[string]string{"operasional.jam_buka": "08:00"},
		bmt:      map[string]string{"operasional.jam_buka": "09:00"},
		cabang:   map[string]string{},
	}

	resolver := settings.NewResolver(store)
	val := resolver.Resolve(context.Background(), uuid.New(), uuid.New(), "operasional.jam_buka")

	// BMT harus menang atas platform
	assert.Equal(t, "09:00", val)
}

func TestResolver_FallbackKePlatform(t *testing.T) {
	store := &MockStore{
		platform: map[string]string{"operasional.jam_buka": "08:00"},
		bmt:      map[string]string{},
		cabang:   map[string]string{},
	}

	resolver := settings.NewResolver(store)
	val := resolver.Resolve(context.Background(), uuid.New(), uuid.New(), "operasional.jam_buka")

	assert.Equal(t, "08:00", val)
}

func TestResolver_TidakAdaSettings_ReturnsEmpty(t *testing.T) {
	store := &MockStore{
		platform: map[string]string{},
		bmt:      map[string]string{},
		cabang:   map[string]string{},
	}

	resolver := settings.NewResolver(store)
	val := resolver.Resolve(context.Background(), uuid.New(), uuid.New(), "tidak.ada")

	assert.Equal(t, "", val)
}

func TestResolver_ResolveInt_Default(t *testing.T) {
	store := &MockStore{platform: map[string]string{}, bmt: map[string]string{}, cabang: map[string]string{}}
	resolver := settings.NewResolver(store)

	val := resolver.ResolveInt(context.Background(), uuid.New(), uuid.New(), "tidak.ada", 5)
	assert.Equal(t, 5, val)
}

func TestResolver_ResolveBool_True(t *testing.T) {
	store := &MockStore{
		platform: map[string]string{"fitur.aktif": "true"},
		bmt:      map[string]string{},
		cabang:   map[string]string{},
	}
	resolver := settings.NewResolver(store)

	val := resolver.ResolveBool(context.Background(), uuid.New(), uuid.New(), "fitur.aktif", false)
	assert.True(t, val)
}

// TestSettings_TidakAdaHardcode_SelaluDariDB memastikan settings engine selalu membaca
// nilai dari DB (via Store interface), tidak pernah dari konstanta hardcode di kode.
//
// Ini adalah test wajib per CLAUDE.md untuk membuktikan prinsip "Settings over hardcode".
func TestSettings_TidakAdaHardcode_SelaluDariDB(t *testing.T) {
	t.Run("jam_buka diambil dari DB bukan konstanta", func(t *testing.T) {
		// Jika diset berbeda di DB, nilai dari DB yang digunakan
		store := &MockStore{
			bmt:      map[string]string{"operasional.jam_buka": "06:30"},
			platform: map[string]string{"operasional.jam_buka": "08:00"},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)
		val := resolver.Resolve(context.Background(), uuid.New(), uuid.New(), "operasional.jam_buka")

		// Nilai dari DB (BMT settings), bukan dari kode
		assert.Equal(t, "06:30", val, "jam buka harus dari DB, bukan hardcode")
		assert.NotEqual(t, "08:00", val, "nilai platform tidak boleh override BMT")
	})

	t.Run("autodebet_tanggal dari DB bukan konstanta", func(t *testing.T) {
		store := &MockStore{
			bmt:      map[string]string{"autodebet.tanggal_simpanan_wajib": "5"},
			platform: map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)
		val := resolver.ResolveInt(context.Background(), uuid.New(), uuid.New(), "autodebet.tanggal_simpanan_wajib", 1)

		// Tanggal autodebet dari DB, bukan dari angka hardcode "1"
		assert.Equal(t, 5, val, "tanggal autodebet harus dari DB")
	})

	t.Run("absensi_metode dari DB bukan array hardcode", func(t *testing.T) {
		store := &MockStore{
			bmt:      map[string]string{"pondok.absensi_metode": `["MANUAL","NFC"]`},
			platform: map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)

		var metode []string
		err := resolver.ResolveJSON(context.Background(), uuid.New(), uuid.New(), "pondok.absensi_metode", &metode)
		assert.NoError(t, err)

		// Metode dari DB, bukan array hardcode ["MANUAL", "NFC", "BIOMETRIK"]
		assert.Len(t, metode, 2, "jumlah metode harus sesuai DB, bukan hardcode")
		assert.Contains(t, metode, "MANUAL")
		assert.Contains(t, metode, "NFC")
		assert.NotContains(t, metode, "BIOMETRIK")
	})

	t.Run("approval_chain dari DB per jenis form", func(t *testing.T) {
		store := &MockStore{
			bmt: map[string]string{
				`approval.FORM_BUKA_REKENING`: `["TELLER","MANAJER_CABANG"]`,
			},
			platform: map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)

		var approvers []string
		err := resolver.ResolveJSON(context.Background(), uuid.New(), uuid.New(), "approval.FORM_BUKA_REKENING", &approvers)
		assert.NoError(t, err)

		// Chain approver dari DB — bisa berbeda per BMT
		assert.Len(t, approvers, 2)
		assert.Equal(t, "TELLER", approvers[0])
		assert.Equal(t, "MANAJER_CABANG", approvers[1])
	})

	t.Run("nilai kosong jika kunci tidak ada di DB — bukan default hardcode", func(t *testing.T) {
		store := &MockStore{
			platform: map[string]string{},
			bmt:      map[string]string{},
			cabang:   map[string]string{},
		}
		resolver := settings.NewResolver(store)

		// Jika kunci tidak ada di DB, Resolve mengembalikan string kosong
		val := resolver.Resolve(context.Background(), uuid.New(), uuid.New(), "kunci.tidak.ada.di.db")
		assert.Equal(t, "", val, "harus kosong jika tidak ada di DB — bukan hardcode default")
	})
}
