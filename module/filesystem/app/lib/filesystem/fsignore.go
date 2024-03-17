package filesystem

import (
	"strings"

	"{{{ .Package }}}/app/util"
)

var defaultIgnore = []string{".DS_Store$", "^.git/", "^.idea/", "^build/", "^client/node_modules", ".html.go$", ".sql.go$"}

func buildIgnore(ign []string) []string {
	ret := util.NewStringSlice(append([]string{}, defaultIgnore...))
	ret.Push(ign...)
	return ret.Slice
}

const (
	keyPrefix = "^"
	keySuffix = "$"
)

func checkIgnore(ignore []string, fp string) bool {
	for _, i := range ignore {
		switch {
		case strings.HasPrefix(i, keyPrefix):
			i = strings.TrimPrefix(i, keyPrefix)
			if fp == strings.TrimSuffix(i, "/") || fp == strings.TrimSuffix(i, "\\") {
				return true
			}
			if strings.HasPrefix(fp, i) {
				return true
			}
		case strings.HasSuffix(i, keySuffix):
			if strings.HasSuffix(fp, strings.TrimSuffix(i, keySuffix)) {
				return true
			}
		case fp == i:
			return true
		}
	}
	return false
}
