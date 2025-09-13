package metaschema

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func ExportEvent(x *model.Event, coll *jsonschema.Collection, arg *metamodel.Args) (*jsonschema.Schema, error) {
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
	if x.PluralOverride != "" {
		ret.AddMetadata("plural", x.PluralOverride)
	}
	if len(x.Tags) > 0 {
		ret.AddMetadata("tags", x.Tags)
	}
	if x.TitleOverride != "" {
		ret.Title = x.TitleOverride
	}
	if len(x.Imports) > 0 {
		ret.AddMetadata("imports", x.Imports)
	}
	return ret, nil
}
