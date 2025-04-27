package jsonschema

import (
	"fmt"
	"github.com/samber/lo"
	"path"
	"projectforge.dev/projectforge/app/util"
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
	little := strings.TrimPrefix(strings.TrimSuffix(id, ".schema.json"), c.BaseURI)
	big := id
	if !strings.HasSuffix(big, ".schema.json") {
		big += ".schema.json"
	}
	if !strings.HasPrefix(big, c.BaseURI) {
		big = path.Join(c.BaseURI, big)
	}
	ret := lo.FindOrElse(c.Schemas, nil, func(x *Schema) bool {
		return x.ID == little || x.ID == big
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
	comment := fmt.Sprintf("managed by %s", util.AppName)
	ret := &Schema{Schema: "https://json-schema.org/draft/2020-12/schema", ID: u, Comment: comment}
	c.AddSchema(ret)
	return ret
}

func (c *Collection) AddSchema(sch *Schema) {
	c.Schemas = append(c.Schemas, sch)
}

func (c *Collection) Extra() Schemas {
	return lo.Filter(c.Schemas, func(x *Schema, _ int) bool {
		return !strings.Contains(x.Comment, util.AppName)
	})
}
