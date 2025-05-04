package metaschema

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
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
		case "enum":
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

func ExportModel(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*model.Model, error) {
	if sch == nil {
		return nil, errors.New("nil schema provided for model")
	}
	_, err := exportGetType(sch, "object")
	if err != nil {
		return nil, err
	}
	n, pkg, grp := parseID(sch.ID)
	ret := &model.Model{Name: n, Package: pkg, Group: grp}
	return ret, nil
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
		grp = parts[len(parts)-3:]
	}
	pid = strings.TrimSuffix(strings.TrimSuffix(pid, ".json"), ".schema")
	return pid, pkg, grp
}

func ExportEnum(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*enum.Enum, error) {
	_, err := exportGetType(sch, "enum")
	if err != nil {
		return nil, err
	}
	n, pkg, grp := parseID(sch.ID)
	ret := &enum.Enum{Name: n, Package: pkg, Group: grp}
	return ret, nil
}

func ExportType(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (types.Type, error) {
	t, err := exportGetType(sch)
	if err != nil {
		return nil, err
	}
	switch t {
	case "string":
		switch sch.Format {
		case "date":
			return types.NewDate(), nil
		case "date-time":
			return types.NewTimestamp(), nil
		case "uuid":
			return types.NewUUID(), nil
		default:
			return nil, errors.Errorf("invalid string format [%s] for schema [%s]", sch.Format, sch.String())
		}
	default:
		return nil, errors.Errorf("unhandled schema type [%s] in schema [%T]", t, sch.String())
	}
}

func exportGetType(sch *jsonschema.Schema, expected ...string) (string, error) {
	if sch == nil {
		return "", errors.New("nil schema provided")
	}
	t, ok := sch.Type.(string)
	if !ok {
		return "", errors.Errorf("invalid schema type [%s], of type [%T]", fmt.Sprint(sch.Type), sch.Type)
	}
	if len(expected) > 0 && !slices.Contains(expected, t) {
		return "", errors.Errorf("invalid schema type [%s], expected [%s]", sch.Type, expected)
	}
	return t, nil
}
