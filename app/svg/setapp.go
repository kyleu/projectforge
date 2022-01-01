package svg

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kyleu/projectforge/app/lib/filesystem"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	pngMsg = "convert -density 1000 -resize %dx%d -define png:exclude-chunks=date,time logo.svg %s"
	noBG   = "convert -density 1000 -background none -resize %dx%d -define png:exclude-chunks=date,time logo.svg %s"
)

func proc(cmd string, path string) error {
	exit, out, err := util.RunProcessSimple(cmd, path)
	if err != nil {
		return errors.Wrap(err, "unable to convert [logo.png]")
	}
	if exit != 0 {
		return errors.Wrapf(errors.Errorf("bad output: %s", out), "unable to convert [logo.png], exit code [%d]", exit)
	}
	return nil
}

func SetAppIcon(prj *project.Project, fs filesystem.FileLoader, x *SVG, logger *zap.SugaredLogger) error {
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

	queue := func(f func(prj *project.Project, orig string, fs filesystem.FileLoader, logger *zap.SugaredLogger) error) {
		err := f(prj, orig, fs, logger)
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

	_, err = Build(fs)
	if err != nil {
		return errors.Wrap(err, "unable to build SVG library")
	}

	return nil
}

func webAssets(prj *project.Project, orig string, fs filesystem.FileLoader, logger *zap.SugaredLogger) error {
	webResize := func(size int, fn string, p string) {
		err := proc(fmt.Sprintf(noBG, size, size, fn), p)
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
	err = proc("convert -density 1000 -background none logo.svg -define icon:auto-resize=128,64,32 favicon.ico", webPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.ico]")
	}
	return nil
}
