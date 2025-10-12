package numeric

import (
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (n NumericMap) LessThan(x NumericMap) bool {
	for k, v := range n {
		amt, ok := x[k]
		if !ok && v.GreaterThan(Zero) {
			return false
		}
		if amt.LessThan(v) {
			return false
		}
	}
	return true
}

func (n NumericMap) Equal(expected NumericMap) bool {
	for k, v := range n {
		amt, ok := expected[k]
		if !ok && v.GreaterThan(Zero) {
			return false
		}
		if !amt.Equals(v) {
			return false
		}
	}
	for k, v := range expected {
		_, ok := n[k]
		if !ok && v.GreaterThan(Zero) {
			return false
		}
	}
	return true
}

func (n NumericMap) Add(x NumericMap) NumericMap {
	return n.merge(x, func(v Numeric, k string) Numeric {
		return v.Add(x.GetOr(k, Zero))
	})
}

func (n NumericMap) AddNumeric(x Numeric) NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.Add(x)
	})
}

func (n NumericMap) Subtract(x NumericMap) NumericMap {
	return n.merge(x, func(v Numeric, k string) Numeric {
		return v.Subtract(x.GetOr(k, Zero))
	})
}

func (n NumericMap) SubtractNumeric(x Numeric) NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.Subtract(x)
	})
}

func (n NumericMap) Multiply(x NumericMap) NumericMap {
	return n.merge(x, func(v Numeric, k string) Numeric {
		return v.Multiply(x.GetOr(k, Zero))
	})
}

func (n NumericMap) MultiplyNumeric(x Numeric) NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.Multiply(x)
	})
}

func (n NumericMap) Divide(x NumericMap) NumericMap {
	return n.merge(x, func(v Numeric, k string) Numeric {
		return v.Divide(x.GetOr(k, One))
	})
}

func (n NumericMap) DivideNumeric(x Numeric) NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.Divide(x)
	})
}

func (n NumericMap) Ceiling() NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.Ceiling()
	})
}

func (n NumericMap) Floor() NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.Floor()
	})
}

func (n NumericMap) Min() Numeric {
	return lo.Reduce(lo.Values(n), func(ret Numeric, v Numeric, _ int) Numeric {
		lt := v.LessThan(ret)
		return util.Choose(lt, v, ret)
	}, PositiveInfinity)
}

func (n NumericMap) Max() Numeric {
	return lo.Reduce(lo.Values(n), func(ret Numeric, v Numeric, _ int) Numeric {
		return util.Choose(v.GreaterThan(ret), v, ret)
	}, NegativeInfinity)
}

func (n NumericMap) LnN() NumericMap {
	return lo.MapValues(n, func(v Numeric, _ string) Numeric {
		return v.LnN()
	})
}
