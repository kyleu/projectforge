package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/telemetry/dbmetrics"
	"{{{ .Package }}}/app/util"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/queries"
	"{{{ .Package }}}/queries/schema"{{{ end }}}
)

type Service struct {
	Key          string  `json:"key"`
	DatabaseName string  `json:"database,omitzero"`
	SchemaName   string  `json:"schema,omitzero"`
	Username     string  `json:"username,omitzero"`
	Debug        bool    `json:"debug,omitzero"`
	Type         *DBType `json:"type"`
	ReadOnly     bool    `json:"readonly,omitzero"`{{{ if .HasModule "databaseui" }}}
	tracing      string{{{ end }}}
	db           *sqlx.DB
	metrics      *dbmetrics.Metrics
	stringRep    string
}

func NewService(
	ctx context.Context, typ *DBType, key string, dbName string, schName string, username string, debug bool, db *sqlx.DB, stringRep string, logger util.Logger,
) (*Service, error) {
	if logger == nil {
		return nil, errors.New("logger must be provided to database service")
	}
	logger = logger.With("database", dbName, "user", username)
	m, err := dbmetrics.NewMetrics(strings.ReplaceAll(key, "-", "_"), db, logger)
	if err != nil {
		logger.Debugf("unable to register database metrics for [%s]: %+v", key, err)
	}

	ret := &Service{Key: key, DatabaseName: dbName, SchemaName: schName, Username: username, Debug: debug, Type: typ, db: db, metrics: m, stringRep: stringRep}
	err = ret.Healthcheck(ctx, dbName, db)
	if err != nil {
		return ret, errors.Wrap(err, "unable to run healthcheck")
	}
	register(ret, logger)
	return ret, nil
}

func (s *Service) String() string {
	return s.stringRep
}

func (s *Service) Healthcheck(ctx context.Context, dbName string, db *sqlx.DB) error {
	q := {{{ if .HasModule "migration" }}}queries.Healthcheck(){{{ else }}}"select 1"{{{ end }}}
	res, err := db.QueryContext(ctx, q)
	if err != nil || res.Err() != nil {
		if err == nil {
			err = res.Err()
		}
		if strings.Contains(err.Error(), "does not exist") {
			{{{ if .HasModule "migration" }}}return errors.Wrapf(err, "database [%s] does not exist; run the following:\n"+schema.CreateDatabase(), dbName){{{ else }}}return errors.Wrapf(err, "database does not exist"){{{ end }}}
		}
		return errors.Wrapf(err, "unable to run healthcheck [%s] for database [%s]", q, dbName)
	}
	defer func() { _ = res.Close() }()
	return nil
}

func (s *Service) StartTransaction(logger util.Logger) (*sqlx.Tx, error) {
	if s.Debug {
		logger.Debug("opening transaction")
	}
	return s.db.Beginx()
}

func (s *Service) Conn(ctx context.Context) (*sql.Conn, error) {
	return s.db.Conn(ctx)
}

func (s *Service) Stats() sql.DBStats {
	return s.db.Stats()
}

func (s *Service) Prepare(ctx context.Context, q string) (*sqlx.Stmt, error) {
	return s.db.PreparexContext(ctx, q)
}

func sqlErrMessage(err error, t string, q string, values []any) error {
	if err == nil {
		return errors.New("nil passed as argument, unable to create error")
	}
	msg := fmt.Sprintf("error [%s] running %s sql [%s] with values [%s]", err.Error(), t, strings.TrimSpace(q), valueStrings(values))
	return errors.Wrap(err, msg)
}

type logFunc func(count int, msg string, err error, output ...any)

func (s *Service) logQuery(ctx context.Context, msg string, q string, logger util.Logger, values []any) logFunc {
	if s.Debug {
		logger.Debugf("%s {\n  SQL: %s\n  Values: %s\n}", msg, strings.TrimSpace(q), valueStrings(values))
	}{{{ if .HasModule "databaseui" }}}
	if s.tracing == "" {
		return func(int, string, error, ...any) {}
	}
	t := util.TimerStart()
	return func(count int, msg string, err error, output ...any) {
		go func() {
			st, stErr := s.newStatement(ctx, q, values, t.End(), logger)
			if stErr == nil {
				st.Complete(count, msg, err, output...)
				s.addDebug(st)
			} else {
				logger.Warnf("error inserting trace history: %+v", stErr)
			}
		}()
	}{{{ else }}}
	return func(count int, msg string, err error, output ...any) {}{{{ end }}}
}

func (s *Service) newSpan(ctx context.Context, name string, q string, logger util.Logger) (time.Time, context.Context, *telemetry.Span, util.Logger) {
	if ctx == nil {
		return util.TimeCurrent(), nil, nil, logger
	}
	if s.metrics != nil {
		s.metrics.IncStmt(q, name)
	}
	nc, span, logger := telemetry.StartSpan(ctx, "database"+name, logger)
	span.Attributes(
		&telemetry.Attribute{Key: "db.statement", Value: q},
		&telemetry.Attribute{Key: "db.system", Value: s.db.DriverName()},
		&telemetry.Attribute{Key: "db.name", Value: s.DatabaseName},
		&telemetry.Attribute{Key: "db.user", Value: s.Username},
	)
	return util.TimeCurrent(), nc, span, logger
}

func (s *Service) complete(q string, op string, span *telemetry.Span, started time.Time, _ util.Logger, err error) {
	if err != nil && span != nil {
		span.OnError(err)
	}
	if span != nil {
		span.Complete()
	}
	if s.metrics != nil {
		s.metrics.CompleteStmt(q, op, started)
	}
}

func (s *Service) Close() error {
	if s.metrics != nil {
		_ = s.metrics.Close()
	}
	unregister(s)
	return s.db.Close()
}
