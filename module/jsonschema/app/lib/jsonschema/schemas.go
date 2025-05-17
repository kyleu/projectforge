package jsonschema

import "github.com/samber/lo"

type Schemas []*Schema

func (s Schemas) IDs() []string {
	return lo.Map(s, func(x *Schema, _ int) string {
		return x.ID()
	})
}
