package action

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
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
	srcFiles = append(srcFiles, readmeFile.Clone())
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
		if len(pm.Prj.ExportArgs.Models.WithService()) == 0 {
			srcFiles = lo.Reject(srcFiles, func(f *file.File, _ int) bool {
				return util.Str(f.FullPath()).HasPrefix("app/lib/svc/")
			})
		}
	}
	settingsJSON, err := parseSettings(pm.FS, pm.Prj.Theme)
	if err != nil {
		pm.Logger.Warnf("unable to generate [%s] for project [%s]: %+v", vscodeSettingsPath, pm.Prj.Key, err)
	}
	if settingsJSON != nil {
		srcFiles = append(srcFiles, settingsJSON)
	}
	launchJSON, err := parseLaunch(pm.FS, pm.Prj.Name, pm.Prj.ExecSafe(), pm.Prj.HasModule("upgrade"))
	if err != nil {
		pm.Logger.Warnf("unable to generate [%s] for project [%s]: %+v", vscodeLaunchPath, pm.Prj.Key, err)
	}
	if launchJSON != nil {
		srcFiles = append(srcFiles, launchJSON)
	}

	configVars, portOffsets := parseConfig(pm)
	lb := util.StringDetectLinebreak(string(pm.File))
	tCtx := template.ToTemplateContext(pm.Prj, configVars, portOffsets, lb)
	err = runTemplates(srcFiles, tCtx, pm.FS)
	if err != nil {
		return nil, nil, err
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

func runTemplates(srcFiles file.Files, tCtx *template.Context, fs filesystem.FileLoader) error {
	for _, f := range srcFiles {
		origPath := f.FullPath()
		if strings.Contains(origPath, delimStart) {
			newPath, err := runTemplate("filename", origPath, tCtx)
			if err != nil {
				return err
			}
			p, n := util.StringCutLast(newPath, '/', true)
			f.Path = strings.Split(p, string(filepath.Separator))
			f.Name = n
		}
		err := file.ReplaceSections(f, fs)
		if err != nil {
			return err
		}
	}
	return nil
}
