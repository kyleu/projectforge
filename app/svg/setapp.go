package svg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const (
	pngMsg = "convert -density 1000 -resize %dx%d -define png:exclude-chunks=date,time logo.svg %s"
	noBG   = "convert -density 1000 -background none -resize %dx%d -define png:exclude-chunks=date,time logo.svg %s"
)

func proc(ctx context.Context, cmd string, path string, logger util.Logger) error {
	exit, out, err := telemetry.RunProcessSimple(ctx, cmd, path, logger)
	if err != nil {
		return errors.Wrap(err, "unable to convert [logo.png]")
	}
	if exit != 0 {
		return errors.Wrapf(errors.Errorf("bad output: %s", out), "unable to convert [logo.png], exit code [%d]", exit)
	}
	return nil
}

func SetAppIcon(ctx context.Context, prj *project.Project, fs filesystem.FileLoader, x *SVG, logger util.Logger) error {
	orig := x.Markup
	for strings.Contains(orig, "<!--") {
		startIdx := strings.Index(orig, "<!--")
		endIdx := strings.Index(orig, "-->")
		if endIdx == -1 {
			break
		}
		orig = strings.TrimPrefix(orig[:startIdx]+orig[endIdx+3:], "\n")
	}

	var wg sync.WaitGroup
	wg.Add(4)
	var errs []error

	queue := func(f func(ctx context.Context, prj *project.Project, orig string, fs filesystem.FileLoader, logger util.Logger) error) {
		err := f(ctx, prj, orig, fs, logger)
		if err != nil {
			errs = append(errs, err)
		}
		wg.Done()
	}
	queue(webAssets)
	queue(androidAssets)
	queue(iOSAssets)
	queue(macOSAssets)

	wg.Wait()
	if len(errs) > 0 {
		return errs[0]
	}

	// app icon
	appIconContent := strings.ReplaceAll(orig, "svg-"+x.Key, "svg-app")
	appIconContent = "<!-- $PF_IGNORE$ -->\n" + appIconContent
	err := fs.WriteFile("client/src/svg/app.svg", []byte(appIconContent), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write initial [logo.svg]")
	}

	_, err = Build(fs, logger)
	if err != nil {
		return errors.Wrap(err, "unable to build SVG library")
	}

	return nil
}

func webAssets(ctx context.Context, prj *project.Project, orig string, fs filesystem.FileLoader, logger util.Logger) error {
	webResize := func(size int, fn string, p string) {
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
	err = proc(ctx, cmd, webPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.ico]")
	}
	return nil
}
