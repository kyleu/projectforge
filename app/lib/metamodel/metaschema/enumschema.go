package metaschema

import (
	"path"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
)

func EnumSchema(x *enum.Enum, sch *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	ret := sch.NewSchema(path.Join(x.PackageWithGroup(""), x.Name))
	ret.Type = "enum"
	ret.Description = x.Description
	ret.Enum = lo.Map(x.Values, func(x *enum.Value, _ int) any {
		return x.Key
	})
	return ret, nil
}
