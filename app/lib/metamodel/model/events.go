package model

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Events []*Event

func (m Events) Get(n string) *Event {
	return lo.FindOrElse(m, nil, func(x *Event) bool {
		return x.Name == n
	})
}

func (m Events) ForGroup(pth ...string) Events {
	return lo.Filter(m, func(x *Event, _ int) bool {
		return slices.Equal(x.Group, pth)
	})
}

func (m Events) Validate(mods []string, groups Groups) error {
	names := util.ValueMap{}
	for _, x := range m {
		if _, ok := names[x.Name]; ok {
			return errors.Errorf("multiple Events found with name [%s]", x.Name)
		}
	}
	return nil
}

func (m Events) WithTag(tag string) Events {
	return lo.Filter(m, func(x *Event, _ int) bool {
		return x.HasTag(tag)
	})
}

func (m Events) WithoutTag(tag string) Events {
	return lo.Reject(m, func(x *Event, _ int) bool {
		return x.HasTag(tag)
	})
}
