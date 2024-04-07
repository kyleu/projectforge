package database

import (
	"context"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

const defaultPostgreSQLSchema = "public"

var typePostgres = &DBType{Key: "postgres", Title: "PostgreSQL", Quote: `"`, Placeholder: "$", SupportsReturning: true}

type PostgresParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Schema   string `json:"schema,omitempty"`
	MaxConns int    `json:"maxConns,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

func PostgresParamsFromEnv(key string, defaultUser string, prefix string) *PostgresParams {
	h := localhost
	if x := util.GetEnv(prefix + "db_host"); x != "" {
		h = x
	}
	p := 0
	if x := util.GetEnv(prefix + "db_port"); x != "" {
		px, _ := strconv.ParseInt(x, 10, 32)
		p = int(px)
	}
	u := defaultUser
	if x := util.GetEnv(prefix + "db_user"); x != "" {
		u = x
	}
	pw := defaultUser
	if x := util.GetEnv(prefix + "db_password"); x != "" {
		pw = x
	}
	d := key
	if x := util.GetEnv(prefix + "db_database"); x != "" {
		d = x
	}
	s := defaultPostgreSQLSchema
	if x := util.GetEnv(prefix + "db_schema"); x != "" {
		s = x
	}
	mc := 16
	if x := util.GetEnv(prefix + "db_max_connections"); x != "" {
		mcx, _ := strconv.ParseInt(x, 10, 32)
		mc = int(mcx)
	}
	debug := false
	if x := util.GetEnv(prefix + "db_debug"); x != "" {
		debug = x != util.BoolFalse
	}
	return &PostgresParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenPostgres(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	envParams := PostgresParamsFromEnv(key, key, prefix)
	return OpenPostgresDatabase(ctx, key, envParams, logger)
}

func OpenDefaultPostgres(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenPostgres(ctx, "", "", logger)
}

func OpenPostgresDatabase(ctx context.Context, key string, params *PostgresParams, logger util.Logger) (*Service, error) {
	host := params.Host
	if host == "" {
		host = localhost
	}
	port := params.Port
	if port == 0 {
		port = 5432
	}
	sch := defaultPostgreSQLSchema
	if params.Schema != "" {
		sch = params.Schema
	}

	const template = "postgres://%s:%s@%s:%d/%s?search_path=%s&application_name=%s"
	url := fmt.Sprintf(template, params.Username, params.Password, host, port, params.Database, sch, key)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(params.MaxConns)
	db.SetMaxIdleConns(0)

	logger = logger.With("svc", "database", "db", key)

	return NewService(typePostgres, key, params.Database, params.Schema, params.Username, params.Debug, db, logger)
}
