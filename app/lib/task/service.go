package task

import (
	"context"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

type Service struct {
	RegisteredTasks Tasks `json:"tasks"`

	path string
	fs   filesystem.FileLoader
}

func NewService(fs filesystem.FileLoader, path string, initialTasks ...*Task) *Service {
	return &Service{RegisteredTasks: initialTasks, path: path, fs: fs}
}

func (s *Service) RegisterTask(t *Task) error {
	if t == nil {
		return errors.New("a task must be provided")
	}
	if s.RegisteredTasks.Get(t.Key) != nil {
		return errors.Errorf("a task with key [%s] is already registered", t.Key)
	}
	s.RegisteredTasks = append(s.RegisteredTasks, t)
	s.RegisteredTasks.Sort()
	return nil
}

func (s *Service) Run(ctx context.Context, task *Task, run string, args util.ValueMap, logger util.Logger, fns ...exec.OutFn) *Result {
	ctx, span, logger := telemetry.StartSpan(ctx, "task:run:"+task.Key, logger)
	defer span.Complete()

	ret := task.Run(ctx, run, args, logger, fns...)
	// if err := s.SaveResult(task, logger); err != nil {
	// 	return ret.CompleteError(errors.Wrapf(err, "unable to save [%s] result for task [%s]", task.Key, task.ID))
	// }
	return ret
}

func (s *Service) RunAll(ctx context.Context, task *Task, run string, argsSet []util.ValueMap, logger util.Logger, fns ...exec.OutFn) (Results, error) {
	call := func(args util.ValueMap) (*Result, error) {
		if task == nil {
			return nil, nil
		}
		return s.Run(ctx, task, run, args, logger, fns...), nil
	}
	maxConcurrent := task.MaxConcurrent
	if maxConcurrent == -1 {
		maxConcurrent = 128
	}
	if maxConcurrent == 0 {
		maxConcurrent = runtime.NumCPU()
	}
	results, errs := util.AsyncRateLimit("run", argsSet, call, maxConcurrent, 12*time.Hour, logger)
	results = util.ArrayRemoveNil(results)
	return results, util.ErrorMerge(errs...)
}

func (s *Service) RemoveTask(key string) bool {
	if s.RegisteredTasks.Get(key) == nil {
		return false
	}
	s.RegisteredTasks = lo.Filter(s.RegisteredTasks, func(x *Task, _ int) bool {
		return x.Key == key
	})
	return true
}
