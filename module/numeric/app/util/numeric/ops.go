package numeric

import "math"

func (n Numeric) Add(x Numeric) Numeric {
	if isZero(n.mantissa) {
		return x
	}
	if isZero(x.mantissa) {
		return n
	}
	if n.IsNaN() || x.IsNaN() || n.IsInf() || x.IsInf() {
		return Numeric{mantissa: n.mantissa + x.mantissa, exponent: 0}
	}
	var bigger, smaller Numeric
	if n.exponent >= x.exponent {
		bigger, smaller = n, x
	} else {
		bigger, smaller = x, n
	}
	if bigger.exponent-smaller.exponent > MaxSignificantDigits {
		return bigger
	}
	return normalize(math.Round(1e14*bigger.mantissa+1e14*smaller.mantissa*powerOf10(smaller.exponent-bigger.exponent)), bigger.exponent-14)
}

func (n Numeric) Subtract(x Numeric) Numeric {
	return n.Add(x.Negate())
}

func (n Numeric) Multiply(x Numeric) Numeric {
	return normalize(n.mantissa*x.mantissa, n.exponent+x.exponent)
}

func (n Numeric) Divide(x Numeric) Numeric {
	return n.Multiply(x.Reciprocate())
}

func (n Numeric) Reciprocate() Numeric {
	return normalize(1.0/n.mantissa, -n.exponent)
}

func (n Numeric) Log(base Numeric) float64 {
	if base.IsZero() {
		return math.NaN()
	}
	return 2.30258509299404568402 / math.Log(base.ToFloat()) * n.Log10()
}

func (n Numeric) Log2() float64 {
	return 3.32192809488736234787 * n.Log10()
}

func (n Numeric) Ln() float64 {
	return 2.30258509299404568402 * n.Log10()
}

func (n Numeric) Log10() float64 {
	return float64(n.exponent) + math.Log10(n.mantissa)
}

func (n Numeric) AbsLog10() float64 {
	return float64(n.exponent) + math.Log10(math.Abs(n.mantissa))
}

func (n Numeric) Pow(other Numeric) Numeric {
	//UN-SAFETY: Accuracy not guaranteed beyond ~9~11 decimal places.

	//Fast track: If (this.exponent*value) is an integer and mantissa^value fits in a Number, we can do a very fast method.
	temp := float64(n.exponent) * other.ToFloat()
	if isInteger(temp) && !math.IsInf(temp, 0) && math.Abs(temp) < ExpLimit {
		newMantissa := math.Pow(n.mantissa, other.ToFloat())
		if !math.IsInf(newMantissa, 0) {
			return normalize(newMantissa, int64(temp))
		}
	}

	//Same speed and usually more accurate. (An arbitrary-precision version of this calculation is used in break_break_infinity.js, sacrificing performance for utter accuracy.)
	var newExponent = math.Trunc(temp)
	var residue = temp - newExponent
	newMantissa := math.Pow(10, other.ToFloat()*math.Log10(n.mantissa)+residue)
	if !math.IsInf(newMantissa, 0) {
		return normalize(newMantissa, int64(newExponent))
	}

	//UN-SAFETY: This should return NaN when mantissa is negative and value is noninteger.
	var result = other.Multiply(FromFloat(n.AbsLog10())).Pow10() //this is 2x faster and gives same values AFAIK
	if n.Sign() == -1 && FromFloat(math.Mod(other.ToFloat(), 2)).Equals(One) {
		return result.Negate()
	}
	return result
}

func (n Numeric) Pow10() Numeric {
	if n.IsInt() {
		return normalize(1, n.ToInt())
	}
	return normalize(math.Pow(10, math.Mod(n.ToFloat(), 1)), n.ToInt())
}

func isInteger(x float64) bool {
	return x == float64(int64(x))
}
