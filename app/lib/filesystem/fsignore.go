package filesystem

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

var defaultIgnore = []string{".DS_Store$", "^.git/", "^.idea/", "^build/", "^client/node_modules", ".html.go$", ".sql.go$"}

func buildIgnore(ign []string) []string {
	if len(ign) == 1 && ign[0] == "-" {
		return nil
	}
	ret := util.NewStringSlice(util.ArrayCopy(defaultIgnore))
	ret.Push(ign...)
	return ret.Slice
}

const (
	keyPrefix   = "^"
	keySuffix   = "$"
	keyWildcard = "*"
)

func checkIgnore(ignore []string, fp string) bool {
	for _, i := range ignore {
		if fp == i {
			return true
		}
		if strings.HasPrefix(i, keyWildcard) && strings.HasSuffix(i, keyWildcard) {
			if strings.Contains(fp, strings.TrimSuffix(strings.TrimPrefix(i, keyWildcard), keyWildcard)) {
				return true
			}
		}
		if strings.HasPrefix(i, keyPrefix) || strings.HasSuffix(i, keyWildcard) {
			i = strings.TrimPrefix(strings.TrimSuffix(i, keyWildcard), keyPrefix)
			if fp == strings.TrimSuffix(i, "/") || fp == strings.TrimSuffix(i, "\\") {
				return true
			}
			if strings.HasPrefix(fp, i) {
				return true
			}
		}
		if strings.HasSuffix(i, keySuffix) || strings.HasPrefix(i, keyWildcard) {
			if strings.HasSuffix(fp, strings.TrimSuffix(strings.TrimPrefix(i, keyWildcard), keySuffix)) {
				return true
			}
		}
	}
	return false
}

func checkPatterns(patterns []string, pth string) bool {
	if len(patterns) == 0 {
		return true
	}
	return checkIgnore(patterns, pth)
}
