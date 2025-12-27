package jsonschema

type DataString struct {
	MaxLength        *uint64 `json:"maxLength,omitzero"`        // maximum length (non-negative integer)
	MinLength        *uint64 `json:"minLength,omitzero"`        // minimum length (non-negative integer, default 0)
	Pattern          string  `json:"pattern,omitzero"`          // ecma 262 regular expression
	Format           string  `json:"format,omitzero"`           // predefined format (e.g., "date-time", "email", "ipv4")
	ContentEncoding  string  `json:"contentEncoding,omitzero"`  // encoding of the string content (e.g., "base64")
	ContentMediaType string  `json:"contentMediaType,omitzero"` // media type of the string content (e.g., "application/json")
	ContentSchema    *Schema `json:"contentSchema,omitzero"`    // schema for the decoded content if contentEncoding/contentMediaType are present
}

func (d DataString) IsEmpty() bool {
	return d.MaxLength == nil && d.MinLength == nil && d.Pattern == "" && d.Format == "" &&
		d.ContentEncoding == "" && d.ContentMediaType == "" && d.ContentSchema == nil
}

func (d DataString) Clone() DataString {
	return DataString{
		MaxLength: d.MaxLength, MinLength: d.MinLength, Pattern: d.Pattern, Format: d.Format,
		ContentEncoding: d.ContentEncoding, ContentMediaType: d.ContentMediaType, ContentSchema: d.ContentSchema.Clone(),
	}
}
