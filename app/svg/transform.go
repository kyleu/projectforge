package svg

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func Transform(tgt string, b []byte, url string) (*SVG, error) {
	x := &xmlNode{}
	if err := xml.Unmarshal(b, x); err != nil {
		return nil, errors.Wrapf(err, "unable to parse XML from [%s]", string(b))
	}
	if x.XMLName.Local != "svg" {
		return nil, errors.New("root element must be [svg]")
	}
	vb := getAttr(x.Attrs, "viewBox")
	if vb == "" {
		return nil, errors.New("no [viewBox] available in <svg> attributes")
	}

	var markup string
	add := func(s string) {
		markup += s + "\n"
	}

	add("<!-- imported from " + url + " -->")
	add(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`)
	add(fmt.Sprintf(`  <symbol id="svg-%s" viewBox="%s">`, tgt, vb))

	content, err := transformNodes(x.Nodes)
	if err != nil {
		return nil, err
	}
	add(content)

	add("  </symbol>")
	add(fmt.Sprintf(`  <use xlink:href="#svg-%s" />`, tgt))
	add("</svg>")

	return &SVG{Key: tgt, Markup: markup}, nil
}

func transformNodes(nodes []xmlNode) (string, error) {
	ret := ""
	for _, node := range nodes {
		findAttr := func(k string) *xml.Attr {
			for _, a := range node.Attrs {
				if a.Name.Local == k {
					return &a
				}
			}
			return nil
		}

		node.Attrs = cleanAttrs(node.Attrs)

		class := findAttr("class")
		var classes []string
		stroke := findAttr("stroke")
		if stroke != nil {
			classes = append(classes, "svg-stroke")
		}
		fill := findAttr("fill")
		if stroke == nil || fill != nil {
			classes = append(classes, "svg-fill")
		}

		if len(classes) > 0 {
			cls := strings.Join(classes, " ")
			if class == nil {
				base := []xml.Attr{{Name: xml.Name{Local: "class"}, Value: cls}}
				attrs := append([]xml.Attr{}, base...)
				attrs = append(attrs, node.Attrs...)
				node.Attrs = attrs
			} else {
				class.Value = class.Value + " " + cls
			}
		}

		b, err := xml.Marshal(node)
		if err != nil {
			return "", err
		}
		ret += strings.ReplaceAll(string(b), ` xmlns="http://www.w3.org/2000/svg"`, "")
	}
	ret = strings.ReplaceAll(ret, "></path>", " />")
	return "    " + ret, nil
}

func cleanAttrs(attrs []xml.Attr) []xml.Attr {
	ret := make([]xml.Attr, 0, len(attrs))
	for _, a := range attrs {
		n := a.Name.Local
		if n != "xmlns" {
			hit := false
			for _, x := range ret {
				if x.Name.Local == a.Name.Local {
					hit = true
				}
			}
			if !hit {
				ret = append(ret, a)
			}
		}
	}
	return ret
}
