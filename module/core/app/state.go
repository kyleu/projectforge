package app

import (
	"fmt"

	"go.uber.org/zap"
{{{ if .HasModule "oauth" }}}
	"{{{ .Package }}}/app/auth"{{{ end }}}
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
	Auth      *auth.Service{{{ end }}}
	Themes    *theme.Service
	Telemetry *telemetry.Service
	Logger    *zap.SugaredLogger
	Services  *Services
}

func NewState(debug bool, bi *BuildInfo, f filesystem.FileLoader, m *telemetry.Metrics, log *zap.SugaredLogger) (*State, error) {
	return &State{
		Debug:     debug,
		BuildInfo: bi,
		Files:     f,
		Auth:      auth.NewService("", log),
		Themes:    theme.NewService(f, log),
		Telemetry: telemetry.NewService(m, log),
		Logger:    log.With(zap.String("service", "router")),
	}, nil
}
