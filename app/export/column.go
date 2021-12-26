package export

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

type Column struct {
	Name string `json:"name"`
	Type *Type  `json:"type"`
	PK   bool   `json:"pk,omitempty"`
}

func (c Column) camel() string {
	return util.StringToLowerCamel(c.Name)
}

func (c Column) proper() string {
	return util.StringToCamel(c.Name)
}

type Columns []*Column

func (c Columns) PKs() Columns {
	var ret Columns
	for _, x := range c {
		if x.PK {
			ret = append(ret, x)
		}
	}
	return ret
}

func (c Columns) Names() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Name)
	}
	return ret
}

func (c Columns) camelNames() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.camel())
	}
	return ret
}

func (c Columns) Types() Types {
	var ret Types
	for _, x := range c {
		ret = append(ret, x.Type)
	}
	return ret
}

func (c Columns) Smushed() string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.proper())
	}
	return strings.Join(ret, "")
}

func (c Columns) Args() string {
	args := make([]string, 0, len(c))
	for _, col := range c {
		args = append(args, fmt.Sprintf("%s %s", col.camel(), col.Type.Go))
	}
	return strings.Join(args, ", ")
}

func (c Columns) Refs() string {
	refs := make([]string, 0, len(c))
	for _, col := range c {
		refs = append(refs, fmt.Sprintf("%s: %s", col.proper(), col.camel()))
	}
	return strings.Join(refs, ", ")
}

func (c Columns) WhereClause() interface{} {
	wc := make([]string, 0, len(c))
	for idx, col := range c {
		wc = append(wc, fmt.Sprintf("%s = $%d", col.Name, idx+1))
	}
	return strings.Join(wc, " and ")
}
