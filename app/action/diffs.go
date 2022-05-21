package action

import (
	"context"
	"path"
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(ctx context.Context, pm *PrjAndMods) (file.Files, diff.Diffs, error) {
	tgt := pm.PSvc.GetFilesystem(pm.Prj)

	srcFiles, err := pm.MSvc.GetFiles(pm.Mods, true)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to get files from [%d] modules", len(pm.Mods))
	}

	if pm.Mods.Get("export") != nil {
		args, err := pm.Prj.Info.ModuleArgExport()
		if err != nil {
			return nil, nil, errors.Wrap(err, "export module arguments are invalid")
		}
		args.Modules = pm.Mods.Keys()
		files, e := pm.ESvc.Files(ctx, args, true, pm.Logger)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to export code")
		}
		srcFiles = append(srcFiles, files...)
		e = pm.ESvc.Inject(ctx, args, srcFiles, pm.Logger)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to inject code")
		}
	}

	configVars, portOffsets := parse(pm)

	tCtx := pm.Prj.ToTemplateContext(configVars, portOffsets)

	for _, f := range srcFiles {
		origPath := f.FullPath()
		if strings.Contains(origPath, delimStart) {
			newPath, e := runTemplate("filename", origPath, tCtx)
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

	for _, m := range pm.Prj.Modules {
		mod := pm.Mods.Get(m)
		for _, src := range mod.ConfigVars {
			var hit bool
			for _, tgt := range configVars {
				if src.Key == tgt.Key {
					hit = true
					break
				}
			}
			if !hit {
				configVars = append(configVars, src)
			}
		}
		configVars.Sort()

		for k, v := range mod.PortOffsets {
			portOffsets[k] = v
		}
	}
	return configVars, portOffsets
}
