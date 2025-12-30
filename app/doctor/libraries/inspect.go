package libraries

import (
	"context"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func onInspect(ctx context.Context, lib *Library, fs filesystem.FileLoader, logger util.Logger) (*Result, error) {
	ret := NewResult(lib, "inspect")
	files, err := fs.ListFilesRecursive(".", nil, logger)
	if err != nil {
		return nil, err
	}
	ret.AddMessage("Inspecting [%d] files in library [%s]...", len(files), lib.Key)
	return ret, nil
}
