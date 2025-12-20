package jsonschema

type dataAnnotations struct {
	Title       string `json:"title,omitzero"`       // a short description of the schema
	Description string `json:"description,omitzero"` // a more detailed explanation
	Default     any    `json:"default,omitzero"`     // default value for the instance
	Deprecated  any    `json:"deprecated,omitzero"`  // indicates the instance is deprecated (default: false)
	ReadOnly    bool   `json:"readOnly,omitzero"`    // indicates the instance should not be modified (default: false)
	WriteOnly   bool   `json:"writeOnly,omitzero"`   // indicates the instance may be set but should not be returned (default: false)
}

func (d dataAnnotations) IsEmpty() bool {
	return d.Title == "" && d.Description == "" && d.Default == nil && d.Deprecated == nil && d.ReadOnly == false && d.WriteOnly == false
}

func (d dataAnnotations) Clone() dataAnnotations {
	return dataAnnotations{
		Title: d.Title, Description: d.Description, Default: d.Default, Deprecated: d.Deprecated, ReadOnly: d.ReadOnly, WriteOnly: d.WriteOnly,
	}
}
