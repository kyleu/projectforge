package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

var DefaultMode = os.FileMode(0o755)

type FileSystem struct {
	root     string
	logger   *zap.SugaredLogger
	children []FileLoader
}

var _ FileLoader = (*FileSystem)(nil)

func NewFileSystem(root string, logger *zap.SugaredLogger) *FileSystem {
	return &FileSystem{root: root, logger: logger.With(zap.String("service", "filesystem"))}
}

func (f *FileSystem) AddChildren(fls ...FileLoader) {
	f.children = append(f.children, fls...)
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

func (f *FileSystem) Stat(path string) (os.FileInfo, error) {
	p := f.getPath(path)
	s, err := os.Stat(p)
	if err == nil {
		return s, nil
	}
	for _, c := range f.children {
		if x, e := c.Stat(path); x != nil && e == nil {
			return x, nil
		}
	}
	return nil, err
}

func (f *FileSystem) Exists(path string) bool {
	x, _ := f.Stat(path)
	for _, c := range f.children {
		if c.Exists(path) {
			return true
		}
	}
	return x != nil
}

func (f *FileSystem) IsDir(path string) bool {
	s, err := f.Stat(path)
	if s == nil || err == nil {
		return false
	}
	return s.IsDir()
}

func (f *FileSystem) CreateDirectory(path string) error {
	p := f.getPath(path)
	if err := os.MkdirAll(p, DefaultMode); err != nil {
		return errors.Wrapf(err, "unable to create data directory [%s]", p)
	}
	return nil
}

func (f *FileSystem) ReadFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(f.getPath(path))
	if err != nil {
		for _, c := range f.children {
			if x, err := c.ReadFile(path); err == nil {
				return x, nil
			}
		}
		return nil, errors.Wrapf(err, "unable to read file [%s]", path)
	}
	return b, nil
}

func (f *FileSystem) WriteFile(path string, content []byte, mode os.FileMode, overwrite bool) error {
	p := f.getPath(path)
	s, err := os.Stat(p)
	if os.IsExist(err) && !overwrite {
		return errors.Errorf("file [%s] exists, will not overwrite", p)
	}
	if mode == 0 {
		if s == nil {
			mode = DefaultMode
		} else {
			mode = s.Mode()
		}
	}
	dd := filepath.Dir(path)
	err = f.CreateDirectory(dd)
	if err != nil {
		return errors.Wrapf(err, "unable to create data directory [%s]", dd)
	}
	file, err := os.Create(p)
	if err != nil {
		return errors.Wrapf(err, "unable to create file [%s]", p)
	}
	err = file.Chmod(mode)
	if err != nil {
		return errors.Wrapf(err, "unable to set mode [%s] for file [%s]", mode.String(), p)
	}
	defer func() { _ = file.Close() }()
	_, err = file.Write(content)
	if err != nil {
		return errors.Wrapf(err, "unable to write content to file [%s]", p)
	}
	return nil
}
