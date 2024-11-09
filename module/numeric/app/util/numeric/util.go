package numeric

import (
	"math"
	"strconv"
)

const (
	Tolerance = 1e-18

	MaxSignificantDigits = 17

	ExpLimit = math.MaxInt64

	FloatExpMin int64 = -324
	FloatExpMax int64 = 308

	Epsilon = float64(7.0)/3 - float64(4.0)/3 - float64(1.0)
)

var (
	Zero             = from(0, 0)
	One              = from(1, 0)
	Ten              = from(1, 1)
	NaN              = from(math.NaN(), 0)
	PositiveInfinity = from(math.Inf(1), 0)
	NegativeInfinity = from(math.Inf(-1), 0)
)

func areSameInfinity(a, b Numeric) bool {
	return (math.IsInf(a.mantissa, 1) && math.IsInf(b.mantissa, 1)) || (math.IsInf(a.mantissa, -1) && math.IsInf(b.mantissa, -1))
}

func areEqual(a, b float64) bool {
	return math.Abs(a-b) < Tolerance
}

func normalize(mantissa float64, exponent int64) Numeric {
	if (mantissa >= 1 && mantissa < 10) || !isFinite(mantissa) {
		return from(mantissa, exponent)
	}
	if isZero(mantissa) {
		return Zero
	}

	tempExponent := int64(math.Floor(math.Log10(math.Abs(mantissa))))

	if tempExponent == FloatExpMin {
		mantissa = mantissa * 10 / 1e-323
	} else {
		mantissa = mantissa / powerOf10(tempExponent)
	}

	return from(mantissa, exponent+tempExponent)
}

func isZero(x float64) bool {
	return math.Abs(x) < Tolerance
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func mapNumeric(n Numeric, zero func() Numeric, fn func(f float64) float64) Numeric {
	if n.IsNaN() {
		return n
	}
	if n.exponent < -1 {
		return zero()
	}
	if n.exponent < MaxSignificantDigits {
		return FromFloat(fn(n.ToFloat()))
	}
	return n
}

var powersOf10 = func() []float64 {
	ret := make([]float64, FloatExpMax-FloatExpMin)
	indexOf0 := -FloatExpMin - 1
	for i := 0; i < len(ret); i++ {
		val, _ := strconv.ParseFloat("1e"+strconv.Itoa(i-int(indexOf0)), 64)
		ret[i] = val
	}
	return ret
}()

func powerOf10(power int64) float64 {
	indexOf0 := -FloatExpMin - 1
	return powersOf10[indexOf0+power]
}
