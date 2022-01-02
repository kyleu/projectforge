package model

import (
	"strings"
)

func (t *Type) Help() string {
	switch t.Key {
	case TypeInt.Key:
		return "Integer"
	case TypeMap.Key:
		return "JSON object"
	case TypeString.Key:
		return "String text"
	case TypeTimestamp.Key:
		return "Date and time, in almost any format"
	case TypeUUID.Key:
		return "UUID in format (00000000-0000-0000-0000-000000000000)"
	default:
		return t.Key
	}
}

func (c *Column) Help() string {
	ret := c.Type.Help()
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
