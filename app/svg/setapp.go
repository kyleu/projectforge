package svg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func SetAppIcon(fs filesystem.FileLoader, x *SVG, logger *zap.SugaredLogger) error {
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

	logoResize := func(size int, fn string, p string) {
		msg := "convert -density 1000 -background none -resize %dx%d logo.svg %s"
		err := proc(fmt.Sprintf(msg, size, size, fn), p)
		if err != nil {
			logger.Warnf("error processing icon [%s]: %+v", fn, err)
		}
	}

	iOSResize := func(size int, fn string, p string) {
		msg := "convert -density 1000 -background #888888 -resize %dx%d logo.svg %s"
		err := proc(fmt.Sprintf(msg, size, size, fn), p)
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
	logoResize(256, "logo.png", webPath)
	logoResize(64, "favicon.png", webPath)
	err = proc("convert -density 1000 -background none logo.svg -define icon:auto-resize=128,64,32 favicon.ico", webPath)
	if err != nil {
		return errors.Wrap(err, "unable to convert [favicon.ico]")
	}

	// Android assets
	androidLogoPath := "tools/android/app/src/main/res/logo.svg"
	androidPath := filepath.Join(fs.Root(), "tools/android/app/src/main/res")
	err = fs.WriteFile(androidLogoPath, []byte(orig), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write temporary Android [logo.svg]")
	}
	logoResize(48, "mipmap-mdpi/ic_launcher.png", androidPath)
	logoResize(72, "mipmap-hdpi/ic_launcher.png", androidPath)
	logoResize(96, "mipmap-xhdpi/ic_launcher.png", androidPath)
	logoResize(144, "mipmap-xxhdpi/ic_launcher.png", androidPath)
	logoResize(192, "mipmap-xxxhdpi/ic_launcher.png", androidPath)
	err = fs.Remove(androidLogoPath)
	if err != nil {
		return errors.Wrap(err, "unable to remove temporary Android [logo.svg]")
	}

	// iOS assets
	iOSLogoPath := "tools/ios/Assets.xcassets/AppIcon.appiconset/logo.svg"
	iOSPath := filepath.Join(fs.Root(), "tools/ios/Assets.xcassets/AppIcon.appiconset")
	err = fs.WriteFile(iOSLogoPath, []byte(orig), filesystem.DefaultMode, true)
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
