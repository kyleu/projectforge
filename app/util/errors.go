// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type unwrappable interface {
	Unwrap() error
}

type ErrorFrame struct {
	Key string `json:"key" xml:"key"`
	Loc string `json:"loc" xml:"loc"`
}

type ErrorDetail struct {
	Type       string            `json:"type" xml:"type"`
	Message    string            `json:"message" xml:"message"`
	Stack      []ErrorFrame      `json:"stack,omitempty" xml:"stack,omitempty"`
	StackTrace errors.StackTrace `json:"-" xml:"-"`
	Cause      *ErrorDetail      `json:"cause,omitempty" xml:"cause,omitempty"`
}

func GetErrorDetail(e error, includeStack bool) *ErrorDetail {
	var stack errors.StackTrace
	var cause *ErrorDetail

	if includeStack {
		t, ok := e.(stackTracer)
		if ok {
			stack = t.StackTrace()
		}
		u, ok := e.(unwrappable)
		if ok {
			cause = GetErrorDetail(u.Unwrap(), includeStack)
		}
	}
	msg := KeyError
	if e != nil {
		msg = e.Error()
	}
	return &ErrorDetail{
		Type:       KeyError,
		Message:    msg,
		Stack:      TraceDetail(stack),
		StackTrace: stack,
		Cause:      cause,
	}
}

func TraceDetail(trace errors.StackTrace) []ErrorFrame {
	s := fmt.Sprintf("%+v", trace)
	lines := StringSplitLines(s)
	var validLines []string

	lo.ForEach(lines, func(line string, _ int) {
		if l := strings.TrimSpace(line); l != "" {
			validLines = append(validLines, l)
		}
	})

	var ret []ErrorFrame
	for i := 0; i < len(validLines)-1; i += 2 {
		f := ErrorFrame{Key: validLines[i], Loc: validLines[i+1]}
		ret = append(ret, f)
	}
	return ret
}

func ErrorMerge(errs ...error) error {
	errs = lo.Filter(errs, func(e error, _ int) bool {
		return e != nil
	})
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}
	msg := lo.Map(errs, func(e error, _ int) string {
		return e.Error()
	})
	return errors.Wrapf(errs[0], strings.Join(msg, ", "))
}
