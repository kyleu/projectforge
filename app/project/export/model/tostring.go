package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
)

func TypeToString(t types.Type, prop string) string {
	switch t.Key() {
	case types.KeyUUID:
		return prop + ".String()"
	default:
		return prop
	}
}

func TypeToViewString(t types.Type, prop string, nullable bool) string {
	ret := ToGoString(t, prop, false)
	switch t.Key() {
	case types.KeyDate:
		if nullable {
			return "{%%= components.DisplayTimestampDay(" + ret + ") %%}"
		}
		return "{%%= components.DisplayTimestampDay(&" + ret + ") %%}"
	case types.KeyTimestamp:
		if nullable {
			return "{%%= components.DisplayTimestamp(" + ret + ") %%}"
		}
		return "{%%= components.DisplayTimestamp(&" + ret + ") %%}"
	default:
		return "{%%s " + ret + " %%}"
	}
}

func DetectLL(msg string, args ...any) string {
	ret := fmt.Sprintf(msg, args...)
	if len(ret) > 160 {
		return ret + " //nolint:lll"
	}
	return ret
}
