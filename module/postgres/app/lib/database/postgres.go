package database

import (
	"context"
	"fmt"
	"strconv"

	// load postgres driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
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
		p, _ = strconv.Atoi(x)
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
		mc, _ = strconv.Atoi(x)
	}
	debug := false
	if x := util.GetEnv(prefix + "db_debug"); x != "" {
		debug = x != falseKey
	}
	return &PostgresParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenPostgres(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	if key == "" {
		key = util.AppKey
	}
	envParams := PostgresParamsFromEnv(key, util.AppKey, prefix)
	if util.GetEnv("db_ssl") == util.BoolTrue {
		serviceParams, err := PostgresParamsFromService()
		if err != nil {
			return nil, err
		}
		return OpenPostgresDatabaseSSL(ctx, key, envParams, serviceParams, logger)
	}
	return OpenPostgresDatabase(ctx, key, envParams, logger)
}

func OpenDefaultPostgres(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenPostgres(ctx, "", "", logger)
}

func OpenPostgresDatabase(ctx context.Context, key string, params *PostgresParams, logger util.Logger) (*Service, error) {
	_, span, logger := telemetry.StartSpan(ctx, "database:open", logger)
	defer span.Complete()
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

func OpenPostgresDatabaseSSL(ctx context.Context, key string, ep *PostgresParams, sp *PostgresServiceParams, logger util.Logger) (*Service, error) {
	_, span, logger := telemetry.StartSpan(ctx, "database:openssl", logger)
	defer span.Complete()

	dbname := sp.Database
	if dbname == "" {
		dbname = ep.Database
	}
	schema := sp.Schema
	if schema == "" {
		schema = ep.Schema
	}
	if schema == "" {
		schema = defaultPostgreSQLSchema
	}

	const template = "postgres://%s:%d/%s?search_path=%s&application_name=%s&user=%s&sslmode=%s&sslcert=%s&sslrootcert=%s&sslkey=%s"
	url := fmt.Sprintf(template, sp.Host, sp.Port, dbname, schema, key, sp.Username, sp.SSLMode, sp.SSLCert, sp.SSLRootCert, sp.SSLKey)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(ep.MaxConns)
	db.SetMaxIdleConns(0)

	logger = logger.With("svc", "database", "db", key)

	return NewService(typePostgres, key, dbname, ep.Schema, sp.Username, ep.Debug, db, logger)
}
