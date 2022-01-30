package model

import (
	"strings"

	"github.com/kyleu/projectforge/app/lib/types"
)

func Help(t types.Type) string {
	switch t.Key() {
	case types.KeyAny:
		return "Interface, could be anything"
	case types.KeyBool:
		return "Value [true] or [false]"
	case types.KeyInt:
		return "Integer"
	case types.KeyMap:
		return "JSON object"
	case types.KeyString:
		return "String text"
	case types.KeyTimestamp:
		return "Date and time, in almost any format"
	case types.KeyUUID:
		return "UUID in format (00000000-0000-0000-0000-000000000000)"
	default:
		return t.Key()
	}
}

func (c *Column) Help() string {
	ret := Help(c.Type)
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
