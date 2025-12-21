package numeric

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type NumericMap map[string]Numeric

func FromMap(m util.ValueMap) (NumericMap, []string) {
	ret := make(NumericMap, len(m))
	var errs []string
	for k, v := range m {
		n, err := FromAny(v)
		if err == nil {
			ret[k] = n
		} else {
			errs = append(errs, errors.Wrapf(err, "unable to parse numeric value from type [%T]", v).Error())
		}
	}
	return ret, errs
}

func (n NumericMap) Clone() NumericMap {
	return util.MapClone(n)
}

func (n NumericMap) Negate() NumericMap {
	ret := make(NumericMap, len(n))
	for k, v := range n {
		ret[k] = v.Negate()
	}
	return ret
}

func (n NumericMap) ToMap() util.ValueMap {
	ret := make(util.ValueMap, len(n))
	for k, v := range n {
		ret[k] = v
	}
	return ret
}

func (n NumericMap) GetOr(key string, dflt Numeric) Numeric {
	if amt, ok := n[key]; ok {
		return amt
	}
	return dflt
}

func (n NumericMap) Matching(tgt NumericMap) (NumericMap, error) {
	ret := make(NumericMap, len(n))
	for k := range n {
		x, ok := tgt[k]
		if !ok {
			return nil, errors.Errorf("can't find [%s] among NumericMap keys [%s]", k, util.StringJoin(util.MapKeysSorted(tgt), ", "))
		}
		ret[k] = x
	}
	return ret, nil
}

func (n NumericMap) Trimmed() NumericMap {
	ret := make(NumericMap, len(n))
	for k, v := range n {
		if !v.IsZero() {
			ret[k] = v
		}
	}
	return ret
}

func (n NumericMap) IsZero() bool {
	for _, v := range n {
		if !v.IsZero() {
			return false
		}
	}
	return true
}

func (n NumericMap) String() string {
	if len(n) == 0 {
		return "[none]"
	}
	ret := util.NewStringSliceWithSize(len(n))
	for k, v := range n {
		ret.Pushf("%s: %s", k, v.String())
	}
	return "[" + ret.JoinCommas() + "]"
}

func (n NumericMap) merge(x NumericMap, fn func(v Numeric, k string) Numeric) NumericMap {
	ret := lo.MapValues(n, fn)
	for k, v := range x {
		if _, ok := n[k]; !ok {
			ret[k] = v
		}
	}
	return ret
}

func ReduceNumericMap(m ...NumericMap) NumericMap {
	return lo.Reduce(m, func(agg NumericMap, i NumericMap, _ int) NumericMap {
		return agg.Add(i)
	}, NumericMap{})
}
