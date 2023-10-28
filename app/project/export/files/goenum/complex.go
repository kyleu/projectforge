package goenum

import (
	"golang.org/x/exp/maps"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func structComplex(e *enum.Enum, g *golang.File) []*golang.Block {
	g.AddImport(helper.ImpAppUtil, helper.ImpDBDriver, helper.ImpErrors, helper.ImpFmt, helper.ImpLo, helper.ImpStrings, helper.ImpXML)
	structBlock := golang.NewBlock(e.Proper(), "struct")
	structBlock.W("type %s struct {", e.Proper())
	structBlock.W("\tKey         string")
	structBlock.W("\tName        string")
	structBlock.W("\tDescription string")
	structBlock.W("\tIcon        string")
	extraKeys := maps.Keys(e.ExtraFields())
	if len(extraKeys) > 0 {
		maxLength := util.StringArrayMaxLength(extraKeys)
		structBlock.WB()
		for _, x := range extraKeys {
			structBlock.W("\t%s string", util.StringPad(util.StringToCamel(x), maxLength))
		}
	}
	structBlock.W("}")

	fnStringBlock := golang.NewBlock(e.Proper()+".String", "method")
	fnStringBlock.W("func (%s %s) String() string {", e.FirstLetter(), e.Proper())
	fnStringBlock.W("\tif %s.Name != \"\" {", e.FirstLetter())
	fnStringBlock.W("\t\treturn %s.Name", e.FirstLetter())
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
	fnXMLOutBlock.W("func (%s %s) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {", e.FirstLetter(), e.Proper())
	fnXMLOutBlock.W("\treturn enc.EncodeElement(%s.Key, start)", e.FirstLetter())
	fnXMLOutBlock.W("}")

	fnXMLInBlock := golang.NewBlock(e.Proper()+".UnmarshalXML", "method")
	fnXMLInBlock.W("func (%s *%s) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {", e.FirstLetter(), e.Proper())
	fnXMLInBlock.W("\tvar key string")
	fnXMLInBlock.W("\tif err := dec.DecodeElement(&key, &start); err != nil {")
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
