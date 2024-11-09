package numeric

func (n Numeric) LessThan(x Numeric) bool {
	if n.IsNaN() || x.IsNaN() {
		return false
	}
	if isZero(n.mantissa) {
		return x.mantissa > 0
	}
	if isZero(x.mantissa) {
		return n.mantissa < 0
	}
	if n.exponent == x.exponent {
		return n.mantissa < x.mantissa
	}
	if n.mantissa > 0 {
		return x.mantissa > 0 && n.exponent < x.exponent
	}
	return x.mantissa > 0 || n.exponent > x.exponent
}

func (n Numeric) GreaterThan(x Numeric) bool {
	if n.IsNaN() || x.IsNaN() {
		return false
	}
	if isZero(n.mantissa) {
		return x.mantissa < 0
	}
	if isZero(x.mantissa) {
		return n.mantissa > 0
	}
	if n.exponent == x.exponent {
		return n.mantissa > x.mantissa
	}
	if n.mantissa > 0 {
		return x.mantissa < 0 || n.exponent > x.exponent
	}
	return x.mantissa < 0 && n.exponent < x.exponent
}

func (n Numeric) Min(x Numeric) Numeric {
	if n.IsNaN() || x.IsNaN() {
		return NaN
	}
	if n.GreaterThan(x) {
		return x
	}
	return n
}

func (n Numeric) Max(x Numeric) Numeric {
	if n.IsNaN() || x.IsNaN() {
		return NaN
	}
	if n.GreaterThan(x) {
		return n
	}
	return x
}
