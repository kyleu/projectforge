package metaschema

import (
	"fmt"
	"slices"
	"strings"

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
	case "", KeyNil, KeyNilString, KeyNull:
		switch sch.Ref {
		case "":
			ret = types.NewAny()
		case types.KeyNumeric:
			ret = types.NewNumeric()
		case types.KeyNumericMap:
			ret = types.NewNumericMap()
		default:
			if strings.Contains(sch.Ref, "/") {
				ret = types.NewReferencePath(sch.Ref, true)
			} else {
				ref := sch.Ref
				ref = strings.TrimSuffix(ref, ".json")
				ref = strings.TrimSuffix(ref, ".schema")
				ret = types.NewEnum(ref)
			}
		}
	case KeyArray:
		ret, err = importArrayType(sch, coll, args)
	case KeyBoolean:
		ret = types.NewBool()
	case KeyEnum:
		ret = types.NewEnum(sch.Ref)
	case KeyInteger:
		ret = types.NewInt(sch.GetMetadata().GetIntOpt("bits"))
	case KeyNumber:
		ret = types.NewFloat(sch.GetMetadata().GetIntOpt("bits"))
	case KeyObject:
		md := sch.GetMetadata()
		if md != nil && md["type"] == util.KeyJSON {
			ret = types.NewJSON()
		} else {
			ret = types.NewStringKeyedMap()
		}
	case KeyString:
		switch sch.Format {
		case KeyDate:
			ret = types.NewDate()
		case KeyDateTime:
			md := sch.GetMetadata()
			if md != nil && md["type"] == types.KeyTimestampZoned {
				ret = types.NewTimestampZoned()
			} else {
				ret = types.NewTimestamp()
			}
		case KeyUUID:
			ret = types.NewUUID()
		case "", KeyString:
			s := &types.String{}
			if sch.MaxLength != nil && *sch.MaxLength > 0 {
				s.MaxLength = int(*sch.MaxLength)
			}
			if sch.MinLength != nil && *sch.MinLength > 0 {
				s.MinLength = int(*sch.MinLength)
			}
			ret = types.Wrap(s)
		default:
			return nil, errors.Errorf("invalid string format [%s] for schema [%s]", sch.Format, sch.String())
		}
	default:
		return nil, errors.Errorf("unhandled schema type [%s] in schema [%s]", t, sch.String())
	}
	return ret, err
}

func importArrayType(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*types.Wrapped, error) {
	if sch.Items == nil {
		return types.NewList(types.NewAny()), nil
	}

	itemT, e := ImportType(sch.Items, coll, args)
	if e != nil {
		return nil, errors.Wrapf(e, "error processing item subschema [%s] for schema [%s]", sch.Items.String(), sch.String())
	}
	return types.NewList(itemT), nil
}

func exportGetType(sch *jsonschema.Schema, expected ...string) (string, error) {
	if sch == nil {
		return "", errors.New("nil schema provided")
	}
	if sch.Type == nil {
		return KeyNil, nil
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
