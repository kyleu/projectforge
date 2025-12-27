package jsonschema

import "projectforge.dev/projectforge/app/util"

type DataArray struct {
	Items            *Schema `json:"items,omitzero"`            // schema for array items (schema or boolean false). applied *after* prefixItems
	PrefixItems      Schemas `json:"prefixItems,omitzero"`      // array of schemas for tuple validation (items at specific indices)
	UnevaluatedItems *Schema `json:"unevaluatedItems,omitzero"` // validation for items not covered by `items` or `prefixItems` (schema or boolean)
	MaxItems         *uint64 `json:"maxItems,omitzero"`         // maximum number of items (non-negative integer)
	MinItems         *uint64 `json:"minItems,omitzero"`         // minimum number of items (non-negative integer, default 0)
	UniqueItems      *bool   `json:"uniqueItems,omitzero"`      // whether all items must be unique (default: false)
	Contains         *Schema `json:"contains,omitzero"`         // schema that at least one item must match
	MaxContains      *uint64 `json:"maxContains,omitzero"`      // maximum number of items matching `contains` (non-negative integer)
	MinContains      *uint64 `json:"minContains,omitzero"`      // minimum number of items matching `contains` (non-negative integer, default 1)
}

func (d DataArray) IsEmpty() bool {
	return d.Items == nil && len(d.PrefixItems) == 0 && d.UnevaluatedItems == nil && d.MaxItems == nil &&
		d.MinItems == nil && d.UniqueItems == nil && d.Contains == nil && d.MaxContains == nil && d.MinContains == nil
}

func (d DataArray) Clone() DataArray {
	return DataArray{
		Items: d.Items, PrefixItems: util.ArrayCopy(d.PrefixItems), UnevaluatedItems: d.UnevaluatedItems,
		MaxItems: d.MaxItems, MinItems: d.MinItems, UniqueItems: d.UniqueItems,
		Contains: d.Contains.Clone(), MaxContains: d.MaxContains, MinContains: d.MinContains,
	}
}
