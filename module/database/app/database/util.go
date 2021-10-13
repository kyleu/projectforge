package database

import (
	"strings"
)

const (
	localhost = "localhost"
	falseKey  = "false"
)

func ArrayToString(a []string) string {
	return "{" + strings.Join(a, ",") + "}"
}

func StringToArray(s string) []string {
	split := strings.Split(strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{"), ",")
	ret := make([]string, 0)

	for _, x := range split {
		y := strings.TrimSpace(x)
		if len(y) > 0 {
			ret = append(ret, y)
		}
	}

	return ret
}
