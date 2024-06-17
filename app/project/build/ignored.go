package build

import (
	"slices"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Ignored(prj *project.Project, fs filesystem.FileLoader, logger util.Logger) ([]string, error) {
	ign := append(slices.Clone(prj.Ignore), "^tmp")
	files, err := fs.ListFilesRecursive(".", ign, logger)
	if err != nil {
		return nil, err
	}
	ret := &util.StringSlice{}
	lo.ForEach(files, func(f string, _ int) {
		if slices.Contains(prj.Info.IgnoredFiles, f) {
			ret.Push(f)
		}
	})
	return ret.Slice, nil
}
