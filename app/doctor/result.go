package doctor

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Result struct {
	Check     *Check   `json:"-"`
	Key       string   `json:"key"`
	Title     string   `json:"title"`
	Status    string   `json:"status,omitzero"`
	Summary   string   `json:"summary,omitzero"`
	Errors    Errors   `json:"errors,omitempty"`
	Duration  int      `json:"duration,omitzero"`
	Solutions []string `json:"solution,omitempty"`
	Logs      []string `json:"logs,omitempty"`
}

func NewResult(check *Check, key string, title string, summary string) *Result {
	return &Result{Check: check, Key: key, Title: title, Status: util.KeyOK, Summary: summary}
}

func (p *Result) AddLog(msg string) *Result {
	p.Logs = append(p.Logs, msg)
	return p
}

func (p *Result) CleanSolutions() []string {
	return lo.Map(p.Solutions, func(s string, _ int) string {
		return strings.TrimPrefix(strings.TrimPrefix(s, "!"), "#")
	})
}

func (p *Result) WithError(err *Error) *Result {
	p.Status = util.KeyError
	p.Errors = append(p.Errors, err)
	return p
}

func (p *Result) AddSolution(msg string) *Result {
	p.Solutions = append(p.Solutions, msg)
	return p
}

func (p *Result) AddPackageSolution(name string, pkg string) *Result {
	msg := fmt.Sprintf("Install [%s] using your platform's package manager", name)
	switch runtime.GOOS {
	case "windows":
		msg += fmt.Sprintf(" by running [choco install %s]", pkg)
	case "darwin":
		msg += fmt.Sprintf(" by running [brew install %s]", pkg)
	case "linux":
		if pkg == "golang" {
			msg += " by running [sudo snap install --classic go] (don't use [apt], it still has an old version)"
		} else {
			msg += fmt.Sprintf(" by running [sudo apt install %s]", pkg)
		}
	}
	return p.AddSolution(msg)
}

func (p *Result) ErrorString() string {
	if len(p.Errors) == 0 && len(p.Solutions) == 0 {
		return ""
	}
	msg := lo.Map(p.Errors, func(e *Error, _ int) string {
		return e.String()
	})
	msg = append(msg, lo.Map(p.Solutions, func(s string, _ int) string {
		return "   Solution: " + s
	})...)
	return util.StringJoin(msg, util.StringDefaultLinebreak)
}

type Results []*Result

func (r Results) Errors() Results {
	return lo.Filter(r, func(x *Result, _ int) bool {
		return x.Status == util.KeyError
	})
}

func (r Results) ErrorSummary() string {
	return util.StringJoin(lo.Map(r.Errors(), func(x *Result, _ int) string {
		return x.Summary
	}), ", ")
}
