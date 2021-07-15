package svg

import (
	"encoding/xml"
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
