package metaschema

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/util"
)

func ExportArgs(coll *jsonschema.Collection, args *metamodel.Args) (*metamodel.Args, error) {
	if coll == nil || len(coll.Schemas) == 0 {
		return nil, errors.New("empty collection provided for arguments")
	}
	ret := &metamodel.Args{}
	for _, sch := range coll.Schemas {
		switch sch.Type {
		case "object":
			x, err := ExportModel(sch, coll, args)
			if err != nil {
				return nil, err
			}
			ret.Models = append(ret.Models, x)
		case "string":
			if len(sch.Enum) == 0 {
				return nil, errors.Errorf("unhandled type [string] without [enum] values for schema [%s]", sch.String())
			}
			x, err := ExportEnum(sch, coll, args)
			if err != nil {
				return nil, err
			}
			ret.Enums = append(ret.Enums, x)
		default:
			return nil, errors.Errorf("invalid type [%v] for schema [%s]", sch.Type, sch.String())
		}
	}
	return ret, nil
}

func ExportAny(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (any, error) {
	switch sch.Type {
	case "object":
		return ExportModel(sch, coll, args)
	case "enum":
		return ExportEnum(sch, coll, args)
	default:
		return nil, errors.Errorf("invalid type [%v] for schema [%s]", sch.Type, sch.String())
	}
}

func parseID(id string) (string, string, []string) {
	parts := util.StringSplitAndTrim(id, "/")
	pid, pkg := "", ""
	var grp []string
	switch len(parts) {
	case 1:
		pid = id
	case 2:
		pid = parts[1]
		pkg = parts[0]
	default:
		pid = parts[len(parts)-1]
		pkg = parts[len(parts)-2]
		grp = parts[:len(parts)-2]
	}
	pid = strings.TrimSuffix(strings.TrimSuffix(pid, ".json"), ".schema")
	return pid, pkg, grp
}
