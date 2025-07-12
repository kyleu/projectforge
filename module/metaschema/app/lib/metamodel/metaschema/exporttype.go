package metaschema

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/lib/types"
)

func ExportType(typ types.Type, coll *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	ret := &jsonschema.Schema{}
	switch t := typ.(type) {
	case *types.Any:
		ret.Type = nil
	case *types.Bool:
		ret.Type = KeyBoolean
	case *types.Date:
		ret.Type = KeyString
		ret.Format = KeyDate
	case *types.Enum:
		ret.Ref = t.Ref + ".schema.json"
	case *types.Float:
		ret.Type = KeyNumber
		if t.Bits != 0 {
			ret.AddMetadata("bits", t.Bits)
		}
	case *types.Int:
		ret.Type = KeyInteger
		if t.Bits != 0 {
			ret.AddMetadata("bits", t.Bits)
		}
	case *types.JSON:
		ret.Type = KeyObject
		ret.AddMetadata("type", "json")
	case *types.List:
		ret.Type = KeyArray
		if t.V.Scalar() && t.V.EnumKey() == "" {
			ret.Items = ToJSONType(t.V)
		} else {
			itemSchema, e := ExportType(t.V, coll, args)
			if e != nil {
				return nil, e
			}
			ret.Items = itemSchema
		}
	case *types.Map:
		ret.Type = KeyObject
	case *types.Numeric:
		ret.Ref = types.KeyNumeric
	case *types.Reference:
		ret.Ref = t.String()
	case *types.String:
		ret.Type = KeyString
		if ml := uint64(t.MinLength); ml > 0 {
			ret.MinLength = &ml
		}
		if ml := uint64(t.MaxLength); ml > 0 {
			ret.MaxLength = &ml
		}
		if t.Pattern != "" {
			ret.Pattern = t.Pattern
		}
	case *types.Timestamp:
		ret.Type = KeyString
		ret.Format = KeyDateTime
	case *types.TimestampZoned:
		ret.Type = KeyString
		ret.Format = KeyDateTime
		ret.AddMetadata("type", types.KeyTimestampZoned)
	case *types.UUID:
		ret.Type = KeyString
		ret.Format = KeyUUID
	case *types.ValueMap:
		ret.Type = KeyObject
	case *types.Wrapped:
		return ExportType(t.T, coll, args)
	default:
		return nil, errors.Errorf("unhandled JSONSchema type [%T]", t)
	}
	return ret, nil
}
