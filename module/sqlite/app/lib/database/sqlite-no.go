// +build aix android dragonfly freebsd js illumos linux,ppc64 linux,mips linux,mipsle linux,mips64 linux,mips64le linux,ppc64 linux,ppc64le linux,s390x netbsd openbsd plan9 solaris windows,arm windows,arm64

package database

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const SQLiteEnabled = false

func OpenSQLiteDatabase(ctx context.Context, key string, params *SQLiteParams, logger *zap.SugaredLogger) (*Service, error) {
	return nil, errors.New("SQLite is not enabled in this build")
}
