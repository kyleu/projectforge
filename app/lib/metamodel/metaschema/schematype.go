package metaschema

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/types"
)

func TypeSchema(typ types.Type, coll *jsonschema.Collection, args *metamodel.Args) (*jsonschema.Schema, error) {
	ret := &jsonschema.Schema{Type: typ.String()}
	switch t := typ.(type) {
	case *types.Any:
		ret.Type = "object" // TODO: Hoo boy...
	case *types.Bool:
		ret.Type = "boolean"
	case *types.Date:
		ret.Type = "string"
		ret.Format = "date"
	case *types.Enum:
		ret.Type = "enum" // TODO: Ref
	case *types.Float:
		ret.Type = "number"
	case *types.Int:
		ret.Type = "integer"
	case *types.List:
		ret.Type = "array" // TODO: Ref
	case *types.Map:
		ret.Type = "object"
	case *types.Reference:
		ret.Type = "object" // TODO: Who knows?
	case *types.String:
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
		ret.Type = "object" // TODO: Ref
	case *types.Wrapped:
		return TypeSchema(t.T, coll, args)
	default:
		return nil, errors.Errorf("unhandled JSON Schema type [%T]", t)
	}
	return ret, nil
}
