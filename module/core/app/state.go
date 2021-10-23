package app

import (
	"fmt"
{{{ if .HasModule "migration" }}}
	"github.com/pkg/errors"{{{ end }}}
	"go.uber.org/zap"
{{{ if .HasModule "oauth" }}}
	"{{{ .Package }}}/app/auth"{{{ end }}}{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/database"{{{ end }}}
	"{{{ .Package }}}/app/filesystem"
	"{{{ .Package }}}/app/telemetry"
	"{{{ .Package }}}/app/theme"
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
	BuildInfo *BuildInfo
	Files     filesystem.FileLoader{{{ if .HasModule "oauth" }}}
	Auth      *auth.Service{{{ end }}}{{{ if .HasModule "migration" }}}
	DB        *database.Service{{{ end }}}
	Themes    *theme.Service
	Logger    *zap.SugaredLogger
	Services  *Services
}

func NewState(debug bool, bi *BuildInfo, f filesystem.FileLoader, logger *zap.SugaredLogger) (*State, error) {
	_ = telemetry.InitializeIfNeeded(true, logger){{{ if .HasModule "oauth" }}}
	as := auth.NewService("", logger){{{ end }}}
	ts := theme.NewService(f, logger){{{ if .HasModule "migration" }}}

	db, err := database.OpenDefaultPostgres(logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}{{{ end }}}

	return &State{
		Debug:     debug,
		BuildInfo: bi,
		Files:     f{{{ if .HasModule "oauth" }}},
		Auth:      as{{{ end }}}{{{ if .HasModule "migration" }}},
		DB:        db{{{ end }}},
		Themes:    ts,
		Logger:    logger,
	}, nil
}
