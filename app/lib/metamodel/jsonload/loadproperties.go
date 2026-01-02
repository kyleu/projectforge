package jsonload

import (
	"fmt"
	"slices"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func parseProperties(sch *jsonschema.Schema) (model.Columns, error) {
	if sch == nil || sch.Properties == nil {
		return nil, nil
	}
	ret := make(model.Columns, 0, len(sch.Properties.Order))
	for _, prop := range sch.Properties.Order {
		req := slices.Contains(sch.Required, prop)
		col, err := parseProperty(sch, prop, req, sch.Properties.GetSimple(prop))
		if err != nil {
			return nil, err
		}
		ret = append(ret, col)
	}
	return ret, nil
}

func parseProperty(sch *jsonschema.Schema, name string, required bool, prop *jsonschema.Schema) (*model.Column, error) {
	if prop == nil {
		return nil, nil
	}
	ret := &model.Column{Name: name, TitleOverride: prop.Title, Type: parseType(prop, required), Nullable: !required}
	ret.Comment = util.Choose(prop.Description != "", prop.Description, prop.Comment)
	md := prop.GetMetadata()
	for k := range md {
		switch k {
		case "pk":
			ret.PK = md.GetBoolOpt(k)
		case "search":
			ret.Search = md.GetBoolOpt(k)
		case "default":
			ret.SQLDefault = md.GetStringOpt(k)
		case "indexed":
			ret.Indexed = md.GetBoolOpt(k)
		case "display":
			ret.Display = md.GetStringOpt(k)
		case "format":
			ret.Format = md.GetStringOpt(k)
		case "json":
			ret.JSON = md.GetStringOpt(k)
		case "sql":
			ret.SQLOverride = md.GetStringOpt(k)
		case "title":
			ret.TitleOverride = md.GetStringOpt(k)
		case "plural":
			ret.PluralOverride = md.GetStringOpt(k)
		case "proper":
			ret.ProperOverride = md.GetStringOpt(k)
		case "example":
			ret.Example = md.GetStringOpt(k)
		case "validation":
			ret.Validation = md.GetStringOpt(k)
		case "values":
			ret.Values = md.GetStringArrayOpt(k)
		case "tags":
			ret.Tags = md.GetStringArrayOpt(k)
		case "help":
			ret.Help = md.GetStringOpt(k)
		default:
			if ret.Metadata == nil {
				ret.Metadata = make(util.ValueMap)
			}
			ret.Metadata[k] = md[k]
		}
	}
	if ret.SQLDefault == "" && prop.Default != nil {
		ret.SQLDefault = fmt.Sprint(prop.Default)
	}
	return ret, nil
}

func parseType(prop *jsonschema.Schema, required bool) *types.Wrapped {
	switch st := prop.DetectSchemaType(); st {
	case jsonschema.SchemaTypeNull:
		return types.NewNil()
	case jsonschema.SchemaTypeBoolean:
		return types.NewBool()
	case jsonschema.SchemaTypeInteger:
		return types.NewInt(0)
	case jsonschema.SchemaTypeNumber:
		return types.NewFloat(0)
	case jsonschema.SchemaTypeString:
		if prop.Format == "uuid" || prop.Format == util.UUIDRegex {
			return types.NewUUID()
		}
		if prop.Format == "date-time" {
			return types.NewTimestamp()
		}
		return types.NewString()
	case jsonschema.SchemaTypeDate:
		return types.NewTimestamp()
	case jsonschema.SchemaTypeObject:
		return parseObj(prop)
	case jsonschema.SchemaTypeEnum:
		return types.NewEnum("TODO")
	case jsonschema.SchemaTypeArray:
		return types.NewList(parseType(prop.Items, true))
	case jsonschema.SchemaTypeRef:
		return parseRef(prop, required)
	case jsonschema.SchemaTypeUnion:
		return types.NewUnion(types.NewAny())
	case jsonschema.SchemaTypeEmpty:
		return types.NewAny()
	case jsonschema.SchemaTypeNot:
		return types.NewError("TODO: not")
	case jsonschema.SchemaTypeUnknown:
		return types.NewUnknown("unknown")
	default:
		return types.NewError(fmt.Sprintf("unhandled type [%s]", st))
	}
}

func parseObj(prop *jsonschema.Schema) *types.Wrapped {
	if !prop.HasProperties() {
		return types.NewStringKeyedMap()
	}
	return types.NewReferenceArgs(util.Pkg{"obj"}, "TODO")
}

func parseRef(prop *jsonschema.Schema, required bool) *types.Wrapped {
	switch prop.Ref {
	case "numeric":
		return types.NewNumeric()
	case "numericMap":
		return types.NewNumericMap()
	default:
		ref := util.Str(prop.Ref)
		var hasRef bool
		if ref.HasPrefix("ref:") {
			ref = ref.TrimPrefix("ref:")
			hasRef = true
		}
		p := util.NewPkgSplit(ref.String(), "/")
		n := p.Last()
		if hasRef && !required {
			n = "*" + n
		}
		return types.NewReferenceArgs(p.Shift(), n)
	}
}
