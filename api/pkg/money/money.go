package money

import (
	"fmt"
	"math"
)

// Money represents Indonesian Rupiah as int64 (in smallest unit)
// NEVER use float64 for money
type Money int64

const (
	Zero Money = 0
)

func New(amount int64) Money {
	return Money(amount)
}

func (m Money) Add(other Money) Money {
	return m + other
}

func (m Money) Sub(other Money) Money {
	return m - other
}

func (m Money) Mul(factor int64) Money {
	return Money(int64(m) * factor)
}

func (m Money) MulPercent(persen float64) Money {
	// Use integer arithmetic to avoid float precision issues
	return Money(math.Round(float64(m) * persen / 100))
}

func (m Money) IsZero() bool {
	return m == 0
}

func (m Money) IsPositive() bool {
	return m > 0
}

func (m Money) IsNegative() bool {
	return m < 0
}

func (m Money) Int64() int64 {
	return int64(m)
}

func (m Money) String() string {
	return fmt.Sprintf("Rp %d", int64(m))
}

// Format formats money as "1.000.000"
func (m Money) Format() string {
	n := int64(m)
	if n < 0 {
		return fmt.Sprintf("-%s", Money(-n).Format())
	}

	s := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}

func Min(a, b Money) Money {
	if a < b {
		return a
	}
	return b
}

func Max(a, b Money) Money {
	if a > b {
		return a
	}
	return b
}
