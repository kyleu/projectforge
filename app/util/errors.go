// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type unwrappable interface {
	Unwrap() error
}

type ErrorFrame struct {
	Key string
	Loc string
}

type ErrorDetail struct {
	Message    string            `json:"message"`
	StackTrace errors.StackTrace `json:"-"`
	Cause      *ErrorDetail      `json:"cause,omitempty"`
}

func GetErrorDetail(e error) *ErrorDetail {
	var stack errors.StackTrace

	// nolint
	t, ok := e.(stackTracer)
	if ok {
		stack = t.StackTrace()
	}

	var cause *ErrorDetail

	// nolint
	u, ok := e.(unwrappable)
	if ok {
		cause = GetErrorDetail(u.Unwrap())
	}

	return &ErrorDetail{
		Message:    e.Error(),
		StackTrace: stack,
		Cause:      cause,
	}
}

func TraceDetail(trace errors.StackTrace) []ErrorFrame {
	s := fmt.Sprintf("%+v", trace)
	lines := strings.Split(s, "\n")
	var validLines []string

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if len(l) > 0 {
			validLines = append(validLines, l)
		}
	}

	var ret []ErrorFrame
	for i := 0; i < len(validLines)-1; i += 2 {
		f := ErrorFrame{Key: validLines[i], Loc: validLines[i+1]}
		ret = append(ret, f)
	}
	return ret
}
