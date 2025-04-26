package filesystem

import (
	"cmp"
	"io/fs"
	"net/url"
	"slices"
	"strings"

	"github.com/samber/lo"
)

var (
	DirectoryMode = FileMode(0o755)
	DefaultMode   = FileMode(0o644)
)

type FileMode uint32

func (m FileMode) ToFS() fs.FileMode {
	return fs.FileMode(m)
}

type FileInfo struct {
	Name  string
	Size  int64
	Mode  FileMode
	IsDir bool
}

func (f *FileInfo) Equal(x *FileInfo) bool {
	return f.Name == x.Name && f.Size == x.Size && f.Mode == x.Mode && f.IsDir == x.IsDir
}

func (f *FileInfo) QueryEscapedPath(pth ...string) string {
	return "/" + strings.Join(lo.Map(append(pth, f.Name), func(x string, _ int) string {
		return url.QueryEscape(x)
	}), "/")
}

func FileInfoFromFS(s fs.FileInfo) *FileInfo {
	return &FileInfo{Name: s.Name(), Size: s.Size(), Mode: FileMode(s.Mode()), IsDir: s.IsDir()}
}

func FileInfoFromDE(e fs.DirEntry) *FileInfo {
	i, _ := e.Info()
	return &FileInfo{Name: e.Name(), Size: i.Size(), Mode: FileMode(i.Mode()), IsDir: e.IsDir()}
}

type FileInfos []*FileInfo

func FileInfosFromFS(x []fs.FileInfo) FileInfos {
	return lo.Map(x, func(x fs.FileInfo, _ int) *FileInfo {
		return FileInfoFromFS(x)
	})
}

func (f FileInfos) Sorted() FileInfos {
	slices.SortFunc(f, func(l *FileInfo, r *FileInfo) int {
		return cmp.Compare(l.Name, r.Name)
	})
	return f
}

func (f FileInfos) Equal(x FileInfos) bool {
	if len(f) != len(x) {
		return false
	}
	for idx, i := range f {
		if !i.Equal(x[idx]) {
			return false
		}
	}
	return true
}
