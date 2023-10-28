package goenum

import (
	"slices"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func structCollection(e *enum.Enum, g *golang.File) ([]*golang.Block, error) {
	tBlock := golang.NewBlock(e.ProperPlural(), "typealias")
	tBlock.W("type %s []%s", e.ProperPlural(), e.Proper())

	ksBlock := golang.NewBlock(e.ProperPlural()+"Keys", "method")
	ksBlock.W("func (%s %s) Keys() []string {", e.FirstLetter(), e.ProperPlural())
	ksBlock.W("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	ksBlock.W("\t\treturn x.Key")
	ksBlock.W("\t})")
	ksBlock.W("}")

	strBlock := golang.NewBlock(e.ProperPlural()+"Strings", "method")
	strBlock.W("func (%s %s) Strings() []string {", e.FirstLetter(), e.ProperPlural())
	strBlock.W("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	strBlock.W("\t\treturn x.String()")
	strBlock.W("\t})")
	strBlock.W("}")

	fnHelpBlock := golang.NewBlock(e.Proper()+".Help", "method")
	fnHelpBlock.W("func (%s %s) Help() string {", e.FirstLetter(), e.ProperPlural())
	fnHelpBlock.W("\treturn \"Available options: [\" + strings.Join(%s.Strings(), \", \") + \"]\"", e.FirstLetter())
	fnHelpBlock.W("}")

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

	gnBlock := golang.NewBlock(e.ProperPlural()+"Get", "method")
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

	ret := []*golang.Block{tBlock, ksBlock, strBlock, fnHelpBlock, gBlock, gnBlock}

	ef := e.ExtraFields()
	efks := maps.Keys(ef)
	slices.Sort(efks)
	for _, efk := range efks {
		prop := util.StringToCamel(efk)
		t := ef[efk]
		dflt := "\"\""
		switch t {
		case "string":
			// noop
		default:
			return nil, errors.Errorf("unable to create enum helper for type [%s]", t)
		}
		gxBlock := golang.NewBlock(e.ProperPlural()+"Get", "method")
		gxBlock.W("func (%s %s) GetBy%s(x %s, logger util.Logger) %s {", e.FirstLetter(), e.ProperPlural(), prop, t, e.Proper())
		gxBlock.W("\tfor _, value := range %s {", e.FirstLetter())
		gxBlock.W("\t\tif value.%s == x {", prop)
		gxBlock.W("\t\t\treturn value")
		gxBlock.W("\t\t}")
		gxBlock.W("\t}")
		if def := e.Values.Default(); def != nil {
			gxBlock.W("\tif x == %s {", dflt)
			gxBlock.W("\t\treturn %s%s", e.Proper(), def.Proper())
			gxBlock.W("\t}")
		}
		gxBlock.W("\tmsg := fmt.Sprintf(\"unable to find [%s] with %s [%%%%s]\", x)", e.Proper(), prop)
		gxBlock.W("\tif logger != nil {")
		gxBlock.W("\t\tlogger.Warn(msg)")
		gxBlock.W("\t}")
		gxBlock.W("\treturn %s{Key: \"_error\", Name: \"error: \" + msg}", e.Proper())
		gxBlock.W("}")
		ret = append(ret, gxBlock)
	}

	rBlock := golang.NewBlock(e.ProperPlural()+"Random", "method")
	rBlock.W("func (%s %s) Random() %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	rBlock.W("\treturn %s[util.RandomInt(len(%s))]", e.FirstLetter(), e.FirstLetter())
	rBlock.W("}")
	ret = append(ret, rBlock)

	return ret, nil
}
