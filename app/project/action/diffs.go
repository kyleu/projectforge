package action

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/project/template"
	"projectforge.dev/projectforge/app/util"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(pm *PrjAndMods) (file.Files, diff.Diffs, error) {
	isPrivate := pm.Prj.Build == nil || pm.Prj.Build.Private
	srcFiles, err := pm.MSvc.GetFiles(pm.Mods, isPrivate, pm.Logger)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to get files from [%d] modules", len(pm.Mods))
	}
	if pm.Prj.Key == util.AppKey {
		for _, x := range pm.MSvc.Modules() {
			if pm.Mods.Get(x.Key) == nil {
				docs, err := pm.MSvc.GetDocFiles(x, isPrivate, pm.Logger)
				if err != nil {
					return nil, nil, err
				}
				srcFiles = append(srcFiles, docs...)
			}
		}
	}
	if pm.Mods.Get("export") != nil && len(srcFiles) > 0 {
		linebreak := util.StringDetectLinebreak(srcFiles[0].Content)
		pm.Prj.ExportArgs.Modules = pm.Mods.Keys()
		files, e := pm.ESvc.Files(pm.Prj, linebreak)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to export code")
		}
		srcFiles = lo.UniqBy(append(files, srcFiles...), func(f *file.File) string {
			return f.FullPath()
		})
	}
	settingsJSON, err := parseSettings(pm.FS, pm.Prj.Theme)
	if err != nil {
		// return nil, nil, errors.Wrapf(err, "unable to generate [" + vscodeSettingsPath + "]")
		pm.Logger.Warnf("unable to generate ["+vscodeSettingsPath+"]: %+v", err)
	}
	if settingsJSON != nil {
		srcFiles = append(srcFiles, settingsJSON)
	}
	launchJSON, err := parseLaunch(pm.FS, pm.Prj.Name, pm.Prj.ExecSafe(), pm.Prj.HasModule("upgrade"))
	if err != nil {
		// return nil, nil, errors.Wrapf(err, "unable to generate [" + vscodeLaunchPath + "]")
		pm.Logger.Warnf("unable to generate ["+vscodeLaunchPath+"]: %+v", err)
	}
	if launchJSON != nil {
		srcFiles = append(srcFiles, launchJSON)
	}

	configVars, portOffsets := parseConfig(pm)
	lb := util.StringDetectLinebreak(string(pm.File))
	tCtx := template.ToTemplateContext(pm.Prj, configVars, portOffsets, lb)
	for _, f := range srcFiles {
		origPath := f.FullPath()
		if strings.Contains(origPath, delimStart) {
			newPath, e := runTemplate("filename", origPath, tCtx)
			if e != nil {
				return nil, nil, e
			}
			p, n := util.StringSplitLast(newPath, '/', true)
			f.Path = strings.Split(p, string(filepath.Separator))
			f.Name = n
		}
		err = file.ReplaceSections(f, pm.FS)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, f := range srcFiles {
		f.Content, err = runTemplateFile(f, tCtx)
		if err != nil {
			return nil, nil, err
		}
	}
	dfs, err := diff.FileLoader(pm.Mods.Keys(), srcFiles, pm.FS, pm.Prj.Info.IgnoredFiles, false, pm.Logger)
	if err != nil {
		return nil, nil, err
	}
	return srcFiles, dfs, nil
}
