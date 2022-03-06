package svg

import (
	"fmt"
	"strings"

	"projectforge.dev/app/util"
)

func template(src string, svgs []*SVG) string {
	out := strings.Builder{}
	w := func(s string) {
		out.WriteString(s)
		out.WriteString("\n")
	}

	maxKeyLength := 0
	var keys []string
	for _, svg := range svgs {
		if len(svg.Key) > maxKeyLength {
			maxKeyLength = len(svg.Key)
		}
		switch svg.Key {
		case "search":
			// noop
		default:
			keys = append(keys, fmt.Sprintf(`%q`, svg.Key))
		}
	}

	w("// Package util $PF_IGNORE$")
	w("// Code generated from files in [client/src/svg]. See " + util.AppURL + " for details. DO NOT EDIT.")
	w("package util")
	w("")
	w("// nolint")
	w("var SVGLibrary = map[string]string{")
	msg := "\t%-" + fmt.Sprintf("%d", maxKeyLength+3) + "s `%s`,"
	for _, fn := range svgs {
		w(fmt.Sprintf(msg, `"`+fn.Key+`":`, fn.Markup))
	}
	w("}")
	w("")
	w("// nolint")
	w("var SVGIconKeys = []string{" + strings.Join(keys, ", ") + "}")

	return out.String()
}
