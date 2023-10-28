package goenum

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func Enum(e *enum.Enum, addHeader bool, linebreak string) (*file.File, error) {
	var m model.Model
	m.Camel()
	g := golang.NewFile(e.Package, []string{"app", e.PackageWithGroup("")}, strings.ToLower(e.Camel()))
	if e.Simple() {
		g.AddBlocks(structSimple(e)...)
	} else {
		g.AddBlocks(structComplex(e, g)...)
		g.AddBlocks(enumStructParse(e))
		coll, err := structCollection(e, g)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(coll...)
		g.AddBlocks(enumValues(e, g))
	}
	return g.Render(addHeader, linebreak)
}

func enumStructParse(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock(e.Proper(), "parse")
	ret.W("func %sParse(logger util.Logger, strings ...string) %s {", e.Proper(), e.ProperPlural())
	ret.W("\treturn lo.Map(strings, func(x string, _ int) %s {", e.Proper())
	ret.W("\t\treturn All%s.Get(x, logger)", e.ProperPlural())
	ret.W("\t})")
	ret.W("}")
	return ret
}
