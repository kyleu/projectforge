package metaschema

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/lib/metamodel/model"
	"{{{ .Package }}}/app/util"
)

func ExportModel(x *model.Model, coll *jsonschema.Collection, arg *metamodel.Args) (*jsonschema.Schema, error) {
	id := util.StringPath(x.PackageWithGroup(""), x.Name)
	ret := coll.NewSchema(id)
	ret.Type = KeyObject
	ret.Description = x.Description
	ret.Properties = util.NewOrderedMap[*jsonschema.Schema](false, len(x.Columns))
	ret.Required = x.Columns.Required().Names()
	for _, col := range x.Columns {
		colSch, err := ExportColumn(col, coll, arg)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse column [%s]", col.String())
		}
		ret.Properties.Set(col.Name, colSch)
	}
	if x.Schema != "" {
		ret.AddMetadata("schema", x.Schema)
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
	if x.SortIndex != 0 {
		ret.AddMetadata("sortIndex", x.SortIndex)
	}
	if x.View != "" {
		ret.AddMetadata("view", x.View)
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
	if x.TableOverride != "" {
		ret.AddMetadata("table", x.TableOverride)
	}
	if len(x.Tags) > 0 {
		ret.AddMetadata("tags", x.Tags)
	}
	if x.TitleOverride != "" {
		ret.Title = x.TitleOverride
	}
	if len(x.Relations) > 0 {
		ret.AddMetadata("relations", x.Relations)
	}
	if len(x.Indexes) > 0 {
		ret.AddMetadata("indexes", x.Indexes)
	}
	if len(x.SeedData) > 0 {
		ret.Examples = util.ArrayToAnyArray(x.SeedData)
	}
	if len(x.Links) > 0 {
		ret.AddMetadata("links", x.Links)
	}
	if len(x.Imports) > 0 {
		ret.AddMetadata("imports", x.Imports)
	}
	return ret, nil
}
