package svg

import (
	"fmt"
	"path/filepath"

	"github.com/kyleu/projectforge/app/lib/filesystem"
	"github.com/kyleu/projectforge/app/project"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func macOSAssets(prj *project.Project, orig string, fs filesystem.FileLoader, logger *zap.SugaredLogger) error {
	if prj.Build == nil || (!prj.Build.Desktop) {
		return nil
	}
	macOSResize := func(size int, fn string, p string) {
		err := proc(fmt.Sprintf(pngMsg, size, size, fn), p)
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
	err = fs.Remove(macOSLogoPath)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary macOS [logo.svg]")
	}
	err = proc("iconutil --convert icns icons.iconset", filepath.Join(fs.Root(), "tools", "desktop", "template", "macos"))
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary macOS [logo.svg]")
	}
	return fs.RemoveRecursive(macOSPath)
}
