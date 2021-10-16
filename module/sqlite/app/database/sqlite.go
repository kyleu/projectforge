package database

import (
	"os"
	"strconv"
)

type SQLiteParams struct {
	File   string `json:"file"`
	Schema string `json:"schema,omitempty"`
	Debug  bool   `json:"debug,omitempty"`
}

func SQLiteParamsFromEnv(key string, defaultUser string, prefix string) *SQLiteParams {
	f := ""
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
