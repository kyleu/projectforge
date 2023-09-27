// Package filesystem - Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import "github.com/pkg/errors"

func (f *FileSystem) PeekFile(path string, maxSize int) ([]byte, error) {
	file, err := f.f.Open(f.getPath(path))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open file [%s]", path)
	}
	defer func() { _ = file.Close() }()

	if stat, _ := file.Stat(); stat.Size() == 0 {
		return nil, nil
	}

	b := make([]byte, maxSize)

	size, err := file.Read(b)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file [%s]", path)
	}
	return b[:size], nil
}

func (f *FileSystem) Size(path string) int {
	file, err := f.f.Open(f.getPath(path))
	if err != nil {
		return -1
	}
	defer func() { _ = file.Close() }()

	stat, _ := file.Stat()
	if stat == nil {
		return -1
	}
	return int(stat.Size())
}

func (f *FileSystem) ReadFile(path string) ([]byte, error) {
	file, err := f.f.Open(f.getPath(path))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open file [%s]", path)
	}
	defer func() { _ = file.Close() }()

	stat, _ := file.Stat()
	if stat.Size() == 0 {
		return nil, nil
	}

	b := make([]byte, stat.Size())

	size, err := file.Read(b)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file [%s]", path)
	}
	if int64(size) != stat.Size() {
		return nil, errors.Wrapf(err, "file [%s] size [%d] does not match expected size [%d]", path, size, stat.Size())
	}
	return b, nil
}
