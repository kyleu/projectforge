package controller

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

// Initialize app-specific system dependencies.
func initApp(*app.State, util.Logger) {
}

// Configure app-specific data for each request.
func initAppRequest(as *app.State, ps *cutil.PageState) error {
	if err := initProjects(ps.Context, as, ps.Logger); err != nil {
		return errors.Wrap(err, "unable to initialize projects")
	}

	root := as.Services.Projects.Default()
	p := string(ps.URI.Path())
	if root.Info == nil && !strings.HasSuffix(p, "/about") && !strings.HasPrefix(p, "/welcome") {
		ps.ForceRedirect = "/welcome"
	}

	return nil
}

// Initialize system dependencies for the marketing site.
func initSite(as *app.State, logger util.Logger) {
	as.Services = &app.Services{
		Modules: module.NewService(context.Background(), as.Files, logger),
	}
}

// Configure marketing site data for each request.
func initSiteRequest(*app.State, *cutil.PageState) error {
	return nil
}

func initProjects(ctx context.Context, as *app.State, logger util.Logger) error {
	prjs, err := as.Services.Projects.Refresh(logger)
	if err != nil {
		return errors.Wrap(err, "can't load projects")
	}
	for _, prj := range prjs {
		var mods project.ModuleDefs
		if prj.Info != nil {
			mods = prj.Info.ModuleDefs
		}
		var keys []string
		for _, mod := range mods {
			k, err := as.Services.Modules.Register(ctx, prj.Path, mod.Key, mod.Path, mod.URL, logger)
			if err != nil {
				return errors.Wrapf(err, "unable to register module definition for module [%s]", mod.Key)
			}
			keys = append(keys, k...)
		}
		if len(keys) > 0 {
			if len(keys) != 1 && keys[0] != "*" {
				logger.Debugf("Loaded modules for [%s]: %s", prj.Key, strings.Join(keys, ", "))
			}
		}
	}
	return nil
}
