package goenum

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func structCollection(e *enum.Enum) ([]*golang.Block, error) {
	tBlock := golang.NewBlock(e.ProperPlural(), "typealias")
	tBlock.W("type %s []%s", e.ProperPlural(), e.Proper())
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
	ksBlock.W("func (%s %s) Keys() []string {", e.FirstLetter(), e.ProperPlural())
	ksBlock.W("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	ksBlock.W("\t\treturn x.Key")
	ksBlock.W("\t})")
	ksBlock.W("}")
	return ksBlock
}

func blockString(e *enum.Enum) *golang.Block {
	strBlock := golang.NewBlock(e.ProperPlural()+"Strings", "method")
	strBlock.W("func (%s %s) Strings() []string {", e.FirstLetter(), e.ProperPlural())
	strBlock.W("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	strBlock.W("\t\treturn x.String()")
	strBlock.W("\t})")
	strBlock.W("}")
	return strBlock
}

func blockHelp(e *enum.Enum) *golang.Block {
	fnHelpBlock := golang.NewBlock(e.Proper()+".Help", "method")
	fnHelpBlock.W("func (%s %s) Help() string {", e.FirstLetter(), e.ProperPlural())
	fnHelpBlock.W("\treturn \"Available options: [\" + strings.Join(%s.Strings(), \", \") + \"]\"", e.FirstLetter())
	fnHelpBlock.W("}")
	return fnHelpBlock
}

func blockGet(e *enum.Enum) *golang.Block {
	gBlock := golang.NewBlock(e.ProperPlural()+"Get", "method")
	gBlock.W("func (%s %s) Get(key string, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	gBlock.W("\tfor _, value := range %s {", e.FirstLetter())
	gBlock.W("\t\tif value.Key == key {")
	gBlock.W("\t\t\treturn value")
	gBlock.W("\t\t}")
	gBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gBlock.W("\tif key == \"\" {")
		gBlock.W("\t\treturn %s%s", e.Proper(), def.Proper())
		gBlock.W("\t}")
	}
	gBlock.W("\tmsg := fmt.Sprintf(\"unable to find [%s] with key [%%%%s]\", key)", e.Proper())
	gBlock.W("\tif logger != nil {")
	gBlock.W("\t\tlogger.Warn(msg)")
	gBlock.W("\t}")
	gBlock.W("\treturn %s{Key: \"_error\", Name: \"error: \" + msg}", e.Proper())
	gBlock.W("}")
	return gBlock
}

func blockGetByName(e *enum.Enum) *golang.Block {
	gnBlock := golang.NewBlock(e.ProperPlural()+"GetByName", "method")
	gnBlock.W("func (%s %s) GetByName(name string, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	gnBlock.W("\tfor _, value := range %s {", e.FirstLetter())
	gnBlock.W("\t\tif value.Name == name {")
	gnBlock.W("\t\t\treturn value")
	gnBlock.W("\t\t}")
	gnBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gnBlock.W("\tif name == \"\" {")
		gnBlock.W("\t\treturn %s%s", e.Proper(), def.Proper())
		gnBlock.W("\t}")
	}
	gnBlock.W("\tmsg := fmt.Sprintf(\"unable to find [%s] with name [%%%%s]\", name)", e.Proper())
	gnBlock.W("\tif logger != nil {")
	gnBlock.W("\t\tlogger.Warn(msg)")
	gnBlock.W("\t}")
	gnBlock.W("\treturn %s{Key: \"_error\", Name: \"error: \" + msg}", e.Proper())
	gnBlock.W("}")
	return gnBlock
}

func blockGetByPropUnique(e *enum.Enum, efk string, t string) (*golang.Block, error) {
	prop := util.StringToCamel(efk)
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
		dflt = "false"
	case types.KeyTimestamp:
		dflt = "nil"
		goType = "*time.Time"
	default:
		return nil, errors.Errorf("unable to create enum helper for type [%s]", t)
	}
	gxBlock := golang.NewBlock(e.ProperPlural()+"GetBy"+efk, "method")
	gxBlock.W("func (%s %s) GetBy%s(input %s, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), prop, goType, e.Proper())
	gxBlock.W("\tfor _, value := range %s {", e.FirstLetter())
	gxBlock.W("\t\tif value.%s == input {", prop)
	gxBlock.W("\t\t\treturn value")
	gxBlock.W("\t\t}")
	gxBlock.W("\t}")
	if def := e.Values.Default(); def != nil {
		gxBlock.W("\tif input == %s {", dflt)
		gxBlock.W("\t\treturn %s%s", e.Proper(), def.Proper())
		gxBlock.W("\t}")
	}
	gxBlock.W("\tmsg := fmt.Sprintf(\"unable to find [%s] with %s [%%%%s]\", input)", e.Proper(), prop)
	gxBlock.W("\tif logger != nil {")
	gxBlock.W("\t\tlogger.Warn(msg)")
	gxBlock.W("\t}")
	gxBlock.W("\treturn %s{Key: \"_error\", Name: \"error: \" + msg}", e.Proper())
	gxBlock.W("}")
	return gxBlock, nil
}

func blockGetByPropShared(e *enum.Enum, efk string, t string) (*golang.Block, error) {
	prop := util.StringToCamel(efk)
	goType := t
	if t == types.KeyTimestamp {
		goType = "*time.Time"
	}
	gxBlock := golang.NewBlock(e.ProperPlural()+"GetBy"+efk, "method")
	gxBlock.W("func (%s %s) GetBy%s(input %s, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), prop, goType, e.ProperPlural())
	gxBlock.W("\tvar ret = %s", e.FirstLetter())
	gxBlock.W("\tfor _, value := range %s {", e.FirstLetter())
	gxBlock.W("\t\tif value.%s == input {", prop)
	gxBlock.W("\t\t\tret = append(ret, value)")
	gxBlock.W("\t\t}")
	gxBlock.W("\t}")
	gxBlock.W("\treturn ret")
	gxBlock.W("}")
	return gxBlock, nil
}

func blockRandom(e *enum.Enum) *golang.Block {
	rBlock := golang.NewBlock(e.ProperPlural()+"Random", "method")
	rBlock.W("func (%s %s) Random() %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	rBlock.W("\treturn %s[util.RandomInt(len(%s))]", e.FirstLetter(), e.FirstLetter())
	rBlock.W("}")
	return rBlock
}
