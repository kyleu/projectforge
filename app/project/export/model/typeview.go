package model

import (
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/util"
)

func ToGoViewString(t *types.Wrapped, prop string, nullable bool, format string, verbose bool, url bool, enums enum.Enums, src string) string {
	switch t.Key() {
	case types.KeyAny:
		if src == util.KeySimple {
			return "{%%v " + prop + " %%}"
		}
		return "{%%= components.JSON(" + prop + ") %%}"
	case types.KeyBool:
		return "{%%v " + prop + " %%}"
	case types.KeyInt:
		switch format {
		case FmtSI.Key:
			return "{%%s util.ByteSizeSI(int64(" + prop + ")) %%}"
		case "":
			return "{%%d " + prop + " %%}"
		default:
			return "INVALID_INT_FORMAT[" + format + "]"
		}
	case types.KeyFloat:
		return "{%%f " + prop + " %%}"
	case types.KeyList:
		if src == util.KeySimple {
			return "{%%v " + prop + " %%}"
		}
		lt := t.ListType()
		if lt == nil {
			lt = types.NewString()
		}
		switch lt.Key() {
		case types.KeyString:
			return "{%%= view.StringArray(" + prop + ") %%}"
		case types.KeyInt:
			return "{%%= view.IntArray(util.ArrayFromAny(" + prop + ")) %%}"
		case types.KeyEnum:
			e, _ := AsEnumInstance(lt, enums)
			if e == nil {
				return "ERROR: invalid enum [" + lt.String() + "]"
			}
			return "{%%= view.StringArray(" + prop + ".Strings()) %%}"
		default:
			return "{%%= components.JSON(" + prop + ") %%}"
		}
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
		if src == util.KeySimple {
			return "{%%v " + prop + " %%}"
		}
		return "{%%= components.JSON(" + prop + ") %%}"
	case types.KeyDate:
		if nullable {
			return "{%%= view.TimestampDay(" + prop + ") %%}"
		}
		return "{%%= view.TimestampDay(&" + prop + ") %%}"
	case types.KeyEnum:
		e, _ := AsEnumInstance(t, enums)
		if e == nil || e.Simple() {
			return "{%%v " + ToGoString(t, nullable, prop, false) + " %%}"
		}
		return "{%%s " + ToGoString(t, nullable, prop, false) + ".String() %%}"
	case types.KeyTimestamp:
		if nullable {
			return "{%%= view.Timestamp(" + prop + ") %%}"
		}
		return "{%%= view.Timestamp(&" + prop + ") %%}"
	case types.KeyUUID:
		if nullable {
			return "{%%= view.UUID(" + prop + ") %%}"
		}
		return "{%%= view.UUID(&" + prop + ") %%}"
	case types.KeyString:
		return goViewStringForString(url, src, t, nullable, prop, format, verbose)
	default:
		return "{%%v " + ToGoString(t, nullable, prop, false) + " %%}"
	}
}
