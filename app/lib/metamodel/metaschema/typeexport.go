package metaschema

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
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

func ExportModel(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*model.Model, error) {
	if sch == nil {
		return nil, errors.New("nil schema provided for model")
	}
	_, err := exportGetType(sch, "object")
	if err != nil {
		return nil, err
	}
	ret := &model.Model{Name: "TODO"}
	return ret, nil
}

func ExportEnum(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*enum.Enum, error) {
	if sch == nil {
		return nil, errors.New("nil schema provided for enum")
	}
	_, err := exportGetType(sch, "enum")
	if err != nil {
		return nil, err
	}
	ret := &enum.Enum{Name: "TODO"}
	return ret, nil
}

func ExportType(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (types.Type, error) {
	if sch == nil {
		return nil, errors.New("nil schema provided for type")
	}
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
	t, ok := sch.Type.(string)
	if !ok {
		return "", errors.Errorf("invalid schema type [%s], of type [%T]", fmt.Sprint(sch.Type), sch.Type)
	}
	if len(expected) > 0 && !slices.Contains(expected, t) {
		return "", errors.Errorf("invalid schema type [%s], expected [%s]", sch.Type, expected)
	}
	return t, nil
}
