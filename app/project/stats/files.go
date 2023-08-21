package stats

import (
	"path/filepath"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type FileStats []*FileStat

func (f FileStats) Count() int {
	return lo.SumBy(f, func(x *FileStat) int {
		return x.Count()
	})
}

func (f FileStats) Largest() *FileStat {
	var ret *FileStat
	lo.ForEach(f, func(x *FileStat, _ int) {
		if ret == nil || x.Size > ret.Size {
			ret = x.Largest()
		}
	})
	return ret
}

func (f FileStats) TotalSize() int64 {
	return lo.Sum(lo.Map(f, func(x *FileStat, _ int) int64 {
		return x.TotalSize()
	}))
}

func (f FileStats) Extensions() map[string]int {
	ret := map[string]int{}
	lo.ForEach(f, func(x *FileStat, _ int) {
		for k, v := range x.Extensions() {
			ret[k] += v
		}
	})
	return ret
}

func GetFileStats(fs filesystem.FileLoader, pth string, logger util.Logger) (FileStats, error) {
	if pth == "" {
		pth = "."
	}
	return listDir(fs, logger, pth)
}

var ignores = []string{"^tmp/", "^node_modules/", "^libs/"}

func listDir(fs filesystem.FileLoader, logger util.Logger, pth ...string) (FileStats, error) {
	curr := fs.ListFiles(filepath.Join(pth...), ignores, logger)
	ret := make(FileStats, 0, len(curr))
	for _, f := range curr {
		s := newFileStat(pth, f.Name, f.IsDir)
		if f.IsDir {
			kids, err := listDir(fs, logger, filepath.Join(append(pth, f.Name)...))
			if err != nil {
				return nil, err
			}
			s.Kids = kids
		} else {
			s.Size = f.Size
		}
		ret = append(ret, s)
	}
	return ret, nil
}
