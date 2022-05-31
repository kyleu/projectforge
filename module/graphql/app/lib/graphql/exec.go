package graphql

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

func (s *Service) Exec(ctx context.Context, key string, q string, op string, vars map[string]any, logger util.Logger) (g *graphql.Response, e error) {
	r, ok := s.schemata[key]
	if !ok {
		return nil, errors.Errorf("no schema available with key [%s]", key)
	}

	ctx, span, logger := telemetry.StartSpan(ctx, "graphql", logger)
	defer span.Complete()

	defer func() {
		if rec := recover(); rec != nil {
			if recoverErr, ok := rec.(error); ok {
				e = errors.Wrap(recoverErr, "panic")
			} else {
				e = errors.Errorf("graphql encountered panic recovery of type [%T]: %s", rec, fmt.Sprint(rec))
			}
		}
	}()
	logger.Debugf("running GraphQL query")
	r.ExecCount += 1
	g = r.Schema.Exec(ctx, q, op, vars)
	return
}

func (s *Service) ExecCount(key string) int {
	r, ok := s.schemata[key]
	if !ok {
		return 0
	}
	return r.ExecCount
}
