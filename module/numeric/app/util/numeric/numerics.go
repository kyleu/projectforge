package numeric

import "github.com/samber/lo"

type Numerics []Numeric

func (n Numerics) Sum() Numeric {
	if len(n) == 0 {
		return Zero
	}
	return lo.Reduce(n[1:], func(agg Numeric, x Numeric, _ int) Numeric {
		return agg.Add(x)
	}, n[0])
}
