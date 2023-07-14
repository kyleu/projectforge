package svg

import (
	"encoding/xml"
	"projectforge.dev/projectforge/app/util"
	"strings"
)

type xmlNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []xmlNode  `xml:",any"`
}

func (n *xmlNode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node xmlNode

	return d.DecodeElement((*node)(n), &start)
}

func getAttr(attrs []xml.Attr, k string) string {
	for _, a := range attrs {
		if a.Name.Local == k {
			return a.Value
		}
	}
	return ""
}

func cleanMarkup(orig string, color string) (string, string) {
	for strings.Contains(orig, "<!--") {
		startIdx := strings.Index(orig, "<!--")
		endIdx := strings.Index(orig, "-->")
		if endIdx == -1 {
			break
		}
		orig = strings.TrimPrefix(orig[:startIdx]+orig[endIdx+3:], util.StringDetectLinebreak(orig))
	}
	origColored := orig
	if color != "" {
		origColored = strings.ReplaceAll(orig, "<path ", "<path fill=\""+color+"\" ")
	}
	return orig, origColored
}
