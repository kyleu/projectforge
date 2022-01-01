package model

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

type Columns []*Column

func (c Columns) Get(name string) *Column {
	for _, x := range c {
		if x.Name == name {
			return x
		}
	}
	return nil
}

func (c Columns) PKs() Columns {
	var ret Columns
	for _, x := range c {
		if x.PK {
			ret = append(ret, x)
		}
	}
	return ret
}

func (c Columns) Searches() Columns {
	var ret Columns
	for _, x := range c {
		if x.Search {
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

func (c Columns) NamesQuoted() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, fmt.Sprintf("%q", x.Name))
	}
	return ret
}

func (c Columns) CamelNames() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Camel())
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

func (c Columns) ZeroVals() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.ZeroVal())
	}
	return ret
}

func (c Columns) Smushed() string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Proper())
	}
	return strings.Join(ret, "")
}

func (c Columns) Args() string {
	args := make([]string, 0, len(c))
	for _, col := range c {
		args = append(args, fmt.Sprintf("%s %s", col.Camel(), col.ToGoType()))
	}
	return strings.Join(args, ", ")
}

func (c Columns) Refs() string {
	refs := make([]string, 0, len(c))
	for _, col := range c {
		refs = append(refs, fmt.Sprintf("%s: %s", col.Proper(), col.Camel()))
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

func (c Columns) GoTypeKeys() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Type.ToGoType(x.Nullable))
	}
	return ret
}

func (c Columns) GoDTOTypeKeys() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Type.ToGoDTOType(x.Nullable))
	}
	return ret
}

func (c Columns) MaxGoKeyLength() int {
	return util.StringArrayMaxLength(c.GoTypeKeys())
}

func (c Columns) MaxGoDTOKeyLength() int {
	return util.StringArrayMaxLength(c.GoDTOTypeKeys())
}
