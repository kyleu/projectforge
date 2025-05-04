package metaschema

import (
	"github.com/pkg/errors"
	"path"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func ModelSchema(x *model.Model, sch *jsonschema.Collection, arg *metamodel.Args) (*jsonschema.Schema, error) {
	ret := sch.NewSchema(path.Join(x.PackageWithGroup(""), x.Name))
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
	return ret, nil
}
