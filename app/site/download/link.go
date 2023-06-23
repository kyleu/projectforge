// Content managed by Project Forge, see [projectforge.md] for details.
package download

import (
	"github.com/samber/lo"
)

type Link struct {
	URL  string `json:"url"`
	Mode string `json:"mode"`
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

func (l *Link) OSString() string {
	switch l.OS {
	case osAndroid:
		return "Android"
	case osDragonfly:
		return "Dragonfly"
	case osFreeBSD:
		return "FreeBSD"
	case osIllumos:
		return "Illumos"
	case osIOS:
		return "iOS"
	case osJS:
		return "JavaScript"
	case osLinux:
		return "Linux"
	case osMac:
		return "macOS"
	case osNetBSD:
		return "NetBSD"
	case osOpenBSD:
		return "OpenBSD"
	case osPlan9:
		return "Plan9"
	case osSolaris:
		return "Solaris"
	case osWindows:
		return "Windows"
	}
	return "Unknown"
}

func (l *Link) OSIcon() string {
	if l.OS == osMac {
		return "apple"
	}
	return l.OS
}

type Links []*Link

func (l Links) Get(mode string, os string, arch string) *Link {
	return lo.FindOrElse(l, nil, func(link *Link) bool {
		return link.Mode == mode && link.OS == os && link.Arch == arch
	})
}

func (l Links) GetByModes(modes ...string) Links {
	return lo.Filter(l, func(link *Link, _ int) bool {
		return lo.Contains(modes, link.Mode)
	})
}

func (l Links) GetByOS(os string) Links {
	return lo.Filter(l, func(link *Link, _ int) bool {
		return link.OS == os
	})
}

var availableLinks Links

const (
	modeServer  = "server"
	modeDesktop = "desktop"
	modeMobile  = "mobile"

	osAndroid   = "android"
	osDragonfly = "dragonfly"
	osFreeBSD   = "freebsd"
	osIllumos   = "illumos"
	osIOS       = "ios"
	osJS        = "js"
	osLinux     = "linux"
	osMac       = "macos"
	osNetBSD    = "netbsd"
	osOpenBSD   = "openbsd"
	osPlan9     = "plan9"
	osSolaris   = "solaris"
	osWindows   = "windows"

	archAMD64        = "x86_64"
	archARM64        = "arm64"
	archARMV5        = "armv5"
	archARMV6        = "armv6"
	archARMV7        = "armv7"
	archI386         = "i386"
	archMIPS64Hard   = "mips64_hardfloat"
	archMIPS64LEHard = "mips64le_hardfloat"
	archMIPS64LESoft = "mips64le_softfloat"
	archMIPS64Soft   = "mips64_softfloat"
	archMIPSHard     = "mips_hardfloat"
	archMIPSLEHard   = "mipsle_hardfloat"
	archMIPSLESoft   = "mipsle_softfloat"
	archMIPSSoft     = "mips_softfloat"
	archPPC64        = "ppc64"
	archRISCV64      = "riscv64"
	archS390X        = "s390x"
	archUniversal    = "all"
	archWASM         = "wasm"
)
