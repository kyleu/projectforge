// Package filesystem - Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import "strings"

var defaultIgnore = []string{".DS_Store$", "^.git/", "^.idea/", "^build/", "^client/node_modules", ".html.go$", ".sql.go$"}

func buildIgnore(ign []string) []string {
	ret := append([]string{}, defaultIgnore...)
	ret = append(ret, ign...)
	return ret
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
