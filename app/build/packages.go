package build

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type Pkg struct {
	Path string `json:"path"`
}

type Pkgs []*Pkg

func Packages(fs filesystem.FileLoader, logger util.Logger) (Pkgs, error) {
	var ret Pkgs
	files, err := fs.ListFilesRecursive(".", nil, logger)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.HasSuffix(f, ".go") || strings.HasSuffix(f, ".html") {
			x := &Pkg{Path: f}
			ret = append(ret, x)
		}
	}
	return ret, nil
}
