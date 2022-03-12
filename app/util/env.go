// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

func ReplaceEnvVars(s string, logger *zap.SugaredLogger) string {
	sIdx := strings.Index(s, "${")
	if sIdx > -1 {
		eIdx := strings.Index(s[sIdx:], "}")
		if eIdx > -1 {
			orig := s[sIdx : sIdx+eIdx+1]

			n := orig[2 : len(orig)-1]
			d := ""

			dIdx := strings.Index(orig, "|")
			if dIdx > -1 {
				n = orig[2:dIdx]
				d = orig[dIdx+1 : len(orig)-1]
			}

			o := GetEnv(n)
			// logger.Debug(fmt.Sprintf("Replacing [%v] in address (original: %v, env[%v]: (%v), default: %v)", s, orig, n, o, d))
			if o == "" {
				o = d
			}
			return ReplaceEnvVars(strings.Replace(s, orig, o, 1), logger)
		}
	}

	return s
}

func GetEnv(name string, defaultValue ...string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	if l := strings.ToLower(name); l != name {
		if v := os.Getenv(l); v != "" {
			return v
		}
	}
	if u := strings.ToUpper(name); u != name {
		if v := os.Getenv(u); v != "" {
			return v
		}
	}
	return strings.Join(defaultValue, "")
}
