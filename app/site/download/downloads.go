package download

import (
	"fmt"

	"github.com/kyleu/projectforge/app/util"
)

func GetLinks(version string) Links {
	if availableLinks == nil {
		availableLinks = calcDownloadLinks(version)
	}
	return availableLinks
}

// nolint
func calcDownloadLinks(version string) Links {
	ret := Links{}
	add := func(url string, mode string, os string, arch string) {
		ret = append(ret, &Link{URL: url, Mode: mode, OS: os, Arch: arch})
	}
	addDefault := func(mode string, os string, arch string) {
		var u string
		switch mode {
		case modeServer, modeMobile:
			u = fmt.Sprintf("%s_%s_%s_%s.zip", util.AppKey, version, os, arch)
		case modeDesktop:
			u = fmt.Sprintf("%s_%s_%s_%s_desktop.zip", util.AppKey, version, os, arch)
		}
		add(u, mode, os, arch)
	}
	addARMs := func(mode string, os string) {
		for _, arm := range []string{archARMV5, archARMV6, archARMV7} {
			addDefault(mode, os, arm)
		}
	}
	addMIPS := func(mode string, os string) {
		for _, weird := range []string{archMIPS64Hard, archMIPS64Soft, archMIPS64LEHard, archMIPS64LESoft, archMIPSHard, archMIPSSoft, archMIPSLEHard, archMIPSLESoft} {
			addDefault(mode, os, weird)
		}
	}

	addDefault(modeDesktop, osMac, archAMD64)
	addDefault(modeServer, osMac, archAMD64)
	addDefault(modeServer, osMac, archARM64)
	addDefault(modeDesktop, osWindows, archAMD64)
	addDefault(modeServer, osWindows, archAMD64)
	addDefault(modeServer, osWindows, archI386)
	addDefault(modeServer, osWindows, archARM64)
	addARMs(modeServer, osWindows)
	addDefault(modeDesktop, osLinux, archAMD64)
	addDefault(modeServer, osLinux, archAMD64)
	addDefault(modeServer, osLinux, archI386)
	addDefault(modeServer, osLinux, archARM64)
	addARMs(modeServer, osLinux)
	addDefault(modeServer, osLinux, archPPC64)
	addDefault(modeServer, osLinux, archRISCV64)
	addDefault(modeServer, osLinux, archS390X)
	addMIPS(modeServer, osLinux)
	addDefault(modeServer, osAIX, archPPC64)
	addDefault(modeMobile, osAndroid, "apk")
	addDefault(modeMobile, osAndroid, "aar")
	addDefault(modeServer, osDragonfly, archAMD64)
	addDefault(modeServer, osFreeBSD, archAMD64)
	addDefault(modeServer, osFreeBSD, archI386)
	addDefault(modeServer, osFreeBSD, archARM64)
	addARMs(modeServer, osFreeBSD)
	addDefault(modeServer, osIllumos, archAMD64)
	addDefault(modeMobile, osIOS, "app")
	addDefault(modeMobile, osIOS, "framework")
	addDefault(modeServer, osJS, archWASM)
	addDefault(modeServer, osNetBSD, archAMD64)
	addDefault(modeServer, osNetBSD, archI386)
	addDefault(modeServer, osNetBSD, archARMV7)
	addDefault(modeServer, osOpenBSD, archAMD64)
	addDefault(modeServer, osOpenBSD, archARM64)
	addDefault(modeServer, osOpenBSD, archI386)
	addARMs(modeServer, osOpenBSD)
	addDefault(modeServer, osPlan9, archAMD64)
	addDefault(modeServer, osPlan9, archI386)
	addARMs(modeServer, osPlan9)
	addDefault(modeServer, osSolaris, archAMD64)

	return ret
}
