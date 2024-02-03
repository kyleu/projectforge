package database

import (
	"strings"

	"github.com/samber/lo"
)

const (
	localhost = "localhost"
)

func ArrayToString(a []string) string {
	return "{" + strings.Join(a, ",") + "}"
}

func StringToArray(s string) []string {
	split := strings.Split(strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{"), ",")
	ret := make([]string, 0)
	lo.ForEach(split, func(x string, _ int) {
		y := strings.TrimSpace(x)
		if len(y) > 0 {
			ret = append(ret, y)
		}
	})
	return ret
}
