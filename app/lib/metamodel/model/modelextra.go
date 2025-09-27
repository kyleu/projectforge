package model

import (
	"fmt"

	"github.com/samber/lo"

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

func (m *Model) ConfigMap() util.ValueMap {
	return m.Config
}

func (m Models) WithRoutes() Models {
	return lo.Reject(m, func(x *Model, _ int) bool {
		return x.SkipRoutes()
	})
}

func (m Models) WithController() Models {
	return lo.Reject(m, func(x *Model, _ int) bool {
		return x.SkipController()
	})
}

func (m Models) WithService() Models {
	return lo.Reject(m, func(x *Model, _ int) bool {
		return x.SkipService()
	})
}

func (m Models) WithDatabase() Models {
	return lo.Reject(m, func(x *Model, _ int) bool {
		return x.SkipDatabase()
	})
}

func (m Models) WithTypeScript() Models {
	return lo.Filter(m, func(x *Model, _ int) bool {
		return x.HasTag("typescript")
	})
}

func (m *Model) RelativePath(rGroup []string, extra ...string) string {
	mGroup := m.GroupAndPackage()
	commonPrefix := 0
	for i := 0; i < len(mGroup) && i < len(rGroup) && mGroup[i] == rGroup[i]; i++ {
		commonPrefix++
	}
	upLevels := len(mGroup) - commonPrefix
	var pathParts []string
	for i := commonPrefix; i < len(rGroup); i++ {
		pathParts = append(pathParts, rGroup[i])
	}
	pathParts = append(pathParts, extra...)
	return fmt.Sprintf("%s%s", util.StringRepeat("../", upLevels), util.StringJoin(pathParts, "/"))
}
