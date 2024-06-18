package filesystem

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func (f *FileSystem) UnzipToDir(src string, dest string) (*util.OrderedMap[int64], error) {
	r, err := f.FileReader(src)
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(r, int64(f.Size(src)))
	if err != nil {
		return nil, err
	}

	err = f.CreateDirectory(dest)
	if err != nil {
		return nil, err
	}

	extractAndWriteFile := func(zf *zip.File) (int64, error) {
		rc, err := zf.Open()
		if err != nil {
			return 0, err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, zf.Name) // #nosec G305

		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return 0, errors.Errorf("illegal file path: %s", path)
		}

		var sz int64
		if zf.FileInfo().IsDir() {
			sz = -1
			_ = f.CreateDirectory(path)
		} else {
			_ = f.CreateDirectory(filepath.Dir(path))
			df, err := f.FileWriter(path, true, false)
			if err != nil {
				return 0, err
			}
			defer func() {
				_ = df.Close()
			}()
			sz, err = io.Copy(df, rc) // #nosec G110
			if err != nil {
				return 0, err
			}
		}
		return sz, nil
	}

	ret := util.NewOrderedMap[int64](false, len(zr.File))
	for _, f := range zr.File {
		fsz, err := extractAndWriteFile(f)
		if err != nil {
			return nil, err
		}
		ret.Append(f.Name, fsz)
	}
	return ret, nil
}
