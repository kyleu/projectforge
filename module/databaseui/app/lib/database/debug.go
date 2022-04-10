package database

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const maxStatements = 100

type DebugStatement struct {
	SQL    string      `json:"sql"`
	Values []any       `json:"values,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
	Timing int         `json:"timing,omitempty"`
}

func NewStatement(ctx context.Context, s *Service, q string, values []any, timing int) (*DebugStatement, error) {
	ret := &DebugStatement{SQL: q, Timing: timing}
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

func (s *Service) Tracing() string {
	return s.tracing
}

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

func GetDebugStatements(key string) DebugStatements {
	statementsMu.Lock()
	defer statementsMu.Unlock()
	return slices.Clone(statements[key])
}

func (s *Service) addDebug(st *DebugStatement) {
	if s.tracing != "" {
		statementsMu.Lock()
		defer statementsMu.Unlock()
		statements[s.Key] = statements[s.Key].Add(st)
	}
}
