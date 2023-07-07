// Content managed by Project Forge, see [projectforge.md] for details.
package download

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var (
	arms = []string{ArchARMV5, ArchARMV6, ArchARMV7}
	mips = []string{ArchMIPS64Hard, ArchMIPS64Soft, ArchMIPS64LEHard, ArchMIPS64LESoft, ArchMIPSHard, ArchMIPSSoft, ArchMIPSLEHard, ArchMIPSLESoft}
)

func GetLinks(version string) Links {
	if availableLinks == nil {
		availableLinks = calcDownloadLinks(version)
	}
	return availableLinks
}

func calcDownloadLinks(version string) Links {
	ret := Links{}
	add := func(url string, mode string, os string, arch string) {
		ret = append(ret, &Link{URL: url, Mode: mode, OS: os, Arch: arch})
	}
	addDefault := func(mode string, os string, arch string) {
		var u string
		switch mode {
		case ModeServer, ModeMobile:
			u = fmt.Sprintf("%s_%s_%s_%s.zip", util.AppKey, version, os, arch)
		case ModeDesktop:
			u = fmt.Sprintf("%s_%s_%s_%s_desktop.zip", util.AppKey, version, os, arch)
		}
		add(u, mode, os, arch)
	}
	addARMs := func(mode string, os string) {
		lo.ForEach(arms, func(arm string, _ int) {
			addDefault(mode, os, arm)
		})
	}
	addMIPS := func(mode string, os string) {
		lo.ForEach(mips, func(weird string, _ int) {
			addDefault(mode, os, weird)
		})
	}

	addDefault(ModeDesktop, OSMac, ArchUniversal)
	addDefault(ModeServer, OSMac, ArchAMD64)
	addDefault(ModeServer, OSMac, ArchARM64)
	addDefault(ModeServer, OSMac, ArchUniversal)
	addDefault(ModeDesktop, OSWindows, ArchAMD64)
	addDefault(ModeServer, OSWindows, ArchAMD64)
	addDefault(ModeServer, OSWindows, Arch386)
	addDefault(ModeServer, OSWindows, ArchARM64)
	addARMs(ModeServer, OSWindows)
	addDefault(ModeDesktop, OSLinux, ArchAMD64)
	addDefault(ModeServer, OSLinux, ArchAMD64)
	addDefault(ModeServer, OSLinux, Arch386)
	addDefault(ModeServer, OSLinux, ArchARM64)
	addARMs(ModeServer, OSLinux)
	addDefault(ModeServer, OSLinux, ArchPPC64)
	addDefault(ModeServer, OSLinux, ArchPPC64LE)
	addDefault(ModeServer, OSLinux, ArchRISCV64)
	addDefault(ModeServer, OSLinux, ArchS390X)
	addMIPS(ModeServer, OSLinux)
	addDefault(ModeMobile, OSAndroid, "apk")
	addDefault(ModeMobile, OSAndroid, "aar")
	addDefault(ModeServer, OSDragonfly, ArchAMD64)
	addDefault(ModeServer, OSFreeBSD, ArchAMD64)
	addDefault(ModeServer, OSFreeBSD, Arch386)
	addDefault(ModeServer, OSFreeBSD, ArchARM64)
	addARMs(ModeServer, OSFreeBSD)
	addDefault(ModeServer, OSIllumos, ArchAMD64)
	addDefault(ModeMobile, OSIOS, "app")
	addDefault(ModeMobile, OSIOS, "framework")
	addDefault(ModeServer, OSJS, ArchWASM)
	addDefault(ModeServer, OSNetBSD, ArchAMD64)
	addDefault(ModeServer, OSNetBSD, Arch386)
	addDefault(ModeServer, OSNetBSD, ArchARMV7)
	addDefault(ModeServer, OSOpenBSD, ArchAMD64)
	addDefault(ModeServer, OSOpenBSD, ArchARM64)
	addDefault(ModeServer, OSOpenBSD, Arch386)
	addARMs(ModeServer, OSOpenBSD)
	addDefault(ModeServer, OSSolaris, ArchAMD64)

	return ret
}
