package metaschema

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
)

func ImportColumn(key string, parent *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*model.Column, error) {
	prop := parent.Properties.Map[key]
	typ, err := ImportType(prop, coll, args)
	if err != nil {
		return nil, err
	}
	nullable := !slices.Contains(parent.Required, key)
	col := &model.Column{Name: key, Type: typ, Nullable: nullable}
	md := prop.GetMetadata()
	if md != nil {
		for _, k := range md.Keys() {
			switch k {
			case "pk":
				col.PK = md.GetBoolOpt(k)
			case "search":
				col.Search = md.GetBoolOpt(k)
			case "indexed":
				col.Indexed = md.GetBoolOpt(k)
			case "display":
				col.Display = md.GetStringOpt(k)
			case "format":
				col.Format = md.GetStringOpt(k)
			case "json":
				col.JSON = md.GetStringOpt(k)
			case "sql":
				col.SQLOverride = md.GetStringOpt(k)
			case "title":
				col.TitleOverride = md.GetStringOpt(k)
			case "plural":
				col.PluralOverride = md.GetStringOpt(k)
			case "proper":
				col.ProperOverride = md.GetStringOpt(k)
			case "example":
				col.Example = md.GetStringOpt(k)
			case "validation":
				col.Validation = md.GetStringOpt(k)
			case "values":
				col.Values = md.GetStringArrayOpt(k)
			case "tags":
				col.Tags = md.GetStringArrayOpt(k)
			case "comment":
				col.Comment = md.GetStringOpt(k)
			case "help":
				col.Help = md.GetStringOpt(k)
			case "bits":
				switch t := col.Type.T.(type) {
				case *types.Int:
					t.Bits = md.GetIntOpt(k)
				case *types.Float:
					t.Bits = md.GetIntOpt(k)
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
