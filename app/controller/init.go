package controller

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var allowedRoutes = []string{"/about", "/admin", "/testbed", "/welcome"}

// Initialize app-specific system dependencies.
func initApp(_ *app.State, _ util.Logger) {
}

// Configure app-specific data for each request.
func initAppRequest(as *app.State, ps *cutil.PageState) error {
	if err := initProjects(ps.Context, as, ps.Logger); err != nil {
		return errors.Wrap(err, "unable to initialize projects")
	}
	root := as.Services.Projects.Default()
	if root.Info == nil && !lo.ContainsBy(allowedRoutes, func(r string) bool {
		return strings.HasSuffix(string(ps.URI.Path()), r)
	}) {
		ps.ForceRedirect = "/welcome"
	}
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
