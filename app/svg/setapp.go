package svg

import (
	"path/filepath"
	"strings"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func SetAppIcon(fs filesystem.FileLoader, x *SVG) error {
	path := filepath.Join(fs.Root(), "assets")

	orig := x.Markup
	for strings.Contains(orig, "<!--") {
		startIdx := strings.Index(orig, "<!--")
		endIdx := strings.Index(orig, "-->")
		if endIdx == -1 {
			break
		}
		orig = strings.TrimPrefix(orig[:startIdx]+orig[endIdx+3:], "\n")
	}

	proc := func(cmd string) error {
		exit, out, err := util.RunProcessSimple(cmd, path)
		if err != nil {
			return errors.Wrap(err, "unable to convert [logo.png]")
		}
		if exit != 0 {
			return errors.Wrapf(errors.Errorf("bad output: %s", out), "unable to convert [logo.png], exit code [%d]", exit)
		}
		return nil
	}

	err := fs.WriteFile("assets/logo.svg", []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write initial [logo.svg]")
	}

	err = proc("convert -density 1000 -background none -resize 256x256 logo.svg logo.png")
	if err != nil {
		return errors.Wrap(err, "unable to convert [logo.png]")
	}

	err = proc("convert -density 1000 -background none -resize 64x64 logo.svg favicon.png")
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.png]")
	}

	err = proc("convert -density 1000 -background none logo.svg -define icon:auto-resize=128,64,32 favicon.ico")
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.ico]")
	}

	appIconContent := strings.ReplaceAll(orig, "svg-"+x.Key, "svg-app")
	appIconContent = "<!-- $PF_IGNORE$ -->\n" + appIconContent
	err = fs.WriteFile("client/src/svg/app.svg", []byte(appIconContent), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write initial [logo.svg]")
	}

	_, err = Build(fs)
	if err != nil {
		return errors.Wrap(err, "unable to build SVG library")
	}

	return nil
}
