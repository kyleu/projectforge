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

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

const defaultSchema = "public"

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
	pw := defaultUser
	if x := os.Getenv(prefix + "DB_PASSWORD"); x != "" {
		pw = x
	}
	d := key
	if x := os.Getenv(prefix + "DB_DATABASE"); x != "" {
		d = x
	}
	s := defaultSchema
	if x := os.Getenv(prefix + "DB_SCHEMA"); x != "" {
		s = x
	}
	mc := 16
	if x := os.Getenv(prefix + "DB_MAX_CONNECTIONS"); x != "" {
		mc, _ = strconv.Atoi(x)
	}
	debug := false
	if x := os.Getenv(prefix + "DB_DEBUG"); x != "" {
		debug = x != falseKey
	}
	return &PostgresParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenPostgres(ctx context.Context, prefix string, logger *zap.SugaredLogger) (*Service, error) {
	envParams := PostgresParamsFromEnv(util.AppKey, util.AppKey, prefix)
	if os.Getenv("DB_SSL") == "true" {
		serviceParams, err := PostgresParamsFromService()
		if err != nil {
			return nil, err
		}
		return OpenPostgresDatabaseSSL(ctx, util.AppKey, envParams, serviceParams, logger)
	}
	return OpenPostgresDatabase(ctx, util.AppKey, envParams, logger)
}

func OpenDefaultPostgres(ctx context.Context, logger *zap.SugaredLogger) (*Service, error) {
	return OpenPostgres(ctx, "", logger)
}

func OpenPostgresDatabase(ctx context.Context, key string, params *PostgresParams, logger *zap.SugaredLogger) (*Service, error) {
	_, span := telemetry.StartSpan(ctx, "database", "open")
	defer span.Complete()
	host := params.Host
	if host == "" {
		host = localhost
	}
	port := params.Port
	if port == 0 {
		port = 5432
	}
	sch := defaultSchema
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

	var log *zap.SugaredLogger
	if params.Debug {
		log = logger.With("svc", "database", "db", key)
	}

	return NewService(typePostgres, key, params.Database, params.Schema, params.Username, db, log)
}

func OpenPostgresDatabaseSSL(ctx context.Context, key string, ep *PostgresParams, sp *PostgresServiceParams, logger *zap.SugaredLogger) (*Service, error) {
	_, span := telemetry.StartSpan(ctx, "database", "openssl")
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
		schema = defaultSchema
	}

	const template = "postgres://%s:%d/%s?search_path=%s&application_name=%s&user=%s&sslmode=%s&sslcert=%s&sslrootcert=%s&sslkey=%s"
	url := fmt.Sprintf(template, sp.Host, sp.Port, dbname, schema, key, sp.Username, sp.SSLMode, sp.SSLCert, sp.SSLRootCert, sp.SSLKey)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(ep.MaxConns)
	db.SetMaxIdleConns(0)

	var log *zap.SugaredLogger
	if ep.Debug {
		log = logger.With("svc", "database", "db", key)
	}

	return NewService(typePostgres, key, dbname, ep.Schema, sp.Username, db, log)
}
