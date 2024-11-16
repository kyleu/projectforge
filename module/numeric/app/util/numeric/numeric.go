package numeric

import (
	"fmt"
	"math"
	"strings"

	"{{{ .Package }}}/app/util"
)

type Numeric struct {
	mantissa float64
	exponent int64
}

func (n Numeric) Mantissa() float64 {
	return n.mantissa
}

func (n Numeric) Exponent() int64 {
	return n.exponent
}

func (n Numeric) IsNaN() bool {
	return math.IsNaN(n.mantissa)
}

func (n Numeric) IsPosInf() bool {
	return math.IsInf(n.mantissa, 1)
}

func (n Numeric) IsNegInf() bool {
	return math.IsInf(n.mantissa, -1)
}

func (n Numeric) IsInf() bool {
	return math.IsInf(n.mantissa, 0)
}

func (n Numeric) IsInt() bool {
	return math.Round(n.mantissa) == n.mantissa
}

func (n Numeric) IsZero() bool {
	return math.Abs(n.mantissa) < Epsilon
}

func (n Numeric) Sign() int {
	return util.Choose(math.Signbit(n.mantissa), -1, 1)
}

func (n Numeric) Abs() Numeric {
	return Numeric{mantissa: math.Abs(n.mantissa), exponent: n.exponent}
}

func (n Numeric) Negate() Numeric {
	return Numeric{mantissa: -n.mantissa, exponent: n.exponent}
}

func (n Numeric) Round() Numeric {
	return mapNumeric(n, func() Numeric {
		return Zero
	}, math.Round)
}

func (n Numeric) Floor() Numeric {
	return mapNumeric(n, func() Numeric {
		if math.Signbit(n.mantissa) {
			return One.Negate()
		}
		return Zero
	}, math.Floor)
}

func (n Numeric) Ceiling() Numeric {
	return mapNumeric(n, func() Numeric {
		if n.mantissa > 0 {
			return One
		}
		return Zero
	}, math.Ceil)
}

func (n Numeric) Truncate() Numeric {
	if n.IsNaN() {
		return n
	}
	if n.exponent < -1 {
		return Zero
	}
	if n.exponent < MaxSignificantDigits {
		return FromFloat(math.Trunc(n.ToFloat()))
	}
	return n
}

func (n Numeric) Clone() Numeric {
	return Numeric{mantissa: n.mantissa, exponent: n.exponent}
}

func (n Numeric) Equals(other Numeric) bool {
	ret := !math.IsNaN(n.mantissa) && !math.IsNaN(other.mantissa)
	return ret && (areSameInfinity(n, other) || (n.exponent == other.exponent && areEqual(n.mantissa, other.mantissa)))
}

func (n Numeric) String() string {
	ret := fmt.Sprintf("%.3f", n.mantissa)
	for strings.HasSuffix(ret, "0") {
		ret = strings.TrimSuffix(ret, "0")
	}
	ret = strings.TrimSuffix(ret, ".")
	return fmt.Sprintf("%se%d", ret, n.exponent)
}

func (n Numeric) Debug() string {
	return fmt.Sprintf("%fe%d", n.mantissa, n.exponent)
}

func (n Numeric) MarshalJSON() ([]byte, error) {
	return []byte(`{"m": ` + fmt.Sprint(n.mantissa) + `, "e": ` + fmt.Sprint(n.exponent) + `}`), nil
}

func (n Numeric) UnmarshalJSON(data []byte) error {
	return nil
}
