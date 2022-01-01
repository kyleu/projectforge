package field

type Metadata struct {
	Description string   `json:"description,omitempty"`
	Comments    []string `json:"comments,omitempty"`
	Source      string   `json:"source,omitempty"`
	Line        int      `json:"line,omitempty"`
	Column      int      `json:"column,omitempty"`
}
