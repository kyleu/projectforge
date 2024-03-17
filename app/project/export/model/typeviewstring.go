package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func goViewStringForString(url bool, src string, t *types.Wrapped, nullable bool, prop string, format string, verbose bool) string {
	key := "s"
	if url {
		key = "u"
	}
	if src == util.KeySimple {
		return "{%%" + key + " " + prop + " %%}"
	}
	switch format {
	case FmtCode.Key:
		return "<pre class=\"prewsw\">{%%s " + ToGoString(t, nullable, prop, false) + " %%}</pre>"
	case FmtLinebreaks.Key:
		return "<div class=\"prewsl\">{%%s " + ToGoString(t, nullable, prop, false) + " %%}</div>"
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
		return "{%%= view.Format(" + ToGoString(t, nullable, prop, false) + ", \"" + format + "\") %%}"
	case FmtURL.Key:
		x := "{%%" + key + " " + ToGoString(t, nullable, prop, false) + " %%}"
		return fmt.Sprintf("<a href=%q target=\"_blank\" rel=\"noopener noreferrer\">%s</a>", x, x)
	case FmtIcon.Key:
		size := "18"
		if src == util.KeyDetail {
			size = "64"
		}
		return "{%%= components.Icon(" + ToGoString(t, nullable, prop, false) + ", " + size + ", \"\", ps) %%}"
	case FmtCountry.Key:
		if verbose {
			x := ToGoString(t, nullable, prop, false)
			return "{%%" + key + " " + x + " %%} {%%s util.CountryFlag(" + x + ") %%}"
		}
		return "{%%" + key + " " + ToGoString(t, nullable, prop, false) + " %%}"
	case FmtSelect.Key:
		return "<strong>{%%" + key + " " + ToGoString(t, nullable, prop, false) + " %%}</strong>"
	case "":
		return "{%%= view.String(" + prop + ") %%}"
	default:
		return "INVALID_STRING_FORMAT[" + format + "]"
	}
}
