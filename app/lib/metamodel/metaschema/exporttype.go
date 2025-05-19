package metaschema

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/types"
)

func ExportType(typ types.Type, coll *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	ret := &jsonschema.Schema{}
	switch t := typ.(type) {
	case *types.Any:
		ret.Type = nil // TODO: Hoo boy...
	case *types.Bool:
		ret.Type = "boolean"
	case *types.Date:
		ret.Type = "string"
		ret.Format = "date"
	case *types.Enum:
		ret.Ref = t.Ref + ".schema.json"
	case *types.Float:
		ret.Type = "number"
		if t.Bits != 0 {
			ret.AddMetadata("bits", t.Bits)
		}
	case *types.Int:
		ret.Type = "integer"
		if t.Bits != 0 {
			ret.AddMetadata("bits", t.Bits)
		}
	case *types.List:
		ret.Type = "array"
		if t.V.Scalar() && t.V.EnumKey() == "" {
			ret.Items = types.ToJSONType(t.V)
		} else {
			itemSchema, e := ExportType(t.V, coll, args)
			if e != nil {
				return nil, e
			}
			ret.Items = itemSchema
		}
	case *types.Map:
		ret.Type = "object"
	case *types.Numeric:
		ret.Ref = "Numeric"
	case *types.Reference:
		ret.Ref = t.String()
	case *types.String:
		ret.Type = typ.String()
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
		ret.Type = "string"
		ret.Format = "date-time"
	case *types.UUID:
		ret.Type = "string"
		ret.Format = "uuid"
	case *types.ValueMap:
		ret.Type = "object"
	case *types.Wrapped:
		return ExportType(t.T, coll, args)
	default:
		return nil, errors.Errorf("unhandled JSONSchema type [%T]", t)
	}
	return ret, nil
}
