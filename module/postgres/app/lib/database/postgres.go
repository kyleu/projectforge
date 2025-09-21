package database

import (
	"context"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

const defaultPostgreSQLSchema = "public"

var TypePostgres = &DBType{Key: "postgres", Title: "PostgreSQL", Quote: `"`, Placeholder: "$", SupportsReturning: true}

type PostgresParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port,omitzero"`
	Username string `json:"username"`
	Password string `json:"password,omitzero"`
	Database string `json:"database,omitzero"`
	Schema   string `json:"schema,omitzero"`
	MaxConns int    `json:"maxConns,omitzero"`
	Debug    bool   `json:"debug,omitzero"`
}

func PostgresParamsFromEnv(key string, defaultUser string, prefix string) *PostgresParams {
	h := localhost
	if x := util.GetEnv(prefix + cfgHost); x != "" {
		h = x
	}
	var p int
	if x := util.GetEnv(prefix + cfgPort); x != "" {
		px, _ := strconv.ParseInt(x, 10, 32)
		p = int(px)
	}
	u := defaultUser
	if x := util.GetEnv(prefix + cfgUser); x != "" {
		u = x
	}
	pw := defaultUser
	if x := util.GetEnv(prefix + cfgPassword); x != "" {
		pw = x
	}
	d := key
	if x := util.GetEnv(prefix + cfgDatabase); x != "" {
		d = x
	}
	s := defaultPostgreSQLSchema
	if x := util.GetEnv(prefix + cfgSchema); x != "" {
		s = x
	}
	mc := 16
	if x := util.GetEnv(prefix + cfgMaxConns); x != "" {
		mcx, _ := strconv.ParseInt(x, 10, 32)
		mc = int(mcx)
	}
	debug := false
	if x := util.GetEnv(prefix + cfgDebug); x != "" {
		debug = x != util.BoolFalse
	}
	return &PostgresParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func PostgresParamsFromMap(m util.ValueMap) *PostgresParams {
	return &PostgresParams{
		Host:     m.GetStringOpt("host"),
		Port:     m.GetIntOpt("port"),
		Username: m.GetStringOpt("username"),
		Password: m.GetStringOpt("password"),
		Database: m.GetStringOpt("database"),
		Schema:   m.GetStringOpt("schema"),
		MaxConns: m.GetIntOpt("maxConns"),
		Debug:    m.GetBoolOpt("debug"),
	}
}

func OpenPostgres(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	if key == "" {
		key = util.AppKey
	}
	envParams := PostgresParamsFromEnv(key, util.AppKey, prefix)
	return OpenPostgresDatabase(ctx, key, envParams, logger)
}

func OpenDefaultPostgres(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenPostgres(ctx, "", "", logger)
}

func OpenPostgresDatabase(ctx context.Context, key string, params *PostgresParams, logger util.Logger) (*Service, error) {
	_, span, logger := telemetry.StartSpan(ctx, "database:open", logger)
	defer span.Complete()
	host := util.OrDefault(params.Host, localhost)
	port := util.OrDefault(params.Port, 5432)
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
	stringRep := fmt.Sprintf("%s@%s:%d", params.Database, host, port)

	return NewService(TypePostgres, key, params.Database, params.Schema, params.Username, params.Debug, db, stringRep, logger)
}
