package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/telemetry"
)

type Service struct {
	DatabaseName string
	SchemaName   string
	Username     string
	telemetry    *telemetry.Service
	logger       *zap.SugaredLogger
	db           *sqlx.DB
}

func NewService(dbName string, schName string, username string, telemetry *telemetry.Service, logger *zap.SugaredLogger, db *sqlx.DB) *Service {
	return &Service{DatabaseName: dbName, SchemaName: schName, Username: username, telemetry: telemetry, logger: logger, db: db}
}

func (s *Service) StartTransaction() (*sqlx.Tx, error) {
	if s.logger != nil {
		s.logger.Info("opening transaction")
	}
	return s.db.Beginx()
}

func errMessage(t string, q string, values []interface{}) string {
	return fmt.Sprintf("error running %s sql [%s] with values [%s]", t, strings.TrimSpace(q), valueStrings(values))
}

func (s *Service) logQuery(msg string, q string, values []interface{}) {
	if s.logger != nil {
		s.logger.Infof("%s {\n  SQL: %s\n  Values: %s\n}", msg, strings.TrimSpace(q), valueStrings(values))
	}
}

func (s *Service) newSpan(ctx context.Context, name string, q string) (context.Context, trace.Span) {
	nc, span := s.telemetry.StartSpan(ctx, "database", name)
	span.SetAttributes(
		semconv.DBStatementKey.String(q),
		semconv.DBSystemPostgreSQL,
		semconv.DBNameKey.String(s.DatabaseName),
		semconv.DBUserKey.String(s.Username),
	)
	return nc, span
}
