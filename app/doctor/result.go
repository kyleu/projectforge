package doctor

import (
	"strings"

	"github.com/samber/lo"
)

type Result struct {
	Check     *Check   `json:"-"`
	Key       string   `json:"key"`
	Title     string   `json:"title"`
	Status    string   `json:"status,omitempty"`
	Summary   string   `json:"summary,omitempty"`
	Errors    Errors   `json:"errors,omitempty"`
	Duration  int      `json:"duration,omitempty"`
	Solutions []string `json:"solution,omitempty"`
	Logs      []string `json:"logs,omitempty"`
}

func NewResult(check *Check, key string, title string, summary string) *Result {
	return &Result{Check: check, Key: key, Title: title, Status: "OK", Summary: summary}
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
	p.Status = "error"
	p.Errors = append(p.Errors, err)
	return p
}

func (p *Result) AddSolution(msg string) *Result {
	p.Solutions = append(p.Solutions, msg)
	return p
}

type Results []*Result

func (r Results) Errors() Results {
	return lo.Filter(r, func(x *Result, _ int) bool {
		return x.Status == "error"
	})
}

func (r Results) ErrorSummary() string {
	return strings.Join(lo.Map(r.Errors(), func(x *Result, _ int) string {
		return x.Summary
	}), ", ")
}
