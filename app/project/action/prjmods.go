package action

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

type PrjAndMods struct {
	Cfg    util.ValueMap
	Prj    *project.Project
	File   json.RawMessage
	Mods   module.Modules
	MSvc   *module.Service
	PSvc   *project.Service
	XSvc   *exec.Service
	SSvc   *websocket.Service
	ESvc   *export.Service
	EArgs  *model.Args
	Logger util.Logger
}

func getPrjAndMods(ctx context.Context, p *Params) (context.Context, *PrjAndMods, error) {
	if p.ProjectKey == "" {
		prj := p.PSvc.ByPath("")
		if prj != nil {
			p.ProjectKey = prj.Key
		}
	}

	f := p.PSvc.GetFile(p.ProjectKey)

	prj, err := p.PSvc.Get(p.ProjectKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to load project [%s]", p.ProjectKey)
	}
	if prj.Info != nil {
		for _, mod := range prj.Info.ModuleDefs {
			_, e := p.MSvc.Register(ctx, prj.Path, mod.Key, mod.Path, mod.URL, p.Logger)
			if e != nil {
				return nil, nil, errors.Wrap(e, "unable to register modules")
			}
		}
	}

	mods, err := p.MSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, nil, err
	}

	args, err := prj.ModuleArgExport(p.PSvc, p.Logger)
	if err != nil {
		return nil, nil, err
	}
	if args != nil {
		args.Modules = mods.Keys()
	}

	pm := &PrjAndMods{Cfg: p.Cfg, File: f, Prj: prj, Mods: mods, MSvc: p.MSvc, PSvc: p.PSvc, XSvc: p.XSvc, ESvc: p.ESvc, EArgs: args, Logger: p.Logger}
	return ctx, pm, nil
}
