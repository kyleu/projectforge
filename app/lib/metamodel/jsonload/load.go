package jsonload

import (
	"context"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

type Validation struct {
	Collection *jsonschema.Collection
}

func (s *Validation) AddSchema(sch ...*jsonschema.Schema) error {
	return s.Collection.AddSchema(true, sch...)
}

func (s *Validation) Export(ctx context.Context, logger util.Logger) ([]string, *metamodel.Args, error) {
	logs := util.NewStringSlice()
	ret := &metamodel.Args{}
	if err := s.Collection.Validate(); err != nil {
		return nil, nil, err
	}
	schemata := s.Collection.Schemas()
	for _, sch := range schemata {
		t, err := sch.Validate()
		if err != nil {
			return nil, nil, err
		}
		switch {
		case t.Matches(jsonschema.SchemaTypeObject):
			if err := exportObject(ctx, sch, ret, s.Collection, logs); err != nil {
				return nil, nil, err
			}
		case t.Matches(jsonschema.SchemaTypeEnum):
			if err := exportEnum(ctx, sch, ret, s.Collection, logs); err != nil {
				return nil, nil, err
			}
		default:
			logs.Pushf("unsupported type [%s] for schema [%s]", t, sch.ID())
		}
	}
	logs.Push("OK!")
	return logs.Slice, ret, nil
}

func exportObject(ctx context.Context, sch *jsonschema.Schema, ret *metamodel.Args, s *jsonschema.Collection, logs *util.StringSlice) error {
	logs.Pushf("exporting object schema [%s]", sch.ID())
	m := &model.Model{
		Name: sch.ID(),
	}
	ret.Models = append(ret.Models, m)
	return nil
}

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

func extractPath(sch *jsonschema.Schema) (string, string, []string, error) {
	id := util.Str(sch.ID())
	var prefix util.RichString
	if idx := id.Index("://"); idx > -1 {
		prefix = id.Substring(0, idx+3)
		id = id.Substring(idx+3, id.Length())
	}
	if id.HasSuffix(".json") {
		id = id.TrimSuffix(".json")
	}
	if id.HasSuffix(".schema") {
		id = id.TrimSuffix(".schema")
	}

	parts := id.Split("/")
	partsLen := len(parts)
	n := parts[partsLen-1]
	var p util.RichString
	var g util.RichStrings
	if partsLen > 1 {
		p = parts[partsLen-2]
		g = parts[:partsLen-2]
		g[0] = prefix + g[0]
	}
	return n.String(), p.String(), g.Strings(), nil
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
