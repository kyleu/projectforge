package app

import (
	"fmt"

	"github.com/fasthttp/router"
	"$PF_PACKAGE$/app/auth"
	"$PF_PACKAGE$/app/filesystem"
	"$PF_PACKAGE$/app/theme"
	"$PF_PACKAGE$/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	ret := &State{
		Debug:     debug,
		BuildInfo: bi,
		Router:    r,
		Files:     f,
		Auth:      auth.NewService("", log),
		Themes:    theme.NewService(f, log),
		Logger:    rl,
	}

	svcs, err := NewServices(ret)
	if err != nil {
		return nil, errors.Wrap(err, "error creating services")
	}
	ret.Services = svcs

	return ret, nil
}
