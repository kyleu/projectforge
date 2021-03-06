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

func macOSAssets(ctx context.Context, prj *project.Project, orig string, fs filesystem.FileLoader, logger util.Logger) error {
	if prj.Build == nil || (!prj.Build.Desktop) {
		return nil
	}
	macOSResize := func(size int, fn string, p string) {
		if x := filepath.Dir(filepath.Join(p, fn)); !fs.Exists(x) {
			_ = fs.CreateDirectory(x)
		}
		err := proc(ctx, fmt.Sprintf(pngMsg, size, size, fn), p, logger)
		if err != nil {
			logger.Warnf("error processing icon [%s]: %+v", fn, err)
		}
	}

	const macOSLogoPath = "tools/desktop/template/macos/icons.iconset/logo.svg"
	macOSPath := filepath.Join(fs.Root(), "tools", "desktop", "template", "macos", "icons.iconset")
	err := fs.WriteFile(macOSLogoPath, []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write temporary macOS [logo.svg]")
	}
	macOSResize(16, "16.png", macOSPath)
	macOSResize(32, "icon_16x16@2x.png", macOSPath)
	macOSResize(32, "32.png", macOSPath)
	macOSResize(64, "icon_32x32@2x.png", macOSPath)
	macOSResize(128, "icon_128x128.png", macOSPath)
	macOSResize(256, "icon_128x128@2x.png", macOSPath)
	macOSResize(256, "icon_256x256.png", macOSPath)
	macOSResize(512, "icon_256x256@2x.png", macOSPath)
	macOSResize(512, "icon_512x512.png", macOSPath)
	macOSResize(1024, "icon_512x512@2x.png", macOSPath)
	err = fs.Remove(macOSLogoPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary macOS [logo.svg]")
	}
	templatePath := filepath.Join(fs.Root(), "tools", "desktop", "template", "macos")
	if !fs.Exists(templatePath) {
		_ = fs.CreateDirectory(templatePath)
	}
	err = proc(ctx, "iconutil --convert icns icons.iconset", templatePath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary macOS [logo.svg]")
	}
	if !fs.Exists(filepath.Join(templatePath, "background.png")) {
		cmdMsg := "magick -background %q -fill black -font %q -pointsize 48 -size 640x400 -gravity center label:%q background.png"
		cmd := fmt.Sprintf(cmdMsg, prj.Theme.Light.NavBackground, "Helvetica-Neue", prj.Title())
		err = proc(ctx, cmd, templatePath, logger)
		if err != nil {
			return errors.Wrap(err, "unable to generate background image")
		}
	}
	if !fs.Exists(filepath.Join(templatePath, "background@2x.png")) {
		cmdMsg := "magick -background %q -fill black -font %q -pointsize 72 -size 1280x800 -gravity center label:%q background@2x.png"
		cmd := fmt.Sprintf(cmdMsg, prj.Theme.Light.NavBackground, "Helvetica-Neue", prj.Title())
		err = proc(ctx, cmd, templatePath, logger)
		if err != nil {
			return errors.Wrap(err, "unable to generate background image")
		}
	}
	return fs.RemoveRecursive(macOSPath, logger)
}
