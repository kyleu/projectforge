package svg

import (
	"path/filepath"
	"strings"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func SetAppIcon(fs filesystem.FileLoader, x *SVG) error {

	orig := x.Markup
	for strings.Contains(orig, "<!--") {
		startIdx := strings.Index(orig, "<!--")
		endIdx := strings.Index(orig, "-->")
		if endIdx == -1 {
			break
		}
		orig = strings.TrimPrefix(orig[:startIdx]+orig[endIdx+3:], "\n")
	}

	proc := func(cmd string, path string) error {
		exit, out, err := util.RunProcessSimple(cmd, path)
		if err != nil {
			return errors.Wrap(err, "unable to convert [logo.png]")
		}
		if exit != 0 {
			return errors.Wrapf(errors.Errorf("bad output: %s", out), "unable to convert [logo.png], exit code [%d]", exit)
		}
		return nil
	}

	// web assets
	webPath := filepath.Join(fs.Root(), "assets")
	err := fs.WriteFile("assets/logo.svg", []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write initial [logo.svg]")
	}
	err = proc("convert -density 1000 -background none -resize 256x256 logo.svg logo.png", webPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [logo.png]")
	}
	err = proc("convert -density 1000 -background none -resize 64x64 logo.svg favicon.png", webPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.png]")
	}
	err = proc("convert -density 1000 -background none logo.svg -define icon:auto-resize=128,64,32 favicon.ico", webPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.ico]")
	}

	// android assets
	tempPath := "tools/android/app/src/main/res/logo.svg"
	androidPath := filepath.Join(fs.Root(), "tools/android/app/src/main/res")
	err = fs.WriteFile(tempPath, []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write temporary [logo.svg]")
	}
	err = proc("convert -density 1000 -background none -resize 48x48 logo.svg mipmap-mdpi/ic_launcher.png", androidPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [mipmap-mdpi]")
	}
	err = proc("convert -density 1000 -background none -resize 72x72 logo.svg mipmap-hdpi/ic_launcher.png", androidPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [mipmap-hdpi]")
	}
	err = proc("convert -density 1000 -background none -resize 96x96 logo.svg mipmap-xhdpi/ic_launcher.png", androidPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [mipmap-xhdpi]")
	}
	err = proc("convert -density 1000 -background none -resize 144x144 logo.svg mipmap-xxhdpi/ic_launcher.png", androidPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [mipmap-xxhdpi]")
	}
	err = proc("convert -density 1000 -background none -resize 192x192 logo.svg mipmap-xxxhdpi/ic_launcher.png", androidPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [mipmap-xxxhdpi]")
	}
	err = fs.Remove(tempPath)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary [logo.svg]")
	}

	// app icon
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
