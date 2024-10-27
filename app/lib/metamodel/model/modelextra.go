package model

import "github.com/samber/lo"

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
