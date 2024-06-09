//go:build darwin || (!android && linux && 386) || (!android && linux && amd64) || (!android && linux && arm) || (!android && linux && arm64) || (!android && linux && riscv64) || (windows && amd64)

package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

const SQLiteEnabled = true

var TypeSQLite = &DBType{Key: "sqlite", Title: "SQLite", Quote: `"`, Placeholder: "$", SupportsReturning: true}

func OpenSQLiteDatabase(ctx context.Context, key string, params *SQLiteParams, logger util.Logger) (*Service, error) {
	_, span, logger := telemetry.StartSpan(ctx, "database:open", logger)
	defer span.Complete()
	if params.File == "" {
		return nil, errors.New("need filename for SQLite database")
	}
	if pth, fn := util.StringSplitLast(params.File, '/', true); fn != "" {
		_ = os.MkdirAll(pth, os.FileMode(filesystem.DirectoryMode))
	}
	conn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=trusted_schema(0)", params.File)
	if params.User != "" && params.Password != "" {
		conn += fmt.Sprintf("&_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha1", params.User, params.Password)
	}
	db, err := sqlx.Open("sqlite", conn)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	logger = logger.With("svc", "database", "db", key)
	return NewService(TypeSQLite, key, key, params.Schema, "sqlite", params.Debug, db, params.File, logger)
}
