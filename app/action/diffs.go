package action

import (
	"path"
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(pm *PrjAndMods) (file.Files, diff.Diffs, error) {
	tgt := pm.PSvc.GetFilesystem(pm.Prj)

	srcFiles, err := pm.MSvc.GetFiles(pm.Mods, true)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to get files from [%d] modules", len(pm.Mods))
	}

	if pm.Mods.Get("export") != nil {
		args := &model.Args{}
		if argsX := pm.Prj.Info.ModuleArg("export"); argsX != nil {
			e := util.CycleJSON(argsX, &args)
			if e != nil {
				return nil, nil, errors.Wrap(e, "export module arguments are invalid")
			}
		}
		args.Modules = pm.Mods.Keys()
		files, e := pm.ESvc.Files(args, true)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to export code")
		}
		srcFiles = append(srcFiles, files...)
		e = pm.ESvc.Inject(args, srcFiles)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to inject code")
		}
	}

	portOffsets := map[string]int{}
	var configVars util.KeyTypeDescs
	for _, m := range pm.Prj.Modules {
		mod := pm.Mods.Get(m)
		configVars = append(configVars, mod.ConfigVars...)
		for k, v := range mod.PortOffsets {
			portOffsets[k] = v
		}
	}

	ctx := pm.Prj.ToTemplateContext(configVars, portOffsets)

	for _, f := range srcFiles {
		origPath := f.FullPath()
		if strings.Contains(origPath, delimStart) {
			newPath, e := runTemplate("filename", origPath, ctx)
			if e != nil {
				return nil, nil, e
			}
			p, n := path.Split(newPath)
			f.Path = strings.Split(p, "/")
			f.Name = n
		}
		err = file.ReplaceSections(f, tgt)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, f := range srcFiles {
		f.Content, err = runTemplateFile(f, ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	dfs := diff.FileLoader(srcFiles, tgt, false, pm.Logger)
	return srcFiles, dfs, nil
}
