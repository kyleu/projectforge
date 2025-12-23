package build

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Packages(prj *project.Project, fs filesystem.FileLoader, showAll bool, logger util.Logger) (Pkgs, error) {
	var ret Pkgs
	files, err := fs.ListFilesRecursive(".", nil, logger)
	if err != nil {
		return nil, err
	}
	root := prj.Package
	for _, f := range files {
		if ((!strings.HasSuffix(f, util.ExtGo)) && (!strings.HasSuffix(f, util.ExtHTML))) || strings.HasPrefix(f, "module/") {
			continue
		}
		dir, fn := util.StringSplitPath(f)
		pth := util.StringPath(root, dir)
		curr := ret.Get(pth)
		if curr == nil {
			curr = &Pkg{Path: pth}
			ret = append(ret, curr)
		}
		curr.Files = append(curr.Files, fn)
		bs, err := fs.ReadFile(f)
		if err != nil {
			return nil, err
		}
		lines := util.StringSplitLines(string(bs))
		imps, _, _, err := processFileImports(fn, lines, root)
		if err != nil {
			return nil, err
		}
		lo.ForEach(imps, func(impRaw string, _ int) {
			imp, typ := util.StringCutLast(impRaw, ':', true)
			if imp != "" && (showAll || typ == "self") {
				curr.AddDep(imp)
			}
		})
	}
	return ret.Sort(), nil
}
