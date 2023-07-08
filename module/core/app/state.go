package app

import (
	"context"
	"fmt"
	"sync"
	"time"
{{{ if .HasModule "oauth" }}}
	"{{{ .Package }}}/app/lib/auth"{{{ end }}}{{{ if .HasDatabaseModule }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}{{{ if .HasModule "filesystem" }}}
	"{{{ .Package }}}/app/lib/filesystem"{{{ end }}}{{{ if .HasModule "graphql" }}}
	"{{{ .Package }}}/app/lib/graphql"{{{ end }}}
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/theme"
	"{{{ .Package }}}/app/util"
)

var once sync.Once

type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (b *BuildInfo) String() string {
	if b.Date == util.KeyUnknown {
		return b.Version
	}
	d, _ := util.TimeFromJS(b.Date)
	return fmt.Sprintf("%s (%s)", b.Version, util.TimeToYMD(d))
}

type State struct {
	Debug     bool
	BuildInfo *BuildInfo{{{ if .HasModule "filesystem" }}}
	Files     filesystem.FileLoader{{{ end }}}{{{ if .HasModule "oauth" }}}
	Auth      *auth.Service{{{ end }}}{{{ if .HasModule "migration" }}}
	DB        *database.Service{{{ end }}}{{{ if .HasModule "readonlydb" }}}
	DBRead    *database.Service{{{ end }}}{{{ if .HasModule "graphql" }}}
	GraphQL   *graphql.Service{{{ end }}}
	Themes    *theme.Service
	Services  *Services
	Started   time.Time
}

func NewState(debug bool, bi *BuildInfo{{{ if .HasModule "filesystem" }}}, f filesystem.FileLoader{{{ end }}}, enableTelemetry bool, {{{ if .HasModule "oauth" }}}port{{{ else }}}_{{{ end }}} uint16, logger util.Logger) (*State, error) {
	var loadLocationError error
	once.Do(func() {
		loc, err := time.LoadLocation("UTC")
		if err != nil {
			loadLocationError = err
			return
		}
		time.Local = loc
	})
	if loadLocationError != nil {
		return nil, loadLocationError
	}

	_ = telemetry.InitializeIfNeeded(enableTelemetry, bi.Version, logger){{{ if .HasModule "oauth" }}}
	as := auth.NewService("", port, logger){{{ end }}}
	ts := theme.NewService({{{ if .HasModule "filesystem" }}}f{{{ end }}}){{{ if .HasModule "graphql" }}}
	gqls := graphql.NewService(){{{ end }}}

	return &State{
		Debug:     debug,
		BuildInfo: bi{{{ if .HasModule "filesystem" }}},
		Files:     f{{{ end }}}{{{ if .HasModule "oauth" }}},
		Auth:      as{{{ end }}}{{{ if .HasModule "graphql" }}},
		GraphQL:   gqls{{{ end }}},
		Themes:    ts,
		Started:   util.TimeCurrent(),
	}, nil
}

func (s State) Close(ctx context.Context, logger util.Logger) error {
	defer func() { _ = telemetry.Close(ctx) }()
	{{{ if .HasModule "migration" }}}if err := s.DB.Close(); err != nil {
		logger.Errorf("error closing database: %+v", err)
	}
	{{{ end }}}{{{ if .HasModule "readonlydb" }}}if err := s.DBRead.Close(); err != nil {
		logger.Errorf("error closing read-only database: %+v", err)
	}
	{{{ end }}}{{{ if .HasModule "graphql" }}}if err := s.GraphQL.Close(); err != nil {
		logger.Errorf("error closing GraphQL service: %+v", err)
	}
	{{{ end }}}return s.Services.Close(ctx, logger)
}
