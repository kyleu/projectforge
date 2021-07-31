package doctor

import (
	"fmt"
	"strings"
)

type Prognosis struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Status  string   `json:"status,omitempty"`
	Summary string   `json:"summary,omitempty"`
	Errors  Errors   `json:"errors,omitempty"`
	Logs    []string `json:"logs,omitempty"`
}

func NewPrognosis(key string, title string, summary string) *Prognosis {
	return &Prognosis{Key: key, Title: title, Status: "ok", Summary: summary}
}

func (p *Prognosis) AddLog(msg string) *Prognosis {
	p.Logs = append(p.Logs, msg)
	return p
}

func (p *Prognosis) WithError(err *Error) *Prognosis {
	p.Status = "error"
	p.Errors = append(p.Errors, err)
	return p
}

func (p *Prognosis) String() string {
	logs := strings.Builder{}
	for _, l := range p.Logs {
		logs.WriteString("\n- ")
		logs.WriteString(l)
	}
	return fmt.Sprintf("%s: %s\n[%s]%s", p.Title, p.Status, p.Summary, logs.String())
}

type Prognoses []*Prognosis
