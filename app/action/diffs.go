package action

import (
	"path"
	"strings"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(pm *PrjAndMods) (file.Files, diff.Diffs, error) {
	tgt := pm.PSvc.GetFilesystem(pm.Prj)

	srcFiles, err := pm.MSvc.GetFiles(pm.Mods)
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
		files, e := pm.ESvc.Files(args)
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
	for _, m := range pm.Prj.Modules {
		for k, v := range pm.Mods.Get(m).PortOffsets {
			portOffsets[k] = v
		}
	}

	ctx := pm.Prj.ToTemplateContext(portOffsets)

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
