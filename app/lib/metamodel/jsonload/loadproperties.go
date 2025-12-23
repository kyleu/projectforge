package jsonload

import (
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
)

func parseProperties(sch *jsonschema.Schema) (model.Columns, error) {
	if sch == nil || sch.Properties == nil {
		return nil, nil
	}
	ret := make(model.Columns, 0, len(sch.Properties.Order))
	for _, prop := range sch.Properties.Order {
		col, err := parseProperty(sch, prop, sch.Properties.GetSimple(prop))
		if err != nil {
			return nil, err
		}
		ret = append(ret, col)
	}
	return ret, nil
}

func parseProperty(sch *jsonschema.Schema, name string, prop *jsonschema.Schema) (*model.Column, error) {
	md := prop.GetMetadata()
	ret := &model.Column{
		Name:           name,
		Type:           types.NewAny(),
		PK:             md.GetBoolOpt("pk"),
		Nullable:       md.GetBoolOpt("nullable"),
		Search:         md.GetBoolOpt("search"),
		SQLDefault:     md.GetStringOpt("default"),
		Indexed:        md.GetBoolOpt("indexed"),
		Display:        md.GetStringOpt("display"),
		Format:         md.GetStringOpt("format"),
		JSON:           md.GetStringOpt("json"),
		SQLOverride:    md.GetStringOpt("sql"),
		TitleOverride:  md.GetStringOpt("title"),
		PluralOverride: md.GetStringOpt("plural"),
		ProperOverride: md.GetStringOpt("proper"),
		Example:        md.GetStringOpt("example"),
		Validation:     md.GetStringOpt("validation"),
		Values:         md.GetStringArrayOpt("values"),
		Tags:           md.GetStringArrayOpt("tags"),
		Comment:        md.GetStringOpt("comment"),
		Help:           md.GetStringOpt("help"),
		Metadata:       md.GetMapOpt("metadata"),
	}
	return ret, nil
}
