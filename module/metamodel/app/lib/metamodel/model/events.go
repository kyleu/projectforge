package model

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Events []*Event

func (e Events) Get(n string) *Event {
	return lo.FindOrElse(e, nil, func(x *Event) bool {
		return x.Name == n || x.Camel() == n || x.Proper() == n
	})
}

func (e Events) ForGroup(pth ...string) Events {
	return lo.Filter(e, func(x *Event, _ int) bool {
		return slices.Equal(x.Group, pth)
	})
}

func (e Events) Validate(mods []string, groups Groups) error {
	names := util.ValueMap{}
	for _, x := range e {
		if _, ok := names[x.Name]; ok {
			return errors.Errorf("multiple Events found with name [%s]", x.Name)
		}
	}
	return nil
}

func (e Events) WithTag(tag string) Events {
	return lo.Filter(e, func(x *Event, _ int) bool {
		return x.HasTag(tag)
	})
}

func (e Events) WithoutTag(tag string) Events {
	return lo.Reject(e, func(x *Event, _ int) bool {
		return x.HasTag(tag)
	})
}
