package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	// load postgres driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/telemetry"
)

type PostgresParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Schema   string `json:"schema,omitempty"`
	MaxConns int    `json:"maxConns,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

func PostgresParamsFromEnv(key string, defaultUser string, prefix string) *PostgresParams {
	h := "localhost"
	if x := os.Getenv(prefix + "DB_HOST"); x != "" {
		h = x
	}
	p := 0
	if x := os.Getenv(prefix + "DB_PORT"); x != "" {
		p, _ = strconv.Atoi(x)
	}
	u := defaultUser
	if x := os.Getenv(prefix + "DB_USER"); x != "" {
		u = x
	}
	pw := ""
	if x := os.Getenv(prefix + "DB_PASSWORD"); x != "" {
		pw = x
	}
	d := key
	if x := os.Getenv(prefix + "DB_DATABASE"); x != "" {
		d = x
	}
	s := "public"
	if x := os.Getenv(prefix + "DB_SCHEMA"); x != "" {
		s = x
	}
	mc := 16
	if x := os.Getenv(prefix + "DB_MAX_CONNECTIONS"); x != "" {
		mc, _ = strconv.Atoi(x)
	}
	debug := false
	if x := os.Getenv(prefix + "DB_DEBUG"); x != "" {
		debug = x != "false"
	}
	return &PostgresParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenPostgresDatabase(ctx context.Context, params *PostgresParams, tel *telemetry.Service, logger *zap.SugaredLogger) (*Service, error) {
	ctx, span := tel.StartSpan(ctx, "database", "open")
	defer span.End()
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

	svc := NewService(params.Database, params.Schema, params.Username, tel, log, db)

	return svc, nil
}