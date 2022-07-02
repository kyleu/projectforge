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

func androidAssets(ctx context.Context, prj *project.Project, orig string, fs filesystem.FileLoader, logger util.Logger) error {
	if prj.Build == nil || (!prj.Build.Android) {
		return nil
	}
	androidResize := func(size int, fn string, p string) {
		if x := filepath.Dir(filepath.Join(p, fn)); !fs.Exists(x) {
			_ = fs.CreateDirectory(x)
		}
		err := proc(ctx, fmt.Sprintf(noBG, size, size, fn), p, logger)
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
	err = fs.Remove(androidLogoPath, logger)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary Android [logo.svg]")
	}
	return nil
}
