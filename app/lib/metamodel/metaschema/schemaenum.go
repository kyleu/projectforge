package metaschema

import (
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/util"
)

func EnumSchema(x *enum.Enum, sch *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	ret := sch.NewSchema(util.StringPath(x.PackageWithGroup(""), x.Name))
	ret.Type = "string"
	ret.Description = x.Description

	vals := util.ValueMap{}
	for _, v := range x.Values {
		ret.Enum = append(ret.Enum, v.Key)
		vals[v.Key] = v.ToOrderedMap(false)
	}
	if len(vals) > 0 {
		ret.AddMetadata("values", vals)
	}

	if len(x.Config) > 0 {
		ret.AddMetadata("config", x.Config)
	}
	if x.Icon != "" {
		ret.AddMetadata("icon", x.Icon)
	}
	if x.ProperOverride != "" {
		ret.AddMetadata("proper", x.ProperOverride)
	}
	if x.RouteOverride != "" {
		ret.AddMetadata("route", x.RouteOverride)
	}
	if len(x.Tags) > 0 {
		ret.AddMetadata("tags", x.Tags)
	}
	if x.TitleOverride != "" {
		ret.AddMetadata("title", x.TitleOverride)
	}

	return ret, nil
}
