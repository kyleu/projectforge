package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/util"
)

func (c Columns) ToGoStrings(prefix string, alwaysString bool, maxLength int) string {
	ret := lo.Map(c, func(x *Column, _ int) string {
		return ToGoString(x.Type, x.Nullable, prefix+x.Proper(), alwaysString)
	})
	if maxLength == 0 {
		return strings.Join(ret, ", ")
	}
	return strings.Join(ret, ", ")
}

func (c Columns) ZeroVals() []string {
	return lo.Map(c, func(x *Column, _ int) string {
		return x.ZeroVal()
	})
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

func (c Columns) GoTypeKeys(pkg string, enums enum.Enums) ([]string, error) {
	ret := util.NewStringSlice(make([]string, 0, len(c)))
	for _, x := range c {
		gt, err := x.ToGoType(pkg, enums)
		if err != nil {
			return nil, err
		}
		ret.Push(gt)
	}
	return ret.Slice, nil
}

func (c Columns) GoTypes(pkg string, enums enum.Enums) ([]string, error) {
	ret := util.NewStringSlice(make([]string, 0, len(c)))
	for _, x := range c {
		gt, err := x.ToGoType(pkg, enums)
		if err != nil {
			return nil, err
		}
		ret.Push(gt)
	}
	return ret.Slice, nil
}

func (c Columns) GoRowTypes(pkg string, enums enum.Enums, database string) ([]string, error) {
	ret := util.NewStringSlice(make([]string, 0, len(c)))
	for _, x := range c {
		gdt, err := x.ToGoRowType(pkg, enums, database)
		if err != nil {
			return nil, err
		}
		ret.Push(gdt)
	}
	return ret.Slice, nil
}

func (c Columns) MaxGoTypeLength(pkg string, enums enum.Enums) int {
	gt, _ := c.GoTypes(pkg, enums)
	return util.StringArrayMaxLength(gt)
}

func (c Columns) MaxGoRowTypeLength(pkg string, enums enum.Enums, database string) int {
	ks, _ := c.GoRowTypes(pkg, enums, database)
	return util.StringArrayMaxLength(ks)
}
