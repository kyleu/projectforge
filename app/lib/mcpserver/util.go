package mcpserver

import "projectforge.dev/projectforge/app/util"

func valToText(x any) string {
	switch t := x.(type) {
	case nil:
		return "<no result>"
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return util.ToJSON(t)
	}
}
