package app

import (
	"context"
	"fmt"
	"sync"
	"time"{{{ if .HasUser }}}

	"github.com/google/uuid"{{{ end }}}
{{{ if .HasAccount }}}
	"{{{ .Package }}}/app/lib/auth"{{{ end }}}{{{ if .HasDatabaseModule }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}{{{ if .HasModule "filesystem" }}}
	"{{{ .Package }}}/app/lib/filesystem"{{{ end }}}{{{ if .HasModule "graphql" }}}
	"{{{ .Package }}}/app/lib/graphql"{{{ end }}}
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/theme"{{{ if .HasUser }}}
	"{{{ .Package }}}/app/user"{{{ end }}}
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
	Files     filesystem.FileLoader{{{ end }}}{{{ if .HasAccount }}}
	Auth      *auth.Service{{{ end }}}{{{ if .HasModule "migration" }}}
	DB        *database.Service{{{ end }}}{{{ if .HasModule "readonlydb" }}}
	DBRead    *database.Service{{{ end }}}{{{ if .HasModule "graphql" }}}
	GraphQL   *graphql.Service{{{ end }}}
	Themes    *theme.Service
	Services  *Services
	Started   time.Time
}

func NewState(debug bool, bi *BuildInfo{{{ if .HasModule "filesystem" }}}, f filesystem.FileLoader{{{ end }}}, enableTelemetry bool, {{{ if .HasAccount }}}port{{{ else }}}_{{{ end }}} uint16, logger util.Logger) (*State, error) {
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

	_ = telemetry.InitializeIfNeeded(enableTelemetry, bi.Version, logger)

	return &State{
		Debug:     debug,
		BuildInfo: bi{{{ if .HasModule "filesystem" }}},
		Files:     f{{{ end }}}{{{ if .HasAccount }}},
		Auth:      auth.NewService("", port, logger){{{ end }}}{{{ if .HasModule "graphql" }}},
		GraphQL:   graphql.NewService(){{{ end }}},
		Themes:    theme.NewService({{{ if .HasModule "filesystem" }}}f{{{ end }}}),
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
}{{{ if .HasUser }}}

func (s State) User(ctx context.Context, id uuid.UUID, logger util.Logger) (*user.User, error) {
	if s.Services == nil || s.Services.User == nil {
		return nil, nil
	}
	return s.Services.User.Get(ctx, nil, id, logger)
}{{{ end }}}
