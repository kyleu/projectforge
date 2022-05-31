// Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	DirectoryMode = os.FileMode(0o755)
	DefaultMode   = os.FileMode(0o644)
)

type FileSystem struct {
	root string
}

var _ FileLoader = (*FileSystem)(nil)

func NewFileSystem(root string) *FileSystem {
	return &FileSystem{root: root}
}

func (f *FileSystem) getPath(ss ...string) string {
	s := filepath.Join(ss...)
	if strings.HasPrefix(s, f.root) {
		return s
	}
	return filepath.Join(f.root, s)
}

func (f *FileSystem) Root() string {
	return f.root
}

func (f *FileSystem) Clone() FileLoader {
	return NewFileSystem(f.root)
}

func (f *FileSystem) Stat(path string) (os.FileInfo, error) {
	p := f.getPath(path)
	s, err := os.Stat(p)
	if err == nil {
		return s, nil
	}
	return nil, err
}

func (f *FileSystem) SetMode(path string, mode os.FileMode) error {
	p := f.getPath(path)
	return os.Chmod(p, mode)
}

func (f *FileSystem) Exists(path string) bool {
	x, _ := f.Stat(path)
	return x != nil
}

func (f *FileSystem) IsDir(path string) bool {
	s, err := f.Stat(path)
	if s == nil || err != nil {
		return false
	}
	return s.IsDir()
}

func (f *FileSystem) CreateDirectory(path string) error {
	p := f.getPath(path)
	if err := os.MkdirAll(p, DirectoryMode); err != nil {
		return errors.Wrapf(err, "unable to create data directory [%s]", p)
	}
	return nil
}
