package controller

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/gql"
	"projectforge.dev/projectforge/app/lib/task"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var allowedRoutes = []string{"/about", "/admin", "/testbed", "/welcome"}

func initApp(_ context.Context, as *app.State, logger util.Logger) error {
	_, err := gql.NewSchema(as, logger)
	if err != nil {
		return err
	}
	tf := func(ctx context.Context, res *task.Result, logger util.Logger) error {
		res.Log("Testbed!")
		return nil
	}
	t := task.NewTask("testbed", "", "utility", "star", "Who knows what it'll do?", tf)
	t.Fields = util.FieldDescs{{Key: "project", Title: "Project"}}
	err = as.Services.Task.RegisterTask(t)
	if err != nil {
		return err
	}
	return nil
}

func initAppRequest(as *app.State, ps *cutil.PageState) error {
	if err := initProjects(ps.Context, as, ps.Logger); err != nil {
		return errors.Wrap(err, "unable to initialize projects")
	}
	root := as.Services.Projects.Default()
	if root.Info == nil {
		allowed := lo.ContainsBy(allowedRoutes, func(r string) bool {
			return strings.HasSuffix(ps.URI.Path, r)
		})
		if !allowed {
			ps.ForceRedirect = "/welcome"
		}
	}
	as.Services.Task.RegisteredTasks.Get("testbed").Fields[0].Choices = as.Services.Projects.Keys()
	return nil
}

var initProjectsMu = sync.Mutex{}

func initProjects(ctx context.Context, as *app.State, logger util.Logger) error {
	initProjectsMu.Lock()
	defer initProjectsMu.Unlock()
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
