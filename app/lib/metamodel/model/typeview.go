package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func ToGoViewString(t *types.Wrapped, prop string, nullable bool, format string, display string, verbose bool, url bool, enums enum.Enums, src string) string {
	switch t.Key() {
	case types.KeyAny, types.KeyJSON:
		if src == util.KeySimple {
			return tmplStartV + prop + tmplEnd
		}
		return jsonGoViewString(prop)
	case types.KeyBool:
		return tmplStartEQ + "view.BoolIcon(" + prop + ", 18, \"\", ps" + tmplEndP
		// return tmplStartV + prop + tmplEnd
	case types.KeyInt:
		switch format {
		case FmtSI.Key:
			return tmplStartS + fmt.Sprintf("util.ByteSizeSI(int64(%s))", prop) + tmplEnd
		case FmtSeconds.Key:
			return tmplStartS + fmt.Sprintf("view.DurationSeconds(float64(%s))", prop) + tmplEnd
		case "":
			return tmplStart + "d " + prop + tmplEnd
		default:
			return "INVALID_INT_FORMAT[" + format + "]"
		}
	case types.KeyFloat:
		switch format {
		case FmtSeconds.Key:
			return tmplStartEQ + fmt.Sprintf("view.DurationSeconds(%s)", prop) + tmplEnd
		case "":
			return tmplStart + "f " + prop + tmplEnd
		default:
			return "INVALID_FLOAT_FORMAT[" + format + "]"
		}
	case types.KeyList:
		return listGoViewString(t, prop, src, enums)
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
		if display == "summary" && src == "table" {
			return tmplStartEQ + "view.MapKeys(" + prop + tmplEndP
		}
		if src == util.KeySimple {
			return tmplStartV + prop + " %%}"
		}
		return jsonGoViewString(prop)
	case types.KeyDate:
		if nullable {
			return tmplStartEQ + "view.TimestampDay(" + prop + tmplEndP
		}
		return tmplStartEQ + "view.TimestampDay(&" + prop + tmplEndP
	case types.KeyEnum:
		e, _ := AsEnumInstance(t, enums)
		if e == nil || e.Simple() {
			return tmplStartV + ToGoString(t, nullable, prop, false) + tmplEnd
		}
		ret := tmplStartS + ToGoString(t, nullable, prop, false) + ".String()" + tmplEnd
		if e.HasValueIcons() {
			ret += " " + tmplStartEQ + "components.SVGIcon(" + ToGoString(t, nullable, prop, false) + ".Icon, ps)" + tmplEnd
		}
		return ret
	case types.KeyTimestamp, types.KeyTimestampZoned:
		if nullable {
			return tmplStartEQ + "view.Timestamp(" + prop + tmplEndP
		}
		return tmplStartEQ + "view.Timestamp(&" + prop + tmplEndP
	case types.KeyUUID:
		if nullable {
			return tmplStartEQ + "view.UUID(" + prop + tmplEndP
		}
		return tmplStartEQ + "view.UUID(&" + prop + tmplEndP
	case types.KeyString:
		return goViewStringForString(url, src, t, nullable, prop, format, verbose)
	default:
		return tmplStartV + ToGoString(t, nullable, prop, false) + tmplEnd
	}
}

func jsonGoViewString(prop string) string {
	return tmplStartEQ + "components.JSON(" + prop + tmplEndP
}

func saGoViewString(x string) string {
	return tmplStartEQ + "view.StringArray(" + x + tmplEndP
}

func listGoViewString(t *types.Wrapped, prop string, src string, enums enum.Enums) string {
	if src == util.KeySimple {
		return tmplStartV + prop + tmplEnd
	}
	lt := t.ListType()
	if lt == nil {
		lt = types.NewString()
	}
	switch lt.Key() {
	case types.KeyString:
		return saGoViewString(prop)
	case types.KeyInt:
		return tmplStartEQ + fmt.Sprintf("view.IntArray(util.ArrayFromAny[any](%s))", prop) + tmplEnd
	case types.KeyEnum:
		e, _ := AsEnumInstance(lt, enums)
		if e == nil {
			return "ERROR: invalid enum [" + lt.String() + "]"
		}
		return saGoViewString(prop + stringsSuffix)
	default:
		return jsonGoViewString(prop)
	}
}
