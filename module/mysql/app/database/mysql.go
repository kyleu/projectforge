package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	// load MySQL driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/telemetry"
	"{{{ .Package }}}/app/util"
)

var typeMySQL = &DBType{Key: "mysql", Title: "MySQL", Quote: "`"}

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
	return &MySQLParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenMySQLDatabase(ctx context.Context, key string, params *MySQLParams, logger *zap.SugaredLogger) (*Service, error) {
	ctx, span := telemetry.StartSpan(ctx, "database", "open")
	defer span.End()
	host := params.Host
	if host == "" {
		host = "localhost"
	}
	port := params.Port
	if port == 0 {
		port = 3306
	}

	const template = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
	url := fmt.Sprintf(template, params.Username, params.Password, host, port, params.Database)

	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	db.SetMaxOpenConns(params.MaxConns)
	db.SetMaxIdleConns(0)

	var log *zap.SugaredLogger
	if params.Debug {
		log = logger.With("svc", "database", "db", key)
	}

	svc := NewService(typeMySQL, key, params.Database, params.Schema, params.Username, db, log)
	return svc, nil
}

func OpenDefaultMySQL(logger *zap.SugaredLogger) (*Service, error) {
	params := MySQLParamsFromEnv(util.AppKey, util.AppKey, "")
	return OpenMySQLDatabase(context.Background(), util.AppKey, params, logger)
}
