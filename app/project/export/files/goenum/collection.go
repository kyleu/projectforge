package goenum

import (
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func structCollection(e *enum.Enum, identityProper string) ([]*golang.Block, error) {
	tBlock := golang.NewBlock(e.ProperPlural(), "typealias")
	tBlock.WF("type %s []%s", e.ProperPlural(), e.Proper())
	ksBlock := blockKeys(e)
	strBlock := blockString(e)
	namesSafeBlock := blockNamesSafe(e)
	fnHelpBlock := blockHelp(e)
	gBlock := blockGet(e, identityProper)
	gnBlock := blockGetByName(e)
	ret := []*golang.Block{tBlock, ksBlock, strBlock, namesSafeBlock, fnHelpBlock, gBlock, gnBlock}

	ef := e.ExtraFields()
	for _, efk := range ef.Keys() {
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

func blockNamesSafe(e *enum.Enum) *golang.Block {
	strBlock := golang.NewBlock(e.ProperPlural()+"NamesSafe", "method")
	strBlock.WF("func (%s %s) NamesSafe() []string {", e.FirstLetter(), e.ProperPlural())
	strBlock.WF("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	strBlock.W("\t\treturn x.NameSafe()")
	strBlock.W("\t})")
	strBlock.W("}")
	return strBlock
}

func blockHelp(e *enum.Enum) *golang.Block {
	fnHelpBlock := golang.NewBlock(e.Proper()+".Help", "method")
	fnHelpBlock.WF("func (%s %s) Help() string {", e.FirstLetter(), e.ProperPlural())
	fnHelpBlock.WF("\treturn \"Available %s options: [\" + util.StringJoin(%s.Strings(), \", \") + \"]\"", e.TitleLower(), e.FirstLetter())
	fnHelpBlock.W("}")
	return fnHelpBlock
}

func blockRandom(e *enum.Enum) *golang.Block {
	rBlock := golang.NewBlock(e.ProperPlural()+"Random", "method")
	rBlock.WF("func (%s %s) Random() %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	rBlock.WF("\treturn util.RandomElement(%s)", e.FirstLetter())
	rBlock.W("}")
	return rBlock
}
