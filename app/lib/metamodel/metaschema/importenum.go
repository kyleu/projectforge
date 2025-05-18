package metaschema

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/util"
)

func ImportEnum(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*enum.Enum, error) {
	_, err := exportGetType(sch, "string")
	if err != nil {
		return nil, err
	}
	if len(sch.Enum) == 0 {
		return nil, errors.Errorf("string type provided without enum values in schema [%s]", sch.String())
	}
	n, pkg, grp := parseID(sch.ID())
	ret := &enum.Enum{Name: n, Package: pkg, Group: grp, Description: sch.Description}

	extra := sch.Metadata.GetMapOpt("values")
	for _, x := range sch.Enum {
		v, e := ValueFromAny(x)
		if e != nil {
			return nil, e
		}
		if v != nil && v.Description == "" && extra != nil && len(extra) > 0 {
			if curr, ok := extra[v.Key]; ok {
				newValue, ee := ValueFromAny(curr)
				if ee != nil {
					return nil, ee
				}
				newValue.Key = v.Key
				v = newValue
			}
		}
		ret.Values = append(ret.Values, v)
	}

	if x := sch.Metadata.GetStringOpt("icon"); x != "" {
		ret.Icon = x
	}
	if x := sch.Metadata.GetMapOpt("config"); len(x) > 0 {
		ret.Config = x
	}
	if x := sch.Metadata.GetStringOpt("proper"); x != "" {
		ret.ProperOverride = x
	}
	if x := sch.Metadata.GetStringOpt("route"); len(x) > 0 {
		ret.RouteOverride = x
	}
	if x := sch.Metadata.GetStringArrayOpt("tags"); len(x) > 0 {
		ret.Tags = x
	}
	if x := sch.Metadata.GetStringOpt("title"); x != "" {
		ret.TitleOverride = x
	}

	return ret, nil
}

func ValueFromAny(x any) (*enum.Value, error) {
	switch t := x.(type) {
	case []byte:
		v := &enum.Value{}
		if err := util.FromJSONStrict(t, v); err != nil {
			return nil, errors.Wrapf(err, "unhandled byte array for enum: %v", string(t))
		}
		return v, nil
	case string:
		return &enum.Value{Key: t}, nil
	case *enum.Value:
		return t, nil
	case map[string]any:
		v := &enum.Value{}
		if err := util.FromJSONStrict(util.ToJSONBytes(t, true), v); err != nil {
			return nil, errors.Wrapf(err, "invalid fields for [%T]: %v", x, x)
		}
		return v, nil
	case *util.OrderedMap[any]:
		v := &enum.Value{}
		if err := util.FromJSONStrict(util.ToJSONBytes(t.Map, true), v); err != nil {
			return nil, errors.Wrapf(err, "invalid fields for [%T]: %v", x, x)
		}
		return v, nil
	default:
		return nil, errors.Errorf("unhandled enum type [%T]: %v", x, x)
	}
}
