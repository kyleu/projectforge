package goenum

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func blockGet(e *enum.Enum, identityProper string) *golang.Block {
	gBlock := golang.NewBlock(e.ProperPlural()+"Get", "method")
	gBlock.WF("func (%s %s) Get(key string, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	gBlock.WF("\tfor _, value := range %s {", e.FirstLetter())
	gBlock.WF("\t\tif strings.EqualFold(value.%s, key) {", identityProper)
	gBlock.W("\t\t\treturn value")
	gBlock.W("\t\t}")
	gBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gBlock.W("\tif key == \"\" {")
		gBlock.WF("\t\treturn %s%s", e.Proper(), def.Proper(e.Acronyms...))
		gBlock.W("\t}")
	}
	missingHandler(gBlock, e, "key", "key")
	gBlock.W("}")
	return gBlock
}

func blockGetByName(e *enum.Enum) *golang.Block {
	gnBlock := golang.NewBlock(e.ProperPlural()+"GetByName", "method")
	gnBlock.WF("func (%s %s) GetByName(name string, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	gnBlock.WF("\tfor _, value := range %s {", e.FirstLetter())
	gnBlock.W("\t\tif strings.EqualFold(value.Name, name) {")
	gnBlock.W("\t\t\treturn value")
	gnBlock.W("\t\t}")
	gnBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gnBlock.W("\tif name == \"\" {")
		gnBlock.WF("\t\treturn %s%s", e.Proper(), def.Proper(e.Acronyms...))
		gnBlock.W("\t}")
	}
	missingHandler(gnBlock, e, "name", "name")
	gnBlock.W("}")
	return gnBlock
}

func blockGetByPropUnique(e *enum.Enum, efk string, t string) (*golang.Block, error) {
	prop := util.StringToProper(efk)
	dflt := "\"\""
	goType := t
	switch t {
	case types.KeyString:
		// noop
	case types.KeyInt:
		dflt = "0"
	case types.KeyFloat:
		dflt = "0"
	case types.KeyBool:
		dflt = util.BoolFalse
	case types.KeyTimestamp, types.KeyTimestampZoned:
		dflt = "nil"
		goType = timePointer
	default:
		return nil, errors.Errorf("unable to create enum helper for type [%s]", t)
	}
	gxBlock := golang.NewBlock(e.ProperPlural()+helper.TextGetBy+efk, "method")
	gxBlock.WF("func (%s %s) GetBy%s(input %s, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), prop, goType, e.Proper())
	gxBlock.WF("\tfor _, value := range %s {", e.FirstLetter())
	gxBlock.WF("\t\tif value.%s == input {", prop)
	gxBlock.W("\t\t\treturn value")
	gxBlock.W("\t\t}")
	gxBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gxBlock.WF("\tif input == %s {", dflt)
		gxBlock.WF("\t\treturn %s%s", e.Proper(), def.Proper(e.Acronyms...))
		gxBlock.W("\t}")
	}
	missingHandler(gxBlock, e, prop, "input")
	gxBlock.W("}")
	return gxBlock, nil
}

func blockGetByPropShared(e *enum.Enum, efk string, t string) (*golang.Block, error) {
	prop := util.StringToProper(efk)
	goType := t
	if t == types.KeyTimestamp || t == types.KeyTimestampZoned {
		goType = timePointer
	}
	gxBlock := golang.NewBlock(e.ProperPlural()+helper.TextGetBy+efk, "method")
	gxBlock.WF("func (%s %s) GetBy%s(input %s) %s {", e.FirstLetter(), e.ProperPlural(), prop, goType, e.ProperPlural())
	gxBlock.WF("\tret := %s", e.FirstLetter())
	gxBlock.WF("\tfor _, value := range %s {", e.FirstLetter())
	gxBlock.WF("\t\tif value.%s == input {", prop)
	gxBlock.W("\t\t\tret = append(ret, value)")
	gxBlock.W("\t\t}")
	gxBlock.W("\t}")
	gxBlock.W("\treturn ret")
	gxBlock.W("}")
	return gxBlock, nil
}

func missingHandler(b *golang.Block, e *enum.Enum, prop string, v string) {
	b.WF(`	msg := fmt.Sprintf("unable to find [%s] with %s [%%%%s]", %s)`, e.Proper(), prop, v)
	b.W("\tif logger != nil {")
	b.W("\t\tlogger.Warn(msg)")
	b.W("\t}")
	if d := e.Values.Default(); d != nil {
		b.WF("\treturn %s%s", e.Proper(), d.Proper(e.Acronyms...))
	} else {
		b.WF(`	return %s{Key: "_error", Name: "error: " + msg}`, e.Proper())
	}
}
