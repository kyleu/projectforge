package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/telemetry"
	"{{{ .Package }}}/app/telemetry/dbmetrics"
)

type DBType struct {
	Key               string `json:"key"`
	Title             string `json:"title"`
	Quote             string `json:"-"`
	Placeholder       string `json:"-"`
	SupportsReturning bool   `json:"-"`
}

type Service struct {
	Key          string  `json:"key"`
	DatabaseName string  `json:"database,omitempty"`
	SchemaName   string  `json:"schema,omitempty"`
	Username     string  `json:"username,omitempty"`
	Type         *DBType `json:"type"`
	db           *sqlx.DB
	metrics      *dbmetrics.Metrics
	logger       *zap.SugaredLogger
}

func NewService(typ *DBType, key string, dbName string, schName string, username string, db *sqlx.DB, logger *zap.SugaredLogger) (*Service, error) {
	m, err := dbmetrics.NewMetrics(key, db)
	if err != nil {
		return nil, err
	}
	return &Service{Key: key, DatabaseName: dbName, SchemaName: schName, Username: username, Type: typ, db: db, metrics: m, logger: logger}, nil
}

func (s *Service) StartTransaction() (*sqlx.Tx, error) {
	if s.logger != nil {
		s.logger.Info("opening transaction")
	}
	return s.db.Beginx()
}

func (s *Service) Stats() sql.DBStats {
	return s.db.Stats()
}

func errMessage(t string, q string, values []interface{}) string {
	return fmt.Sprintf("error running %s sql [%s] with values [%s]", t, strings.TrimSpace(q), valueStrings(values))
}

func (s *Service) logQuery(msg string, q string, values []interface{}) {
	if s.logger != nil {
		s.logger.Infof("%s {\n  SQL: %s\n  Values: %s\n}", msg, strings.TrimSpace(q), valueStrings(values))
	}
}

func (s *Service) newSpan(ctx context.Context, name string, q string) (time.Time, context.Context, trace.Span) {
	s.metrics.IncStmt(q, name)
	nc, span := telemetry.StartSpan(ctx, "database", name)
	span.SetAttributes(
		semconv.DBStatementKey.String(q),
		semconv.DBSystemPostgreSQL,
		semconv.DBNameKey.String(s.DatabaseName),
		semconv.DBUserKey.String(s.Username),
	)
	return time.Now(), nc, span
}

func (s *Service) complete(q string, op string, span trace.Span, started time.Time, err error) {
	span.End()
	s.metrics.CompleteStmt(q, op, started, err)
}

func (s *Service) Close() error {
	_ = s.metrics.Close()
	return s.db.Close()
}
