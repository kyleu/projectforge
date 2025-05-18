package metaschema

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func ImportType(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*types.Wrapped, error) {
	t, err := exportGetType(sch)
	if err != nil {
		return nil, err
	}
	var ret *types.Wrapped
	switch t {
	case "", "nil", "<nil>", "null":
		switch sch.Ref {
		case "":
			ret = types.NewAny()
		case "Numeric":
			ret = types.NewNumeric()
		default:
			ret = types.NewReferencePath(sch.Ref, true)
		}
	case "array":
		if sch.Items == nil {
			ret = types.NewList(types.NewAny())
		} else {
			switch itemType := sch.Items.(type) {
			case string:
				ret = types.NewList(types.FromJSONType(itemType, sch.Ref))
			case map[string]any:
				b := util.ToJSONBytes(itemType, true)
				itemSch, e := jsonschema.FromJSON(b)
				if e != nil {
					return nil, errors.Wrapf(e, "invalid array item subschema [%s] for schema [%s]", util.ToJSON(itemType), sch.String())
				}
				itemT, e := ImportType(itemSch, coll, args)
				if e != nil {
					return nil, errors.Wrapf(e, "error processing item subschema [%s] for schema [%s]", util.ToJSON(itemType), sch.String())
				}
				ret = types.NewList(itemT)
			case *jsonschema.Schema:
				itemT, e := ImportType(itemType, coll, args)
				if e != nil {
					return nil, errors.Wrapf(e, "error processing item subschema [%s] for schema [%s]", itemType.String(), sch.String())
				}
				ret = types.NewList(itemT)
			default:
				return nil, errors.Errorf("invalid array item type [%T] for schema [%s]", itemType, sch.String())
			}
		}
	case "boolean":
		ret = types.NewBool()
	case "enum":
		ret = types.NewEnum(sch.Ref)
	case "integer":
		ret = types.NewInt(sch.Metadata.GetIntOpt("bits"))
	case "number":
		ret = types.NewFloat(sch.Metadata.GetIntOpt("bits"))
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
		return nil, errors.Errorf("unhandled schema type [%s] in schema [%s]", t, sch.String())
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
