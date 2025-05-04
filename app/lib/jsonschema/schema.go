package jsonschema

import "projectforge.dev/projectforge/app/util"

type Schema struct {
	// Core Vocabulary & Metadata
	Schema        string                    `json:"$schema,omitempty"`        // URI identifying the schema dialect (e.g., "https://json-schema.org/draft/2020-12/schema")
	ID            string                    `json:"$id,omitempty"`            // Base URI for the schema
	Anchor        string                    `json:"$anchor,omitempty"`        // An identifier for this subschema
	Ref           string                    `json:"$ref,omitempty"`           // Reference to another schema (URI or JSON Pointer)
	DynamicRef    string                    `json:"$dynamicRef,omitempty"`    // Reference that resolves dynamically (requires $dynamicAnchor)
	DynamicAnchor string                    `json:"$dynamicAnchor,omitempty"` // Anchor for dynamic resolution
	Vocabulary    *util.OrderedMap[bool]    `json:"$vocabulary,omitempty"`    // Declares vocabularies used (keys are URIs, values must be true)
	Comment       string                    `json:"$comment,omitempty"`       // A comment string, ignored by validators
	Defs          *util.OrderedMap[*Schema] `json:"$defs,omitempty"`          // Definitions for reusable subschemas

	// Annotations (Metadata Keywords)
	Title       string `json:"title,omitempty"`       // A short description of the schema
	Description string `json:"description,omitempty"` // A more detailed explanation
	Default     any    `json:"default,omitempty"`     // Default value for the instance
	Deprecated  bool   `json:"deprecated,omitempty"`  // Indicates the instance is deprecated (default: false)
	ReadOnly    bool   `json:"readOnly,omitempty"`    // Indicates the instance should not be modified (default: false)
	WriteOnly   bool   `json:"writeOnly,omitempty"`   // Indicates the instance may be set but should not be returned (default: false)
	Examples    []any  `json:"examples,omitempty"`    // Array of example values that validate against the schema

	// Generic Validation Keywords
	Type  any   `json:"type,omitempty"`  // Expected data type(s) (string or array of strings: "string", "number", "integer", "object", "array", "boolean", "null")
	Enum  []any `json:"enum,omitempty"`  // Array of allowed values
	Const any   `json:"const,omitempty"` // Exact required value

	// Validation Keywords for Numbers (number and integer)
	MultipleOf       *float64 `json:"multipleOf,omitempty"`       // Instance must be divisible by this number (strictly positive)
	Maximum          *float64 `json:"maximum,omitempty"`          // Maximum inclusive value
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"` // Maximum exclusive value
	Minimum          *float64 `json:"minimum,omitempty"`          // Minimum inclusive value
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"` // Minimum exclusive value

	// Validation Keywords for Strings
	MaxLength        *uint64 `json:"maxLength,omitempty"`        // Maximum length (non-negative integer)
	MinLength        *uint64 `json:"minLength,omitempty"`        // Minimum length (non-negative integer, default 0)
	Pattern          string  `json:"pattern,omitempty"`          // ECMA 262 regular expression
	Format           string  `json:"format,omitempty"`           // Predefined format (e.g., "date-time", "email", "ipv4")
	ContentEncoding  string  `json:"contentEncoding,omitempty"`  // Encoding of the string content (e.g., "base64")
	ContentMediaType string  `json:"contentMediaType,omitempty"` // Media type of the string content (e.g., "application/json")
	ContentSchema    *Schema `json:"contentSchema,omitempty"`    // Schema for the decoded content if contentEncoding/MediaType are present

	// Validation Keywords for Arrays
	Items            any     `json:"items,omitempty"`            // Schema for array items (schema or boolean false). Applied *after* prefixItems. Use `UnevaluatedItems` for more control.
	PrefixItems      Schemas `json:"prefixItems,omitempty"`      // Array of schemas for tuple validation (items at specific indices)
	UnevaluatedItems any     `json:"unevaluatedItems,omitempty"` // Validation for items not covered by `items` or `prefixItems` (schema or boolean)
	MaxItems         *uint64 `json:"maxItems,omitempty"`         // Maximum number of items (non-negative integer)
	MinItems         *uint64 `json:"minItems,omitempty"`         // Minimum number of items (non-negative integer, default 0)
	UniqueItems      *bool   `json:"uniqueItems,omitempty"`      // Whether all items must be unique (default: false)
	Contains         *Schema `json:"contains,omitempty"`         // Schema that at least one item must match
	MaxContains      *uint64 `json:"maxContains,omitempty"`      // Maximum number of items matching `contains` (non-negative integer, requires `contains`)
	MinContains      *uint64 `json:"minContains,omitempty"`      // Minimum number of items matching `contains` (non-negative integer, default 1, requires `contains`)

	// Validation Keywords for Objects
	Properties            *util.OrderedMap[*Schema]  `json:"properties,omitempty"`            // Schemas for named properties
	PatternProperties     *util.OrderedMap[*Schema]  `json:"patternProperties,omitempty"`     // Schemas for properties matching regex patterns
	AdditionalProperties  any                        `json:"additionalProperties,omitempty"`  // Controls handling of properties not explicitly listed or matched by patterns (schema or boolean). Use `UnevaluatedProperties` for more control.
	UnevaluatedProperties any                        `json:"unevaluatedProperties,omitempty"` // Validation for properties not covered by `properties`, `patternProperties`, or `additionalProperties` (schema or boolean)
	Required              []string                   `json:"required,omitempty"`              // Array of required property names
	PropertyNames         *Schema                    `json:"propertyNames,omitempty"`         // Schema for property names
	MaxProperties         *uint64                    `json:"maxProperties,omitempty"`         // Maximum number of properties (non-negative integer)
	MinProperties         *uint64                    `json:"minProperties,omitempty"`         // Minimum number of properties (non-negative integer, default 0)
	DependentRequired     *util.OrderedMap[[]string] `json:"dependentRequired,omitempty"`     // Properties required based on the presence of other properties
	DependentSchemas      *util.OrderedMap[*Schema]  `json:"dependentSchemas,omitempty"`      // Schemas applied based on the presence of other properties

	// Conditional Applicators
	If   *Schema `json:"if,omitempty"`   // If this schema validates, `then` must also validate
	Then *Schema `json:"then,omitempty"` // Schema applied if `if` validates
	Else *Schema `json:"else,omitempty"` // Schema applied if `if` does not validate

	// Boolean Logic Applicators
	AllOf Schemas `json:"allOf,omitempty"` // Instance must validate against all of these schemas
	AnyOf Schemas `json:"anyOf,omitempty"` // Instance must validate against at least one of these schemas
	OneOf Schemas `json:"oneOf,omitempty"` // Instance must validate against exactly one of these schemas
	Not   *Schema `json:"not,omitempty"`   // Instance must not validate against this schema
}
