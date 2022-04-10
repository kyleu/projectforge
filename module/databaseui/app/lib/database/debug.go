package database

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const maxStatements = 100

type DebugStatement struct {
	SQL    string      `json:"sql"`
	Values []any       `json:"values"`
	Extra  interface{} `json:"extra"`
}

func NewStatement(ctx context.Context, s *Service, q string, values []any) (*DebugStatement, error) {
	ret := &DebugStatement{SQL: q}
	if s.tracing == "values" || s.tracing == "analyze" {
		ret.Values = values
	}
	if s.tracing == "analyze" {
		ret.Extra = "TODO"
	}
	return ret, nil
}

type DebugStatements []*DebugStatement

func (d DebugStatements) Add(st *DebugStatement) DebugStatements {
	if len(d) > maxStatements {
		d = d[1:]
	}
	d = append(d, st)
	return d
}

var (
	statements   = map[string]DebugStatements{}
	statementsMu = sync.Mutex{}
)

func (s *Service) EnableTracing(includeValues bool, analyze bool, logger *zap.SugaredLogger) error {
	if analyze {
		s.tracing = "analyze"
	} else if includeValues {
		s.tracing = "values"
	} else {
		s.tracing = "statement"
	}
	logger.Infof("database [%s] has tracing enabled in [%s] mode", s.Key, s.tracing)
	return nil
}

func (s *Service) DisableTracing(logger *zap.SugaredLogger) error {
	s.tracing = ""
	statementsMu.Lock()
	defer statementsMu.Unlock()
	delete(statements, s.Key)
	logger.Infof("database [%s] no longer has tracing enabled", s.Key)
	return nil
}

func GetDebugStatements(key string) (DebugStatements, error) {
	curr, ok := statements[key]
	if !ok {
		return nil, errors.Errorf("no statements tracked for database [%s]", key)
	}
	return curr, nil
}

func (s *Service) addDebug(st *DebugStatement) {
	if s.tracing != "" {
		statementsMu.Lock()
		defer statementsMu.Unlock()
		statements[s.Key] = statements[s.Key].Add(st)
	}
}
