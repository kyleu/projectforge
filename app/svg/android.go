package svg

import (
	"fmt"
	"path/filepath"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func androidAssets(prj *project.Project, orig string, fs filesystem.FileLoader, logger *zap.SugaredLogger) error {
	if prj.Build == nil || (!prj.Build.Android) {
		return nil
	}
	androidResize := func(size int, fn string, p string) {
		err := proc(fmt.Sprintf(noBG, size, size, fn), p)
		if err != nil {
			logger.Warnf("error processing icon [%s]: %+v", fn, err)
		}
	}

	const androidLogoPath = "tools/android/app/src/main/res/logo.svg"
	androidPath := filepath.Join(fs.Root(), "tools", "android", "app", "src", "main", "res")
	err := fs.WriteFile(androidLogoPath, []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write temporary Android [logo.svg]")
	}
	androidResize(48, "mipmap-mdpi/ic_launcher.png", androidPath)
	androidResize(72, "mipmap-hdpi/ic_launcher.png", androidPath)
	androidResize(96, "mipmap-xhdpi/ic_launcher.png", androidPath)
	androidResize(144, "mipmap-xxhdpi/ic_launcher.png", androidPath)
	androidResize(192, "mipmap-xxxhdpi/ic_launcher.png", androidPath)
	err = fs.Remove(androidLogoPath)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary Android [logo.svg]")
	}
	return nil
}
