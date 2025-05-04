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
	configVars, portOffsets := parse(pm)
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

func parse(pm *PrjAndMods) (util.KeyTypeDescs, map[string]int) {
	var configVars util.KeyTypeDescs
	portOffsets := map[string]int{}

	lo.ForEach(pm.Prj.Modules, func(m string, _ int) {
		mod := pm.Mods.Get(m)
		lo.ForEach(mod.ConfigVars, func(src *util.KeyTypeDesc, _ int) {
			hit := lo.ContainsBy(configVars, func(tgt *util.KeyTypeDesc) bool {
				return src.Matches(tgt)
			})
			if !hit {
				configVars = append(configVars, src)
			}
		})
		for k, v := range mod.PortOffsets {
			portOffsets[k] = v
		}
	})
	return configVars.Sort(), portOffsets
}
