package action

import (
	"path"
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func getGeneratedFiles(tgt filesystem.FileLoader, ignore []string, logger util.Logger) ([]string, error) {
	filenames, err := tgt.ListFilesRecursive("", ignore, logger)
	if err != nil {
		return nil, err
	}
	ret := &util.StringSlice{}
	for _, fn := range filenames {
		b, e := tgt.PeekFile(fn, 1024)
		if e != nil {
			return nil, e
		}
		if file.ContainsHeader(string(b)) {
			ret.Push(fn)
		}
	}
	return ret.Slice, nil
}

func getModuleFiles(pm *PrjAndMods) ([]string, error) {
	ret, err := pm.MSvc.GetFilenames(pm.Mods, pm.Logger)
	if err != nil {
		return nil, err
	}
	if pm.Mods.Get("export") != nil {
		args, err := pm.Prj.ModuleArgExport(pm.PSvc, pm.Logger)
		if err != nil {
			return nil, err
		}
		args.Modules = pm.Mods.Keys()
		lb := util.StringDefaultLinebreak
		pfs, err := pm.PSvc.GetFilesystem(pm.Prj)
		if err != nil {
			return nil, err
		}
		modContent, _ := pfs.ReadFile("go.mod")
		if modContent != nil {
			lb = util.StringDetectLinebreak(string(modContent))
		}
		files, e := pm.ESvc.Files(pm.Prj, args, true, lb)
		if e != nil {
			return nil, err
		}
		lo.ForEach(files, func(f *file.File, _ int) {
			ret = append(ret, f.FullPath())
		})
	}
	return ret, nil
}

func getEmptyFolders(tgt filesystem.FileLoader, ignore []string, logger util.Logger, pth ...string) ([]string, error) {
	ret := &util.StringSlice{}
	pStr := path.Join(pth...)
	fc := len(tgt.ListFiles(pStr, nil, logger))
	ds := tgt.ListDirectories(pStr, ignore, logger)
	if fc == 0 && len(ds) == 0 {
		ret.Push(pStr)
	}
	for _, d := range ds {
		p := append(slices.Clone(pth), d)
		childRes, err := getEmptyFolders(tgt, ignore, logger, p...)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get empty folders for [%s/%s]", path.Join(p...), d)
		}
		ret.Push(childRes...)
	}
	return ret.Slice, nil
}
