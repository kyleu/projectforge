package model

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/util"
)

type Columns []*Column

func (c Columns) Get(name string) *Column {
	return lo.FindOrElse(c, nil, func(x *Column) bool {
		return x.Name == name
	})
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
	return lo.Filter(c, func(col *Column, _ int) bool {
		return col.HasTag(t)
	})
}

func (c Columns) WithoutTag(t string) Columns {
	return lo.Reject(c, func(col *Column, _ int) bool {
		return col.HasTag(t)
	})
}

func (c Columns) PKs() Columns {
	return lo.Filter(c, func(col *Column, _ int) bool {
		return col.PK
	})
}

func (c Columns) NonPKs() Columns {
	return lo.Reject(c, func(col *Column, _ int) bool {
		return col.PK
	})
}

func (c Columns) Searches() Columns {
	return lo.Filter(c, func(col *Column, _ int) bool {
		return col.Search
	})
}

func (c Columns) Names() []string {
	return lo.Map(c, func(x *Column, _ int) string {
		return x.Name
	})
}

func (c Columns) NamesQuoted() []string {
	return util.StringArrayQuoted(c.Names())
}

func (c Columns) CamelNames() []string {
	return lo.Map(c, func(x *Column, _ int) string {
		return x.Camel()
	})
}

func (c Columns) ProperNames() []string {
	return lo.Map(c, func(x *Column, _ int) string {
		return x.Proper()
	})
}

func (c Columns) TitlesLower() []string {
	return lo.Map(c, func(x *Column, _ int) string {
		return x.TitleLower()
	})
}

func (c Columns) ToGoStrings(prefix string) string {
	ret := lo.Map(c, func(x *Column, _ int) string {
		return ToGoString(x.Type, prefix+x.Proper(), false)
	})
	return strings.Join(ret, ", ")
}

func (c Columns) ToRefs(prefix string, relCols ...*Column) string {
	ret := lo.Map(c, func(x *Column, idx int) string {
		r := prefix + x.Proper()
		if len(relCols) > idx {
			tc := relCols[idx]
			if tc.Nullable && !x.Nullable {
				r = "&" + r
			}
		}
		return r
	})
	return strings.Join(ret, ", ")
}

func (c Columns) Types() types.Types {
	return lo.Map(c, func(x *Column, _ int) types.Type {
		return x.Type
	})
}

func (c Columns) ZeroVals() []string {
	return lo.Map(c, func(x *Column, _ int) string {
		return x.ZeroVal()
	})
}

func (c Columns) Smushed() string {
	ret := lo.Map(c, func(x *Column, _ int) string {
		return x.Proper()
	})
	return strings.Join(ret, "")
}

func (c Columns) Args(pkg string, enums enum.Enums) (string, error) {
	args := make([]string, 0, len(c))
	for _, col := range c {
		gt, err := col.ToGoType(pkg, enums)
		if err != nil {
			return "", err
		}
		args = append(args, fmt.Sprintf("%s %s", col.Camel(), gt))
	}
	return strings.Join(args, ", "), nil
}

func (c Columns) Refs() string {
	refs := lo.Map(c, func(col *Column, _ int) string {
		return fmt.Sprintf("%s: %s", col.Proper(), col.Camel())
	})
	return strings.Join(refs, ", ")
}

func (c Columns) WhereClause(offset int, placeholder string) string {
	wc := make([]string, 0, len(c))
	lo.ForEach(c, func(col *Column, idx int) {
		switch placeholder {
		case "$", "":
			wc = append(wc, fmt.Sprintf("%q = $%d", col.Name, idx+offset+1))
		case "?":
			wc = append(wc, fmt.Sprintf("%q = ?", col.Name))
		case "@":
			wc = append(wc, fmt.Sprintf("%q = @p%d", col.Name, idx+offset+1))
		}
	})
	return strings.Join(wc, " and ")
}

func (c Columns) GoTypeKeys(pkg string, enums enum.Enums) ([]string, error) {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		gt, err := x.ToGoType(pkg, enums)
		if err != nil {
			return nil, err
		}

		ret = append(ret, gt)
	}
	return ret, nil
}

func (c Columns) GoTypes(pkg string, enums enum.Enums) ([]string, error) {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		gt, err := x.ToGoType(pkg, enums)
		if err != nil {
			return nil, err
		}
		ret = append(ret, gt)
	}
	return ret, nil
}

func (c Columns) GoRowTypes(pkg string, enums enum.Enums, database string) ([]string, error) {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		gdt, err := x.ToGoRowType(pkg, enums, database)
		if err != nil {
			return nil, err
		}
		ret = append(ret, gdt)
	}
	return ret, nil
}

func (c Columns) MaxCamelLength() int {
	return util.StringArrayMaxLength(c.CamelNames())
}

func (c Columns) MaxGoTypeLength(pkg string, enums enum.Enums) int {
	gt, _ := c.GoTypes(pkg, enums)
	return util.StringArrayMaxLength(gt)
}

func (c Columns) MaxGoRowTypeLength(pkg string, enums enum.Enums, database string) int {
	ks, _ := c.GoRowTypes(pkg, enums, database)
	return util.StringArrayMaxLength(ks)
}

func (c Columns) ForDisplay(k string) Columns {
	return lo.Filter(c, func(col *Column, _ int) bool {
		return col.ShouldDisplay(k)
	})
}

func (c Columns) HasFormat(f string) bool {
	return lo.ContainsBy(c, func(col *Column) bool {
		return col.Format == f
	})
}
