package database

import (
	"context"
	"strings"

	"{{{ .Package }}}/app/util"
)

type SQLiteParams struct {
	File     string `json:"file"`
	User     string `json:"user,omitzero"`
	Password string `json:"password,omitzero"`
	Schema   string `json:"schema,omitzero"`
	Debug    bool   `json:"debug,omitzero"`
}

func OpenSQLite(ctx context.Context, key string, prefix string, logger util.Logger) (*Service, error) {
	envParams := SQLiteParamsFromEnv(key, prefix)
	if key == "" {
		key = util.AppKey
	}
	return OpenSQLiteDatabase(ctx, key, envParams, logger)
}

func OpenDefaultSQLite(ctx context.Context, logger util.Logger) (*Service, error) {
	return OpenSQLite(ctx, "", "", logger)
}

func SQLiteParamsFromEnv(key string, prefix string) *SQLiteParams {
	if key == "" {
		key = util.StringPath(util.ConfigDir, util.AppKey)
	}
	if !strings.HasSuffix(key, ".sqlite") {
		key += ".sqlite"
	}
	if x := util.GetEnv(prefix + cfgFile); x != "" {
		key = x
	}
	var u string
	if x := util.GetEnv(prefix + cfgUser); x != "" {
		u = x
	}
	var p string
	if x := util.GetEnv(prefix + cfgPassword); x != "" {
		p = x
	}
	s := "public"
	if x := util.GetEnv(prefix + cfgSchema); x != "" {
		s = x
	}
	debug := false
	if x := util.GetEnv(prefix + cfgDebug); x != "" {
		debug = x != util.BoolFalse
	}
	return &SQLiteParams{File: key, Schema: s, User: u, Password: p, Debug: debug}
}

func SQLiteParamsFromMap(m util.ValueMap) *SQLiteParams {
	return &SQLiteParams{
		File:     m.GetStringOpt("file"),
		User:     m.GetStringOpt("user"),
		Password: m.GetStringOpt("password"),
		Schema:   m.GetStringOpt("schema"),
		Debug:    m.GetBoolOpt("debug"),
	}
}
