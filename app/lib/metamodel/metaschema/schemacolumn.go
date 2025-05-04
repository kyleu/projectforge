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
	return ret, nil
}
