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

func (l *Link) ArchString() string {
	switch l.Arch {
	case ArchAMD64:
		return "Intel 64-bit"
	case ArchARM64:
		return "ARM64"
	case ArchARMV5:
		return "ARM v5"
	case ArchARMV6:
		return "ARM v6"
	case ArchARMV7:
		return "ARM v7"
	case Arch386:
		return "Intel 32-bit"
	case ArchMIPS64Hard:
		return "MIPS64 (hardfloat)"
	case ArchMIPS64LEHard:
		return "MIPS64 (le, hardfloat)"
	case ArchMIPS64LESoft:
		return "MIPS64 (le, softfloat)"
	case ArchMIPS64Soft:
		return "MIPS64 (softfloat)"
	case ArchMIPSHard:
		return "MIPS (hardfloat)"
	case ArchMIPSLEHard:
		return "MIPS (le, hardfloat)"
	case ArchMIPSLESoft:
		return "MIPS (le, softfloat)"
	case ArchMIPSSoft:
		return "MIPS (softfloat)"
	case ArchLoong64:
		return "Loong64"
	case ArchPPC64:
		return "PPC64"
	case ArchPPC64LE:
		return "PPC64 (le)"
	case ArchRISCV64:
		return "RISC-V 64-bit"
	case ArchS390X:
		return "S390x"
	case ArchUniversal:
		return "All (Universal)"
	case ArchWASM:
		return "WASM"
	case "apk", "app":
		return "Application"
	case "aar":
		return "AAR Library"
	case "framework":
		return "Xcode Framework"
	default:
		return l.Arch
	}
}

func (l *Link) OSString() string {
	switch l.OS {
	case OSAIX:
		return "AIX"
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

func (l *Link) OSIcon() string {
	switch l.OS {
	case OSAndroid:
		return OSAndroid
	case OSLinux:
		return OSLinux
	case OSMac, OSIOS:
		return "apple"
	case OSWindows:
		return OSWindows
	default:
		return "star"
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

	OSAIX       = "aix"
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
	ArchLoong64      = "loong64"
	ArchPPC64        = "ppc64"
	ArchPPC64LE      = "ppc64le"
	ArchRISCV64      = "riscv64"
	ArchS390X        = "s390x"
	ArchUniversal    = "all"
	ArchWASM         = "wasm"
)
