package jsonschema

import (
	"projectforge.dev/projectforge/app/util"
)

type Schema struct {
	// core vocabulary & metadata
	Schema        string                    `json:"$schema,omitempty"`        // uri identifying the schema dialect (e.g., "https://json-schema.org/draft/2020-12/schema")
	MetaID        string                    `json:"$id,omitempty"`            // base uri for the schema
	ExplicitID    string                    `json:"id,omitempty"`             // base uri for the schema
	Anchor        string                    `json:"$anchor,omitempty"`        // an identifier for this subschema
	Ref           string                    `json:"$ref,omitempty"`           // reference to another schema (uri or json pointer)
	DynamicRef    string                    `json:"$dynamicRef,omitempty"`    // reference that resolves dynamically (requires $dynamicanchor)
	DynamicAnchor string                    `json:"$dynamicAnchor,omitempty"` // anchor for dynamic resolution
	Vocabulary    *util.OrderedMap[bool]    `json:"$vocabulary,omitempty"`    // declares vocabularies used (keys are uris, values must be true)
	Comment       string                    `json:"$comment,omitempty"`       // a comment string, ignored by validators
	MetaDefs      *util.OrderedMap[*Schema] `json:"$defs,omitempty"`          // definitions for reusable subschemas
	ExplicitDefs  *util.OrderedMap[*Schema] `json:"definitions,omitempty"`    // definitions for reusable subschemas

	// annotations (metadata keywords)
	Title       string `json:"title,omitempty"`       // a short description of the schema
	Description string `json:"description,omitempty"` // a more detailed explanation
	Default     any    `json:"default,omitempty"`     // default value for the instance
	Deprecated  bool   `json:"deprecated,omitempty"`  // indicates the instance is deprecated (default: false)
	ReadOnly    bool   `json:"readOnly,omitempty"`    // indicates the instance should not be modified (default: false)
	WriteOnly   bool   `json:"writeOnly,omitempty"`   // indicates the instance may be set but should not be returned (default: false)
	Examples    []any  `json:"examples,omitempty"`    // array of example values that validate against the schema

	// generic validation keywords
	Type  any   `json:"type,omitempty"`  // expected data type(s) (string or array of strings: "string", "number", "integer", "object", "array", "boolean", "null")
	Enum  []any `json:"enum,omitempty"`  // array of allowed values
	Const any   `json:"const,omitempty"` // exact required value

	// validation keywords for numbers (number and integer)
	MultipleOf       *float64 `json:"multipleOf,omitempty"`       // instance must be divisible by this number (strictly positive)
	Maximum          *float64 `json:"maximum,omitempty"`          // maximum inclusive value
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"` // maximum exclusive value
	Minimum          *float64 `json:"minimum,omitempty"`          // minimum inclusive value
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"` // minimum exclusive value

	// validation keywords for strings
	MaxLength        *uint64 `json:"maxLength,omitempty"`        // maximum length (non-negative integer)
	MinLength        *uint64 `json:"minLength,omitempty"`        // minimum length (non-negative integer, default 0)
	Pattern          string  `json:"pattern,omitempty"`          // ecma 262 regular expression
	Format           string  `json:"format,omitempty"`           // predefined format (e.g., "date-time", "email", "ipv4")
	ContentEncoding  string  `json:"contentEncoding,omitempty"`  // encoding of the string content (e.g., "base64")
	ContentMediaType string  `json:"contentMediaType,omitempty"` // media type of the string content (e.g., "application/json")
	ContentSchema    *Schema `json:"contentSchema,omitempty"`    // schema for the decoded content if contentencoding/mediatype are present

	// validation keywords for arrays
	Items            any     `json:"items,omitempty"`            // schema for array items (schema or boolean false). applied *after* prefixitems. use `unevaluateditems` for more control.
	PrefixItems      Schemas `json:"prefixItems,omitempty"`      // array of schemas for tuple validation (items at specific indices)
	UnevaluatedItems any     `json:"unevaluatedItems,omitempty"` // validation for items not covered by `items` or `prefixitems` (schema or boolean)
	MaxItems         *uint64 `json:"maxItems,omitempty"`         // maximum number of items (non-negative integer)
	MinItems         *uint64 `json:"minItems,omitempty"`         // minimum number of items (non-negative integer, default 0)
	UniqueItems      *bool   `json:"uniqueItems,omitempty"`      // whether all items must be unique (default: false)
	Contains         *Schema `json:"contains,omitempty"`         // schema that at least one item must match
	MaxContains      *uint64 `json:"maxContains,omitempty"`      // maximum number of items matching `contains` (non-negative integer, requires `contains`)
	MinContains      *uint64 `json:"minContains,omitempty"`      // minimum number of items matching `contains` (non-negative integer, default 1, requires `contains`)

	// validation keywords for objects
	Properties            *util.OrderedMap[*Schema]  `json:"properties,omitempty"`            // schemas for named properties
	PatternProperties     *util.OrderedMap[*Schema]  `json:"patternProperties,omitempty"`     // schemas for properties matching regex patterns
	AdditionalProperties  any                        `json:"additionalProperties,omitempty"`  // controls handling of properties not explicitly listed or matched by patterns (schema or boolean). use `unevaluatedproperties` for more control.
	AllowTrailingCommas   bool                       `json:"allowTrailingCommas,omitempty"`   // indicates that trailing commas are allowed
	UnevaluatedProperties any                        `json:"unevaluatedProperties,omitempty"` // validation for properties not covered by `properties`, `patternproperties`, or `additionalproperties` (schema or boolean)
	Required              []string                   `json:"required,omitempty"`              // array of required property names
	PropertyNames         *Schema                    `json:"propertyNames,omitempty"`         // schema for property names
	MaxProperties         *uint64                    `json:"maxProperties,omitempty"`         // maximum number of properties (non-negative integer)
	MinProperties         *uint64                    `json:"minProperties,omitempty"`         // minimum number of properties (non-negative integer, default 0)
	DependentRequired     *util.OrderedMap[[]string] `json:"dependentRequired,omitempty"`     // properties required based on the presence of other properties
	DependentSchemas      *util.OrderedMap[*Schema]  `json:"dependentSchemas,omitempty"`      // schemas applied based on the presence of other properties

	// conditional applicators
	If   *Schema `json:"if,omitempty"`   // if this schema validates, `then` must also validate
	Then *Schema `json:"then,omitempty"` // schema applied if `if` validates
	Else *Schema `json:"else,omitempty"` // schema applied if `if` does not validate

	// boolean logic applicators
	AllOf Schemas `json:"allOf,omitempty"` // instance must validate against all of these schemas
	AnyOf Schemas `json:"anyOf,omitempty"` // instance must validate against at least one of these schemas
	OneOf Schemas `json:"oneOf,omitempty"` // instance must validate against exactly one of these schemas
	Not   *Schema `json:"not,omitempty"`   // instance must not validate against this schema

	// extensions
	Metadata util.ValueMap `json:"metadata,omitempty"` // additional information about the schema
}

func (s *Schema) ID() string {
	if s.MetaID != "" {
		return s.MetaID
	}
	return s.ExplicitID
}
