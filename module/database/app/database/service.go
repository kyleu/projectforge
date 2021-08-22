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

type Service struct {
	DatabaseName string
	SchemaName   string
	Username     string
	db           *sqlx.DB
	metrics      *dbmetrics.Metrics
	logger       *zap.SugaredLogger
}

func NewService(dbName string, schName string, username string, db *sqlx.DB, logger *zap.SugaredLogger) *Service {
	m := dbmetrics.NewMetrics(dbName, db)
	return &Service{DatabaseName: dbName, SchemaName: schName, Username: username, db: db, metrics: m, logger: logger}
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
