package jsonload

import (
	"context"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
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
