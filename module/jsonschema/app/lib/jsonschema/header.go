package jsonschema

type SchemaHeader struct {
	Schema               string `json:"$schema,omitzero"`
	AdditionalProperties any    `json:"additionalProperties,omitzero"`
}
