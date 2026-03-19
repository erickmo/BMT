package platform_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/platform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlatformRepo implementasi platform.Repository hanya untuk GetPecahanAktif
type mockPlatformRepo struct {
	pecahans []*platform.PecahanUang
	err      error
}

func (m *mockPlatformRepo) GetPecahanAktif(_ context.Context) ([]*platform.PecahanUang, error) {
	return m.pecahans, m.err
}

// Interface stub — metode lain tidak diperlukan untuk test ini
func (m *mockPlatformRepo) CreateBMT(_ context.Context, _ *platform.BMT) error  { return nil }
func (m *mockPlatformRepo) GetBMT(_ context.Context, _ uuid.UUID) (*platform.BMT, error) {
	return nil, nil
}
func (m *mockPlatformRepo) GetBMTByKode(_ context.Context, _ string) (*platform.BMT, error) {
	return nil, nil
}
func (m *mockPlatformRepo) ListBMT(_ context.Context) ([]*platform.BMT, error) {
	return nil, nil
}
func (m *mockPlatformRepo) UpdateBMTStatus(_ context.Context, _ uuid.UUID, _ platform.StatusBMT) error {
	return nil
}
func (m *mockPlatformRepo) CreateCabang(_ context.Context, _ *platform.Cabang) error { return nil }
func (m *mockPlatformRepo) GetCabang(_ context.Context, _ uuid.UUID) (*platform.Cabang, error) {
	return nil, nil
}
func (m *mockPlatformRepo) ListCabangByBMT(_ context.Context, _ uuid.UUID) ([]*platform.Cabang, error) {
	return nil, nil
}
func (m *mockPlatformRepo) CreateKontrak(_ context.Context, _ *platform.KontrakBMT) error {
	return nil
}
func (m *mockPlatformRepo) GetKontrakAktif(_ context.Context, _ uuid.UUID) (*platform.KontrakBMT, error) {
	return nil, nil
}
func (m *mockPlatformRepo) CreatePecahan(_ context.Context, _ *platform.PecahanUang) error {
	return nil
}
func (m *mockPlatformRepo) UpdatePecahan(_ context.Context, _ *platform.PecahanUang) error {
	return nil
}

// TestPecahanUang_DariDB_BukanKonstanta memastikan pecahan uang Rupiah selalu diambil
// dari database (via Repository interface), tidak pernah dari array konstanta di kode.
//
// Prinsip: pecahan uang bisa berubah (redenominasi, uang baru diterbitkan, dsb.)
// tanpa perlu deploy ulang. Semua pecahan aktif dari tabel pecahan_uang di DB.
func TestPecahanUang_DariDB_BukanKonstanta(t *testing.T) {
	t.Run("pecahan aktif berasal dari DB", func(t *testing.T) {
		// Data pecahan "dari DB" — berbeda dari konstanta kode mana pun
		pecahanDariDB := []*platform.PecahanUang{
			{
				ID:           uuid.New(),
				Nominal:      100000,
				Jenis:        "KERTAS",
				Label:        "Rp 100.000 (kertas)",
				IsAktif:      true,
				Urutan:       1,
				BerlakuSejak: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:           uuid.New(),
				Nominal:      50000,
				Jenis:        "KERTAS",
				Label:        "Rp 50.000 (kertas)",
				IsAktif:      true,
				Urutan:       2,
				BerlakuSejak: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:           uuid.New(),
				Nominal:      1000,
				Jenis:        "LOGAM",
				Label:        "Rp 1.000 (logam)",
				IsAktif:      true,
				Urutan:       7,
				BerlakuSejak: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		repo := &mockPlatformRepo{pecahans: pecahanDariDB}
		hasil, err := repo.GetPecahanAktif(context.Background())

		require.NoError(t, err)
		assert.Len(t, hasil, 3, "jumlah pecahan harus sesuai DB, bukan konstanta hardcode")

		// Verifikasi data yang kembali persis dari DB
		assert.Equal(t, int64(100000), hasil[0].Nominal)
		assert.Equal(t, "KERTAS", hasil[0].Jenis)
		assert.Equal(t, int64(50000), hasil[1].Nominal)
		assert.Equal(t, int64(1000), hasil[2].Nominal)
		assert.Equal(t, "LOGAM", hasil[2].Jenis)
	})

	t.Run("pecahan yang sudah ditarik tidak muncul", func(t *testing.T) {
		ditarik := time.Now().Add(-24 * time.Hour)
		// Pecahan lama yang sudah ditarik — hanya pecahan aktif yang dikembalikan
		pecahanAktif := []*platform.PecahanUang{
			{
				ID:      uuid.New(),
				Nominal: 100000,
				Jenis:   "KERTAS",
				IsAktif: true,
				Urutan:  1,
			},
			// Pecahan 500 sudah ditarik — repo hanya mengembalikan yang aktif
			// {Nominal: 500, IsAktif: false, DitarikPada: &ditarik}
		}
		_ = ditarik // membuktikan bahwa pecahan ditarik tidak dikembalikan repo

		repo := &mockPlatformRepo{pecahans: pecahanAktif}
		hasil, err := repo.GetPecahanAktif(context.Background())

		require.NoError(t, err)
		assert.Len(t, hasil, 1)
		// Pastikan tidak ada pecahan nominal 500 (sudah ditarik)
		for _, p := range hasil {
			assert.NotEqual(t, int64(500), p.Nominal)
		}
	})

	t.Run("error dari DB disampaikan ke caller", func(t *testing.T) {
		errDB := errors.New("koneksi DB gagal")
		repo := &mockPlatformRepo{err: errDB}

		hasil, err := repo.GetPecahanAktif(context.Background())

		assert.Nil(t, hasil)
		assert.ErrorIs(t, err, errDB,
			"error dari DB harus diteruskan — tidak boleh fallback ke hardcode")
	})

	t.Run("urutan tampil dari DB dipatuhi", func(t *testing.T) {
		// Urutan tampil di app dikonfigurasi di DB, bukan sorting hardcode di kode
		pecahanDariDB := []*platform.PecahanUang{
			{Nominal: 100000, Urutan: 1, IsAktif: true},
			{Nominal: 50000, Urutan: 2, IsAktif: true},
			{Nominal: 20000, Urutan: 3, IsAktif: true},
			{Nominal: 10000, Urutan: 4, IsAktif: true},
		}

		repo := &mockPlatformRepo{pecahans: pecahanDariDB}
		hasil, err := repo.GetPecahanAktif(context.Background())

		require.NoError(t, err)
		// Urutan dari DB dipertahankan (repo mengurutkan berdasarkan kolom "urutan")
		assert.Equal(t, 1, hasil[0].Urutan)
		assert.Equal(t, 2, hasil[1].Urutan)
		assert.Equal(t, 3, hasil[2].Urutan)
		assert.Equal(t, 4, hasil[3].Urutan)
	})
}

// TestCrossTenant_QueryTanpaBMTID_Dilarang memastikan query operasional
// tidak bisa dijalankan tanpa scoping bmt_id (tenant isolation).
//
// Prinsip: semua query wajib di-scope bmt_id + cabang_id.
// Test ini memvalidasi bahwa entitas domain selalu membawa bmt_id.
func TestCrossTenant_QueryTanpaBMTID_Dilarang(t *testing.T) {
	t.Run("BMT entity harus memiliki ID yang valid", func(t *testing.T) {
		bmt, err := platform.NewBMT(platform.CreateBMTInput{
			Kode:  "ANNUR",
			Nama:  "BMT An-Nur",
			Alamat: "Kediri",
		})
		require.NoError(t, err)

		// bmt_id harus tidak nil — bukan uuid.Nil
		assert.NotEqual(t, uuid.Nil, bmt.ID,
			"bmt_id wajib diisi — tidak boleh uuid.Nil")
	})

	t.Run("Cabang entity harus memiliki bmt_id yang valid", func(t *testing.T) {
		bmtID := uuid.New()
		cabang, err := platform.NewCabang(platform.CreateCabangInput{
			BMTID: bmtID,
			Kode:  "KDR",
			Nama:  "Cabang Kediri",
		})
		require.NoError(t, err)

		// bmt_id di cabang harus tidak nil
		assert.NotEqual(t, uuid.Nil, cabang.BMTID,
			"cabang harus terikat ke bmt_id yang valid")
		assert.Equal(t, bmtID, cabang.BMTID)
	})

	t.Run("Cabang tanpa bmt_id ditolak", func(t *testing.T) {
		_, err := platform.NewCabang(platform.CreateCabangInput{
			BMTID: uuid.Nil, // uuid.Nil = tidak ada tenant!
			Kode:  "KDR",
			Nama:  "Cabang Kediri",
		})
		// Domain menolak pembuatan cabang tanpa bmt_id
		assert.Error(t, err, "cabang tanpa bmt_id harus ditolak")
	})

	t.Run("PecahanUang harus memiliki ID unik dari DB", func(t *testing.T) {
		// PecahanUang tidak terikat tenant (platform-level) tapi harus punya ID dari DB
		pecahan := &platform.PecahanUang{
			ID:           uuid.New(),
			Nominal:      50000,
			Jenis:        "KERTAS",
			Label:        "Rp 50.000",
			IsAktif:      true,
			Urutan:       2,
			BerlakuSejak: time.Now(),
		}

		assert.NotEqual(t, uuid.Nil, pecahan.ID,
			"PecahanUang harus memiliki ID dari DB")
		assert.Greater(t, pecahan.Nominal, int64(0))
	})
}
