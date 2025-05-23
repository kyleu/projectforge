package goenum

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func structComplex(e *enum.Enum, identityProper string, identityFn string, g *golang.File) []*golang.Block {
	g.AddImport(helper.ImpAppUtil, helper.ImpDBDriver, helper.ImpErrors, helper.ImpFmt, helper.ImpLo, helper.ImpStrings, helper.ImpXML)
	structBlock := golang.NewBlock(e.Proper(), "struct")
	structBlock.WF("type %s struct {", e.Proper())
	structBlock.W("\tKey         string")
	structBlock.W("\tName        string")
	structBlock.W("\tDescription string")
	structBlock.W("\tIcon        string")
	ef := e.ExtraFields()
	extraKeys := ef.Order
	if len(extraKeys) > 0 {
		extraKeyNames := lo.Map(extraKeys, func(x string, _ int) string {
			return util.StringToProper(x)
		})
		maxLength := util.StringArrayMaxLength(extraKeyNames)
		structBlock.WB()
		for _, x := range extraKeys {
			t := ef.GetSimple(x)
			if t == types.KeyTimestamp || t == types.KeyTimestampZoned {
				t = timePointer
			}
			structBlock.WF("\t%s %s", util.StringPad(util.StringToProper(x), maxLength), t)
		}
	}
	structBlock.W("}")

	fnStringBlock := golang.NewBlock(e.Proper()+".String", "method")
	fnStringBlock.WF("func (%s %s) String() string {", e.FirstLetter(), e.Proper())
	fnStringBlock.WF("\treturn %s.%s", e.FirstLetter(), identityProper)
	fnStringBlock.W("}")

	fnNameSafeBlock := golang.NewBlock(e.Proper()+".NameSafe", "method")
	fnNameSafeBlock.WF("func (%s %s) NameSafe() string {", e.FirstLetter(), e.Proper())
	fnNameSafeBlock.WF("\tif %s.Name != \"\" {", e.FirstLetter())
	fnNameSafeBlock.WF("\t\treturn %s.Name", e.FirstLetter())
	fnNameSafeBlock.W("\t}")
	fnNameSafeBlock.WF("\treturn %s.String()", e.FirstLetter())
	fnNameSafeBlock.W("}")

	fnMatchBlock := golang.NewBlock(e.ProperPlural()+"Matches", "method")
	fnMatchBlock.WF("func (%s %s) Matches(xx %s) bool {", e.FirstLetter(), e.Proper(), e.Proper())
	fnMatchBlock.WF("\treturn %s.Key == xx.Key", e.FirstLetter())
	fnMatchBlock.W("}")

	fnJSONOutBlock := golang.NewBlock(e.Proper()+".MarshalJSON", "method")
	fnJSONOutBlock.WF("func (%s %s) MarshalJSON() ([]byte, error) {", e.FirstLetter(), e.Proper())
	fnJSONOutBlock.WF("\treturn util.ToJSONBytes(%s.%s, false), nil", e.FirstLetter(), identityProper)
	fnJSONOutBlock.W("}")

	fnJSONInBlock := golang.NewBlock(e.Proper()+".UnmarshalJSON", "method")
	fnJSONInBlock.WF("func (%s *%s) UnmarshalJSON(data []byte) error {", e.FirstLetter(), e.Proper())
	fnJSONInBlock.W("\tkey, err := util.FromJSONString(data)")
	fnJSONInBlock.W("\tif err != nil {")
	fnJSONInBlock.W("\t\treturn err")
	fnJSONInBlock.W("\t}")
	fnJSONInBlock.WF("\t*%s = All%s.Get%s(key, nil)", e.FirstLetter(), e.ProperPlural(), identityFn)
	fnJSONInBlock.W("\treturn nil")
	fnJSONInBlock.W("}")

	fnXMLOutBlock := golang.NewBlock(e.Proper()+".MarshalXML", "method")
	fnXMLOutBlock.WF("func (%s %s) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {", e.FirstLetter(), e.Proper())
	fnXMLOutBlock.WF("\treturn enc.EncodeElement(%s.%s, start)", e.FirstLetter(), identityProper)
	fnXMLOutBlock.W("}")

	fnXMLInBlock := golang.NewBlock(e.Proper()+".UnmarshalXML", "method")
	fnXMLInBlock.WF("func (%s *%s) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {", e.FirstLetter(), e.Proper())
	fnXMLInBlock.W("\tvar key string")
	fnXMLInBlock.W("\tif err := dec.DecodeElement(&key, &start); err != nil {")
	fnXMLInBlock.W("\t\treturn err")
	fnXMLInBlock.W("\t}")
	fnXMLInBlock.WF("\t*%s = All%s.Get%s(key, nil)", e.FirstLetter(), e.ProperPlural(), identityFn)
	fnXMLInBlock.W("\treturn nil")
	fnXMLInBlock.W("}")

	fnValueBlock := golang.NewBlock(e.Proper()+".Value", "method")
	fnValueBlock.WF("func (%s %s) Value() (driver.Value, error) {", e.FirstLetter(), e.Proper())
	fnValueBlock.WF("\treturn %s.%s, nil", e.FirstLetter(), identityProper)
	fnValueBlock.W("}")

	fnScanBlock := golang.NewBlock(e.Proper()+".Scan", "method")
	fnScanBlock.WF("func (%s *%s) Scan(value any) error {", e.FirstLetter(), e.Proper())
	fnScanBlock.W("\tif value == nil {")
	fnScanBlock.W("\t\treturn nil")
	fnScanBlock.W("\t}")
	fnScanBlock.W("\tif converted, err := driver.String.ConvertValue(value); err == nil {")
	fnScanBlock.W("\t\tif str, err := util.Cast[string](converted); err == nil {")
	fnScanBlock.WF("\t\t\t*%s = All%s.Get%s(str, nil)", e.FirstLetter(), e.ProperPlural(), identityFn)
	fnScanBlock.W("\t\t\treturn nil")
	fnScanBlock.W("\t\t}")
	fnScanBlock.W("\t}")
	fnScanBlock.WF("\treturn errors.Errorf(\"failed to scan %s enum from value [%%%%v]\", value)", e.Proper())
	fnScanBlock.W("}")
	return []*golang.Block{
		structBlock, fnStringBlock, fnNameSafeBlock, fnMatchBlock, fnJSONOutBlock, fnJSONInBlock, fnXMLOutBlock, fnXMLInBlock, fnValueBlock, fnScanBlock,
	}
}
