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
	ret := &Schema{data: data{dataCore: dataCore{Schema: CurrentSchemaVersion, MetaID: u, Comment: comment}}}
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

func (c *Collection) AddSchemaExpanded(sch ...*Schema) error {
	for _, x := range sch {
		exp, err := expandSchema(x, x.ID())
		if err != nil {
			return err
		}
		if len(exp) > 0 {
			if err := c.AddSchema(true, exp[0]); err != nil {
				return err
			}
			if err := c.AddSchema(false, exp[1:]...); err != nil {
				return err
			}
		}
	}
	return nil
}

func shouldExpand(k string, sch *Schema, key string) *Schema {
	if sch.Type == "object" {
		if sch.ID() == "" {
			sch.MetaID = k
		}
		ret := NewRefSchema(key)
		ret.Key = key + "--expanded-ref"
		return ret
	}
	return nil
}

// "https://json.schemastore.org/cloudify.json#definitions/nodeTypeCloudifyAzureNodesNetworkLoadBalancerProbeInterfaces#properties/cloudify.interfaces.lifecycle#properties/delete--enum"
func expandSchema(sch *Schema, key string) (Schemas, error) {
	if len(key) > 1024 {
		return nil, errors.Errorf("recursion limit reached for key [%s]", key)
	}
	orig := sch.Clone()
	if orig.Key == "" {
		orig.Key = key
	}
	t := orig.DetectSchemaType()
	if t.Matches(SchemaTypeEnum) {
		orig.Key += "--enum"
		ref := NewRefSchema(orig.Key)
		ref.Key = key + "--expanded-ref"
		ret := Schemas{ref, orig}
		return ret, nil
	}
	ret := Schemas{orig}
	process := func(x *Schema, key string) error {
		x.Key = key
		exp, err := expandSchema(x, key)
		if err != nil {
			return err
		}
		ret = append(ret, exp...)
		return nil
	}
	for _, k := range orig.Properties.Keys() {
		x := orig.Properties.GetSimple(k)
		newKey := key + "#properties/" + k
		if n := shouldExpand(k, x, newKey); n != nil {
			orig.Properties.Set(k, n)
			if err := process(x, newKey); err != nil {
				return nil, err
			}
		}
	}
	defs := orig.Definitions()
	for _, k := range defs.Keys() {
		x := defs.GetSimple(k)
		newKey := key + "#definitions/" + k
		if n := shouldExpand(k, x, newKey); n != nil {
			defs.Set(k, n)
			if err := process(x, newKey); err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}
