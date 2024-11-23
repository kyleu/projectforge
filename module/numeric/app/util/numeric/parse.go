package numeric

import (
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func from(mantissa float64, exponent int64) Numeric {
	return Numeric{mantissa: mantissa, exponent: exponent}
}

func From(mantissa float64, exponent int64) Numeric {
	return normalize(mantissa, exponent)
}

func Random(maxExponent int) Numeric {
	return normalize(util.RandomFloat(1), int64(util.RandomInt(maxExponent+1)))
}

func FromInt(value int) Numeric {
	return FromFloat(float64(value))
}

func FromInt64(value int64) Numeric {
	return FromFloat(float64(value))
}

func FromFloat(value float64) Numeric {
	if math.IsNaN(value) {
		return NaN
	} else if math.IsInf(value, 1) {
		return PositiveInfinity
	} else if math.IsInf(value, -1) {
		return NegativeInfinity
	} else if isZero(value) {
		return Zero
	}
	return normalize(value, 0)
}

func FromString(s string) (Numeric, error) {
	if strings.Contains(s, " ") {
		split := util.StringSplitAndTrim(s, " ")
		if len(split) == 2 {
			amt, err := strconv.ParseFloat(split[0], 64)
			if err != nil {
				return Numeric{}, errors.Errorf("invalid initial amount [%s]: %s", split[0], s)
			}
			pow10, err := Pow10FromEnglish(split[1])
			if err != nil {
				return Numeric{}, errors.Errorf("invalid English number [%s]: %s", split[1], s)
			}
			ret := normalize(amt, int64(pow10))
			return ret, nil
		}
	}
	if strings.Contains(s, "e") {
		parts := strings.Split(s, "e")
		if len(parts) != 2 {
			return Numeric{}, errors.Errorf("invalid scientific notation: %s", s)
		}
		mantissa, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return Numeric{}, errors.Errorf("invalid mantissa: %s", parts[0])
		}
		exponent, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return Numeric{}, errors.Errorf("invalid exponent: %s", parts[1])
		}
		return normalize(mantissa, exponent), nil
	}

	if s == "NaN" {
		return NaN, nil
	}

	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Numeric{}, errors.Errorf("invalid number: %s", s)
	}

	result := FromFloat(floatVal)
	if math.IsNaN(result.mantissa) {
		return Numeric{}, errors.Errorf("invalid argument: %s", s)
	}

	return result, nil
}

func FromStringOK(s string) Numeric {
	ret, _ := FromString(s)
	return ret
}
