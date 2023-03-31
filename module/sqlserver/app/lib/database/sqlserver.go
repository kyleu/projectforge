package database

import (
	"context"
	"fmt"
	"strconv"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

const defaultSQLServerSchema = "dbo"

var typeSQLServer = &DBType{Key: "sqlserver", Title: "SQL Server", Quote: `"`, Placeholder: "@", SupportsReturning: true}

type SQLServerParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Schema   string `json:"schema,omitempty"`
	MaxConns int    `json:"maxConns,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

func SQLServerParamsFromEnv(key string, defaultUser string, prefix string) *SQLServerParams {
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
	s := defaultSQLServerSchema
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
	return &SQLServerParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenSQLServer(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	if key == "" {
		key = util.AppKey
	}
	envParams := SQLServerParamsFromEnv(key, util.AppKey, prefix)
	return OpenSQLServerDatabase(ctx, key, envParams, logger)
}

func OpenDefaultSQLServer(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenSQLServer(ctx, "", "", logger)
}

func OpenSQLServerDatabase(ctx context.Context, key string, params *SQLServerParams, logger util.Logger) (*Service, error) {
	_, span, logger := telemetry.StartSpan(ctx, "database:open", logger)
	defer span.Complete()

	host := params.Host
	if host == "" {
		host = localhost
	}
	port := params.Port
	if port == 0 {
		port = 1433
	}

	const template = "sqlserver://%s:%s@%s:%d?database=%s&app%%20name=%s"
	url := fmt.Sprintf(template, params.Username, params.Password, host, port, params.Database, util.AppKey)

	db, err := sqlx.Open("sqlserver", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(params.MaxConns)
	db.SetMaxIdleConns(0)

	logger = logger.With("svc", "database", "db", key)

	return NewService(typeSQLServer, key, params.Database, params.Schema, params.Username, params.Debug, db, logger)
}

func UUIDFromGUID(x *any) *uuid.UUID {
	if x == nil {
		return nil
	}
	switch t := (*x).(type) {
	case []byte:
		ret := mssql.UniqueIdentifier{}
		_ = ret.Scan(t)
		return util.UUIDFromString(ret.String())
	case mssql.UniqueIdentifier:
		return util.UUIDFromString(t.String())
	default:
		return nil
	}
}
