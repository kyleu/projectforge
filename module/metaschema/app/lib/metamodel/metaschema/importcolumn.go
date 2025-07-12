package metaschema

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/lib/metamodel/model"
	"{{{ .Package }}}/app/lib/types"
)

func ImportColumn(key string, parent *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*model.Column, error) {
	prop := parent.Properties.Map[key]
	typ, err := ImportType(prop, coll, args)
	if err != nil {
		return nil, err
	}
	nullable := !slices.Contains(parent.Required, key)
	col := &model.Column{Name: key, Type: typ, Nullable: nullable}
	if prop.Metadata != nil {
		for _, k := range prop.Metadata.Keys() {
			switch k {
			case "pk":
				col.PK = prop.Metadata.GetBoolOpt(k)
			case "search":
				col.Search = prop.Metadata.GetBoolOpt(k)
			case "indexed":
				col.Indexed = prop.Metadata.GetBoolOpt(k)
			case "display":
				col.Display = prop.Metadata.GetStringOpt(k)
			case "format":
				col.Format = prop.Metadata.GetStringOpt(k)
			case "json":
				col.JSON = prop.Metadata.GetStringOpt(k)
			case "sql":
				col.SQLOverride = prop.Metadata.GetStringOpt(k)
			case "title":
				col.TitleOverride = prop.Metadata.GetStringOpt(k)
			case "plural":
				col.PluralOverride = prop.Metadata.GetStringOpt(k)
			case "proper":
				col.ProperOverride = prop.Metadata.GetStringOpt(k)
			case "example":
				col.Example = prop.Metadata.GetStringOpt(k)
			case "validation":
				col.Validation = prop.Metadata.GetStringOpt(k)
			case "values":
				col.Values = prop.Metadata.GetStringArrayOpt(k)
			case "tags":
				col.Tags = prop.Metadata.GetStringArrayOpt(k)
			case "help":
				col.HelpString = prop.Metadata.GetStringOpt(k)
			case "bits":
				switch t := col.Type.T.(type) {
				case *types.Int:
					t.Bits = prop.Metadata.GetIntOpt(k)
				case *types.Float:
					t.Bits = prop.Metadata.GetIntOpt(k)
				default:
					return nil, errors.Errorf("incorrect [bits] metadata for type [%s]", col.Type.String())
				}
			case "type":
				// handled above
			default:
				return nil, errors.Errorf("unhandled metadata key [%s]", k)
			}
		}
	}
	if prop.Default != nil {
		col.SQLDefault = fmt.Sprint(prop.Default)
	}
	return col, nil
}
