package cmd

import (
	"context"{{{ if .HasModule "readonlydb" }}}
	"strings"{{{ end }}}

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"{{{ if .HasDatabaseModule }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}
	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

func buildDefaultAppState(flags *Flags, logger util.Logger) (*app.State, error) {
	f := filesystem.NewFileSystem(flags.ConfigDir, logger)
	telemetryDisabled := util.GetEnvBool("disable_telemetry", false)
	st, err := app.NewState(flags.Debug, _buildInfo, f, !telemetryDisabled, logger)
	if err != nil {
		return nil, err
	}

	ctx, span, logger := telemetry.StartSpan(context.Background(), "app:init", logger)
	defer span.Complete(){{{ if .HasModule "migration" }}}{{{ if .HasModule "postgres" }}}

	db, err := database.OpenDefaultPostgres(ctx, logger){{{ else }}}{{{ if .HasModule "sqlite" }}}

	db, err := database.OpenDefaultSQLite(ctx, logger){{{ end }}}{{{ end }}}
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}
	st.DB = db{{{ end }}}{{{ if .HasModule "readonlydb" }}}
	roSuffix := "_readonly"
	rKey := util.AppKey + roSuffix
	if x := util.GetEnv("read_db_host", ""); x != "" {
		paramsR := database.PostgresParamsFromEnv(rKey, rKey, "read_")
		logger.Infof("using [%s:%s] for read-only database pool", paramsR.Host, paramsR.Database)
		st.DBRead, err = database.OpenPostgresDatabase(ctx, rKey, paramsR, logger)
	} else {
		paramsR := database.PostgresParamsFromEnv(rKey, util.AppKey, "")
		if strings.HasSuffix(paramsR.Database, roSuffix) {
			paramsR.Database = util.AppKey
		}
		logger.Infof("using default database as read-only database pool")
		st.DBRead, err = database.OpenPostgresDatabase(ctx, rKey, paramsR, logger)
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to open read-only database")
	}
	st.DBRead.ReadOnly = true{{{ end }}}

	svcs, err := app.NewServices(ctx, st)
	if err != nil {
		return nil, errors.Wrap(err, "error creating services")
	}
	st.Services = svcs

	return st, nil
}
