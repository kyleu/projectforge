package filter

import "{{{ .Package }}}/app/util"

var OrderingFieldDescs = util.FieldDescs{
	{Key: "column", Title: "Column", Description: "The name of the column to sort by"},
	{Key: "asc", Title: "Ascending", Description: "Determines if this ordering is applied ascending or descending"},
}

type Ordering struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc,omitempty"`
}

func (o Ordering) String() string {
	if o.Asc {
		return o.Column
	}
	return o.Column + ":desc"
}

type Orderings []*Ordering
