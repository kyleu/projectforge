package jsonschema

import (
	"github.com/samber/lo"
	"path"
	"strings"
)

type Collection struct {
	BaseURI string  `json:"baseURI,omitempty"`
	Schemas Schemas `json:"schemas,omitempty"`
}

func NewCollection(baseURI string) *Collection {
	return &Collection{BaseURI: baseURI}
}

func (c *Collection) GetSchema(id string) *Schema {
	u := strings.TrimPrefix(strings.TrimSuffix(id, ".schema.json"), c.BaseURI)
	ret := lo.FindOrElse(c.Schemas, nil, func(x *Schema) bool {
		return x.ID == id || x.ID == u
	})
	return ret
}

func (c *Collection) NewSchema(id string) *Schema {
	u := id
	if !strings.HasSuffix(u, ".schema.json") {
		u += ".schema.json"
	}
	if !strings.HasPrefix(u, c.BaseURI) {
		u = path.Join(c.BaseURI, u)
	}
	ret := &Schema{Schema: "https://json-schema.org/draft/2020-12/schema", ID: u}
	c.AddSchema(ret)
	return ret
}

func (c *Collection) AddSchema(sch *Schema) {
	c.Schemas = append(c.Schemas, sch)
}

func (c *Collection) Extra() []string {
	return c.Schemas.IDs()
}
