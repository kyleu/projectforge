package datschema

type Schema struct {
	Version      int          `json:"version"`
	CreatedAt    int          `json:"createdAt"`
	Tables       Tables       `json:"tables,omitempty"`
	Enumerations Enumerations `json:"enumerations,omitempty"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns Columns  `json:"columns"`
	Tags    []string `json:"tags,omitempty"`
}

type Tables []*Table

type Column struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Array       bool   `json:"array,omitempty"`
	Type        string `json:"type"`
	Unique      bool   `json:"unique,omitempty"`
	Localized   bool   `json:"localized,omitempty"`
	References  any    `json:"references,omitempty"`
	Until       any    `json:"until,omitempty"`
	File        any    `json:"file,omitempty"`
	Files       any    `json:"files,omitempty"`
}

func (c *Column) SafeName() string {
	if c.Name == "" {
		return "<unknown>"
	}
	return c.Name
}

type Columns []*Column

type Enumeration struct {
	Name        string   `json:"name"`
	Indexing    int      `json:"indexing"`
	Enumerators []string `json:"enumerators,omitempty"`
}

type Enumerations []*Enumeration
