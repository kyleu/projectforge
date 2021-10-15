package doctor

import (
	"strings"

	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

type Result struct {
	Check    *Check   `json:"-"`
	Key      string   `json:"key"`
	Title    string   `json:"title"`
	Status   string   `json:"status,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	Errors   Errors   `json:"errors,omitempty"`
	Duration int      `json:"duration,omitempty"`
	Solution string   `json:"solution,omitempty"`
	Logs     []string `json:"logs,omitempty"`
}

func NewResult(check *Check, key string, title string, summary string) *Result {
	return &Result{Check: check, Key: key, Title: title, Status: "OK", Summary: summary}
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

type Results []*Result

func SimpleOut(path string, cmd string, args []string, outCheck func(r *Result, out string) *Result) func(r *Result, logger *zap.SugaredLogger) *Result {
	return func(r *Result, logger *zap.SugaredLogger) *Result {
		fullCmd := strings.Join(append([]string{cmd}, args...), " ")
		exitCode, out, err := util.RunProcessSimple(fullCmd, path)
		if err != nil {
			msg := "[%s] is not present on your computer"
			return r.WithError(NewError("missing", msg, cmd))
		}
		if exitCode != 0 {
			return r.WithError(NewError("exitcode", "[%s] returned [%d] as an exit code", fullCmd, exitCode))
		}
		return outCheck(r, out)
	}
}
