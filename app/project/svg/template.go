package svg

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func template(svgs []*SVG, linebreak string) string {
	out := strings.Builder{}
	w := func(s string) {
		out.WriteString(s)
		out.WriteString(linebreak)
	}

	maxKeyLength := 0
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
	w("var SVGIconKeys = []string{" + strings.Join(keys, ", ") + "}")

	w("")
	w("func RandomIcon() string {")
	w("\treturn SVGIconKeys[RandomInt(len(SVGIconKeys))]")
	w("}")

	return out.String()
}

func cstemplate(svgs []*SVG, namespace string) string {
	out := strings.Builder{}
	w := func(s string, args ...any) {
		out.WriteString(fmt.Sprintf(s, args...))
		out.WriteString("\n")
	}

	w("// Code generated from files in [wwwroot/svg]. See " + util.AppURL + " for details. DO NOT EDIT.")
	w("namespace %s.Util;", namespace)
	w("")
	w("public static class Icons")
	w("{")
	for _, x := range svgs {
		w("    private const string %sSVG = %q;", x.Key, x.Markup)
	}
	w("")
	w("    public static Dictionary<string, string> Library = new Dictionary<string, string>")
	w("    {")
	for idx, x := range svgs {
		suffix := ""
		if idx < len(svgs)-1 {
			suffix = ","
		}
		w("        { %q, %sSVG }%s", x.Key, x.Key, suffix)
	}
	w("    };")
	w("}")

	return out.String()
}

/*



    {
        { "sign_in", signInSVG }
    };
}
*/
