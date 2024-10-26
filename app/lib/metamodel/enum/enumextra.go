package enum

import "github.com/samber/lo"

func (e *Enum) SkipDatabase() bool {
	return e.HasTag("no-database")
}

func (e Enums) WithDatabase() Enums {
	return lo.Filter(e, func(x *Enum, _ int) bool {
		return !x.SkipDatabase()
	})
}
