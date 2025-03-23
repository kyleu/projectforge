package goenum

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Enum(e *enum.Enum, linebreak string) (*file.File, error) {
	var m model.Model
	m.Camel()
	g := golang.NewFile(e.Package, []string{"app", e.PackageWithGroup("")}, strings.ToLower(e.Camel()))
	g.AddBlocks(enumValues(e))
	if e.Simple() {
		g.AddBlocks(structSimple(e)...)
	} else {
		identityProper := "Key"
		var identityFn string
		if i := e.Config.GetStringOpt("identity"); i != "" {
			identityProper = util.StringToProper(i)
			identityFn = "By" + identityProper
		}
		g.AddBlocks(structComplex(e, identityProper, identityFn, g)...)
		g.AddBlocks(enumStructParse(e, identityFn))
		coll, err := structCollection(e, identityProper)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(coll...)
	}
	return g.Render(linebreak)
}

func enumStructParse(e *enum.Enum, identityFn string) *golang.Block {
	ret := golang.NewBlock(e.Proper(), "parse")
	ret.WF("func %sParse(logger util.Logger, keys ...string) %s {", e.Proper(), e.ProperPlural())
	ret.W("\tif len(keys) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")
	ret.WF("\treturn lo.Map(keys, func(x string, _ int) %s {", e.Proper())
	ret.WF("\t\treturn All%s.Get%s(x, logger)", e.ProperPlural(), identityFn)
	ret.W("\t})")
	ret.W("}")
	return ret
}
