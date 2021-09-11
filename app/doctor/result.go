package doctor

import (
	"fmt"
	"strings"
)

type Result struct {
	Check   *Check   `json:"-"`
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Status  string   `json:"status,omitempty"`
	Summary string   `json:"summary,omitempty"`
	Errors  Errors   `json:"errors,omitempty"`
	Logs    []string `json:"logs,omitempty"`
}

func NewResult(check *Check, key string, title string, summary string) *Result {
	return &Result{Check: check, Key: key, Title: title, Status: "ok", Summary: summary}
}

func (p *Result) AddLog(msg string) *Result {
	p.Logs = append(p.Logs, msg)
	return p
}

func (p *Result) WithError(err *Error) *Result {
	p.Status = "error"
	p.Errors = append(p.Errors, err)
	return p
}

func (p *Result) String() string {
	logs := strings.Builder{}
	for _, l := range p.Logs {
		logs.WriteString("\n- ")
		logs.WriteString(l)
	}
	return fmt.Sprintf("%s: %s\n[%s]%s", p.Title, p.Status, p.Summary, logs.String())
}

type Results []*Result
