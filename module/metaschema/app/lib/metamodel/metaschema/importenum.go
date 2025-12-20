package metaschema

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/lib/metamodel/enum"
	"{{{ .Package }}}/app/util"
)

func ImportEnum(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*enum.Enum, error) {
	_, err := exportGetType(sch, KeyString)
	if err != nil {
		return nil, err
	}
	if len(sch.Enum) == 0 {
		return nil, errors.Errorf("string type provided without enum values in schema [%s]", sch.String())
	}
	n, pkg, grp := parseID(sch.ID())
	ret := &enum.Enum{Name: n, Package: pkg, Group: grp, Description: sch.Description}

	md := sch.GetMetadata()
	extra := md.GetMapOpt("values")
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

	if x := md.GetStringOpt("schema"); x != "" {
		ret.Schema = x
	}
	if x := md.GetStringOpt("icon"); x != "" {
		ret.Icon = x
	}
	if x := md.GetStringArrayOpt("tags"); len(x) > 0 {
		ret.Tags = x
	}
	if x := md.GetStringOpt("title"); x != "" {
		ret.TitleOverride = x
	}
	if x := md.GetStringOpt("proper"); x != "" {
		ret.ProperOverride = x
	}
	if x := md.GetStringOpt("route"); x != "" {
		ret.RouteOverride = x
	}
	if x := md.GetMapOpt("config"); len(x) > 0 {
		ret.Config = x
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
