package filesystem

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type Reader interface {
	io.Reader
	io.ReaderAt
	io.Closer
}

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
	size, err := func() (int64, error) {
		file, err := f.f.Open(f.getPath(path))
		if err != nil {
			return 0, errors.Wrapf(err, "unable to open file [%s]", path)
		}
		defer func() { _ = file.Close() }()

		stat, _ := file.Stat()
		return stat.Size(), nil
	}()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(f.getPath(path))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file [%s]", path)
	}
	if int64(len(b)) != size {
		return nil, errors.Errorf("file [%s] size [%d] does not match expected size [%d]", path, len(b), size)
	}
	return b, nil
}

func (f *FileSystem) FileReader(fn string) (Reader, error) {
	p := f.getPath(fn)
	return f.f.Open(p)
}
