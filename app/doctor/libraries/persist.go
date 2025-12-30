package libraries

import (
	"context"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func onPersist(ctx context.Context, lib *Library, libFS filesystem.FileLoader, logger util.Logger) (*Result, error) {
	moduleFS, err := fsForModules()
	if err != nil {
		return nil, err
	}
	ret := NewResult(lib, "persist")
	for src, tgt := range lib.Files {
		if err := persistFiles(moduleFS, src, libFS, tgt, ret, logger); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func persistFiles(modFS filesystem.FileLoader, src string, libFS filesystem.FileLoader, tgt string, r *Result, logger util.Logger) error {
	r.AddMessage("persisting files from [%s] to [%s]", src, tgt)

	srcFile, err := modFS.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "file [%s] not found", src)
	}
	if srcFile.IsDir {
		kids, err := modFS.ListFilesRecursive(src, nil, logger)
		if err != nil {
			return errors.Wrapf(err, "failed to read directory [%s]", src)
		}
		for _, kid := range kids {
			if err := persistFiles(modFS, util.StringFilePath(src, kid), libFS, util.StringFilePath(tgt, kid), r, logger); err != nil {
				return err
			}
		}
	} else {
		r.AddMessage("persisting file [%s] (%s) to [%s]", src, util.ByteSizeSI(srcFile.Size), tgt)
		b, err := modFS.ReadFile(src)
		if err != nil {
			return errors.Wrapf(err, "failed to read file [%s]", src)
		}
		if err := libFS.WriteFile(tgt, b, srcFile.Mode, true); err != nil {
			return errors.Wrapf(err, "failed to write file [%s]", tgt)
		}
	}
	return nil
}
