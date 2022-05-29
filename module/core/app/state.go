package app

import (
	"context"
	"fmt"
	"time"
{{{ if .HasModule "oauth" }}}
	"{{{ .Package }}}/app/lib/auth"{{{ end }}}{{{ if .HasDatabaseModule }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}{{{ if .HasModule "filesystem" }}}
	"{{{ .Package }}}/app/lib/filesystem"{{{ end }}}
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/theme"
	"{{{ .Package }}}/app/util"
)

type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (b *BuildInfo) String() string {
	if b.Date == "unknown" {
	} else {
		d, _ := util.TimeFromJS(b.Date)
		return fmt.Sprintf("%s (%s)", b.Version, util.TimeToYMD(d))
	}
	return b.Version
}

type State struct {
	Debug     bool
	BuildInfo *BuildInfo{{{ if .HasModule "filesystem" }}}
	Files     filesystem.FileLoader{{{ end }}}{{{ if .HasModule "oauth" }}}
	Auth      *auth.Service{{{ end }}}{{{ if .HasModule "migration" }}}
	DB        *database.Service{{{ end }}}{{{ if .HasModule "readonlydb" }}}
	DBRead    *database.Service{{{ end }}}
	Themes    *theme.Service
	Logger    util.Logger
	Services  *Services
	Started   time.Time
}

func (s State) Close(ctx context.Context, logger util.Logger) error {
	{{{ if .HasModule "migration" }}}if err := s.DB.Close(); err != nil {
		logger.Errorf("error closing database: %+v", err)
	}
	{{{ end }}}{{{ if .HasModule "readonlydb" }}}if err := s.DBRead.Close(); err != nil {
		logger.Errorf("error closing read-only database: %+v", err)
	}
	{{{ end }}}return s.Services.Close(ctx, logger)
}

func NewState(debug bool, bi *BuildInfo{{{ if .HasModule "filesystem" }}}, f filesystem.FileLoader{{{ end }}}, enableTelemetry bool, logger util.Logger) (*State, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	time.Local = loc

	_ = telemetry.InitializeIfNeeded(enableTelemetry, bi.Version, logger){{{ if .HasModule "oauth" }}}
	as := auth.NewService("", logger){{{ end }}}
	ts := theme.NewService({{{ if .HasModule "filesystem" }}}f{{{ end }}})

	return &State{
		Debug:     debug,
		BuildInfo: bi{{{ if .HasModule "filesystem" }}},
		Files:     f{{{ end }}}{{{ if .HasModule "oauth" }}},
		Auth:      as{{{ end }}},
		Themes:    ts,
		Logger:    logger,
		Started:   time.Now(),
	}, nil
}
