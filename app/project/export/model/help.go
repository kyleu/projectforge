package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
)

func Help(t types.Type, f string, nullable bool, enums enum.Enums) (string, error) {
	q := func(s string) string {
		if nullable {
			s += " (optional)"
		}
		return "\"" + s + "\""
	}
	switch t.Key() {
	case types.KeyAny:
		return q("Interface, could be anything"), nil
	case types.KeyBool:
		return q("Value [true] or [false]"), nil
	case types.KeyEnum:
		e, err := AsEnumInstance(t, enums)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.All%s.Help()", e.Package, e.ProperPlural()), nil
	case types.KeyInt:
		return q("Integer"), nil
	case types.KeyFloat:
		return q("Floating-point number"), nil
	case types.KeyList:
		return q("Comma-separated list of values"), nil
	case types.KeyMap, types.KeyValueMap:
		return q("JSON object"), nil
	case types.KeyReference:
		return q("[" + asRefK(t) + "], as a JSON object"), nil
	case types.KeyString:
		switch f {
		case FmtURL.Key:
			return q("URL in string form"), nil
		case FmtCountry.Key:
			return q("Two-digit country code"), nil
		default:
			return q("String text"), nil
		}
	case types.KeyDate:
		return q("Calendar date"), nil
	case types.KeyTimestamp:
		return q("Date and time, in almost any format"), nil
	case types.KeyUUID:
		return q("UUID in format (00000000-0000-0000-0000-000000000000)"), nil
	default:
		return q(t.Key()), nil
	}
}

func (c *Column) Help(enums enum.Enums) (string, error) {
	if c.HelpString != "" {
		return "\"" + c.HelpString + "\"", nil
	}
	ret, err := Help(c.Type, c.Format, c.Nullable, enums)
	if err != nil {
		return "", err
	}
	return ret, nil
}
