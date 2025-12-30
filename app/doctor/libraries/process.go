package libraries

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const libraryPath = "tmp/libraries"

func Process(ctx context.Context, lib *Library, act string, logger util.Logger) (*Result, error) {
	ret := NewResult(lib, act)

	fs, err := fsFor(ctx, lib, true, logger)
	if err != nil {
		return nil, err
	}

	ret.Output = []any{"TODO", fs.Exists(".")}
	ret.Outcome = "WIP"

	return ret, nil
}

func fsFor(ctx context.Context, lib *Library, clean bool, logger util.Logger) (filesystem.FileLoader, error) {
	root, err := filesystem.NewFileSystem(libraryPath, false, "")
	if err != nil {
		return nil, err
	}
	if !root.Exists(".") {
		return nil, errors.Errorf("library path [%s] does not exist", libraryPath)
	}
	dirExists := root.Exists(lib.Key)
	if clean && dirExists {
		err = root.RemoveRecursive(lib.Key, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to clean library directory [%s]", lib.Key)
		}
	}
	if !dirExists {
		err = root.CreateDirectory("./" + lib.Key)
		if err != nil {
			return nil, err
		}
	}
	return filesystem.NewFileSystem(root.Root()+"/"+lib.Key, false, "")
}
