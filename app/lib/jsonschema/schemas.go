package jsonschema

import "github.com/samber/lo"

var (
	trueSchemaData  = &Schema{}
	falseSchemaData = &Schema{data: data{dataApplicators: dataApplicators{Not: &Schema{}}}}
)

func NewTrueSchema() *Schema {
	return trueSchemaData.Clone()
}

func NewFalseSchema() *Schema {
	return falseSchemaData.Clone()
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
