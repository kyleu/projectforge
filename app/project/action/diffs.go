package action

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(pm *PrjAndMods) (file.Files, diff.Diffs, error) {
	srcFiles, err := pm.MSvc.GetFiles(pm.Mods, true, pm.Logger)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to get files from [%d] modules", len(pm.Mods))
	}

	if pm.Mods.Get("export") != nil && len(srcFiles) > 0 {
		linebreak := util.StringDetectLinebreak(srcFiles[0].Content)
		args, e := pm.Prj.ModuleArgExport(pm.PSvc, pm.Logger)
		if e != nil {
			return nil, nil, errors.Wrap(e, "export module arguments are invalid")
		}
		args.Modules = pm.Mods.Keys()
		files, e := pm.ESvc.Files(pm.Prj, args, true, linebreak)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to export code")
		}
		srcFiles = append(srcFiles, files...)
		e = pm.ESvc.Inject(args, srcFiles, linebreak)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to inject code")
		}
	}
	if pm.Mods.Get("datadog") != nil {
		svc := ServiceDefinition(pm.Prj)
		f := file.NewFile("doc/service.json", filesystem.DefaultMode, util.ToJSONBytes(svc, true), false, pm.Logger)
		srcFiles = append(srcFiles, f)
	}

	configVars, portOffsets := parse(pm)

	pm.Prj.ExportArgs, _ = pm.Prj.ModuleArgExport(pm.PSvc, pm.Logger)

	lb := util.StringDetectLinebreak(string(pm.File))
	tCtx := pm.Prj.ToTemplateContext(configVars, portOffsets, lb)

	tgt, err := pm.PSvc.GetFilesystem(pm.Prj)
	if err != nil {
		return nil, nil, err
	}

	for _, f := range srcFiles {
		origPath := f.FullPath()
		if strings.Contains(origPath, delimStart) {
			newPath, e := runTemplate("filename", origPath, tCtx)
			if e != nil {
				return nil, nil, e
			}
			p, n := path.Split(newPath)
			f.Path = strings.Split(p, string(filepath.ListSeparator))
			f.Name = n
		}
		err = file.ReplaceSections(f, tgt)
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

	dfs, err := diff.FileLoader(pm.Mods.Keys(), srcFiles, tgt, false, pm.Logger)
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
				return src.Key == tgt.Key
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
