package action

import (
	"strings"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func onDebug(prj *project.Project, mods module.Modules, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	if removePath, _ := cfg.GetString("remove", true); removePath != "" {
		if removePath == "*" {
			return removeAll(prj, mods, cfg, mSvc, pSvc, logger)
		}
		return remove(removePath, prj, mods, cfg, pSvc, logger)
	}

	return calculateMissing(prj, mods, cfg, mSvc, pSvc, logger)
}

func removeAll(prj *project.Project, mods module.Modules, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	start := util.TimerStart()
	removeCount := 0
	ret.AddLog("removed [%d] %s", removeCount, util.PluralMaybe("file", removeCount))
	mr := &module.Result{Keys: mods.Keys(), Status: "OK", Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func remove(path string, prj *project.Project, mods module.Modules, cfg util.ValueMap, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	start := util.TimerStart()

	tgt := pSvc.GetFilesystem(prj)
	err := tgt.Remove(path)
	if err != nil {
		return errorResult(err, cfg, logger)
	}
	ret.AddLog("removed [%s]", path)
	mr := &module.Result{Keys: mods.Keys(), Status: "OK", Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func calculateMissing(prj *project.Project, mods module.Modules, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	start := util.TimerStart()

	tgt := pSvc.GetFilesystem(prj)
	filenames, err := tgt.ListFilesRecursive("", prj.Ignore)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	var generated []string
	for _, fn := range filenames {
		var b []byte
		b, err = tgt.PeekFile(fn, 1024)
		if err != nil {
			return errorResult(err, cfg, logger)
		}
		if file.ContainsHeader(string(b)) {
			generated = append(generated, fn)
		}
	}

	src, err := mSvc.GetFilenames(mods)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	var audits []*diff.Diff
	for _, g := range generated {
		if !util.StringArrayContains(src, g) {
			if (!strings.HasSuffix(g, "client.js.map")) && (!strings.HasSuffix(g, "file/header.go")) {
				audits = append(audits, &diff.Diff{Path: g, Status: diff.StatusMissing})
			}
		}
	}

	var actions module.Resolutions
	if len(audits) > 0 {
		resFor := func(title string, path string) *module.Resolution {
			return &module.Resolution{Title: title, Project: prj.Key, Action: "debug", Args: map[string]string{"remove": path}}
		}
		for _, a := range audits {
			actions = append(actions, resFor("Remove [" + a.Path + "]", a.Path))
		}
		actions = append(actions, resFor("Remove all invalid files", "*"))
	}

	mr := &module.Result{Keys: mods.Keys(), Status: "OK", Diffs: audits, Actions: actions, Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

