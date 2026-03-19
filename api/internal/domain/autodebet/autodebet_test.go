package autodebet_test

import (
	"testing"

	"github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/stretchr/testify/assert"
)

func TestAutodebet_SaldoCukup_FullDebit(t *testing.T) {
	saldo := money.New(500000)
	target := money.New(50000)

	hasil := autodebet.EksekusiAutodebet(saldo, target)

	assert.Equal(t, target, hasil.NominalDidebit)
	assert.Equal(t, money.Zero, hasil.NominalTunggakan)
	assert.False(t, hasil.IsPartial)
}

func TestAutodebet_SaldoKurang_PartialDebit(t *testing.T) {
	saldo := money.New(30000)
	target := money.New(50000)

	hasil := autodebet.EksekusiAutodebet(saldo, target)

	assert.Equal(t, money.New(30000), hasil.NominalDidebit)
	assert.Equal(t, money.New(20000), hasil.NominalTunggakan)
	assert.True(t, hasil.IsPartial)
}

func TestAutodebet_SaldoNol_SemuaTunggakan(t *testing.T) {
	saldo := money.Zero
	target := money.New(50000)

	hasil := autodebet.EksekusiAutodebet(saldo, target)

	assert.Equal(t, money.Zero, hasil.NominalDidebit)
	assert.Equal(t, target, hasil.NominalTunggakan)
	assert.True(t, hasil.IsPartial)
}

func TestAutodebet_SaldoLebih_FullDebit(t *testing.T) {
	saldo := money.New(1000000)
	target := money.New(50000)

	hasil := autodebet.EksekusiAutodebet(saldo, target)

	assert.Equal(t, target, hasil.NominalDidebit)
	assert.Equal(t, money.Zero, hasil.NominalTunggakan)
	assert.False(t, hasil.IsPartial)
}

// TestAutodebet_SaldoKurang_PartialDebitDanTunggakan adalah test wajib per CLAUDE.md.
// Memastikan ketika saldo tidak mencukupi, sistem melakukan partial debit (sebesar saldo)
// dan sisanya disimpan sebagai tunggakan — bukan skip atau error.
func TestAutodebet_SaldoKurang_PartialDebitDanTunggakan(t *testing.T) {
	t.Run("saldo 30000, target 50000: debit 30000, tunggakan 20000", func(t *testing.T) {
		saldo := money.New(30000)
		target := money.New(50000)

		hasil := autodebet.EksekusiAutodebet(saldo, target)

		// Debit hanya sebesar saldo yang ada
		assert.Equal(t, money.New(30000), hasil.NominalDidebit,
			"partial debit harus sebesar saldo tersedia")

		// Sisa menjadi tunggakan
		assert.Equal(t, money.New(20000), hasil.NominalTunggakan,
			"sisa kewajiban harus masuk tunggakan, bukan dibuang atau di-skip")

		// IsPartial harus true — ini signal untuk INSERT ke tabel tunggakan_autodebet
		assert.True(t, hasil.IsPartial,
			"IsPartial harus true sebagai signal untuk mencatat tunggakan di DB")
	})

	t.Run("saldo 1, target 100000: debit 1, tunggakan 99999", func(t *testing.T) {
		saldo := money.New(1)
		target := money.New(100000)

		hasil := autodebet.EksekusiAutodebet(saldo, target)

		assert.Equal(t, money.New(1), hasil.NominalDidebit)
		assert.Equal(t, money.New(99999), hasil.NominalTunggakan)
		assert.True(t, hasil.IsPartial)
	})

	t.Run("saldo 0, target berapapun: tidak ada debit, semua jadi tunggakan", func(t *testing.T) {
		saldo := money.Zero
		target := money.New(75000)

		hasil := autodebet.EksekusiAutodebet(saldo, target)

		assert.Equal(t, money.Zero, hasil.NominalDidebit,
			"saldo 0 tidak boleh menghasilkan debit")
		assert.Equal(t, target, hasil.NominalTunggakan,
			"jika saldo 0, seluruh nominal menjadi tunggakan")
		assert.True(t, hasil.IsPartial)
	})

	t.Run("nominal target harus tetap sama meski partial", func(t *testing.T) {
		saldo := money.New(20000)
		target := money.New(50000)

		hasil := autodebet.EksekusiAutodebet(saldo, target)

		// Debit + tunggakan harus sama dengan target asli
		total := hasil.NominalDidebit.Add(hasil.NominalTunggakan)
		assert.Equal(t, target, total,
			"debit + tunggakan harus = nominal target asli")
	})
}
