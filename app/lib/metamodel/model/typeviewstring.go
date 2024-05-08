// Package model - Content managed by Project Forge, see [projectforge.md] for details.
package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

const (
	tmplStart   = "{%%"
	tmplStartEQ = tmplStart + "= "
	tmplStartS  = tmplStart + "s "
	tmplStartV  = tmplStart + "v "
	tmplEnd     = " %%}"
	tmplEndP    = ")" + tmplEnd
)

func goViewStringForString(url bool, src string, t *types.Wrapped, nullable bool, prop string, format string, verbose bool) string {
	key := "s"
	if url {
		key = "u"
	}
	if src == util.KeySimple {
		return tmplStart + key + " " + prop + tmplEnd
	}
	componentEnd := ", \"\", ps)"
	switch format {
	case FmtCode.Key:
		return "<pre class=\"prewsw\">" + tmplStartS + ToGoString(t, nullable, prop, false) + tmplEnd + "</pre>"
	case FmtLinebreaks.Key:
		return "<div class=\"prewsl\">" + tmplStartS + ToGoString(t, nullable, prop, false) + tmplEnd + "</div>"
	case FmtCodeHidden.Key:
		_, p := util.StringSplitLast(prop, '.', true)
		ret := util.NewStringSlice(make([]string, 0, 30))
		ret.Push(
			"<ul class=\"accordion\">",
			"<li>",
			"<input id=\"accordion-"+p+"\" type=\"checkbox\" hidden />",
			"<label class=\"no-padding\" for=\"accordion-"+p+"\"><em>(click to show)</em></label>",
			"<div class=\"bd\"><pre>{%%s "+ToGoString(t, nullable, prop, false)+" %%}</pre></div>",
			"</li>",
			"</ul>",
		)
		return ret.Join("")
	case FmtJSON.Key, FmtHTML.Key, FmtSQL.Key:
		return tmplStartEQ + "view.Format(" + ToGoString(t, nullable, prop, false) + ", \"" + format + "\"" + tmplEndP
	case FmtColor.Key:
		return tmplStartEQ + "view.Color(" + ToGoString(t, nullable, prop, false) + componentEnd + tmplEnd
	case FmtCountry.Key:
		if verbose {
			x := ToGoString(t, nullable, prop, false)
			return tmplStart + key + " " + x + tmplEnd + " " + tmplStartS + "util.CountryFlag(" + x + tmplEndP
		}
		return tmplStart + key + " " + ToGoString(t, nullable, prop, false) + tmplEnd
	case FmtIcon.Key:
		size := "18"
		if src == util.KeyDetail {
			size = "64"
		}
		return tmplStartEQ + "components.Icon(" + ToGoString(t, nullable, prop, false) + ", " + size + componentEnd + tmplEnd
	case FmtImage.Key:
		size := "128px"
		msg := `<img style="max-width: %s; max-height: %s" src="{%%%%s %s %%%%}" />`
		return fmt.Sprintf(msg, size, size, ToGoString(t, nullable, prop, false))
	case FmtSelect.Key:
		return "<strong>" + tmplStart + key + " " + ToGoString(t, nullable, prop, false) + tmplEnd + "</strong>"
	case FmtURL.Key:
		x := tmplStart + key + " " + ToGoString(t, nullable, prop, false) + tmplEnd
		return fmt.Sprintf("<a href=%q target=\"_blank\" rel=\"noopener noreferrer\">%s</a>", x, x)
	case "":
		return tmplStartEQ + "view.String(" + prop + tmplEndP
	default:
		return "INVALID_STRING_FORMAT[" + format + "]"
	}
}
