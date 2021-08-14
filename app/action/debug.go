package action

import (
	"strings"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/util"
)

func onDebug(pm *PrjAndMods) *Result {
	if removePath, _ := pm.Cfg.GetString("remove", true); removePath != "" {
		if removePath == "*" {
			return removeAll(pm)
		}
		return remove(removePath, pm)
	}

	return calculateMissing(pm)
}

func removeAll(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	start := util.TimerStart()
	removeCount := 0
	ret.AddLog("removed [%d] %s", removeCount, util.PluralMaybe("file", removeCount))
	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func remove(path string, pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	start := util.TimerStart()

	tgt := pm.PSvc.GetFilesystem(pm.Prj)
	if err := tgt.Remove(path); err != nil {
		return errorResult(err, pm.Cfg, pm.Logger)
	}
	ret.AddLog("removed [%s]", path)
	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func calculateMissing(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	start := util.TimerStart()

	tgt := pm.PSvc.GetFilesystem(pm.Prj)
	filenames, err := tgt.ListFilesRecursive("", pm.Prj.Ignore)
	if err != nil {
		return errorResult(err, pm.Cfg, pm.Logger)
	}

	var generated []string
	for _, fn := range filenames {
		var b []byte
		b, err = tgt.PeekFile(fn, 1024)
		if err != nil {
			return errorResult(err, pm.Cfg, pm.Logger)
		}
		if file.ContainsHeader(string(b)) {
			generated = append(generated, fn)
		}
	}

	src, err := pm.MSvc.GetFilenames(pm.Mods)
	if err != nil {
		return errorResult(err, pm.Cfg, pm.Logger)
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
			return &module.Resolution{Title: title, Project: pm.Prj.Key, Action: "debug", Args: map[string]string{"remove": path}}
		}
		for _, a := range audits {
			actions = append(actions, resFor("Remove ["+a.Path+"]", a.Path))
		}
		actions = append(actions, resFor("Remove all invalid files", "*"))
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: audits, Actions: actions, Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
