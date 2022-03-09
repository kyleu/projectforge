package svg

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
)

func iOSAssets(prj *project.Project, orig string, fs filesystem.FileLoader, logger *zap.SugaredLogger) error {
	if prj.Build == nil || (!prj.Build.IOS) {
		return nil
	}
	iOSResize := func(size int, fn string, p string) {
		err := proc(fmt.Sprintf(pngMsg, size, size, fn), p)
		if err != nil {
			logger.Warnf("error processing icon [%s]: %+v", fn, err)
		}
	}

	const iOSLogoPath = "tools/ios/Assets.xcassets/AppIcon.appiconset/logo.svg"
	iOSPath := filepath.Join(fs.Root(), "tools", "ios", "Assets.xcassets", "AppIcon.appiconset")
	err := fs.WriteFile(iOSLogoPath, []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write temporary iOS [logo.svg]")
	}
	iOSResize(16, "16.png", iOSPath)
	iOSResize(20, "20.png", iOSPath)
	iOSResize(29, "29.png", iOSPath)
	iOSResize(32, "32.png", iOSPath)
	iOSResize(40, "40.png", iOSPath)
	iOSResize(48, "48.png", iOSPath)
	iOSResize(50, "50.png", iOSPath)
	iOSResize(55, "55.png", iOSPath)
	iOSResize(57, "57.png", iOSPath)
	iOSResize(58, "58.png", iOSPath)
	iOSResize(60, "60.png", iOSPath)
	iOSResize(64, "64.png", iOSPath)
	iOSResize(72, "72.png", iOSPath)
	iOSResize(76, "76.png", iOSPath)
	iOSResize(80, "80.png", iOSPath)
	iOSResize(87, "87.png", iOSPath)
	iOSResize(88, "88.png", iOSPath)
	iOSResize(100, "100.png", iOSPath)
	iOSResize(114, "114.png", iOSPath)
	iOSResize(120, "120.png", iOSPath)
	iOSResize(128, "128.png", iOSPath)
	iOSResize(144, "144.png", iOSPath)
	iOSResize(152, "152.png", iOSPath)
	iOSResize(167, "167.png", iOSPath)
	iOSResize(172, "172.png", iOSPath)
	iOSResize(180, "180.png", iOSPath)
	iOSResize(196, "196.png", iOSPath)
	iOSResize(216, "216.png", iOSPath)
	iOSResize(256, "256.png", iOSPath)
	iOSResize(512, "512.png", iOSPath)
	iOSResize(1024, "1024.png", iOSPath)
	err = fs.Remove(iOSLogoPath)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary iOS [logo.svg]")
	}
	return nil
}
