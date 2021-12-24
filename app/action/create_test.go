package action_test

import (
	"context"
	"testing"

	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/codegen"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/log"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
)

func prj(key string, modules ...string) *project.Project {
	return &project.Project{
		Key:     key,
		Name:    "[" + key + "] Test",
		Version: "0.0.1",
		Package: "github.com/kyleu/projectforge/test/" + key,
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
	mSvc := module.NewService(fs, logger)
	pSvc := project.NewService(logger)
	cSvc := codegen.NewService(logger)

	for _, c := range cases {
		t.Log("Testing [" + c.project.Name + "]")
		cfg := c.project.ToMap()
		params := &action.Params{Span: nil, ProjectKey: c.project.Key, T: action.TypeCreate, Cfg: cfg, MSvc: mSvc, PSvc: pSvc, CSvc: cSvc, Logger: logger}
		res := action.Apply(context.Background(), params)
		t.Log(res.Status)
	}
}
