package schedule

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type Service struct {
	Engine   gocron.Scheduler      `json:"-"`
	Started  time.Time             `json:"started"`
	Results  map[uuid.UUID]*Result `json:"results,omitempty"`
	resultMu sync.Mutex
}

func NewService() *Service {
	engine, _ := gocron.NewScheduler()
	if engine != nil {
		engine.Start()
	}
	return &Service{Engine: engine, Started: time.Now(), Results: map[uuid.UUID]*Result{}}
}

func (s *Service) NewJob(
	ctx context.Context, name string, sched gocron.JobDefinition, f JobFunc, singleton bool, logger util.Logger, tags ...string,
) (*Job, error) {
	opts := []gocron.JobOption{gocron.WithName(name), gocron.WithTags(tags...)}
	if singleton {
		opts = append(opts, gocron.WithSingletonMode(gocron.LimitModeReschedule))
	}
	var id uuid.UUID
	wrapped := func(ctx context.Context, logger util.Logger) {
		t := util.TimerStart()
		var sp *telemetry.Span
		ctx, sp, logger = telemetry.StartSpan(context.Background(), "job."+id.String(), logger)
		defer sp.Complete()
		logger.Debugf("running scheduled job [%s]", id.String())
		res := &Result{Job: id, Occurred: time.Now()}
		ret, err := f(ctx, logger)
		res.DurationMicro = t.End()
		res.Returned = ret
		if err != nil {
			res.Error = err.Error()
			logger.Warnf("error running scheduled job [%s]: %+v", id.String(), err)
		}
		logger.Debugf("completed scheduled job [%s] in [%s]: returned [%T]", id.String(), util.MicrosToMillis(res.DurationMicro), res.Returned)
		s.resultMu.Lock()
		defer s.resultMu.Unlock()
		s.Results[id] = res
	}
	x, err := s.Engine.NewJob(sched, gocron.NewTask(wrapped, ctx, logger), opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to start Schedule engine")
	}
	id = x.ID()
	return jobFromGoCron(x), nil
}

func (s *Service) ListJobs() Jobs {
	return lo.Map(s.Engine.Jobs(), func(x gocron.Job, _ int) *Job {
		return jobFromGoCron(x)
	})
}

func (s *Service) GetJob(id uuid.UUID) *Job {
	return lo.FindOrElse(s.ListJobs(), nil, func(j *Job) bool {
		return j.ID == id
	})
}

func (s *Service) Testbed(ctx context.Context, logger util.Logger) (any, error) {
	ret := "TODO"
	return ret, nil
}