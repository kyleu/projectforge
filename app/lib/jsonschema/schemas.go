package jsonschema

import "github.com/samber/lo"

var (
	trueSchema  = &Schema{}
	falseSchema = &Schema{Data: Data{DataApplicators: DataApplicators{Not: &Schema{}}}}
)

func NewTrueSchema() *Schema {
	return trueSchema.Clone()
}

func NewFalseSchema() *Schema {
	return falseSchema.Clone()
}

type Schemas []*Schema

func (s Schemas) IDs() []string {
	return lo.Map(s, func(x *Schema, _ int) string {
		return x.ID()
	})
}

func (s Schemas) Clone() Schemas {
	return lo.Map(s, func(x *Schema, _ int) *Schema {
		return x.Clone()
	})
}
