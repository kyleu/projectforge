package svg

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func Transform(tgt string, b []byte, url string) (*SVG, error) {
	c := string(b)
	x := &xmlNode{}
	if err := xml.Unmarshal(b, x); err != nil {
		return nil, errors.Wrapf(err, "unable to parse XML from [%s]", c)
	}
	if x.XMLName.Local != "svg" {
		return nil, errors.New("root element must be [svg]")
	}
	vb := getAttr(x.Attrs, "viewBox")
	linebreak := util.StringDetectLinebreak(c)
	var markup string
	add := func(s string) {
		markup += s + linebreak
	}

	if idx := strings.LastIndex(tgt, "/"); idx > -1 {
		tgt = tgt[idx+1:]
	}
	if idx := strings.LastIndex(tgt, "."); idx > -1 {
		tgt = tgt[:idx]
	}

	if url == "" {
		add("<!-- imported from inline content -->")
	} else {
		add("<!-- imported from " + url + " -->")
	}
	add(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 32 32">`)
	if vb == "" {
		add(fmt.Sprintf(`  <symbol id="svg-%s">`, tgt))
	} else {
		add(fmt.Sprintf(`  <symbol id="svg-%s" viewBox=%q>`, tgt, vb))
	}

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
	var ret string
	for _, node := range nodes {
		findAttr := func(k string) *xml.Attr {
			for _, a := range node.Attrs {
				if a.Name.Local == k {
					return &a
				}
			}
			return nil
		}

		if node.XMLName.Local == "title" {
			continue
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
			cls := util.StringJoin(classes, " ")
			if class == nil {
				base := []xml.Attr{{Name: xml.Name{Local: "class"}, Value: cls}}
				node.Attrs = append(util.ArrayCopy(base), node.Attrs...)
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
			hit := lo.ContainsBy(ret, func(x xml.Attr) bool {
				return x.Name.Local == a.Name.Local
			})
			if !hit {
				ret = append(ret, a)
			}
		}
	}
	return ret
}
