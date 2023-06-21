package download

import (
	"fmt"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

var (
	arms = []string{archARMV5, archARMV6, archARMV7}
	mips = []string{archMIPS64Hard, archMIPS64Soft, archMIPS64LEHard, archMIPS64LESoft, archMIPSHard, archMIPSSoft, archMIPSLEHard, archMIPSLESoft}
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
		case modeServer, modeMobile:
			u = fmt.Sprintf("%s_%s_%s_%s.zip", util.AppKey, version, os, arch)
		case modeDesktop:
			u = fmt.Sprintf("%s_%s_%s_%s_desktop.zip", util.AppKey, version, os, arch)
		}
		add(u, mode, os, arch)
	}{{{ if .Build.HasArm }}}
	addARMs := func(mode string, os string) {
		lo.ForEach(arms, func(arm string, index int) {
			addDefault(mode, os, arm)
		})
	}{{{ end }}}{{{ if .Build.LinuxMIPS }}}
	addMIPS := func(mode string, os string) {
		lo.ForEach(mips, func(weird string, _ int) {
			addDefault(mode, os, weird)
		})
	}{{{ end }}}
{{{ if .Build.Desktop }}}
	addDefault(modeDesktop, osMac, archAMD64){{{ end }}}
	addDefault(modeServer, osMac, archAMD64)
	addDefault(modeServer, osMac, archARM64)
	addDefault(modeServer, osMac, archUniversal){{{ if .Build.Desktop }}}
	addDefault(modeDesktop, osWindows, archAMD64){{{ end }}}
	addDefault(modeServer, osWindows, archAMD64)
	addDefault(modeServer, osWindows, archI386){{{ if .Build.WindowsARM }}}
	addDefault(modeServer, osWindows, archARM64)
	addARMs(modeServer, osWindows){{{ end }}}{{{ if .Build.Desktop }}}
	addDefault(modeDesktop, osLinux, archAMD64){{{ end }}}
	addDefault(modeServer, osLinux, archAMD64)
	addDefault(modeServer, osLinux, archI386){{{ if .Build.LinuxARM }}}
	addDefault(modeServer, osLinux, archARM64)
	addARMs(modeServer, osLinux){{{ end }}}{{{ if .Build.LinuxOdd }}}
	addDefault(modeServer, osLinux, archPPC64)
	addDefault(modeServer, osLinux, archRISCV64)
	addDefault(modeServer, osLinux, archS390X){{{ end }}}{{{ if .Build.LinuxMIPS }}}
	addMIPS(modeServer, osLinux){{{ end }}}{{{ if .Build.Android }}}
	addDefault(modeMobile, osAndroid, "apk")
	addDefault(modeMobile, osAndroid, "aar"){{{ end }}}{{{ if .Build.Dragonfly }}}
	addDefault(modeServer, osDragonfly, archAMD64){{{ end }}}{{{ if .Build.FreeBSD }}}
	addDefault(modeServer, osFreeBSD, archAMD64)
	addDefault(modeServer, osFreeBSD, archI386)
	addDefault(modeServer, osFreeBSD, archARM64)
	addARMs(modeServer, osFreeBSD){{{ end }}}{{{ if .Build.Illumos }}}
	addDefault(modeServer, osIllumos, archAMD64){{{ end }}}{{{ if .Build.IOS }}}
	addDefault(modeMobile, osIOS, "app")
	addDefault(modeMobile, osIOS, "framework"){{{ end }}}{{{ if .Build.WASM }}}
	addDefault(modeServer, osJS, archWASM){{{ end }}}{{{ if .Build.NetBSD }}}
	addDefault(modeServer, osNetBSD, archAMD64)
	addDefault(modeServer, osNetBSD, archI386)
	addDefault(modeServer, osNetBSD, archARMV7){{{ end }}}{{{ if .Build.OpenBSD }}}
	addDefault(modeServer, osOpenBSD, archAMD64)
	addDefault(modeServer, osOpenBSD, archARM64)
	addDefault(modeServer, osOpenBSD, archI386)
	addARMs(modeServer, osOpenBSD){{{ end }}}{{{ if .Build.Plan9 }}}
	addDefault(modeServer, osPlan9, archAMD64)
	addDefault(modeServer, osPlan9, archI386)
	addARMs(modeServer, osPlan9){{{ end }}}{{{ if .Build.Solaris }}}
	addDefault(modeServer, osSolaris, archAMD64){{{ end }}}

	return ret
}
