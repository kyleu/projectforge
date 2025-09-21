package filter

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var OrderingFieldDescs = util.FieldDescs{
	{Key: "column", Title: "Column", Description: "The name of the column to sort by"},
	{Key: "asc", Title: "Ascending", Description: "Determines if this ordering is applied ascending or descending"},
}

type Ordering struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc,omitzero"`
}

func (o Ordering) String() string {
	if o.Asc {
		return o.Column
	}
	return o.Column + ":desc"
}

type Orderings []*Ordering

func (o Orderings) Get(key string) *Ordering {
	return lo.FindOrElse(o, nil, func(x *Ordering) bool {
		return x.Column == key
	})
}
