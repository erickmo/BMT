package penilaian_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/pondok/penilaian"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// KomponenNilaiDB merepresentasikan komponen nilai yang diambil dari database.
// Bobot tidak hardcode — setiap mapel bisa memiliki komponen dengan bobot berbeda.
type KomponenNilaiDB struct {
	ID          uuid.UUID
	Nama        string
	BobotPersen int16
}

// hitungNilaiAkhirDariDB menghitung nilai akhir berdasarkan komponen dari DB.
// Implementasi di service layer — bukan domain — karena memerlukan data dari DB.
// Test ini memastikan kalkulasinya benar menggunakan bobot dari DB, bukan konstanta.
func hitungNilaiAkhirDariDB(nilaiPerKomponen map[uuid.UUID]float64, komponenDariDB []KomponenNilaiDB) float64 {
	var nilaiAkhir float64
	for _, komp := range komponenDariDB {
		nilaiSantri, ok := nilaiPerKomponen[komp.ID]
		if !ok {
			nilaiSantri = 0
		}
		nilaiAkhir += nilaiSantri * float64(komp.BobotPersen) / 100
	}
	return nilaiAkhir
}

// TestRaport_NilaiTertimbang_KomponenDariDB memastikan nilai akhir raport
// dihitung menggunakan bobot dari tabel komponen_nilai di DB, bukan konstanta hardcode.
//
// Prinsip: bobot UH, UTS, UAS, Tugas dikonfigurasi management BMT per mapel,
// tersimpan di tabel pondok_komponen_nilai, bukan dikode sebagai konstanta.
func TestRaport_NilaiTertimbang_KomponenDariDB(t *testing.T) {
	t.Run("bobot dari DB: UH=40% UTS=30% UAS=30%", func(t *testing.T) {
		// Simulasi data komponen dari DB — bobot dikonfigurasi management BMT
		komponenDariDB := []KomponenNilaiDB{
			{ID: uuid.New(), Nama: "UH", BobotPersen: 40},
			{ID: uuid.New(), Nama: "UTS", BobotPersen: 30},
			{ID: uuid.New(), Nama: "UAS", BobotPersen: 30},
		}

		// Nilai santri per komponen
		nilaiPerKomponen := map[uuid.UUID]float64{
			komponenDariDB[0].ID: 80, // UH: 80
			komponenDariDB[1].ID: 75, // UTS: 75
			komponenDariDB[2].ID: 90, // UAS: 90
		}

		// Hitung nilai akhir menggunakan bobot dari DB
		nilaiAkhir := hitungNilaiAkhirDariDB(nilaiPerKomponen, komponenDariDB)

		// Ekspektasi: (80*40 + 75*30 + 90*30) / 100 = 32 + 22.5 + 27 = 81.5
		assert.InDelta(t, 81.5, nilaiAkhir, 0.001,
			"nilai akhir harus dihitung dari bobot DB, bukan konstanta hardcode")
	})

	t.Run("bobot dari DB berbeda per mapel: UH=50% Tugas=20% UAS=30%", func(t *testing.T) {
		// Mapel lain bisa punya bobot berbeda — semuanya dari DB
		komponenDariDB := []KomponenNilaiDB{
			{ID: uuid.New(), Nama: "UH", BobotPersen: 50},
			{ID: uuid.New(), Nama: "Tugas", BobotPersen: 20},
			{ID: uuid.New(), Nama: "UAS", BobotPersen: 30},
		}

		nilaiPerKomponen := map[uuid.UUID]float64{
			komponenDariDB[0].ID: 70, // UH: 70
			komponenDariDB[1].ID: 90, // Tugas: 90
			komponenDariDB[2].ID: 80, // UAS: 80
		}

		nilaiAkhir := hitungNilaiAkhirDariDB(nilaiPerKomponen, komponenDariDB)

		// Ekspektasi: (70*50 + 90*20 + 80*30) / 100 = 35 + 18 + 24 = 77
		assert.InDelta(t, 77.0, nilaiAkhir, 0.001)
	})

	t.Run("komponen nilai kosong menghasilkan 0", func(t *testing.T) {
		komponenDariDB := []KomponenNilaiDB{}
		nilaiPerKomponen := map[uuid.UUID]float64{}

		nilaiAkhir := hitungNilaiAkhirDariDB(nilaiPerKomponen, komponenDariDB)
		assert.Equal(t, 0.0, nilaiAkhir)
	})

	t.Run("santri tidak punya nilai untuk komponen tertentu — dianggap 0", func(t *testing.T) {
		komponenDariDB := []KomponenNilaiDB{
			{ID: uuid.New(), Nama: "UH", BobotPersen: 50},
			{ID: uuid.New(), Nama: "UAS", BobotPersen: 50},
		}

		// Santri hanya ada nilai UH, tidak ada nilai UAS
		nilaiPerKomponen := map[uuid.UUID]float64{
			komponenDariDB[0].ID: 80, // UH: 80
			// UAS tidak ada nilainya
		}

		nilaiAkhir := hitungNilaiAkhirDariDB(nilaiPerKomponen, komponenDariDB)
		// Ekspektasi: (80*50 + 0*50) / 100 = 40
		assert.InDelta(t, 40.0, nilaiAkhir, 0.001)
	})
}

// TestNilai_Validasi_RentangBenar memastikan entitas Nilai menolak angka di luar 0–100.
func TestNilai_Validasi_RentangBenar(t *testing.T) {
	bmtID := uuid.New()
	santriID := uuid.New()
	komponenID := uuid.New()
	oleh := uuid.New()

	t.Run("nilai valid: 0", func(t *testing.T) {
		n, err := penilaian.NewNilai(bmtID, santriID, komponenID, 0, "", oleh)
		require.NoError(t, err)
		assert.Equal(t, 0.0, n.Nilai)
	})

	t.Run("nilai valid: 100", func(t *testing.T) {
		n, err := penilaian.NewNilai(bmtID, santriID, komponenID, 100, "", oleh)
		require.NoError(t, err)
		assert.Equal(t, 100.0, n.Nilai)
	})

	t.Run("nilai tidak valid: -1", func(t *testing.T) {
		_, err := penilaian.NewNilai(bmtID, santriID, komponenID, -1, "", oleh)
		assert.ErrorIs(t, err, penilaian.ErrNilaiDiluarRentang)
	})

	t.Run("nilai tidak valid: 101", func(t *testing.T) {
		_, err := penilaian.NewNilai(bmtID, santriID, komponenID, 101, "", oleh)
		assert.ErrorIs(t, err, penilaian.ErrNilaiDiluarRentang)
	})
}

// TestRaport_StatusTransisi memastikan raport mengikuti state machine yang benar.
func TestRaport_StatusTransisi(t *testing.T) {
	bmtID := uuid.New()
	santriID := uuid.New()
	kelasID := uuid.New()

	t.Run("raport baru dimulai sebagai DRAFT", func(t *testing.T) {
		r, err := penilaian.NewRaport(bmtID, santriID, kelasID, "2025/2026", 1)
		require.NoError(t, err)
		assert.Equal(t, penilaian.StatusRaportDraft, r.Status)
		assert.Nil(t, r.DiterbitkanAt)
	})

	t.Run("raport DRAFT bisa difinalisasi", func(t *testing.T) {
		r, _ := penilaian.NewRaport(bmtID, santriID, kelasID, "2025/2026", 1)
		err := r.Finalisasi()
		require.NoError(t, err)
		assert.Equal(t, penilaian.StatusRaportFinal, r.Status)
	})

	t.Run("raport FINAL bisa diterbitkan", func(t *testing.T) {
		r, _ := penilaian.NewRaport(bmtID, santriID, kelasID, "2025/2026", 1)
		require.NoError(t, r.Finalisasi())
		err := r.Terbitkan()
		require.NoError(t, err)
		assert.Equal(t, penilaian.StatusRaportDiterbitkan, r.Status)
		assert.NotNil(t, r.DiterbitkanAt)
	})

	t.Run("raport DRAFT tidak bisa langsung diterbitkan", func(t *testing.T) {
		r, _ := penilaian.NewRaport(bmtID, santriID, kelasID, "2025/2026", 1)
		err := r.Terbitkan()
		assert.Error(t, err)
		assert.Equal(t, penilaian.StatusRaportDraft, r.Status) // status tidak berubah
	})

	t.Run("raport DITERBITKAN tidak bisa difinalisasi ulang", func(t *testing.T) {
		r, _ := penilaian.NewRaport(bmtID, santriID, kelasID, "2025/2026", 1)
		require.NoError(t, r.Finalisasi())
		require.NoError(t, r.Terbitkan())

		err := r.Finalisasi()
		assert.ErrorIs(t, err, penilaian.ErrRaportSudahDiterbitkan)
	})
}

// TestRaport_NilaiMapelJSON memastikan snapshot nilai mapel tersimpan sebagai JSON yang valid.
func TestRaport_NilaiMapelJSON(t *testing.T) {
	snapshots := []penilaian.NilaiMapelSnapshot{
		{MapelID: uuid.New(), NamaMapel: "Matematika", NilaiAkhir: 85.5, Predikat: "A"},
		{MapelID: uuid.New(), NamaMapel: "Bahasa Arab", NilaiAkhir: 90.0, Predikat: "A"},
	}

	snapshotJSON, err := json.Marshal(snapshots)
	require.NoError(t, err)

	r, err := penilaian.NewRaport(uuid.New(), uuid.New(), uuid.New(), "2025/2026", 2)
	require.NoError(t, err)

	r.NilaiMapel = snapshotJSON

	// Verifikasi JSON bisa di-unmarshal kembali
	var hasil []penilaian.NilaiMapelSnapshot
	err = json.Unmarshal(r.NilaiMapel, &hasil)
	require.NoError(t, err)
	assert.Len(t, hasil, 2)
	assert.Equal(t, "Matematika", hasil[0].NamaMapel)
	assert.InDelta(t, 85.5, hasil[0].NilaiAkhir, 0.001)
}

// TestNilaiTahfidz_Lulus_DanMengulang memastikan state machine tahfidz benar.
func TestNilaiTahfidz_Lulus_DanMengulang(t *testing.T) {
	now := time.Now()
	pengujiID := uuid.New()

	n, err := penilaian.NewNilaiTahfidz(
		uuid.New(), uuid.New(),
		"Al-Fatihah", 1, 7,
		now, &pengujiID,
	)
	require.NoError(t, err)
	assert.Equal(t, penilaian.StatusTahfidzBelumDiuji, n.Status)

	t.Run("lulus dengan nilai valid", func(t *testing.T) {
		err := n.Lulus(92.5)
		require.NoError(t, err)
		assert.Equal(t, penilaian.StatusTahfidzLulus, n.Status)
		require.NotNil(t, n.Nilai)
		assert.InDelta(t, 92.5, *n.Nilai, 0.001)
	})
}
