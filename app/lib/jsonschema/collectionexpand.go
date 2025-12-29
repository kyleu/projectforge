package jsonschema

import "github.com/pkg/errors"

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
	t := sch.DetectSchemaType()
	makeRef := func() *Schema {
		if sch.ID() == "" {
			sch.MetaID = k
		}
		ret := NewRefSchema(key)
		ret.Key = key + "--expanded-ref"
		return ret
	}
	if t.Matches(SchemaTypeObject) && sch.Properties.Length() > 0 {
		return makeRef()
	}
	if t.Matches(SchemaTypeEnum) && len(sch.Enum) > 0 {
		return makeRef()
	}
	return nil
}

func simplify(sch *Schema) *Schema {
	if len(sch.AnyOf) == 1 && len(sch.AllOf) == 0 && len(sch.OneOf) == 0 {
		return sch.AnyOf[0]
	}
	if len(sch.AnyOf) == 0 && len(sch.AllOf) == 1 && len(sch.OneOf) == 0 {
		return sch.AllOf[0]
	}
	if len(sch.AnyOf) == 0 && len(sch.AllOf) == 0 && len(sch.OneOf) == 1 {
		return sch.OneOf[0]
	}
	return sch
}

func expandSchema(sch *Schema, key string) (Schemas, error) {
	if len(key) > 1024 {
		return nil, errors.Errorf("recursion limit reached for key [%s]", key)
	}
	orig := simplify(sch.Clone())
	if orig.Key == "" {
		orig.Key = key
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
