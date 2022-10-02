// Content managed by Project Forge, see [projectforge.md] for details.
package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/util"
)

type writer struct {
	Key string
	fn  func(key string, b []byte) error
}

func (w *writer) Write(p []byte) (int, error) {
	err := w.fn(w.Key, p)
	return len(p), err
}

type Exec struct {
	Key       string        `json:"key"`
	Idx       int           `json:"idx"`
	Cmd       string        `json:"cmd"`
	Env       []string      `json:"env,omitempty"`
	Path      string        `json:"path"`
	Started   *time.Time    `json:"started"`
	PID       int           `json:"pid"`
	Completed *time.Time    `json:"completed"`
	ExitCode  int           `json:"exitCode"`
	Buffer    *bytes.Buffer `json:"-"`
	execCmd   *exec.Cmd
}

func NewExec(key string, idx int, cmd string, path string, envvars ...string) *Exec {
	return &Exec{Key: key, Idx: idx, Cmd: cmd, Env: envvars, Path: path, Buffer: &bytes.Buffer{}}
}

func (e *Exec) WebPath() string {
	return fmt.Sprintf("/admin/exec/%s/%d", e.Key, e.Idx)
}

func (e *Exec) Start(ctx context.Context, fn func(key string, b []byte) error, logger util.Logger) error {
	if e.Started != nil {
		return errors.New("process already started")
	}
	var w io.Writer = e.Buffer
	if fn != nil {
		w = io.MultiWriter(e.Buffer, &writer{Key: e.Key, fn: fn})
	}
	e.Started = util.NowPointer()
	cmd, err := util.StartProcess(e.Cmd, e.Path, nil, w, w, e.Env...)
	if err != nil {
		return err
	}
	e.execCmd = cmd
	e.PID = cmd.Process.Pid
	defer func() {
		go func() {
			_ = e.Wait()
		}()
	}()
	return nil
}

func (e *Exec) Kill() error {
	if e.execCmd == nil {
		return errors.New("not started")
	}
	return e.execCmd.Process.Kill()
}

func (e *Exec) Wait() error {
	if e.execCmd == nil {
		return errors.New("not started")
	}
	exit, err := e.execCmd.Process.Wait()
	if err != nil {
		return err
	}

	e.Completed = util.NowPointer()
	e.ExitCode = exit.ExitCode()
	return nil
}

func (e *Exec) String() string {
	return fmt.Sprintf("%s:%d", e.Key, e.Idx)
}

type Execs []*Exec

func (m Execs) Get(key string, idx int) *Exec {
	for _, x := range m {
		if x.Key == key && x.Idx == idx {
			return x
		}
	}
	return nil
}

func (m Execs) GetByKey(key string) Execs {
	var ret Execs
	for _, x := range m {
		if x.Key == key {
			ret = append(ret, x)
		}
	}
	return ret
}

func (m Execs) Running() int {
	var ret int
	for _, x := range m {
		if x.Completed == nil {
			ret++
		}
	}
	return ret
}

func (m Execs) Sort() {
	slices.SortFunc(m, func(l *Exec, r *Exec) bool {
		if l.Key != r.Key {
			return l.Key < r.Key
		}
		return l.Idx < r.Idx
	})
}
