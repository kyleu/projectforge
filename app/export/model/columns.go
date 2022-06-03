package model

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
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

func (c Columns) OneWithTag(t string) (*Column, error) {
	ret := c.WithTag(t)
	if len(ret) == 0 {
		return nil, errors.Errorf("must have one [%s], but found none", t)
	}
	if len(ret) > 1 {
		return nil, errors.Errorf("may only have one [%s], but found [%d]", t, len(ret))
	}
	return ret[0], nil
}

func (c Columns) WithTag(t string) Columns {
	var ret Columns
	for _, col := range c {
		if col.HasTag(t) {
			ret = append(ret, col)
		}
	}
	return ret
}

func (c Columns) WithoutTag(t string) Columns {
	var ret Columns
	for _, col := range c {
		if !col.HasTag(t) {
			ret = append(ret, col)
		}
	}
	return ret
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

func (c Columns) NonPKs() Columns {
	var ret Columns
	for _, x := range c {
		if !x.PK {
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
	return util.StringArrayQuoted(c.Names())
}

func (c Columns) CamelNames() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Camel())
	}
	return ret
}

func (c Columns) ProperNames() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.Proper())
	}
	return ret
}

func (c Columns) TitlesLower() []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.TitleLower())
	}
	return ret
}

func (c Columns) ToGoStrings(prefix string) string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, ToGoString(x.Type, prefix+x.Proper()))
	}
	return strings.Join(ret, ", ")
}

func (c Columns) ToRefs(prefix string) string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, prefix+x.Proper())
	}
	return strings.Join(ret, ", ")
}

func (c Columns) Types() types.Types {
	var ret types.Types
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

func (c Columns) Args(pkg string) string {
	args := make([]string, 0, len(c))
	for _, col := range c {
		args = append(args, fmt.Sprintf("%s %s", col.Camel(), col.ToGoType(pkg)))
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

func (c Columns) WhereClause(offset int) string {
	wc := make([]string, 0, len(c))
	for idx, col := range c {
		wc = append(wc, fmt.Sprintf("%q = $%d", col.Name, idx+offset+1))
	}
	return strings.Join(wc, " and ")
}

func (c Columns) GoTypeKeys(pkg string) []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, ToGoType(x.Type, x.Nullable, pkg))
	}
	return ret
}

func (c Columns) GoTypes(pkg string) []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, x.ToGoType(pkg))
	}
	return ret
}

func (c Columns) GoDTOTypeKeys(pkg string) []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		ret = append(ret, ToGoDTOType(x.Type, x.Nullable, pkg))
	}
	return ret
}

func (c Columns) MaxCamelLength() int {
	return util.StringArrayMaxLength(c.CamelNames())
}

func (c Columns) MaxGoTypeLength(pkg string) int {
	return util.StringArrayMaxLength(c.GoTypes(pkg))
}

func (c Columns) MaxGoKeyLength(pkg string) int {
	return util.StringArrayMaxLength(c.GoTypeKeys(pkg))
}

func (c Columns) MaxGoDTOKeyLength(pkg string) int {
	return util.StringArrayMaxLength(c.GoDTOTypeKeys(pkg))
}

func (c Columns) ForDisplay(k string) Columns {
	ret := make(Columns, 0, len(c))
	for _, x := range c {
		if x.ShouldDisplay(k) {
			ret = append(ret, x)
		}
	}
	return ret
}

func (c Columns) HasFormat(f string) bool {
	for _, col := range c {
		if col.Format == f {
			return true
		}
	}
	return false
}
