package doctor

import (
	"fmt"

	"github.com/kyleu/projectforge/app/util"
)

type Dependency struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Summary string   `json:"summary,omitempty"`
	Cmd     string   `json:"cmd"`
	Args    []string `json:"args,omitempty"`
}

func (d *Dependency) Check() *Prognosis {
	ret := NewPrognosis(d.Key, d.Title, d.Summary)
	cmd := d.Cmd
	for _, arg := range d.Args {
		cmd += " " + arg
	}
	ret.AddLog("running [" + cmd + "]")
	exitCode, out, err := util.RunProcessSimple(cmd, ".")
	if err != nil {
		msg := fmt.Sprintf("[%s] is not present on your computer", cmd)
		return ret.WithError(NewError("missing", msg, "rsvg"))
	}
	if exitCode != 0 {
		msg := fmt.Sprintf("[%s] returned [%d] as an exit code", cmd, exitCode)
		return ret.WithError(NewError("exitcode", msg, "rsvg", fmt.Sprint(exitCode), out))
	}
	return ret.AddLog("out: " + out)
}

type Dependencies = []*Dependency
