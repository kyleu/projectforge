package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func goViewStringForString(url bool, src string, t *types.Wrapped, prop string, format string, verbose bool) string {
	key := "s"
	if url {
		key = "u"
	}
	if src == util.KeySimple {
		return "{%%" + key + " " + prop + " %%}"
	}
	switch format {
	case FmtCode.Key:
		return "<pre>{%%s " + ToGoString(t, prop, false) + " %%}</pre>"
	case FmtLinebreaks.Key:
		return "<div class=\"prewsl\">{%%s " + ToGoString(t, prop, false) + " %%}</div>"
	case FmtCodeHidden.Key:
		_, p := util.StringSplitLast(prop, '.', true)
		ret := make([]string, 0, 30)
		ret = append(ret,
			"<ul class=\"accordion\">",
			"<li>",
			"<input id=\"accordion-"+p+"\" type=\"checkbox\" hidden />",
			"<label class=\"no-padding\" for=\"accordion-"+p+"\"><em>(click to show)</em></label>",
			"<div class=\"bd\"><pre>{%%s "+ToGoString(t, prop, false)+" %%}</pre></div>",
			"</li>",
			"</ul>",
		)
		return strings.Join(ret, "")
	case FmtHTML.Key:
		return "{%%= components.Format(" + ToGoString(t, prop, false) + ", \"html\") %%}</pre>"
	case FmtJSON.Key:
		return "{%%= components.Format(" + ToGoString(t, prop, false) + ", \"json\") %%}</pre>"
	case FmtURL.Key:
		x := "{%%" + key + " " + ToGoString(t, prop, false) + " %%}"
		return fmt.Sprintf("<a href=%q target=\"_blank\" rel=\"noopener noreferrer\">%s</a>", x, x)
	case FmtIcon.Key:
		size := "18"
		if src == util.KeyDetail {
			size = "64"
		}
		return "{%%= components.Icon(" + ToGoString(t, prop, false) + ", " + size + ", \"\", ps) %%}"
	case FmtCountry.Key:
		if verbose {
			return "{%%" + key + " " + ToGoString(t, prop, false) + " %%} {%%s util.CountryFlag(" + ToGoString(t, prop, false) + ") %%}"
		}
		return "{%%" + key + " " + ToGoString(t, prop, false) + " %%}"
	case FmtSelect.Key:
		return "<strong>{%%" + key + " " + ToGoString(t, prop, false) + " %%}</strong>"
	case "":
		return "{%%" + key + " " + ToGoString(t, prop, false) + " %%}"
	default:
		return "INVALID_STRING_FORMAT[" + format + "]"
	}
}
