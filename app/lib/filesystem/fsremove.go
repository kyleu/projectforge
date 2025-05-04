package filesystem

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func (f *FileSystem) Remove(path string, logger util.Logger) error {
	p := f.getPath(path)
	if logger != nil {
		logger.Warnf("removing file at path [%s]", p)
	}
	if err := f.f.Remove(p); err != nil {
		return errors.Wrapf(err, "error removing file [%s]", path)
	}
	return nil
}

func (f *FileSystem) RemoveRecursive(path string, logger util.Logger) error {
	if !f.Exists(path) {
		return nil
	}
	p := f.getPath(path)
	s, err := f.f.Stat(p)
	if err != nil {
		return errors.Wrapf(err, "unable to stat file [%s]", path)
	}
	if s.IsDir() {
		dir, e := f.f.Open(p)
		if e != nil && logger != nil {
			logger.Warnf("cannot open path [%s] for removal: %+v", path, e)
			return nil
		}
		if dir == nil {
			return errors.Errorf("path [%s] does not exist", path)
		}
		files, e := dir.Readdir(0)
		if e != nil && logger != nil {
			logger.Warnf("cannot read path [%s] for removal: %+v", path, e)
			return nil
		}
		for _, file := range files {
			err = f.RemoveRecursive(util.StringFilePath(path, file.Name()), logger)
			if err != nil {
				return err
			}
		}
	}
	err = f.f.Remove(p)
	if err != nil {
		return errors.Wrapf(err, "unable to remove file [%s]", path)
	}
	return nil
}
