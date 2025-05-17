package jsonschema

type SchemaHeader struct {
	Schema               string `json:"$schema,omitempty"`
	AdditionalProperties any    `json:"additionalProperties,omitempty"`
}
