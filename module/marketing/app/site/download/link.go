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
	case OSAndroid:
		return "Android"
	case OSDragonfly:
		return "Dragonfly"
	case OSFreeBSD:
		return "FreeBSD"
	case OSIllumos:
		return "Illumos"
	case OSIOS:
		return "iOS"
	case OSJS:
		return "JavaScript"
	case OSLinux:
		return "Linux"
	case OSMac:
		return "macOS"
	case OSNetBSD:
		return "NetBSD"
	case OSOpenBSD:
		return "OpenBSD"
	case OSPlan9:
		return "Plan9"
	case OSSolaris:
		return "Solaris"
	case OSWindows:
		return "Windows"
	}
	return "Unknown"
}

func (l *Link) Caveat(plat string) string {
	switch plat {
	case "desktop":
		switch l.OS {
		case OSMac:
			return "(Universal)"
		case OSWindows, OSLinux:
			return "(64-bit Intel only)"
		}
	}
	return ""
}

func (l *Link) OSIcon() string {
	switch l.OS {
	case OSMac, OSIOS:
		return "apple"
	default:
		return l.OS
	}
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
	ModeServer  = "server"
	ModeDesktop = "desktop"
	ModeMobile  = "mobile"

	OSAndroid   = "android"
	OSDragonfly = "dragonfly"
	OSFreeBSD   = "freebsd"
	OSIllumos   = "illumos"
	OSIOS       = "ios"
	OSJS        = "js"
	OSLinux     = "linux"
	OSMac       = "darwin"
	OSNetBSD    = "netbsd"
	OSOpenBSD   = "openbsd"
	OSPlan9     = "plan9"
	OSSolaris   = "solaris"
	OSWindows   = "windows"

	ArchAMD64        = "amd64"
	ArchARM64        = "arm64"
	ArchARMV5        = "armv5"
	ArchARMV6        = "armv6"
	ArchARMV7        = "armv7"
	Arch386          = "386"
	ArchMIPS64Hard   = "mips64_hardfloat"
	ArchMIPS64LEHard = "mips64le_hardfloat"
	ArchMIPS64LESoft = "mips64le_softfloat"
	ArchMIPS64Soft   = "mips64_softfloat"
	ArchMIPSHard     = "mips_hardfloat"
	ArchMIPSLEHard   = "mipsle_hardfloat"
	ArchMIPSLESoft   = "mipsle_softfloat"
	ArchMIPSSoft     = "mips_softfloat"
	ArchPPC64        = "ppc64"
	ArchRISCV64      = "riscv64"
	ArchS390X        = "s390x"
	ArchUniversal    = "all"
	ArchWASM         = "wasm"
)
