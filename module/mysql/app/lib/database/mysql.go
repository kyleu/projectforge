package database

import (
	"context"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

var TypeMySQL = &DBType{Key: "mysql", Title: "MySQL", Quote: "`", Placeholder: "?", SupportsReturning: false}

type MySQLParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Schema   string `json:"schema,omitempty"`
	MaxConns int    `json:"maxConns,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

func MySQLParamsFromEnv(key string, defaultUser string, prefix string) *MySQLParams {
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
	pw := ""
	if x := util.GetEnv(prefix + cfgPassword); x != "" {
		pw = x
	}
	d := key
	if x := util.GetEnv(prefix + cfgDatabase); x != "" {
		d = x
	}
	s := ""
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
	return &MySQLParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func MySQLParamsFromMap(m util.ValueMap) *MySQLParams {
	return &MySQLParams{
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

func OpenMySQL(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	if key == "" {
		key = util.AppKey
	}
	envParams := MySQLParamsFromEnv(key, key, prefix)
	return OpenMySQLDatabase(ctx, key, envParams, logger)
}

func OpenDefaultMySQL(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenMySQL(ctx, "", "", logger)
}

func OpenMySQLDatabase(ctx context.Context, key string, params *MySQLParams, logger util.Logger) (*Service, error) {
	_, span, logger := telemetry.StartSpan(ctx, "database:open", logger)
	defer span.Complete()
	host := util.OrDefault(params.Host, localhost)
	port := util.OrDefault(params.Port, 3306)

	const template = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
	url := fmt.Sprintf(template, params.Username, params.Password, host, port, params.Database)

	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(params.MaxConns)
	db.SetMaxIdleConns(0)

	logger = logger.With("svc", "database", "db", key)
	stringRep := fmt.Sprintf("%s@%s:%d", params.Database, host, port)

	return NewService(TypeMySQL, key, params.Database, params.Schema, params.Username, params.Debug, db, stringRep, logger)
}
