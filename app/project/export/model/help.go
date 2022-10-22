package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
)

func Help(t types.Type, f string, enums enum.Enums) (string, error) {
	switch t.Key() {
	case types.KeyAny:
		return "Interface, could be anything", nil
	case types.KeyBool:
		return "Value [true] or [false]", nil
	case types.KeyEnum:
		e, err := AsEnumInstance(t, enums)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Available options: [%s]", strings.Join(e.Values, ", ")), nil
	case types.KeyInt:
		return "Integer", nil
	case types.KeyFloat:
		return "Floating-point number", nil
	case types.KeyList:
		return "Comma-separated list of values", nil
	case types.KeyMap, types.KeyValueMap:
		return "JSON object", nil
	case types.KeyReference:
		return "[" + asRefK(t) + "], as a JSON object", nil
	case types.KeyString:
		switch f {
		case FmtURL:
			return "URL in string form", nil
		case FmtCountry:
			return "Two-digit country code", nil
		default:
			return "String text", nil
		}
	case types.KeyTimestamp:
		return "Date and time, in almost any format", nil
	case types.KeyUUID:
		return "UUID in format (00000000-0000-0000-0000-000000000000)", nil
	default:
		return t.Key(), nil
	}
}

func (c *Column) Help(enums enum.Enums) (string, error) {
	if c.HelpString != "" {
		return c.HelpString, nil
	}
	ret, err := Help(c.Type, c.Format, enums)
	if err != nil {
		return "", err
	}
	if c.Nullable {
		ret += " (optional)"
	}
	return ret, nil
}

func (m *Model) Help(enums enum.Enums) (string, error) {
	ret := make([]string, 0, len(m.Columns)+1)
	ret = append(ret, "")
	for _, x := range m.Columns {
		ch, err := x.Help(enums)
		if err != nil {
			return "", err
		}
		ret = append(ret, " - "+ch)
	}
	return strings.Join(ret, "\n"), nil
}
