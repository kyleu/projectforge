package util

import (
	"strings"

	"go.uber.org/zap"
)

type Logger = *zap.SugaredLogger

var RootLogger Logger

var logSanitizer = strings.NewReplacer("\r", "", "\n", "")

func SanitizeLogValue(s string) string {
	if s == "" {
		return s
	}
	return logSanitizer.Replace(s)
}

func SanitizeLogArgs(args ...any) []any {
	if len(args) == 0 {
		return nil
	}
	ret := make([]any, len(args))
	for i, arg := range args {
		switch t := arg.(type) {
		case string:
			ret[i] = SanitizeLogValue(t)
		case []byte:
			ret[i] = SanitizeLogValue(string(t))
		default:
			ret[i] = arg
		}
	}
	return ret
}
