package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReplaceEnvVars(s string, logger Logger) string {
	sIdx := strings.Index(s, "${")
	if sIdx > -1 {
		eIdx := strings.Index(s[sIdx:], "}")
		if eIdx > -1 {
			orig := s[sIdx : sIdx+eIdx+1]

			n := orig[2 : len(orig)-1]
			var d string

			dIdx := strings.Index(orig, "|")
			if dIdx > -1 {
				n = orig[2:dIdx]
				d = orig[dIdx+1 : len(orig)-1]
			}

			o := OrDefault(GetEnv(n), d)
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
	return StringJoin(defaultValue, "")
}

func GetEnvBool(name string, defaultValue bool) bool {
	return GetEnv(name, fmt.Sprint(defaultValue)) == BoolTrue
}

func GetEnvBoolAny(defaultValue bool, names ...string) bool {
	for _, name := range names {
		if v := GetEnv(name, ""); v != "" {
			return v == BoolTrue
		}
	}
	return defaultValue
}

func GetEnvInt(name string, defaultValue int) int {
	v := GetEnv(name, "")
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(i)
}

func GetEnvDuration(name string, defaultValue time.Duration) time.Duration {
	v := GetEnv(name, "")
	ret, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue
	}
	return ret
}
