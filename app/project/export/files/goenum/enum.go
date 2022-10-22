package goenum

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Enum(e *enum.Enum, args *model.Args, addHeader bool) (*file.File, error) {
	var m model.Model
	m.Camel()
	g := golang.NewFile(e.Package, []string{"app", e.PackageWithGroup("")}, strings.ToLower(e.Camel()))
	g.AddBlocks(enumStruct(e))
	return g.Render(addHeader)
}

func enumStruct(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock(e.Proper(), "struct")
	ret.W("type %s string", e.Proper())
	ret.W("")
	ret.W("const (")
	maxColLength := util.StringArrayMaxLength(e.Values) + len(e.Proper())
	for _, v := range e.Values {
		ret.W("\t%s %s = %q", util.StringPad(e.Proper()+util.StringToCamel(v), maxColLength), e.Proper(), v)
	}
	ret.W(")")
	return ret
}
