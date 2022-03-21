package action_test

import (
	"context"
	"testing"

	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/export"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
)

func prj(key string, modules ...string) *project.Project {
	return &project.Project{
		Key:     key,
		Name:    "[" + key + "] Test",
		Version: "0.0.1",
		Package: "projectforge.dev/projectforge/test/" + key,
		Port:    11100,
		Modules: modules,
		Info:    &project.Info{},
		Build: &project.Build{
			Publish: true,
			Private: true,
		},
		Path: "../../tmp/test/projects/" + key,
	}
}

func TestFoo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		project *project.Project
	}{
		{project: prj("core", "core")},
	}

	logger, _ := log.InitLogging(true)
	fs := filesystem.NewFileSystem("../../tmp/test/cfg", logger)
	mSvc := module.NewService(context.Background(), fs, logger)
	pSvc := project.NewService(logger)
	eSvc := export.NewService(logger)

	for _, c := range cases {
		t.Log("Testing [" + c.project.Name + "]")
		cfg := c.project.ToMap()
		params := &action.Params{ProjectKey: c.project.Key, T: action.TypeCreate, Cfg: cfg, MSvc: mSvc, PSvc: pSvc, ESvc: eSvc, Logger: logger}
		res := action.Apply(context.Background(), params)
		t.Log(res.Status)
	}
}
