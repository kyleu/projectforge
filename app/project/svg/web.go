package svg

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func webAssets(ctx context.Context, prj *project.Project, orig string, fs filesystem.FileLoader, logger util.Logger) error {
	webResize := func(size int, fn string, p string) {
		if x := filepath.Dir(filepath.Join(p, fn)); !fs.Exists(x) {
			_ = fs.CreateDirectory(x)
		}
		err := proc(ctx, fmt.Sprintf(noBG, size, size, fn), p, logger)
		if err != nil {
			logger.Warnf("error processing icon [%s]: %+v", fn, err)
		}
	}

	// web assets
	webPath := filepath.Join(fs.Root(), "assets")
	err := fs.WriteFile("assets/logo.svg", []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write initial [logo.svg]")
	}
	webResize(256, "logo.png", webPath)
	webResize(64, "favicon.png", webPath)
	cmd := "convert -density 1000 -background none logo.svg -define icon:auto-resize=128,64,32 favicon.ico"
	if !fs.Exists(webPath) {
		_ = fs.CreateDirectory(webPath)
	}
	err = proc(ctx, cmd, webPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.ico]")
	}
	return nil
}
