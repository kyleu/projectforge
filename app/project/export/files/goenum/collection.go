package goenum

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func structCollection(e *enum.Enum) ([]*golang.Block, error) {
	tBlock := golang.NewBlock(e.ProperPlural(), "typealias")
	tBlock.WF("type %s []%s", e.ProperPlural(), e.Proper())
	ksBlock := blockKeys(e)
	strBlock := blockString(e)
	fnHelpBlock := blockHelp(e)
	gBlock := blockGet(e)
	gnBlock := blockGetByName(e)
	ret := []*golang.Block{tBlock, ksBlock, strBlock, fnHelpBlock, gBlock, gnBlock}

	ef := e.ExtraFields()
	for _, efk := range ef.Order {
		f := ef.GetSimple(efk)
		_, unique := e.ExtraFieldValues(efk)
		if unique {
			efkBlock, err := blockGetByPropUnique(e, efk, f)
			if err != nil {
				return nil, err
			}
			ret = append(ret, efkBlock)
		} else {
			efkBlock, err := blockGetByPropShared(e, efk, f)
			if err != nil {
				return nil, err
			}
			ret = append(ret, efkBlock)
		}
	}

	rBlock := blockRandom(e)
	ret = append(ret, rBlock)
	return ret, nil
}

func blockKeys(e *enum.Enum) *golang.Block {
	ksBlock := golang.NewBlock(e.ProperPlural()+"Keys", "method")
	ksBlock.WF("func (%s %s) Keys() []string {", e.FirstLetter(), e.ProperPlural())
	ksBlock.WF("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	ksBlock.W("\t\treturn x.Key")
	ksBlock.W("\t})")
	ksBlock.W("}")
	return ksBlock
}

func blockString(e *enum.Enum) *golang.Block {
	strBlock := golang.NewBlock(e.ProperPlural()+"Strings", "method")
	strBlock.WF("func (%s %s) Strings() []string {", e.FirstLetter(), e.ProperPlural())
	strBlock.WF("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	strBlock.W("\t\treturn x.String()")
	strBlock.W("\t})")
	strBlock.W("}")
	return strBlock
}

func blockHelp(e *enum.Enum) *golang.Block {
	fnHelpBlock := golang.NewBlock(e.Proper()+".Help", "method")
	fnHelpBlock.WF("func (%s %s) Help() string {", e.FirstLetter(), e.ProperPlural())
	fnHelpBlock.WF("\treturn \"Available %s options: [\" + strings.Join(%s.Strings(), \", \") + \"]\"", e.TitleLower(), e.FirstLetter())
	fnHelpBlock.W("}")
	return fnHelpBlock
}

func blockGet(e *enum.Enum) *golang.Block {
	gBlock := golang.NewBlock(e.ProperPlural()+"Get", "method")
	gBlock.WF("func (%s %s) Get(key string, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	gBlock.WF("\tfor _, value := range %s {", e.FirstLetter())
	gBlock.W("\t\tif strings.EqualFold(value.Key, key) {")
	gBlock.W("\t\t\treturn value")
	gBlock.W("\t\t}")
	gBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gBlock.W("\tif key == \"\" {")
		gBlock.WF("\t\treturn %s%s", e.Proper(), def.Proper())
		gBlock.W("\t}")
	}
	gBlock.WF("\tmsg := fmt.Sprintf(\"unable to find [%s] with key [%%%%s]\", key)", e.Proper())
	gBlock.W("\tif logger != nil {")
	gBlock.W("\t\tlogger.Warn(msg)")
	gBlock.W("\t}")
	gBlock.W(retUnknown(e))
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
		gnBlock.WF("\t\treturn %s%s", e.Proper(), def.Proper())
		gnBlock.W("\t}")
	}
	gnBlock.WF("\tmsg := fmt.Sprintf(\"unable to find [%s] with name [%%%%s]\", name)", e.Proper())
	gnBlock.W("\tif logger != nil {")
	gnBlock.W("\t\tlogger.Warn(msg)")
	gnBlock.W("\t}")
	gnBlock.W(retUnknown(e))
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
		gxBlock.WF("\t\treturn %s%s", e.Proper(), def.Proper())
		gxBlock.W("\t}")
	}
	possibleLogger(gxBlock, e, prop)
	gxBlock.W(retUnknown(e))
	gxBlock.W("}")
	return gxBlock, nil
}

func possibleLogger(gxBlock *golang.Block, e *enum.Enum, prop string) {
	gxBlock.W("\tif logger != nil {")
	gxBlock.WF("\t\tmsg := fmt.Sprintf(\"unable to find [%s] with %s [%%%%s]\", input)", e.Proper(), prop)
	gxBlock.W("\t\tlogger.Warn(msg)")
	gxBlock.W("\t}")
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

func blockRandom(e *enum.Enum) *golang.Block {
	rBlock := golang.NewBlock(e.ProperPlural()+"Random", "method")
	rBlock.WF("func (%s %s) Random() %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	rBlock.WF("\treturn util.RandomElement(%s)", e.FirstLetter())
	rBlock.W("}")
	return rBlock
}

func retUnknown(e *enum.Enum) string {
	if d := e.Values.Default(); d != nil {
		return fmt.Sprintf("\treturn %s%s", e.Proper(), d.Proper())
	}
	return fmt.Sprintf("\treturn %s{Key: \"_error\", Name: \"error: \" + msg}", e.Proper())
}
