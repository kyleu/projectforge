package exec

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type OutFn func(key string, b []byte) error

func WriteOutFns(key string, b []byte, fns ...OutFn) error {
	for _, fn := range fns {
		if err := fn(key, b); err != nil {
			return err
		}
	}
	return nil
}

func WriteOutFnsString(s string, fns ...OutFn) {
	_ = WriteOutFns("default", []byte(s), fns...)
}

type writer struct {
	Key string
	fn  OutFn
}

func (w *writer) Write(p []byte) (int, error) {
	err := w.fn(w.Key, p)
	return len(p), err
}

type Exec struct {
	Key       string        `json:"key"`
	Idx       int           `json:"idx,omitempty"`
	Cmd       string        `json:"cmd,omitempty"`
	Env       []string      `json:"env,omitempty"`
	Path      string        `json:"path,omitempty"`
	Debug     bool          `json:"debug,omitempty"`
	Started   *time.Time    `json:"started,omitempty"`
	PID       int           `json:"pid,omitempty"`
	Completed *time.Time    `json:"completed,omitempty"`
	ExitCode  int           `json:"exitCode,omitempty"`
	Link      string        `json:"link,omitempty"`
	Buffer    *bytes.Buffer `json:"-"`
	execCmd   *exec.Cmd
	writer    io.Writer
}

func NewExec(key string, idx int, cmd string, path string, debug bool, envvars ...string) *Exec {
	if x := strings.Index(key, "/"); x > -1 && x < len(key) {
		key = key[x+1:]
	}
	return &Exec{Key: key, Idx: idx, Cmd: cmd, Env: envvars, Path: path, Debug: debug, Buffer: &bytes.Buffer{}}
}

func (e *Exec) WebPath() string {
	return fmt.Sprintf("/admin/exec/%s/%d", url.QueryEscape(e.Key), e.Idx)
}

func (e *Exec) Start(fns ...OutFn) error {
	if e.Started != nil {
		return errors.New("process already started")
	}
	e.writer = lo.Reduce(fns, func(agg io.Writer, fn OutFn, _ int) io.Writer {
		return io.MultiWriter(agg, &writer{Key: e.String(), fn: fn})
	}, io.Writer(e.Buffer))
	e.Started = util.TimeCurrentP()
	cmd, err := util.StartProcess(e.Cmd, e.Path, nil, e.writer, e.writer, e.Env...)
	if err != nil {
		return err
	}
	e.execCmd = cmd
	e.PID = cmd.Process.Pid
	return nil
}

func (e *Exec) Run(fns ...OutFn) error {
	t := util.TimerStart()
	err := e.Start(fns...)
	if err != nil {
		return err
	}
	err = e.Wait()
	if err != nil {
		_, _ = fmt.Fprintf(e.writer, " ::: error while wating for process to terminate %s", err.Error())
	}
	if e.Debug {
		_, _ = fmt.Fprintf(e.writer, " ::: process completed in [%s] with exit code [%d]", t.EndString(), e.ExitCode)
	}
	return nil
}

func (e *Exec) Kill() error {
	if e.execCmd == nil {
		return errors.New("not started")
	}
	e.Completed = util.TimeCurrentP()
	return e.execCmd.Process.Kill()
}

func (e *Exec) Wait() error {
	if e.execCmd == nil {
		return errors.New("not started")
	}
	exit, err := e.execCmd.Process.Wait()
	if err != nil {
		if err2 := e.execCmd.Wait(); err2 == nil {
			for (e.execCmd.ProcessState == nil || !e.execCmd.ProcessState.Exited()) && time.Since(util.TimeCurrent()) < (4*time.Second) {
				time.Sleep(500 * time.Millisecond)
			}
		}
		exit = e.execCmd.ProcessState
	}

	e.Completed = util.TimeCurrentP()
	if exit != nil {
		e.ExitCode = exit.ExitCode()
	}
	return nil
}

func (e *Exec) String() string {
	return fmt.Sprintf("%s:%d", e.Key, e.Idx)
}

type Execs []*Exec

func (m Execs) Get(key string, idx int) *Exec {
	return lo.FindOrElse(m, nil, func(x *Exec) bool {
		return x.Key == key && x.Idx == idx
	})
}

func (m Execs) GetByKey(key string) Execs {
	return lo.Filter(m, func(x *Exec, _ int) bool {
		return x.Key == key
	})
}

func (m Execs) Running() int {
	return lo.CountBy(m, func(x *Exec) bool {
		return x.Completed == nil
	})
}

func (m Execs) Sort() {
	slices.SortFunc(m, func(l *Exec, r *Exec) int {
		lk, rk := strings.ToLower(l.Key), strings.ToLower(r.Key)
		if lk != rk {
			return cmp.Compare(lk, rk)
		}
		return cmp.Compare(l.Idx, r.Idx)
	})
}
