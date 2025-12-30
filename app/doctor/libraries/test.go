package libraries

import (
	"context"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func onTest(ctx context.Context, lib *Library, fs filesystem.FileLoader, logger util.Logger) (*Result, error) {
	ret := NewResult(lib, "test")
	return ret, nil
}
