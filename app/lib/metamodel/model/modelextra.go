package model

import (
	"fmt"
	"slices"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func (m *Model) SkipRoutes() bool {
	return m.HasTag("no-routes") || m.SkipController()
}

func (m *Model) SkipController() bool {
	return m.HasTag("no-controller") || m.SkipService()
}

func (m *Model) SkipService() bool {
	return m.HasTag("no-service") || m.SkipDatabase()
}

func (m *Model) SkipDatabase() bool {
	return m.HasTag("no-database") || m.SkipGolang()
}

func (m *Model) SkipGolang() bool {
	return m.HasTag("no-golang")
}

func (m Models) WithRoutes() Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return !x.SkipRoutes()
	})
}

func (m Models) WithController() Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return !x.SkipController()
	})
}

func (m Models) WithService() Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return !x.SkipService()
	})
}

func (m Models) WithDatabase() Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return !x.SkipDatabase()
	})
}

func (m Models) WithTypeScript() Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return x.HasTag("typescript")
	})
}

func (m *Model) AllSearches(db string) []string {
	if !m.HasTag("search") {
		return m.Search
	}
	ret := util.NewStringSlice(slices.Clone(m.Search))
	lo.ForEach(m.Columns, func(c *Column, _ int) {
		if c.Search {
			x := fmt.Sprintf("%q", c.SQL())
			if !types.IsString(c.Type) {
				switch db {
				case dbSQLServer:
					x = fmt.Sprintf("cast(%q as nvarchar(2048))", c.SQL())
				case dbSQLite:
					x = c.SQL()
				default:
					x = fmt.Sprintf("%q::text", c.SQL())
				}
			}
			ret.Push(fmt.Sprintf("lower(%s)", x))
		}
	})
	return ret.Slice
}

func (m *Model) HasSearches() bool {
	return len(m.AllSearches("")) > 0
}
