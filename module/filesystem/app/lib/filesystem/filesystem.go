package filesystem

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"{{{ .Package }}}/app/util"
)

var (
	memFS     afero.Fs
	osFS      afero.Fs
	memFSOnce sync.Once
	osFSOnce  sync.Once
)

func MemFS() afero.Fs {
	memFSOnce.Do(func() {
		memFS = afero.NewMemMapFs()
	})
	return memFS
}

func OSFS() afero.Fs {
	osFSOnce.Do(func() {
		osFS = afero.NewOsFs()
	})
	return osFS
}

type FileSystem struct {
	Mode     string `json:"mode,omitzero"`
	ReadOnly bool   `json:"readOnly,omitzero"`
	root     string
	f        afero.Fs
}

var _ FileLoader = (*FileSystem)(nil)

func NewFileSystem(root string, readonly bool, mode string) (*FileSystem, error) {
	root = strings.TrimPrefix(root, "./")
	var f afero.Fs
	switch mode {
	case "memory":
		f = MemFS()
	case "file":
		f = OSFS()
	case "":
		switch runtime.GOOS {
		case "js":
			f = MemFS()
		default:
			f = OSFS()
		}
	default:
		return nil, errors.Errorf("invalid filesystem mode [%s]", mode)
	}
	if readonly {
		f = afero.NewReadOnlyFs(f)
	}
	return &FileSystem{root: root, f: f}, nil
}

func (f *FileSystem) getPath(ss ...string) string {
	s := util.StringFilePath(ss...)
	if strings.HasPrefix(s, f.root) {
		return s
	}
	return util.StringFilePath(f.root, s)
}

func (f *FileSystem) Root() string {
	return f.root
}

func (f *FileSystem) Clone() FileLoader {
	ret, _ := NewFileSystem(f.root, f.ReadOnly, f.Mode)
	return ret
}

func (f *FileSystem) Stat(path string) (*FileInfo, error) {
	p := f.getPath(path)
	s, err := f.f.Stat(p)
	if err != nil {
		return nil, err
	}
	return FileInfoFromFS(s), nil
}

func (f *FileSystem) SetMode(path string, mode FileMode) error {
	p := f.getPath(path)
	return f.f.Chmod(p, mode.ToFS())
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
	if err := f.f.MkdirAll(p, DirectoryMode.ToFS()); err != nil {
		return errors.Wrapf(err, "unable to create data directory [%s]", p)
	}
	return nil
}

func (f *FileSystem) String() string {
	return fmt.Sprintf("fs://%s", f.root)
}
