// Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

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
