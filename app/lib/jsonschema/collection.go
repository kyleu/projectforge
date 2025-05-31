package jsonschema

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Collection struct {
	Schemas Schemas `json:"schemas,omitempty"`
}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) GetSchema(id string) *Schema {
	little := strings.TrimPrefix(strings.TrimSuffix(id, KeyExtension), "/")
	big := id
	if !strings.HasSuffix(big, KeyExtension) {
		big += KeyExtension
	}
	medium := strings.TrimPrefix(big, "/")
	ret := lo.FindOrElse(c.Schemas, nil, func(x *Schema) bool {
		return x.ID() == little || x.ID() == medium || x.ID() == big
	})
	return ret
}

func (c *Collection) NewSchema(id string) *Schema {
	u := id
	if !strings.HasSuffix(u, KeyExtension) {
		u += KeyExtension
	}
	comment := fmt.Sprintf("managed by %s", util.AppName)
	ret := &Schema{Schema: "https://json-schema.org/draft/2020-12/schema", MetaID: u, Comment: comment}
	c.AddSchema(ret)
	return ret
}

func (c *Collection) AddSchema(sch *Schema) {
	c.Schemas = append(c.Schemas, sch)
}

func (c *Collection) Extra() Schemas {
	return lo.Reject(c.Schemas, func(x *Schema, _ int) bool {
		return strings.Contains(x.Comment, util.AppName)
	})
}
