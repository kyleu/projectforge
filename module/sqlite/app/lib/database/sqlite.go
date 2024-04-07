package database

import (
	"context"

	"{{{ .Package }}}/app/util"
)

type SQLiteParams struct {
	File   string `json:"file"`
	Schema string `json:"schema,omitempty"`
	Debug  bool   `json:"debug,omitempty"`
}

func OpenSQLite(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	if key == "" {
		key = util.AppKey
	}
	envParams := SQLiteParamsFromEnv(key, prefix)
	return OpenSQLiteDatabase(ctx, key, envParams, logger)
}

func OpenDefaultSQLite(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenSQLite(ctx, "", "", logger)
}

func SQLiteParamsFromEnv(key string, prefix string) *SQLiteParams {
	f := key + ".sqlite"
	if x := util.GetEnv(prefix + "db_file"); x != "" {
		f = x
	}
	s := "public"
	if x := util.GetEnv(prefix + "db_schema"); x != "" {
		s = x
	}
	debug := false
	if x := util.GetEnv(prefix + "db_debug"); x != "" {
		debug = x != util.BoolFalse
	}
	return &SQLiteParams{File: f, Schema: s, Debug: debug}
}
