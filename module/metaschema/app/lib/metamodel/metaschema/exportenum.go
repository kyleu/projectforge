package metaschema

import (
	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/lib/metamodel/enum"
	"{{{ .Package }}}/app/util"
)

func ExportEnum(x *enum.Enum, sch *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	id := util.StringPath(x.PackageWithGroup(""), x.Name)
	ret := sch.NewSchema(id)
	ret.Type = KeyString
	ret.Description = x.Description

	vals := util.ValueMap{}
	for _, v := range x.Values {
		ret.Enum = append(ret.Enum, v.Key)
		m := v.ToOrderedMap(false)
		if len(m.Order) > 0 {
			vals[v.Key] = m
		}
	}
	if len(vals) > 0 {
		ret.AddMetadata("values", vals)
	}

	if x.Schema != "" {
		ret.AddMetadata("schema", x.Schema)
	}
	if x.Icon != "" {
		ret.AddMetadata("icon", x.Icon)
	}
	if len(x.Tags) > 0 {
		ret.AddMetadata("tags", x.Tags)
	}
	if x.TitleOverride != "" {
		ret.AddMetadata("title", x.TitleOverride)
	}
	if x.ProperOverride != "" {
		ret.AddMetadata("proper", x.ProperOverride)
	}
	if x.RouteOverride != "" {
		ret.AddMetadata("route", x.RouteOverride)
	}
	if len(x.Config) > 0 {
		ret.AddMetadata("config", x.Config)
	}

	return ret, nil
}
