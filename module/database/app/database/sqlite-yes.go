// +build darwin !android,linux,386 !android,linux,amd64 !android,linux,arm !android,linux,arm64 windows,386 windows,amd64

package database

import (
	"go.uber.org/zap"

	"github.com/pkg/errors"

	// load sqlite driver.
	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
)

const SQLiteEnabled = true

func OpenSQLiteDatabase(params *SQLiteParams, logger *zap.SugaredLogger) (*Service, error) {
	if params.File == "" {
		return nil, errors.New("need filename for SQLite database")
	}
	db, err := sqlx.Open("sqlite", params.File)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	var log *zap.SugaredLogger
	if params.Debug {
		log = logger
	}

	svc := NewService(params.File, params.Schema, log, db)

	return svc, nil
}
