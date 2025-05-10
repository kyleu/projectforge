package metaschema

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/types"
)

func ExportType(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*types.Wrapped, error) {
	t, err := exportGetType(sch)
	if err != nil {
		return nil, err
	}
	var ret *types.Wrapped
	switch t {
	case "", "nil", "<nil>", "null":
		ret = types.NewNil()
	case "array":
		ret = types.NewList(types.NewAny())
	case "boolean":
		ret = types.NewBool()
	case "enum":
		ret = types.NewEnum("?")
	case "integer":
		ret = types.NewInt(0)
	case "number":
		ret = types.NewFloat(0)
	case "object":
		ret = types.NewMap(types.NewString(), types.NewAny())
	case "string":
		switch sch.Format {
		case "date":
			ret = types.NewDate()
		case "date-time":
			ret = types.NewTimestamp()
		case "uuid":
			ret = types.NewUUID()
		case "":
			ret = types.NewString()
		default:
			return nil, errors.Errorf("invalid string format [%s] for schema [%s]", sch.Format, sch.String())
		}
	default:
		return nil, errors.Errorf("unhandled schema type [%s] in schema [%T]", t, sch.String())
	}
	return ret, nil
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
