package metaschema

import (
	"fmt"
	"github.com/pkg/errors"
	"slices"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/types"
)

func ExportType(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*types.Wrapped, error) {
	t, err := exportGetType(sch)
	if err != nil {
		return nil, err
	}
	switch t {
	case "", "nil", "<nil>", "null":
		return types.NewNil(), nil
	case "array":
		return types.NewList(types.NewAny()), nil
	case "enum":
		return types.NewEnum("?"), nil
	case "integer":
		return types.NewInt(0), nil
	case "number":
		return types.NewFloat(0), nil
	case "object":
		return types.NewValueMap(), nil
	case "string":
		switch sch.Format {
		case "date":
			return types.NewDate(), nil
		case "date-time":
			return types.NewTimestamp(), nil
		case "uuid":
			return types.NewUUID(), nil
		case "":
			return types.NewString(), nil
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
	if sch.Type == nil {
		return "nil", nil
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
