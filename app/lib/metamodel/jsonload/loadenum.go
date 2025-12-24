package jsonload

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/util"
)

func exportEnum(ctx context.Context, sch *jsonschema.Schema, ret *metamodel.Args, s *jsonschema.Collection, logs *util.StringSlice) error {
	logs.Pushf("exporting enum schema [%s]", sch.ID())
	vals := make(enum.Values, 0, len(sch.Enum))

	md := sch.GetMetadata()
	valsExtra := md.GetMapOpt("values")
	for _, v := range sch.Enum {
		e, err := exportEnumValue(v, valsExtra)
		if err != nil {
			return err
		}
		vals = append(vals, e)
	}
	name, pkg, grp, err := extractPath(sch)
	if err != nil {
		return err
	}
	e := &enum.Enum{
		Name: name, Package: pkg, Group: grp,
		Schema:         md.GetStringOpt("schema"),
		Description:    sch.Description,
		Icon:           md.GetStringOpt("icon"),
		Values:         vals,
		Tags:           md.GetStringArrayOpt("tags"),
		TitleOverride:  md.GetStringOpt("title"),
		ProperOverride: md.GetStringOpt("proper"),
		RouteOverride:  md.GetStringOpt("route"),
		Config:         md.GetMapOpt("config"),
	}
	ret.Enums = append(ret.Enums, e)
	return nil
}

func exportEnumValue(v any, valsExtra util.ValueMap) (*enum.Value, error) {
	switch t := v.(type) {
	case string:
		ex := valsExtra.GetMapOpt(t)
		var extra *util.OrderedMap[any]
		if m := ex.GetMapOpt("extra"); len(m) > 0 {
			extra = util.NewOMap[any]()
			extra.SetAll(m)
		}
		return &enum.Value{
			Key:         t,
			Name:        ex.GetStringOpt("name"),
			Description: ex.GetStringOpt("description"),
			Icon:        ex.GetStringOpt("icon"),
			Extra:       extra,
			Default:     ex.GetBoolOpt("default"),
			Simple:      false,
		}, nil
	default:
		return nil, errors.Errorf("unsupported enum value type [%T]", v)
	}
}
