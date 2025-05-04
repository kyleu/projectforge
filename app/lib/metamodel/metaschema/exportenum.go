package metaschema

import (
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
)

func ExportEnum(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*enum.Enum, error) {
	_, err := exportGetType(sch, "enum")
	if err != nil {
		return nil, err
	}
	n, pkg, grp := parseID(sch.ID)
	ret := &enum.Enum{Name: n, Package: pkg, Group: grp}
	return ret, nil
}
