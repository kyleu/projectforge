//go:build aix || android || dragonfly || freebsd || js || illumos || (linux && ppc64) || (linux && mips) || (linux && mipsle) || (linux && mips64) || (linux && mips64le) || (linux && ppc64) || (linux && ppc64le) || (linux && s390x) || (linux && loong64) || netbsd || openbsd || plan9 || solaris || (windows && 386) || (windows && arm) || (windows && arm64)

package database

import (
	"context"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

const SQLiteEnabled = false

func OpenSQLiteDatabase(ctx context.Context, key string, params *SQLiteParams, logger util.Logger) (*Service, error) {
	return nil, errors.New("SQLite is not enabled in this build")
}
