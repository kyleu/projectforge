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
	logs.Push("OK!")
	return logs.Slice, ret, nil
}
