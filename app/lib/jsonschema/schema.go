package jsonschema

import "projectforge.dev/projectforge/app/util"

type Schema struct {
	// core vocabulary & metadata
	Schema        string                    `json:"$schema,omitzero"`        // uri identifying the dialect ("https://json-schema.org/draft/2020-12/schema")
	MetaID        string                    `json:"$id,omitzero"`            // base uri for the schema
	ExplicitID    string                    `json:"id,omitzero"`             // base uri for the schema
	Anchor        string                    `json:"$anchor,omitzero"`        // an identifier for this subschema
	Ref           string                    `json:"$ref,omitzero"`           // reference to another schema (uri or json pointer)
	DynamicRef    string                    `json:"$dynamicRef,omitzero"`    // reference that resolves dynamically (requires $dynamicanchor)
	DynamicAnchor string                    `json:"$dynamicAnchor,omitzero"` // anchor for dynamic resolution
	Vocabulary    *util.OrderedMap[bool]    `json:"$vocabulary,omitzero"`    // declares vocabularies used (keys are uris, values must be true)
	Comment       string                    `json:"$comment,omitzero"`       // a comment string, ignored by validators
	MetaDefs      *util.OrderedMap[*Schema] `json:"$defs,omitzero"`          // definitions for reusable subschemas
	ExplicitDefs  *util.OrderedMap[*Schema] `json:"definitions,omitzero"`    // definitions for reusable subschemas

	// annotations (metadata keywords)
	Title       string `json:"title,omitzero"`       // a short description of the schema
	Description string `json:"description,omitzero"` // a more detailed explanation
	Default     any    `json:"default,omitzero"`     // default value for the instance
	Deprecated  bool   `json:"deprecated,omitzero"`  // indicates the instance is deprecated (default: false)
	ReadOnly    bool   `json:"readOnly,omitzero"`    // indicates the instance should not be modified (default: false)
	WriteOnly   bool   `json:"writeOnly,omitzero"`   // indicates the instance may be set but should not be returned (default: false)

	// generic validation keywords
	Type  any   `json:"type,omitzero"`  // expected data type(s) (string or array of "string", "number", "integer", "object", "array", "boolean", "null")
	Enum  []any `json:"enum,omitempty"` // array of allowed values
	Const any   `json:"const,omitzero"` // exact required value

	// validation keywords for numbers (number and integer)
	MultipleOf       *float64 `json:"multipleOf,omitzero"`       // instance must be divisible by this number (strictly positive)
	Maximum          *float64 `json:"maximum,omitzero"`          // maximum inclusive value
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitzero"` // maximum exclusive value
	Minimum          *float64 `json:"minimum,omitzero"`          // minimum inclusive value
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitzero"` // minimum exclusive value

	// validation keywords for strings
	MaxLength        *uint64 `json:"maxLength,omitzero"`        // maximum length (non-negative integer)
	MinLength        *uint64 `json:"minLength,omitzero"`        // minimum length (non-negative integer, default 0)
	Pattern          string  `json:"pattern,omitzero"`          // ecma 262 regular expression
	Format           string  `json:"format,omitzero"`           // predefined format (e.g., "date-time", "email", "ipv4")
	ContentEncoding  string  `json:"contentEncoding,omitzero"`  // encoding of the string content (e.g., "base64")
	ContentMediaType string  `json:"contentMediaType,omitzero"` // media type of the string content (e.g., "application/json")
	ContentSchema    *Schema `json:"contentSchema,omitzero"`    // schema for the decoded content if contentencoding/mediatype are present

	// validation keywords for arrays
	Items            any     `json:"items,omitzero"`            // schema for array items (schema or boolean false). applied *after* prefixitems
	PrefixItems      Schemas `json:"prefixItems,omitzero"`      // array of schemas for tuple validation (items at specific indices)
	UnevaluatedItems any     `json:"unevaluatedItems,omitzero"` // validation for items not covered by `items` or `prefixitems` (schema or boolean)
	MaxItems         *uint64 `json:"maxItems,omitzero"`         // maximum number of items (non-negative integer)
	MinItems         *uint64 `json:"minItems,omitzero"`         // minimum number of items (non-negative integer, default 0)
	UniqueItems      *bool   `json:"uniqueItems,omitzero"`      // whether all items must be unique (default: false)
	Contains         *Schema `json:"contains,omitzero"`         // schema that at least one item must match
	MaxContains      *uint64 `json:"maxContains,omitzero"`      // maximum number of items matching `contains` (non-negative integer)
	MinContains      *uint64 `json:"minContains,omitzero"`      // minimum number of items matching `contains` (non-negative integer, default 1)

	// validation keywords for objects
	Properties            *util.OrderedMap[*Schema]  `json:"properties,omitzero"`            // schemas for named properties
	PatternProperties     *util.OrderedMap[*Schema]  `json:"patternProperties,omitzero"`     // schemas for properties matching regex patterns
	AdditionalProperties  any                        `json:"additionalProperties,omitzero"`  // schema or boolean
	AllowTrailingCommas   bool                       `json:"allowTrailingCommas,omitzero"`   // indicates that trailing commas are allowed
	UnevaluatedProperties any                        `json:"unevaluatedProperties,omitzero"` // validation for properties not covered by `properties`
	Required              []string                   `json:"required,omitempty"`             // array of required property names
	PropertyNames         *Schema                    `json:"propertyNames,omitzero"`         // schema for property names
	MaxProperties         *uint64                    `json:"maxProperties,omitzero"`         // maximum number of properties (non-negative integer)
	MinProperties         *uint64                    `json:"minProperties,omitzero"`         // minimum number of properties (non-negative integer, default 0)
	DependentRequired     *util.OrderedMap[[]string] `json:"dependentRequired,omitzero"`     // properties required based on the presence of other properties
	DependentSchemas      *util.OrderedMap[*Schema]  `json:"dependentSchemas,omitzero"`      // schemas applied based on the presence of other properties

	// conditional applicators
	If   *Schema `json:"if,omitzero"`   // if this schema validates, `then` must also validate
	Then *Schema `json:"then,omitzero"` // schema applied if `if` validates
	Else *Schema `json:"else,omitzero"` // schema applied if `if` does not validate

	// boolean logic applicators
	AllOf Schemas `json:"allOf,omitempty"` // instance must validate against all of these schemas
	AnyOf Schemas `json:"anyOf,omitempty"` // instance must validate against at least one of these schemas
	OneOf Schemas `json:"oneOf,omitempty"` // instance must validate against exactly one of these schemas
	Not   *Schema `json:"not,omitempty"`   // instance must not validate against this schema

	// examples
	Examples []any `json:"examples,omitempty"` // array of example values that validate against the schema

	// extensions
	Metadata util.ValueMap `json:"metadata,omitempty"` // additional information about the schema
}

func (s *Schema) ID() string {
	if s.MetaID != "" {
		return s.MetaID
	}
	return s.ExplicitID
}
