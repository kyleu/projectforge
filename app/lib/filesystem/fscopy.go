// Package filesystem - Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"path/filepath"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func (f *FileSystem) CopyFile(src string, tgt string) error {
	if targetExists := f.Exists(tgt); targetExists {
		return errors.Errorf("file [%s] exists, will not overwrite", tgt)
	}

	input, err := f.ReadFile(src)
	if err != nil {
		return err
	}

	var mode FileMode
	if stat, e := f.Stat(src); stat != nil && e == nil {
		mode = stat.Mode
	}

	err = f.WriteFile(tgt, input, mode, false)
	return err
}

func (f *FileSystem) CopyRecursive(src string, tgt string, ignore []string, logger util.Logger) error {
	srcFiles, err := f.ListFilesRecursive(src, ignore, logger)
	if err != nil {
		return err
	}

	for _, path := range srcFiles {
		err := f.CopyFile(filepath.Join(src, path), filepath.Join(tgt, path))
		if err != nil {
			return errors.Wrapf(err, "error copying [%s]", path)
		}
	}
	return nil
}

func (f *FileSystem) Move(src string, tgt string) error {
	sp := f.getPath(src)
	if sourceExists := f.Exists(sp); !sourceExists {
		return errors.Errorf("source file [%s] does not exist, can't move", sp)
	}

	tp := f.getPath(tgt)
	if targetExists := f.Exists(tp); targetExists {
		return errors.Errorf("target file [%s] exists, will not overwrite", tp)
	}

	if err := f.f.Rename(sp, tp); err != nil {
		return errors.Wrapf(err, "error renaming [%s] to [%s]", sp, tp)
	}

	return nil
}
