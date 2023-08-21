package database

import (
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

const (
	MaxTracedStatements = 100
	MaxValueCount       = 20
	DebugLevelStatement = "statement"
	DebugLevelValues    = "values"
	DebugLevelAnalyze   = "analyze"
)

var (
	statements   = map[string]DebugStatements{}
	statementsMu = sync.Mutex{}
	lastIndex    = 0
	debugExample = &DebugStatement{
		Index: -1, SQL: "select * from test where a = {{{ .Placeholder 1 }}} and b = {{{ .Placeholder 2 }}}",
		Values: []any{util.ValueMap{"a": true}, util.ValueMap{"b": true}, util.ValueMap{"c": true}},
		Extra:  []util.ValueMap{{"example": "[example plan]"}}, Timing: 1,
		Message: "test query run without issue", Count: 2, Out: []any{1, 2},
	}
)

type DebugStatement struct {
	Index   int             `json:"index"`
	SQL     string          `json:"sql"`
	Values  []any           `json:"values,omitempty"`
	Extra   []util.ValueMap `json:"extra,omitempty"`
	Timing  int             `json:"timing,omitempty"`
	Error   string          `json:"error"`
	Message string          `json:"message,omitempty"`
	Count   int             `json:"count,omitempty"`
	Out     []any           `json:"out,omitempty"`
}

func (s *DebugStatement) SQLTrimmed(maxLength int) string {
	if len(s.SQL) > maxLength {
		return s.SQL[:maxLength] + "..."
	}
	return s.SQL
}

func (s *DebugStatement) Complete(count int, msg string, err error, output ...any) {
	s.Count = count
	s.Message = msg
	if err != nil {
		s.Error = err.Error()
	}
	if len(output) > MaxValueCount {
		s.Out = output[:MaxValueCount]
	} else {
		s.Out = output
	}
}

type DebugStatements []*DebugStatement

func (d DebugStatements) Add(st *DebugStatement) DebugStatements {
	if len(d) > MaxTracedStatements {
		d = d[1:]
	}
	d = append(d, st)
	return d
}

func GetDebugStatements(key string) DebugStatements {
	statementsMu.Lock()
	defer statementsMu.Unlock()
	return slices.Clone(statements[key])
}

func GetDebugStatement(key string, idx int) *DebugStatement {
	if idx == -1 {
		return debugExample
	}
	statementsMu.Lock()
	defer statementsMu.Unlock()
	return lo.FindOrElse(statements[key], nil, func(st *DebugStatement) bool {
		return st.Index == idx
	})
}

func (s *Service) newStatement(ctx context.Context, q string, values []any, timing int, logger util.Logger) (*DebugStatement, error) {
	lastIndex++
	ret := &DebugStatement{Index: lastIndex, SQL: q, Timing: timing}
	if s.tracing == DebugLevelValues || s.tracing == DebugLevelAnalyze {
		ret.Values = values
	}
	q = strings.TrimSpace(q)
	if s.tracing == DebugLevelAnalyze && !strings.HasPrefix(q, "explain") {
		ret.Extra = []util.ValueMap{{"status": "[explain in progress]"}}
		go func() {
			if a, err := s.Explain(ctx, q, values, logger); err == nil {
				ret.Extra = a
			} else {
				ret.Extra = []util.ValueMap{{"error": "[explain error] " + err.Error()}}
			}
		}()
	}
	return ret, nil
}

func (s *Service) Tracing() string {
	return s.tracing
}

func (s *Service) EnableTracing(v string, logger util.Logger) error {
	switch v {
	case DebugLevelStatement, DebugLevelValues, DebugLevelAnalyze, "":
		s.tracing = v
	default:
		return errors.Errorf("invalid tracing level [%s] must be [analyze], [values], or [statement]", v)
	}
	logger.Infof("database [%s] has tracing enabled in [%s] mode", s.Key, s.tracing)
	return nil
}

func (s *Service) DisableTracing(logger util.Logger) error {
	s.tracing = ""
	statementsMu.Lock()
	defer statementsMu.Unlock()
	delete(statements, s.Key)
	logger.Infof("database [%s] no longer has tracing enabled", s.Key)
	return nil
}

func (s *Service) addDebug(st *DebugStatement) {
	if s.tracing != "" {
		statementsMu.Lock()
		defer statementsMu.Unlock()
		statements[s.Key] = statements[s.Key].Add(st)
	}
}
