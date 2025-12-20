package stats

import (
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type FileStat struct {
	Name      string    `json:"name"`
	IsDir     bool      `json:"isDir,omitzero"`
	Kids      FileStats `json:"kids,omitempty"`
	Size      int64     `json:"size,omitzero"`
	fullPath  string
	extension string
}

func newFileStat(pth []string, name string, isDir bool) *FileStat {
	fp := util.StringFilePath(append(util.ArrayCopy(pth), name)...)
	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	return &FileStat{Name: name, IsDir: isDir, fullPath: fp, extension: ext}
}

func (f *FileStat) FullPath() string {
	return f.fullPath
}

func (f *FileStat) Count() int {
	return 1 + len(f.Kids)
}

func (f *FileStat) Largest() *FileStat {
	if len(f.Kids) == 0 {
		return f
	}
	l := f.Kids.Largest()
	if l != nil && l.Size > f.Size {
		return l
	}
	return f
}

func (f *FileStat) TotalSize() int64 {
	return f.Size + f.Kids.TotalSize()
}

func (f *FileStat) Extensions() map[string]int {
	ret := map[string]int{}
	if f.extension == "" {
		if f.IsDir {
			ret["<folder>"] = 1
		} else {
			ret["<none>"] = 1
		}
	} else {
		ret[f.extension] = 1
	}
	lo.ForEach(f.Kids, func(x *FileStat, _ int) {
		for k, v := range x.Extensions() {
			ret[k] += v
		}
	})
	return ret
}
