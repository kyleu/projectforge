package download

type Link struct {
	URL  string `json:"url"`
	Mode string `json:"mode"`
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type Links []*Link

func (l Links) Get(mode string, os string, arch string) *Link {
	for _, link := range l {
		if link.Mode == mode && link.OS == os && link.Arch == arch {
			return link
		}
	}
	return nil
}

func (l Links) GetByOS(os string) Links {
	var ret Links
	for _, link := range l {
		if link.OS == os {
			ret = append(ret, link)
		}
	}
	return ret
}

var availableLinks Links

const (
	modeServer  = "server"
	modeDesktop = "desktop"

	osAIX       = "aix"
	osDragonfly = "dragonfly"
	osFreeBSD   = "freebsd"
	osIllumos   = "illumos"
	osJS        = "js"
	osLinux     = "linux"
	osMac       = "macos"
	osMobile    = "mobile"
	osNetBSD    = "netbsd"
	osOpenBSD   = "openbsd"
	osPlan9     = "plan9"
	osSolaris   = "solaris"
	osWindows   = "windows"

	archAMD64        = "x86_64"
	archAndroid      = "android"
	archARM64        = "arm64"
	archARMV5        = "armv5"
	archARMV6        = "armv6"
	archARMV7        = "armv7"
	archI386         = "i386"
	archIOS          = "ios"
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
	archWASM         = "wasm"
)
