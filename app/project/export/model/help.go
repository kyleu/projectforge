package model

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
)

func Help(t types.Type, f string) string {
	switch t.Key() {
	case types.KeyAny:
		return "Interface, could be anything"
	case types.KeyBool:
		return "Value [true] or [false]"
	case types.KeyInt:
		return "Integer"
	case types.KeyFloat:
		return "Floating-point number"
	case types.KeyList:
		return "Comma-separated list of values"
	case types.KeyMap, types.KeyValueMap:
		return "JSON object"
	case types.KeyReference:
		return "[" + asRefK(t) + "], as a JSON object"
	case types.KeyString:
		switch f {
		case FmtURL:
			return "URL in string form"
		case FmtCountry:
			return "Two-digit country code"
		default:
			return "String text"
		}
	case types.KeyTimestamp:
		return "Date and time, in almost any format"
	case types.KeyUUID:
		return "UUID in format (00000000-0000-0000-0000-000000000000)"
	default:
		return t.Key()
	}
}

func (c *Column) Help() string {
	if c.HelpString != "" {
		return c.HelpString
	}
	ret := Help(c.Type, c.Format)
	if c.Nullable {
		ret += " (optional)"
	}
	return ret
}

func (m *Model) Help() string {
	ret := make([]string, 0, len(m.Columns)+1)
	ret = append(ret, "")
	for _, x := range m.Columns {
		ret = append(ret, " - "+x.Help())
	}
	return strings.Join(ret, "\n")
}
