package download

import (
	"fmt"{{{ if .Build.HasArm }}}

	"github.com/samber/lo"{{{ end }}}

	"{{{ .Package }}}/app/util"
)
{{{ if .IsArmAndMips }}}
var (
	arms = []string{ArchARMV5, ArchARMV6, ArchARMV7}
	mips = []string{ArchMIPS64Hard, ArchMIPS64Soft, ArchMIPS64LEHard, ArchMIPS64LESoft, ArchMIPSHard, ArchMIPSSoft, ArchMIPSLEHard, ArchMIPSLESoft}
)
{{{ else if .Build.HasArm }}}var arms = []string{ArchARMV5, ArchARMV6, ArchARMV7}
{{{ else if .Build.LinuxMIPS }}}var mips = []string{ArchMIPS64Hard, ArchMIPS64Soft, ArchMIPS64LEHard, ArchMIPS64LESoft, ArchMIPSHard, ArchMIPSSoft, ArchMIPSLEHard, ArchMIPSLESoft}
{{{ end }}}
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
		case ModeServer:
			u = fmt.Sprintf("%s_%s_%s_%s.zip", util.AppKey, version, os, arch)
		case ModeMobile:
			u = fmt.Sprintf("%s_%s_%s_%s.zip", util.AppKey, version, os, arch)
		case ModeDesktop:
			u = fmt.Sprintf("%s_%s_%s_%s_desktop.zip", util.AppKey, version, os, arch)
		}
		add(u, mode, os, arch)
	}{{{ if .Build.HasArm }}}
	addARMs := func(mode string, os string) {
		lo.ForEach(arms, func(arm string, _ int) {
			addDefault(mode, os, arm)
		})
	}{{{ end }}}{{{ if .Build.LinuxMIPS }}}
	addMIPS := func(mode string, os string) {
		lo.ForEach(mips, func(weird string, _ int) {
			addDefault(mode, os, weird)
		})
	}{{{ end }}}
{{{ if .BuildDesktop }}}
	addDefault(ModeDesktop, OSMac, ArchUniversal){{{ end }}}
	addDefault(ModeServer, OSMac, ArchAMD64)
	addDefault(ModeServer, OSMac, ArchARM64)
	addDefault(ModeServer, OSMac, ArchUniversal){{{ if .BuildDesktop }}}
	addDefault(ModeDesktop, OSWindows, ArchAMD64){{{ end }}}
	addDefault(ModeServer, OSWindows, ArchAMD64)
	addDefault(ModeServer, OSWindows, Arch386){{{ if .Build.WindowsARM }}}
	addDefault(ModeServer, OSWindows, ArchARM64)
	addARMs(ModeServer, OSWindows){{{ end }}}{{{ if .BuildDesktop }}}
	addDefault(ModeDesktop, OSLinux, ArchAMD64){{{ end }}}
	addDefault(ModeServer, OSLinux, ArchAMD64)
	addDefault(ModeServer, OSLinux, Arch386){{{ if .Build.LinuxARM }}}
	addDefault(ModeServer, OSLinux, ArchARM64)
	addARMs(ModeServer, OSLinux){{{ end }}}{{{ if .Build.LinuxOdd }}}
	addDefault(ModeServer, OSLinux, ArchPPC64)
	addDefault(ModeServer, OSLinux, ArchPPC64LE)
	addDefault(ModeServer, OSLinux, ArchRISCV64)
	addDefault(ModeServer, OSLinux, ArchS390X)
	addDefault(ModeServer, OSLinux, ArchLoong64){{{ end }}}{{{ if .Build.LinuxMIPS }}}
	addMIPS(ModeServer, OSLinux){{{ end }}}{{{ if .BuildAndroid }}}
	addDefault(ModeMobile, OSAndroid, "apk")
	addDefault(ModeMobile, OSAndroid, "aar"){{{ end }}}{{{ if .Build.AIX }}}
	addDefault(ModeServer, OSAIX, ArchPPC64){{{ end }}}{{{ if .Build.Dragonfly }}}
	addDefault(ModeServer, OSDragonfly, ArchAMD64){{{ end }}}{{{ if .Build.FreeBSD }}}
	addDefault(ModeServer, OSFreeBSD, ArchAMD64)
	addDefault(ModeServer, OSFreeBSD, Arch386)
	addDefault(ModeServer, OSFreeBSD, ArchARM64)
	addARMs(ModeServer, OSFreeBSD){{{ end }}}{{{ if .Build.Illumos }}}
	addDefault(ModeServer, OSIllumos, ArchAMD64){{{ end }}}{{{ if .BuildIOS }}}
	addDefault(ModeMobile, OSIOS, "app")
	addDefault(ModeMobile, OSIOS, "framework"){{{ end }}}{{{ if .BuildWASM }}}
	addDefault(ModeServer, OSJS, ArchWASM){{{ end }}}{{{ if .Build.NetBSD }}}
	addDefault(ModeServer, OSNetBSD, ArchAMD64)
	addDefault(ModeServer, OSNetBSD, Arch386)
	addDefault(ModeServer, OSNetBSD, ArchARMV7){{{ end }}}{{{ if .Build.OpenBSD }}}
	addDefault(ModeServer, OSOpenBSD, ArchAMD64)
	addDefault(ModeServer, OSOpenBSD, ArchARM64)
	addDefault(ModeServer, OSOpenBSD, Arch386)
	addARMs(ModeServer, OSOpenBSD){{{ end }}}{{{ if .Build.Plan9 }}}
	addDefault(ModeServer, OSPlan9, ArchAMD64)
	addDefault(ModeServer, OSPlan9, Arch386)
	addARMs(ModeServer, OSPlan9){{{ end }}}{{{ if .Build.Solaris }}}
	addDefault(ModeServer, OSSolaris, ArchAMD64){{{ end }}}

	return ret
}
