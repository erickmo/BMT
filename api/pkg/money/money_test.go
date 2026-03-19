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
