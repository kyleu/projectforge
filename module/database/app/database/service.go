package database

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	DatabaseName string
	SchemaName   string
	logger       *zap.SugaredLogger
	db           *sqlx.DB
}

func NewService(dbName string, schName string, logger *zap.SugaredLogger, db *sqlx.DB) *Service {
	return &Service{DatabaseName: dbName, SchemaName: schName, logger: logger, db: db}
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
