package svg

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const (
	bgMsg   = "magick -background %q -fill black -font %q -pointsize 48 -size 640x400 -gravity center label:%q background.png"
	bg2XMsg = "magick -background %q -fill black -font %q -pointsize 72 -size 1280x800 -gravity center label:%q background@2x.png"
	icoMsg  = "magick -density 1000 -background none logo.png -define icon:auto-resize=128,64,32 favicon.ico"
	noBG    = "inkscape -w %d -h %d logo.svg -o %s"
	pngMsg  = "inkscape -w %d -h %d logo.svg -o %s"
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

func RefreshAppIcon(ctx context.Context, prj *project.Project, fs filesystem.FileLoader, logger util.Logger) error {
	origB, err := fs.ReadFile(svgRoot + prj.Icon + util.ExtSVG)
	if err != nil {
		return errors.Wrap(err, "unable to read initial ["+prj.Icon+util.ExtSVG+"]")
	}
	x := &SVG{Key: prj.Icon, Markup: string(origB)}
	return SetAppIcon(ctx, prj, fs, x, logger)
}

func SetAppIcon(ctx context.Context, prj *project.Project, fs filesystem.FileLoader, x *SVG, logger util.Logger) error {
	orig, origColored := cleanMarkup(x.Markup, prj.Theme.Dark.NavBackground)

	var wg sync.WaitGroup
	wg.Add(4)
	var errs []error

	queue := func(f func(ctx context.Context, prj *project.Project, orig string, fs filesystem.FileLoader, logger util.Logger) error) {
		err := f(ctx, prj, origColored, fs, logger)
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
	err := fs.WriteFile("client/src/svg/app.svg", []byte(appIconContent), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write initial [app.svg]")
	}

	_, err = Build(fs, logger, prj)
	if err != nil {
		return errors.Wrap(err, "unable to build SVG library")
	}

	return nil
}
