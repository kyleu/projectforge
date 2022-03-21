package app

import (
	"fmt"
	"time"

	"go.uber.org/zap"
{{{ if .HasModule "oauth" }}}
	"{{{ .Package }}}/app/lib/auth"{{{ end }}}{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}
	"{{{ .Package }}}/app/lib/filesystem"
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
	BuildInfo *BuildInfo
	Files     filesystem.FileLoader{{{ if .HasModule "oauth" }}}
	Auth      *auth.Service{{{ end }}}{{{ if .HasModule "migration" }}}
	DB        *database.Service{{{ end }}}
	Themes    *theme.Service
	Logger    *zap.SugaredLogger
	Services  *Services
	Started   time.Time
}

func NewState(debug bool, bi *BuildInfo, f filesystem.FileLoader, logger *zap.SugaredLogger) (*State, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	time.Local = loc

	_ = telemetry.InitializeIfNeeded(true, bi.Version, logger){{{ if .HasModule "oauth" }}}
	as := auth.NewService("", logger){{{ end }}}
	ts := theme.NewService(f, logger)

	return &State{
		Debug:     debug,
		BuildInfo: bi,
		Files:     f{{{ if .HasModule "oauth" }}},
		Auth:      as{{{ end }}},
		Themes:    ts,
		Logger:    logger,
		Started:   time.Now(),
	}, nil
}
