package goenum

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Enum(e *enum.Enum, addHeader bool, linebreak string) (*file.File, error) {
	var m model.Model
	m.Camel()
	g := golang.NewFile(e.Package, []string{"app", e.PackageWithGroup("")}, strings.ToLower(e.Camel()))
	g.AddBlocks(enumStruct(e))
	// g.AddBlocks(enumAll(e), enumAllStrings(e))
	return g.Render(addHeader, linebreak)
}

func enumStruct(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock(e.Proper(), "struct")
	ret.W("type %s string", e.Proper())
	ret.WB()
	ret.W("const (")
	max := util.StringArrayMaxLength(e.ValuesCamel())
	pl := len(e.Proper())
	maxColLength := max + pl
	lo.ForEach(e.Values, func(v string, _ int) {
		ret.W("\t%s %s = %q", util.StringPad(e.Proper()+util.StringToCamel(v), maxColLength), e.Proper(), v)
	})
	ret.W(")")
	return ret
}
