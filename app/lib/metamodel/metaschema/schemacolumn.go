package metaschema

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
)

func ColumnSchema(col *model.Column, coll *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	ret, err := TypeSchema(col.Type, coll, args)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render type [%s] for column [%s]", col.Type.String(), col.Name)
	}
	if col.PK {
		ret.AddMetadata("pk", col.PK)
	}
	if col.Search {
		ret.AddMetadata("search", col.Search)
	}
	if col.SQLDefault != "" {
		ret.Default = col.SQLDefault
	}
	if col.Indexed {
		ret.AddMetadata("indexed", col.Indexed)
	}
	if col.Display != "" {
		ret.AddMetadata("display", col.Display)
	}
	if col.Format != "" {
		ret.AddMetadata("format", col.Format)
	}
	if col.JSON != "" {
		ret.AddMetadata("json", col.JSON)
	}
	if col.SQLOverride != "" {
		ret.AddMetadata("sql", col.SQLOverride)
	}
	if col.TitleOverride != "" {
		ret.AddMetadata("title", col.TitleOverride)
	}
	if col.PluralOverride != "" {
		ret.AddMetadata("plural", col.PluralOverride)
	}
	if col.ProperOverride != "" {
		ret.AddMetadata("proper", col.ProperOverride)
	}
	if col.Example != "" {
		ret.AddMetadata("example", col.Example)
	}
	if col.Validation != "" {
		ret.AddMetadata("validation", col.Validation)
	}
	if len(col.Values) > 0 {
		ret.AddMetadata("values", col.Values)
	}
	if len(col.Tags) > 0 {
		ret.AddMetadata("tags", col.Tags)
	}
	if col.HelpString != "" {
		ret.AddMetadata("help", col.HelpString)
	}
	if col.Metadata != nil {
		for k, v := range col.Metadata {
			ret.AddMetadata(k, v)
		}
	}
	return ret, nil
}
