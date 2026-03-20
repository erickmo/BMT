package money_test

import (
	"testing"

	"github.com/bmt-saas/api/pkg/money"
	"github.com/stretchr/testify/assert"
)

func TestMoney_Add(t *testing.T) {
	a := money.New(100000)
	b := money.New(50000)
	assert.Equal(t, money.New(150000), a.Add(b))
}

func TestMoney_Sub(t *testing.T) {
	a := money.New(100000)
	b := money.New(30000)
	assert.Equal(t, money.New(70000), a.Sub(b))
}

func TestMoney_MulPercent(t *testing.T) {
	m := money.New(1000000)
	result := m.MulPercent(2.5)
	assert.Equal(t, money.New(25000), result)
}

func TestMoney_Min(t *testing.T) {
	a := money.New(100000)
	b := money.New(50000)
	assert.Equal(t, b, money.Min(a, b))
}

func TestMoney_Format(t *testing.T) {
	m := money.New(1500000)
	assert.Equal(t, "1.500.000", m.Format())
}

func TestMoney_Mul(t *testing.T) {
	m := money.New(100000)
	assert.Equal(t, money.New(300000), m.Mul(3))
	assert.Equal(t, money.New(0), m.Mul(0))
	assert.Equal(t, money.New(-100000), m.Mul(-1))
}

func TestMoney_Max(t *testing.T) {
	a := money.New(100000)
	b := money.New(200000)
	assert.Equal(t, b, money.Max(a, b))
	assert.Equal(t, a, money.Max(a, money.New(50000)))
	assert.Equal(t, a, money.Max(a, a)) // sama besar
}

func TestMoney_IsZero(t *testing.T) {
	assert.True(t, money.Zero.IsZero())
	assert.True(t, money.New(0).IsZero())
	assert.False(t, money.New(1).IsZero())
	assert.False(t, money.New(-1).IsZero())
}

func TestMoney_IsPositive(t *testing.T) {
	assert.True(t, money.New(1).IsPositive())
	assert.True(t, money.New(1000000).IsPositive())
	assert.False(t, money.Zero.IsPositive())
	assert.False(t, money.New(-1).IsPositive())
}

func TestMoney_IsNegative(t *testing.T) {
	assert.True(t, money.New(-1).IsNegative())
	assert.True(t, money.New(-1000000).IsNegative())
	assert.False(t, money.Zero.IsNegative())
	assert.False(t, money.New(1).IsNegative())
}

func TestMoney_String(t *testing.T) {
	assert.Equal(t, "Rp 500000", money.New(500000).String())
	assert.Equal(t, "Rp 0", money.Zero.String())
	assert.Equal(t, "Rp -100000", money.New(-100000).String())
}

func TestMoney_Format_Negatif(t *testing.T) {
	assert.Equal(t, "-500.000", money.New(-500000).Format())
	assert.Equal(t, "-1.000.000", money.New(-1000000).Format())
	assert.Equal(t, "0", money.New(0).Format())
	assert.Equal(t, "999", money.New(999).Format())
	assert.Equal(t, "1.000", money.New(1000).Format())
}
