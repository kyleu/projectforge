package svg

import (
	"bytes"
	"context"
	"fmt"
	"image"
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

	const macOSLogoPath = "tools/desktop/template/darwin/icons.iconset/logo.svg"
	macOSPath := filepath.Join(fs.Root(), "tools", "desktop", "template", "darwin", "icons.iconset")
	err := fs.WriteFile(macOSLogoPath, []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write temporary macOS [logo.svg]")
	}
	macOSResize(1024, "logo.png", macOSPath)
	pngBytes, err := fs.ReadFile(filepath.Join(macOSPath, "logo.png"))
	if err != nil {
		return errors.Wrap(err, "unable to read generated macOS [logo.png]")
	}
	err = fs.Remove(macOSLogoPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary macOS [logo.svg]")
	}
	templatePath := filepath.Join(fs.Root(), "tools", "desktop", "template", "darwin")
	if !fs.Exists(templatePath) {
		_ = fs.CreateDirectory(templatePath)
	}
	img, _, err := image.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		return errors.Wrap(err, "unable to decode macOS [logo.png]")
	}
	icnsBytes, err := iconset(img)
	if err != nil {
		return errors.Wrap(err, "unable to create [icons.iconset]")
	}
	err = fs.WriteFile(filepath.Join(templatePath, "icons.icns"), icnsBytes, filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write [icons.icns]")
	}
	err = macOSBackgrounds(ctx, fs, templatePath, prj, logger)
	if err != nil {
		return err
	}
	return fs.RemoveRecursive(macOSPath, logger)
}

func macOSBackgrounds(ctx context.Context, fs filesystem.FileLoader, templatePath string, prj *project.Project, logger util.Logger) error {
	if !fs.Exists(filepath.Join(templatePath, "background.png")) {
		cmd := fmt.Sprintf(bgMsg, prj.Theme.Light.NavBackground, "Helvetica-Neue", prj.Title())
		err := proc(ctx, cmd, templatePath, logger)
		if err != nil {
			return errors.Wrap(err, "unable to generate background image")
		}
	}
	if !fs.Exists(filepath.Join(templatePath, "background@2x.png")) {
		cmd := fmt.Sprintf(bg2XMsg, prj.Theme.Light.NavBackground, "Helvetica-Neue", prj.Title())
		err := proc(ctx, cmd, templatePath, logger)
		if err != nil {
			return errors.Wrap(err, "unable to generate background image")
		}
	}
	return nil
}
