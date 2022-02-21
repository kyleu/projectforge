package database

import (
	"context"
	"os"

	"go.uber.org/zap"

	"{{{ .Package }}}/app/util"
)

type SQLiteParams struct {
	File   string `json:"file"`
	Schema string `json:"schema,omitempty"`
	Debug  bool   `json:"debug,omitempty"`
}

func OpenSQLite(ctx context.Context, prefix string, logger *zap.SugaredLogger) (*Service, error) {
	envParams := SQLiteParamsFromEnv(util.AppKey, util.AppKey, prefix)
	return OpenSQLiteDatabase(ctx, util.AppKey, envParams, logger)
}

func OpenDefaultSQLite(ctx context.Context, logger *zap.SugaredLogger) (*Service, error) {
	return OpenSQLite(ctx, "", logger)
}

func SQLiteParamsFromEnv(key string, defaultUser string, prefix string) *SQLiteParams {
	f := util.AppKey + ".sqlite"
	if x := os.Getenv(prefix + "DB_FILE"); x != "" {
		f = x
	}
	s := "public"
	if x := os.Getenv(prefix + "DB_SCHEMA"); x != "" {
		s = x
	}
	debug := false
	if x := os.Getenv(prefix + "DB_DEBUG"); x != "" {
		debug = x != falseKey
	}
	return &SQLiteParams{File: f, Schema: s, Debug: debug}
}
