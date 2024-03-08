package cmd

import (
	"context"{{{ if .HasModule "readonlydb" }}}
	"strings"{{{ end }}}

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"{{{ if .HasModules "migration" "readonlydb" }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}{{{ if .HasModule "filesystem" }}}
	"{{{ .Package }}}/app/lib/filesystem"{{{ end }}}
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

func buildDefaultAppState(flags *Flags, logger util.Logger) (*app.State, error) {
	{{{ if .HasModule "filesystem" }}}fs, err := filesystem.NewFileSystem(flags.ConfigDir, false, "")
	if err != nil {
		return nil, err
	}

	{{{ end }}}telemetryDisabled := util.GetEnvBool("disable_telemetry", false)
	st, err := app.NewState(flags.Debug, _buildInfo{{{ if .HasModule "filesystem" }}}, fs{{{ end }}}, !telemetryDisabled, flags.Port, logger)
	if err != nil {
		return nil, err
	}

	ctx, span, logger := telemetry.StartSpan(context.Background(), "app:init", logger)
	defer span.Complete()
	t := util.TimerStart(){{{ if .HasModule "migration" }}}{{{ if .PostgreSQL }}}

	db, err := database.OpenDefaultPostgres(ctx, logger){{{ else }}}{{{ if .SQLite }}}

	db, err := database.OpenDefaultSQLite(ctx, logger){{{ else }}}{{{ if .SQLServer }}}

	db, err := database.OpenDefaultSQLServer(ctx, logger){{{ end }}}{{{ end }}}{{{ end }}}
	if err != nil {
		logger.Errorf("unable to open default database: %+v", err)
	}
	st.DB = db{{{ end }}}{{{ if .HasModule "readonlydb" }}}
	roSuffix := "_readonly"
	rKey := util.AppKey + roSuffix
	if x := util.GetEnv("read_db_host", ""); x != "" {
		paramsR := database.PostgresParamsFromEnv(rKey, rKey, "read_")
		logger.Infof("using [%s:%s] for read-only database pool", paramsR.Host, paramsR.Database){{{ if .PostgreSQL }}}
		st.DBRead, err = database.OpenPostgresDatabase(ctx, rKey, paramsR, logger){{{ else }}}{{{ if .SQLite }}}
		st.DBRead, err = database.OpenSQLiteDatabase(ctx, rKey, paramsR, logger){{{ else }}}{{{ if .SQLServer }}}
		st.DBRead, err = database.OpenSQLServerDatabase(ctx, rKey, paramsR, logger){{{ end }}}{{{ end }}}{{{ end }}}
	} else {
		paramsR := database.PostgresParamsFromEnv(rKey, util.AppKey, "")
		if strings.HasSuffix(paramsR.Database, roSuffix) {
			paramsR.Database = util.AppKey
		}
		logger.Infof("using default database as read-only database pool"){{{ if .PostgreSQL }}}
		st.DBRead, err = database.OpenPostgresDatabase(ctx, rKey, paramsR, logger){{{ else }}}{{{ if .SQLite }}}
		st.DBRead, err = database.OpenSQLiteDatabase(ctx, rKey, paramsR, logger){{{ else }}}{{{ if .SQLServer }}}
		st.DBRead, err = database.OpenSQLServerDatabase(ctx, rKey, paramsR, logger){{{ end }}}{{{ end }}}{{{ end }}}
	}
	if err != nil {
		logger.Errorf("unable to open default read-only database: %v", err)
	}
	st.DBRead.ReadOnly = true{{{ end }}}
	svcs, err := app.NewServices(ctx, st, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating services")
	}
	logger.Debugf("created app state in [%s]", util.MicrosToMillis(t.End()))
	st.Services = svcs

	return st, nil
}
