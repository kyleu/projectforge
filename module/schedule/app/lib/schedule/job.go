package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

type JobFunc func(ctx context.Context, logger util.Logger) (any, error)

type Job struct {
	ID   uuid.UUID  `json:"id"`
	Name string     `json:"name"`
	Tags []string   `json:"tags"`
	Last *time.Time `json:"last"`
	Next *time.Time `json:"next"`
	Func JobFunc    `json:"-"`
}

func (j *Job) String() string {
	if j.Name == "" {
		return j.ID.String()
	}
	return fmt.Sprintf("%s (%s)", j.Name, j.ID)
}

func jobFromGoCron(j gocron.Job) *Job {
	var l, n *time.Time
	last, _ := j.LastRun()
	if !last.IsZero() {
		l = &last
	}
	next, _ := j.NextRun()
	if !next.IsZero() {
		n = &next
	}
	return &Job{ID: j.ID(), Name: j.Name(), Tags: j.Tags(), Last: l, Next: n}
}

type Jobs []*Job
