package jsonschema

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

const (
	CurrentSchemaVersion = "https://json-schema.org/draft/2020-12/schema"
	KeyExtension         = ".schema.json"
)

type Collection struct {
	SchemaMap map[string]*Schema `json:"schemas,omitempty"`
	Roots     []string           `json:"roots,omitempty"`
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
	ret := NewSchema(DataCore{Schema: CurrentSchemaVersion, MetaID: u, Comment: comment})
	_ = c.AddSchema(true, ret)
	return ret
}

func (c *Collection) AddSchema(root bool, sch ...*Schema) error {
	for _, x := range sch {
		id := x.ID()
		if _, ok := c.SchemaMap[id]; ok {
			return errors.Errorf("schema with id [%s] already exists", id)
		}
		c.SchemaMap[id] = x
		if root {
			c.Roots = append(c.Roots, id)
		}
	}
	return nil
}
