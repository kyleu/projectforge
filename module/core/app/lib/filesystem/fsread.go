package filesystem

import (
	"os"

	"github.com/pkg/errors"
)

func (f *FileSystem) PeekFile(path string, maxSize int) ([]byte, error) {
	file, err := os.Open(f.getPath(path))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open file [%s]", path)
	}
	defer func() { _ = file.Close() }()

	b := make([]byte, maxSize)

	size, err := file.Read(b)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file [%s]", path)
	}
	return b[:size], nil
}

func (f *FileSystem) ReadFile(path string) ([]byte, error) {
	b, err := os.ReadFile(f.getPath(path))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file [%s]", path)
	}
	return b, nil
}
