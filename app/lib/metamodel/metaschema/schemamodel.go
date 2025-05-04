package metaschema

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func ModelSchema(x *model.Model, sch *jsonschema.Collection, arg *metamodel.Args) (*jsonschema.Schema, error) {
	ret := sch.NewSchema(util.StringPath(x.PackageWithGroup(""), x.Name))
	ret.Type = "object"
	ret.Description = x.Description
	ret.Properties = util.NewOrderedMap[*jsonschema.Schema](false, len(x.Columns))
	ret.Required = x.Columns.Required().CamelNames()
	for _, col := range x.Columns {
		colSch, err := ColumnSchema(col, sch, arg)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse column [%s]", col.String())
		}
		ret.Properties.Set(col.Name, colSch)
	}
	if len(x.Config) > 0 {
		ret.AddMetadata("config", x.Config)
	}
	if x.Icon != "" {
		ret.AddMetadata("icon", x.Icon)
	}
	if len(x.Ordering) > 0 {
		ret.AddMetadata("ordering", x.Ordering)
	}
	if x.PluralOverride != "" {
		ret.AddMetadata("plural", x.PluralOverride)
	}
	if x.RouteOverride != "" {
		ret.AddMetadata("route", x.RouteOverride)
	}
	if len(x.Search) > 0 {
		ret.AddMetadata("search", x.Search)
	}
	if x.SortIndex != 0 {
		ret.AddMetadata("sortIndex", x.SortIndex)
	}
	if x.TableOverride != "" {
		ret.AddMetadata("table", x.TableOverride)
	}
	if len(x.Tags) > 0 {
		ret.AddMetadata("tags", x.Tags)
	}
	if x.TitleOverride != "" {
		ret.AddMetadata("title", x.TitleOverride)
	}
	if x.View != "" {
		ret.AddMetadata("view", x.View)
	}
	return ret, nil
}
