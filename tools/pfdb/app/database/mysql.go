package database

import (
	"context"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

var typeMySQL = &DBType{Key: "mysql", Title: "MySQL", Quote: "`", Placeholder: "?", SupportsReturning: false}

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
	pw := ""
	if x := util.GetEnv(prefix + "db_password"); x != "" {
		pw = x
	}
	d := key
	if x := util.GetEnv(prefix + "db_database"); x != "" {
		d = x
	}
	s := ""
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
	return &MySQLParams{Host: h, Port: p, Username: u, Password: pw, Database: d, Schema: s, MaxConns: mc, Debug: debug}
}

func OpenMySQLDatabase(ctx context.Context, key string, params *MySQLParams, logger util.Logger) (*Service, error) {
	host := params.Host
	if host == "" {
		host = localhost
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

	return NewService(typeMySQL, key, params.Database, params.Schema, params.Username, params.Debug, db, logger)
}

func OpenDefaultMySQL(logger util.Logger) (*Service, error) {
	params := MySQLParamsFromEnv(util.AppKey, util.AppKey, "")
	return OpenMySQLDatabase(context.Background(), util.AppKey, params, logger)
}
