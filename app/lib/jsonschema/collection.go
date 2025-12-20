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
	Expanded bool    `json:"expand,omitempty"`
	Schemas  Schemas `json:"schemas,omitempty"`
}

func NewCollection(expand bool) *Collection {
	return &Collection{Expanded: expand}
}

func (c *Collection) GetSchema(id string) *Schema {
	little := strings.TrimPrefix(strings.TrimSuffix(id, KeyExtension), "/")
	big := id
	if !strings.HasSuffix(big, KeyExtension) {
		big += KeyExtension
	}
	medium := strings.TrimPrefix(big, "/")
	ret := lo.FindOrElse(c.Schemas, nil, func(x *Schema) bool {
		chk := x.ID()
		return chk == little || chk == medium || chk == big
	})
	return ret
}

func (c *Collection) NewSchema(id string) *Schema {
	u := id
	if !strings.HasSuffix(u, KeyExtension) {
		u += KeyExtension
	}
	comment := fmt.Sprintf("managed by %s", util.AppName)
	ret := &Schema{data: data{dataCore: dataCore{Schema: CurrentSchemaVersion, MetaID: u, Comment: comment}}}
	_ = c.AddSchema(ret)
	return ret
}

func (c *Collection) AddSchema(sch ...*Schema) error {
	if c.Expanded {
		ret := make(Schemas, 0, len(sch))
		for _, x := range sch {
			exp, err := expandSchema(x, 0)
			if err != nil {
				return err
			}
			ret = append(ret, exp...)
		}
		c.Schemas = ret
	} else {
		c.Schemas = append(c.Schemas, sch...)
	}
	return nil
}

func expandSchema(sch *Schema, recursionLevel int) (Schemas, error) {
	if recursionLevel > 100 {
		return nil, errors.Errorf("recursion limit reached")
	}
	orig := sch.Clone()
	ret := Schemas{orig}
	for _, k := range orig.Properties.Keys() {
		x := orig.Properties.GetSimple(k)
		if n := shouldExpand(k, x); n != nil {
			orig.Properties.Set(k, n)
			exp, err := expandSchema(x, recursionLevel+1)
			if err != nil {
				return nil, err
			}
			ret = append(ret, exp...)
		}
	}
	defs := orig.Definitions()
	for _, k := range defs.Keys() {
		x := defs.GetSimple(k)
		if n := shouldExpand(k, x); n != nil {
			defs.Set(k, n)
			exp, err := expandSchema(x, recursionLevel+1)
			if err != nil {
				return nil, err
			}
			ret = append(ret, exp...)
		}
	}
	return ret, nil
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

func (c *Collection) Extra() Schemas {
	return lo.Reject(c.Schemas, func(x *Schema, _ int) bool {
		return strings.Contains(x.Comment, util.AppName)
	})
}
