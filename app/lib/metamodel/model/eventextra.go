package model

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func (e *Event) ConfigMap() util.ValueMap {
	return e.Config
}

func (e Events) WithTypeScript() Events {
	return lo.Filter(e, func(x *Event, _ int) bool {
		return !x.HasTag("no-typescript")
	})
}
