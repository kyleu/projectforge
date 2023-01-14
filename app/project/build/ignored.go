package build

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Ignored(prj *project.Project, fs filesystem.FileLoader, logger util.Logger) ([]string, error) {
	var ret []string
	ign := append(slices.Clone(prj.Ignore), "^tmp")
	files, err := fs.ListFilesRecursive(".", ign, logger)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		content, err := fs.ReadFile(f)
		if err != nil {
			ret = append(ret, fmt.Sprintf("ERROR[%s]: %s", err.Error(), f))
		}
		if strings.Contains(string(content), file.IgnorePattern) {
			ret = append(ret, f)
		}
	}
	return ret, nil
}
