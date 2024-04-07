package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"

	"projectforge.dev/projectforge/app/util"
)

type Service struct {
	Key          string  `json:"key"`
	DatabaseName string  `json:"database,omitempty"`
	SchemaName   string  `json:"schema,omitempty"`
	Username     string  `json:"username,omitempty"`
	Debug        bool    `json:"debug,omitempty"`
	Type         *DBType `json:"type"`
	ReadOnly     bool    `json:"readonly,omitempty"`
	db           *sqlx.DB
}

func NewService(typ *DBType, key string, dbName string, schName string, username string, debug bool, db *sqlx.DB, logger util.Logger) (*Service, error) {
	if logger == nil {
		return nil, errors.New("logger must be provided to database service")
	}
	logger = logger.With("database", dbName, "user", username)

	ret := &Service{Key: key, DatabaseName: dbName, SchemaName: schName, Username: username, Debug: debug, Type: typ, db: db}
	err := ret.Healthcheck(dbName, db)
	if err != nil {
		return ret, errors.Wrap(err, "unable to run healthcheck")
	}
	return ret, nil
}

func (s *Service) Healthcheck(dbName string, db *sqlx.DB) error {
	q := "select 1"
	res, err := db.Query(q)
	if err != nil || res.Err() != nil {
		if err == nil {
			err = res.Err()
		}
		if strings.Contains(err.Error(), "does not exist") {
			return errors.Wrapf(err, "database does not exist")
		}
		return errors.Wrapf(err, "unable to run healthcheck [%s]", q)
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
	}
	return func(count int, msg string, err error, output ...any) {}
}

func (s *Service) Close() error {
	return s.db.Close()
}
