package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type ResultLogFn func(key string, data any)

type Result struct {
	ID       uuid.UUID     `json:"id"`
	Task     *Task         `json:"task"`
	Run      string        `json:"run,omitempty"`
	Args     util.ValueMap `json:"args,omitempty"`
	Started  time.Time     `json:"started,omitempty"`
	Elapsed  int           `json:"elapsed,omitempty"`
	Status   string        `json:"status,omitempty"`
	Summary  string        `json:"summary,omitempty"`
	Logs     []string      `json:"logs,omitempty"`
	Data     any           `json:"data,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Metadata util.ValueMap `json:"metadata,omitempty"`
	Success  bool          `json:"success"`
	Error    string        `json:"error,omitempty"`
	syncMu   *sync.Mutex
	fns      []ResultLogFn
}

func NewResult(task *Task, run string, args util.ValueMap, fns ...ResultLogFn) *Result {
	return &Result{
		ID: util.UUID(), Task: task, Run: run, Args: args, Started: time.Now(), Status: "ok",
		Metadata: util.ValueMap{}, syncMu: &sync.Mutex{}, fns: fns,
	}
}

func CompletedResult(task *Task, run string, args util.ValueMap, data any, err error, logs ...string) *Result {
	ret := NewResult(task, run, args)
	ret.AddLogs(logs...)
	ret.Complete(data, err)
	return ret
}

func (r *Result) IsOK() bool {
	return r.Status == util.OK
}

func (r *Result) Log(msg string, args ...any) {
	r.AddLogs(fmt.Sprintf(msg, args...))
}

func (r *Result) AddLogs(msgs ...string) {
	r.syncMu.Lock()
	defer r.syncMu.Unlock()
	r.Logs = append(r.Logs, msgs...)
	for _, fn := range r.fns {
		for _, msg := range msgs {
			fn("log", msg)
		}
	}
}

func (r *Result) AddTags(tags ...string) {
	r.syncMu.Lock()
	defer r.syncMu.Unlock()
	r.Tags = util.ArraySorted(lo.Uniq(append(r.Tags, tags...)))
}

func (r *Result) HasTag(tag string) bool {
	return lo.Contains(r.Tags, tag)
}

func (r *Result) Complete(data any, errs ...error) *Result {
	r.Data = data
	if err := util.ErrorMerge(errs...); err != nil {
		r.Error = err.Error()
		r.Status = util.Error
	} else if r.Status == "" {
		r.Status = util.OK
	}
	r.Elapsed = int(time.Since(r.Started).Microseconds())
	r.Log("task [%s] completed in [%s]", r.Task.TitleSafe(), util.MicrosToMillis(r.Elapsed))
	return r
}

func (r *Result) CompleteSimple(data any) *Result {
	return r.Complete(data)
}

func (r *Result) CompleteError(err error) *Result {
	return r.Complete(nil, nil, err)
}

func (r *Result) EndTime() time.Time {
	return r.Started.Add(time.Duration(r.Elapsed) * time.Microsecond)
}

func (r *Result) DataMap() util.ValueMap {
	if r == nil || r.Data == nil {
		return nil
	}
	ret, err := util.ParseMap(r.Data, "", true)
	if ret != nil && err == nil {
		return ret
	}
	ret, _ = util.FromJSONMap(util.ToJSONBytes(r.Data, true))
	return ret
}

func (r *Result) String() string {
	if r == nil {
		return "pending"
	}
	if r.Status == "" {
		return util.KeyUnknown
	}
	return r.Status
}

func (r *Result) Summarize() string {
	if r == nil {
		return "missing"
	}
	if r.Summary != "" {
		return r.Summary
	}
	return r.Status
}
