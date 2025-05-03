package enum

import "github.com/samber/lo"

func (e *Enum) SkipDatabase() bool {
	return e.HasTag("no-database") || e.SkipGolang()
}

func (e *Enum) SkipGolang() bool {
	return e.HasTag("no-database")
}

func (e Enums) WithDatabase() Enums {
	return lo.Reject(e, func(x *Enum, _ int) bool {
		return x.SkipDatabase()
	})
}

func (e Enums) WithTypeScript() Enums {
	return lo.Filter(e, func(x *Enum, _ int) bool {
		return x.HasTag("typescript")
	})
}
