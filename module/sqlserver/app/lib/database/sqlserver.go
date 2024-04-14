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

var TypeSQLServer = &DBType{Key: "sqlserver", Title: "SQL Server", Quote: `"`, Placeholder: "@", SupportsReturning: true}

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
	if x := util.GetEnv(prefix + cfgHost); x != "" {
		h = x
	}
	p := 0
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
	s := defaultSQLServerSchema
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
	return &SQLServerParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func SQLServerParamsFromMap(m util.ValueMap) *SQLServerParams {
	return &SQLServerParams{
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

func OpenSQLServer(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	if key == "" {
		key = util.AppKey
	}
	envParams := SQLServerParamsFromEnv(key, key, prefix)
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

	var url string
	switch params.Username {
	case "native", "local", "":
		const template = "sqlserver://%s:%d?database=%s&app%%20name=%s"
		url = fmt.Sprintf(template, host, port, params.Database, util.AppKey)
	default:
		const template = "sqlserver://%s:%s@%s:%d?database=%s&app%%20name=%s"
		url = fmt.Sprintf(template, params.Username, params.Password, host, port, params.Database, util.AppKey)
	}

	db, err := sqlx.Open("sqlserver", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(params.MaxConns)
	db.SetMaxIdleConns(0)

	logger = logger.With("svc", "database", "db", key)
	stringRep := fmt.Sprintf("%s@%s:%d", params.Database, host, port)

	return NewService(TypeSQLServer, key, params.Database, params.Schema, params.Username, params.Debug, db, stringRep, logger)
}

func UUIDFromGUID(x any) *uuid.UUID {
	if x == nil {
		return nil
	}
	switch t := (x).(type) {
	case *mssql.UniqueIdentifier:
		return util.UUIDFromString(t.String())
	case mssql.UniqueIdentifier:
		return util.UUIDFromString(t.String())
	case uuid.UUID:
		return &t
	case *uuid.UUID:
		return t
	case string:
		return util.UUIDFromString(t)
	case []byte:
		ret := mssql.UniqueIdentifier{}
		_ = ret.Scan(t)
		return util.UUIDFromString(ret.String())
	default:
		return nil
	}
}
