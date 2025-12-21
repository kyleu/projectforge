package jsonschema

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const (
	CurrentSchemaVersion = "https://json-schema.org/draft/2020-12/schema"
	KeyExtension         = ".schema.json"
)

type Collection struct {
	SchemaMap map[string]*Schema `json:"schemas,omitempty"`
}

func NewCollection() *Collection {
	return &Collection{SchemaMap: map[string]*Schema{}}
}

func (c *Collection) GetSchema(id string) *Schema {
	little := strings.TrimPrefix(strings.TrimSuffix(id, KeyExtension), "/")
	big := id
	if !strings.HasSuffix(big, KeyExtension) {
		big += KeyExtension
	}
	medium := strings.TrimPrefix(big, "/")

	if ret, ok := c.SchemaMap[id]; ok {
		return ret
	}
	if ret, ok := c.SchemaMap[little]; ok {
		return ret
	}
	if ret, ok := c.SchemaMap[medium]; ok {
		return ret
	}
	if ret, ok := c.SchemaMap[big]; ok {
		return ret
	}
	return nil
}

func (c *Collection) Keys() []string {
	return util.ArraySorted(lo.Keys(c.SchemaMap))
}

func (c *Collection) Schemas() Schemas {
	return lo.Map(util.ArraySorted(lo.Keys(c.SchemaMap)), func(x string, _ int) *Schema {
		return c.GetSchema(x)
	})
}

func (c *Collection) NewSchema(id string) *Schema {
	u := id
	if !strings.HasSuffix(u, KeyExtension) {
		u += KeyExtension
	}
	comment := fmt.Sprintf("managed by %s", util.AppName)
	ret := &Schema{data: data{dataCore: dataCore{Schema: CurrentSchemaVersion, MetaID: u, Comment: comment}}}
	c.AddSchema(ret)
	return ret
}

func (c *Collection) AddSchema(sch ...*Schema) {
	for _, x := range sch {
		c.SchemaMap[x.ID()] = x
	}
}

func (c *Collection) AddSchemaExpanded(sch ...*Schema) error {
	for _, x := range sch {
		exp, err := expandSchema(x, x.ID())
		if err != nil {
			return err
		}
		c.AddSchema(exp...)
	}
	return nil
}

func shouldExpand(k string, sch *Schema) *Schema {
	if sch.Type == "object" {
		if sch.ID() == "" {
			sch.MetaID = k
		}
		return NewRefSchema(k)
	}
	return nil
}

func expandSchema(sch *Schema, path ...string) (Schemas, error) {
	if len(path) > 100 {
		return nil, errors.Errorf("recursion limit reached")
	}
	orig := sch.Clone()
	ret := Schemas{orig}
	process := func(x *Schema, n ...string) error {
		p := append(util.ArrayCopy(path), n...)
		x.Key = strings.Join(p, "/")
		exp, err := expandSchema(x, p...)
		if err != nil {
			return err
		}
		ret = append(ret, exp...)
		return nil
	}
	for _, k := range orig.Properties.Keys() {
		x := orig.Properties.GetSimple(k)
		if n := shouldExpand(k, x); n != nil {
			orig.Properties.Set(k, n)
			if err := process(x, "properties", k); err != nil {
				return nil, err
			}
		}
	}
	defs := orig.Definitions()
	for _, k := range defs.Keys() {
		x := defs.GetSimple(k)
		if n := shouldExpand(k, x); n != nil {
			defs.Set(k, n)
			if err := process(x, "definitions", k); err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}
