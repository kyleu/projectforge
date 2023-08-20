{{{ if .HasModule "wasmserver" }}}//go:build !js
{{{ end }}}package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	DirectoryMode = fs.FileMode(0o755)
	DefaultMode   = fs.FileMode(0o644)
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

func (f *FileSystem) Stat(path string) (*FileInfo, error) {
	p := f.getPath(path)
	s, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	return &FileInfo{Name: s.Name(), Size: s.Size(), Mode: s.Mode(), ModTime: s.ModTime(), IsDir: s.IsDir(), Sys: s.Sys()}, nil
}

func (f *FileSystem) SetMode(path string, mode fs.FileMode) error {
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
	return s.IsDir
}

func (f *FileSystem) CreateDirectory(path string) error {
	p := f.getPath(path)
	if err := os.MkdirAll(p, DirectoryMode); err != nil {
		return errors.Wrapf(err, "unable to create data directory [%s]", p)
	}
	return nil
}
