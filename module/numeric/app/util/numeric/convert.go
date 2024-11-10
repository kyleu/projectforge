package numeric

import "math"

func (n Numeric) ToFloat() float64 {
	if math.IsNaN(n.mantissa) {
		return math.NaN()
	}
	if n.exponent > FloatExpMax {
		if n.mantissa > 0 {
			return math.Inf(1)
		}
		return math.Inf(-1)
	}
	if n.exponent < FloatExpMin {
		return 0.0
	}
	if n.exponent == FloatExpMin {
		if n.mantissa > 0 {
			return 5e-324
		}
		return -5e-324
	}
	result := n.mantissa * math.Pow10(int(n.exponent))
	if !isFinite(result) || n.exponent < 0 {
		return result
	}
	rounded := math.Round(result)
	if math.Abs(rounded-result) < 1e-10 {
		return rounded
	}
	return result
}

func (n Numeric) ToInt() int {
	return int(n.ToFloat())
}

func (n Numeric) ToInt64() int64 {
	return int64(n.ToFloat())
}
