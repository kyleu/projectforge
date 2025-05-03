package metaschema

import (
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
)

func ColumnSchema(col *model.Column, sch *jsonschema.Collection, arg *metamodel.Args) (*jsonschema.Schema, error) {
	return &jsonschema.Schema{Type: col.Type.String()}, nil
}
