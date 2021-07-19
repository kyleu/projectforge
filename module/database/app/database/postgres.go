package database

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	// load postgres driver.
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
)

type PostgresParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Schema   string `json:"schema,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

func OpenPostgresDatabase(params *PostgresParams, logger *zap.SugaredLogger) (*Service, error) {
	host := params.Host
	if host == "" {
		host = "localhost"
	}
	port := params.Port
	if port == 0 {
		port = 5432
	}

	template := "postgres://%s:%s@%s:%d/%s"
	url := fmt.Sprintf(template, params.Username, params.Password, host, port, params.Database)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	var log *zap.SugaredLogger
	if params.Debug {
		log = logger
	}

	svc := NewService(params.Database, params.Schema, log, db)

	return svc, nil
}
