package svg

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func template(svgs SVGs, linebreak string) string {
	out := strings.Builder{}
	w := func(s string) {
		out.WriteString(s)
		out.WriteString(linebreak)
	}

	var maxKeyLength int
	var keys []string
	lo.ForEach(svgs, func(svg *SVG, _ int) {
		if len(svg.Key) > maxKeyLength {
			maxKeyLength = len(svg.Key)
		}
		switch svg.Key {
		case "search":
			// noop
		default:
			keys = append(keys, fmt.Sprintf(`%q`, svg.Key))
		}
	})

	w("// Code generated from files in [client/src/svg]. See " + util.AppURL + " for details. DO NOT EDIT.")
	w("package util")
	w("")
	w("var SVGLibrary = map[string]string{")
	msg := "\t%-" + fmt.Sprintf("%d", maxKeyLength+3) + "s `%s`,"
	lo.ForEach(svgs, func(fn *SVG, _ int) {
		w(fmt.Sprintf(msg, `"`+fn.Key+`":`, fn.Markup))
	})
	w("}")
	w("")
	w("//nolint:lll")
	w("var SVGIconKeys = []string{" + util.StringJoin(keys, ", ") + "}")

	w("")
	w("func RandomIcon() string {")
	w("\treturn SVGIconKeys[RandomInt(len(SVGIconKeys))]")
	w("}")

	return out.String()
}
