package build

import (
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
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
		content, err := fs.ReadFile(f)
		if err != nil {
			ret.Pushf("ERROR[%s]: %s", err.Error(), f)
		}
		if strings.Contains(string(content), file.IgnorePattern) {
			ret.Push(f)
		}
	})
	return ret.Slice, nil
}
