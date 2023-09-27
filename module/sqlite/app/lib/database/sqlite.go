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

func OpenSQLite(ctx context.Context, prefix string, logger util.Logger) (*Service, error) {
	envParams := SQLiteParamsFromEnv(util.AppKey, prefix)
	return OpenSQLiteDatabase(ctx, util.AppKey, envParams, logger)
}

func OpenDefaultSQLite(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenSQLite(ctx, "", logger)
}

func SQLiteParamsFromEnv(_ string, prefix string) *SQLiteParams {
	f := util.AppKey + ".sqlite"
	if x := util.GetEnv(prefix + "db_file"); x != "" {
		f = x
	}
	s := "public"
	if x := util.GetEnv(prefix + "db_schema"); x != "" {
		s = x
	}
	debug := false
	if x := util.GetEnv(prefix + "db_debug"); x != "" {
		debug = x != falseKey
	}
	return &SQLiteParams{File: f, Schema: s, Debug: debug}
}
