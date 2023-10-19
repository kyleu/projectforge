package goenum

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Enum(e *enum.Enum, addHeader bool, linebreak string) (*file.File, error) {
	var m model.Model
	m.Camel()
	g := golang.NewFile(e.Package, []string{"app", e.PackageWithGroup("")}, strings.ToLower(e.Camel()))
	if e.Simple() {
		g.AddBlocks(enumStructSimple(e)...)
	} else {
		g.AddBlocks(enumStructComplex(e, g)...)
		g.AddBlocks(enumStructCollection(e, g)...)
		g.AddBlocks(enumStructValues(e, g))
	}
	return g.Render(addHeader, linebreak)
}

func enumStructSimple(e *enum.Enum) []*golang.Block {
	tBlock := golang.NewBlock(e.Proper(), "typealias")
	tBlock.W("type %s string", e.Proper())

	cBlock := golang.NewBlock(e.Proper(), "constvar")
	cBlock.W("const (")
	maxCount := util.StringArrayMaxLength(e.ValuesCamel())
	pl := len(e.Proper())
	maxColLength := maxCount + pl
	lo.ForEach(e.Values, func(v *enum.Value, _ int) {
		cBlock.W("\t%s %s = %q", util.StringPad(e.Proper()+util.StringToCamel(v.Key), maxColLength), e.Proper(), v.Key)
	})
	cBlock.W(")")
	return []*golang.Block{tBlock, cBlock}
}

func enumStructComplex(e *enum.Enum, g *golang.File) []*golang.Block {
	g.AddImport(helper.ImpAppUtil, helper.ImpDBDriver, helper.ImpErrors, helper.ImpFmt, helper.ImpLo, helper.ImpStrings, helper.ImpXML)
	structBlock := golang.NewBlock(e.Proper(), "struct")
	structBlock.W("type %s struct {", e.Proper())
	structBlock.W("\tKey         string")
	structBlock.W("\tTitle       string")
	structBlock.W("\tDescription string")
	structBlock.W("}")

	fnStringBlock := golang.NewBlock(e.Proper()+".String", "method")
	fnStringBlock.W("func (%s %s) String() string {", e.FirstLetter(), e.Proper())
	fnStringBlock.W("\tif %s.Title != \"\" {", e.FirstLetter())
	fnStringBlock.W("\t\treturn %s.Title", e.FirstLetter())
	fnStringBlock.W("\t}")
	fnStringBlock.W("\treturn %s.Key", e.FirstLetter())
	fnStringBlock.W("}")

	fnJSONOutBlock := golang.NewBlock(e.Proper()+".MarshalJSON", "method")
	fnJSONOutBlock.W("func (%s *%s) MarshalJSON() ([]byte, error) {", e.FirstLetter(), e.Proper())
	fnJSONOutBlock.W("\treturn util.ToJSONBytes(%s.Key, false), nil", e.FirstLetter())
	fnJSONOutBlock.W("}")

	fnJSONInBlock := golang.NewBlock(e.Proper()+".UnmarshalJSON", "method")
	fnJSONInBlock.W("func (%s *%s) UnmarshalJSON(data []byte) error {", e.FirstLetter(), e.Proper())
	fnJSONInBlock.W("\tvar key string")
	fnJSONInBlock.W("\tif err := util.FromJSON(data, &key); err != nil {")
	fnJSONInBlock.W("\t\treturn err")
	fnJSONInBlock.W("\t}")
	fnJSONInBlock.W("\t*%s = All%s.Get(key, nil)", e.FirstLetter(), e.ProperPlural())
	fnJSONInBlock.W("\treturn nil")
	fnJSONInBlock.W("}")

	fnXMLOutBlock := golang.NewBlock(e.Proper()+".MarshalXML", "method")
	fnXMLOutBlock.W("func (%s %s) MarshalXML(e *xml.Encoder, start xml.StartElement) error {", e.FirstLetter(), e.Proper())
	fnXMLOutBlock.W("\treturn e.EncodeElement(%s.Key, start)", e.FirstLetter())
	fnXMLOutBlock.W("}")

	fnXMLInBlock := golang.NewBlock(e.Proper()+".UnmarshalXML", "method")
	fnXMLInBlock.W("func (%s *%s) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {", e.FirstLetter(), e.Proper())
	fnXMLInBlock.W("\tvar key string")
	fnXMLInBlock.W("\tif err := d.DecodeElement(&key, &start); err != nil {")
	fnXMLInBlock.W("\t\treturn err")
	fnXMLInBlock.W("\t}")
	fnXMLInBlock.W("\t*%s = All%s.Get(key, nil)", e.FirstLetter(), e.ProperPlural())
	fnXMLInBlock.W("\treturn nil")
	fnXMLInBlock.W("}")

	fnValueBlock := golang.NewBlock(e.Proper()+".Value", "method")
	fnValueBlock.W("func (%s %s) Value() (driver.Value, error) {", e.FirstLetter(), e.Proper())
	fnValueBlock.W("\treturn %s.Key, nil", e.FirstLetter())
	fnValueBlock.W("}")

	fnScanBlock := golang.NewBlock(e.Proper()+".Scan", "method")
	fnScanBlock.W("func (%s *%s) Scan(value any) error {", e.FirstLetter(), e.Proper())
	fnScanBlock.W("\tif value == nil {")
	fnScanBlock.W("\t\treturn nil")
	fnScanBlock.W("\t}")
	fnScanBlock.W("\tif sv, err := driver.String.ConvertValue(value); err == nil {")
	fnScanBlock.W("\t\tif v, ok := sv.(string); ok {")
	fnScanBlock.W("\t\t\t*%s = All%s.Get(v, nil)", e.FirstLetter(), e.ProperPlural())
	fnScanBlock.W("\t\t\treturn nil")
	fnScanBlock.W("\t\t}")
	fnScanBlock.W("\t}")
	fnScanBlock.W("\treturn errors.Errorf(\"failed to scan %s enum from value [%%%%v]\", value)", e.Proper())
	fnScanBlock.W("}")
	return []*golang.Block{structBlock, fnStringBlock, fnJSONOutBlock, fnJSONInBlock, fnXMLOutBlock, fnXMLInBlock, fnValueBlock, fnScanBlock}
}

func enumStructCollection(e *enum.Enum, g *golang.File) []*golang.Block {
	tBlock := golang.NewBlock(e.ProperPlural(), "typealias")
	tBlock.W("type %s []%s", e.ProperPlural(), e.Proper())

	ksBlock := golang.NewBlock(e.ProperPlural()+"Keys", "method")
	ksBlock.W("func (%s %s) Keys() []string {", e.FirstLetter(), e.ProperPlural())
	ksBlock.W("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	ksBlock.W("\t\treturn x.Key")
	ksBlock.W("\t})")
	ksBlock.W("}")

	tsBlock := golang.NewBlock(e.ProperPlural()+"Titles", "method")
	tsBlock.W("func (%s %s) Titles() []string {", e.FirstLetter(), e.ProperPlural())
	tsBlock.W("\treturn lo.Map(%s, func(x %s, _ int) string {", e.FirstLetter(), e.Proper())
	tsBlock.W("\t\treturn x.Title")
	tsBlock.W("\t})")
	tsBlock.W("}")

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
	gBlock.W("\tmsg := fmt.Sprintf(\"unable to find [%s] enum with key [%%%%s]\", key)", e.Proper())
	gBlock.W("\tif logger != nil {")
	gBlock.W("\t\tlogger.Warn(msg)")
	gBlock.W("\t}")
	gBlock.W("\treturn %s{Key: \"_error\", Title: \"error: \" + msg}", e.Proper())
	gBlock.W("}")

	rBlock := golang.NewBlock(e.ProperPlural()+"Random", "method")
	rBlock.W("func (%s %s) Random() %s {", e.FirstLetter(), e.ProperPlural(), e.Proper())
	rBlock.W("\treturn %s[util.RandomInt(len(%s))]", e.FirstLetter(), e.FirstLetter())
	rBlock.W("}")

	return []*golang.Block{tBlock, ksBlock, tsBlock, strBlock, fnHelpBlock, gBlock, rBlock}
}

func enumStructValues(e *enum.Enum, g *golang.File) *golang.Block {
	b := golang.NewBlock(e.Proper(), "vars")
	b.W("var (")

	maxCount := util.StringArrayMaxLength(e.ValuesCamel())
	names := make([]string, 0, len(e.Values))
	pl := len(e.Proper())
	maxColLength := maxCount + pl
	lo.ForEach(e.Values, func(v *enum.Value, _ int) {
		n := e.Proper() + util.StringToCamel(v.Key)
		names = append(names, n)
		msg := fmt.Sprintf("\t%s = %s{Key: %q", util.StringPad(n, maxColLength), e.Proper(), v.Key)
		if v.Title != "" {
			msg += fmt.Sprintf(", Title: %q", v.Title)
		}
		if v.Description != "" {
			msg += fmt.Sprintf(", Description: %q", v.Description)
		}
		b.W(msg + "}")
	})

	b.W("")
	b.W("\tAll%s = %s{%s}", e.ProperPlural(), e.ProperPlural(), strings.Join(names, ", "))
	b.W(")")
	return b
}
