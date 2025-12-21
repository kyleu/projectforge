package jsonschema

import "projectforge.dev/projectforge/app/util"

// JSON fields that represent core vocabulary & metadata
type dataCore struct {
	Schema        string                    `json:"$schema,omitzero"`        // uri identifying the dialect ("https://json-schema.org/draft/2020-12/schema")
	MetaID        string                    `json:"$id,omitzero"`            // base uri for the schema
	ExplicitID    string                    `json:"id,omitzero"`             // older [id] key
	Anchor        string                    `json:"$anchor,omitzero"`        // an identifier for this subschema
	Ref           string                    `json:"$ref,omitzero"`           // reference to another schema (uri or json pointer)
	DynamicRef    string                    `json:"$dynamicRef,omitzero"`    // reference that resolves dynamically (requires $dynamicAnchor)
	DynamicAnchor string                    `json:"$dynamicAnchor,omitzero"` // anchor for dynamic resolution
	Vocabulary    *util.OrderedMap[bool]    `json:"$vocabulary,omitzero"`    // declares vocabularies used (keys are uris, values must be true)
	Comment       string                    `json:"$comment,omitzero"`       // a comment string, ignored by validators
	Defs          *util.OrderedMap[*Schema] `json:"$defs,omitzero"`          // definitions for reusable subschemas
	ExplicitDefs  *util.OrderedMap[*Schema] `json:"definitions,omitzero"`    // older [definition] key
}

func (d dataCore) IsEmpty() bool {
	return d.Schema == "" && d.MetaID == "" && d.ExplicitID == "" && d.Anchor == "" &&
		d.Ref == "" && d.DynamicRef == "" && d.DynamicAnchor == "" && d.Vocabulary.Empty() &&
		d.Comment == "" && d.Defs.Empty() && d.ExplicitDefs.Empty()
}

func (d dataCore) Clone() dataCore {
	return dataCore{
		Schema: d.Schema, MetaID: d.MetaID, ExplicitID: d.ExplicitID, Anchor: d.Anchor,
		Ref: d.Ref, DynamicRef: d.DynamicRef, DynamicAnchor: d.DynamicAnchor, Vocabulary: d.Vocabulary.Clone(),
		Comment: d.Comment, Defs: d.Defs.Clone(), ExplicitDefs: d.ExplicitDefs.Clone(),
	}
}
