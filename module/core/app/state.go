package app

import (
	"fmt"

	"github.com/fasthttp/router"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/auth"
	"{{{ .Package }}}/app/filesystem"
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
	Router    *router.Router
	Files     filesystem.FileLoader
	Auth      *auth.Service
	Themes    *theme.Service
	Logger    *zap.SugaredLogger
	Services  *Services
}

func NewState(debug bool, bi *BuildInfo, r *router.Router, f filesystem.FileLoader, log *zap.SugaredLogger) (*State, error) {
	rl := log.With(zap.String("service", "router"))
	a := auth.NewService("", log)
	t := theme.NewService(f, log)
	return &State{Debug: debug, BuildInfo: bi, Router: r, Files: f, Auth: a, Themes: t, Logger: rl}, nil
}
