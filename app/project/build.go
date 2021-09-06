package project

import (
	"fmt"

	"github.com/kyleu/projectforge/app/util"
)

type Build struct {
	Publish bool `json:"publish,omitempty"`

	Desktop  bool `json:"desktop,omitempty"`
	Notarize bool `json:"notarize,omitempty"`
	Signing  bool `json:"signing,omitempty"`

	Android bool `json:"android,omitempty"`
	IOS     bool `json:"iOS,omitempty"`
	WASM    bool `json:"wasm,omitempty"`

	WindowsARM bool `json:"windowsARM,omitempty"`
	LinuxARM   bool `json:"linuxARM,omitempty"`
	LinuxMIPS  bool `json:"linuxMIPS,omitempty"`
	LinuxOdd   bool `json:"linuxOdd,omitempty"`

	AIX       bool `json:"aix,omitempty"`
	Dragonfly bool `json:"dragonfly,omitempty"`
	Illumos   bool `json:"illumos,omitempty"`
	FreeBSD   bool `json:"freeBSD,omitempty"`
	NetBSD    bool `json:"netBSD,omitempty"`
	OpenBSD   bool `json:"openBSD,omitempty"`
	Plan9     bool `json:"plan9,omitempty"`
	Solaris   bool `json:"solaris,omitempty"`

	Homebrew  bool `json:"homebrew,omitempty"`
	NFPMS     bool `json:"nfpms,omitempty"`
	Snapcraft bool `json:"snapcraft,omitempty"`
}

func (b *Build) Empty() bool {
	return !(b.Publish || b.Desktop || b.Notarize || b.Signing || b.Android || b.IOS || b.WASM || b.WindowsARM ||
		b.LinuxARM || b.LinuxMIPS || b.LinuxOdd || b.AIX || b.Dragonfly || b.Illumos || b.FreeBSD || b.NetBSD ||
		b.OpenBSD || b.Plan9 || b.Solaris || b.Homebrew || b.NFPMS || b.Snapcraft)
}

func (b *Build) ToMap() map[string]bool {
	return map[string]bool{
		"publish": b.Publish, "desktop": b.Desktop, "notarize": b.Notarize, "signing": b.Signing,
		"android": b.Android, "ios": b.IOS, "wasm": b.WASM, "windows-arm": b.WindowsARM,
		"linux-arm": b.LinuxARM, "linux-mips": b.LinuxMIPS, "linux-odd": b.LinuxOdd,
		"aix": b.AIX, "dragonfly": b.Dragonfly, "illumos": b.Illumos, "freebsd": b.FreeBSD,
		"netbsd": b.NetBSD, "openbsd": b.OpenBSD, "plan9": b.Plan9, "solaris": b.Solaris,
		"homebrew": b.Homebrew, "nfpms": b.NFPMS, "snapcraft": b.Snapcraft,
	}
}

func BuildFromMap(frm util.ValueMap) *Build {
	x := func(k string) bool {
		v := fmt.Sprint(frm["build-"+k])
		return v == "true"
	}
	return &Build{
		Publish: x("publish"), Desktop: x("desktop"), Notarize: x("notarize"), Signing: x("signing"),
		Android: x("android"), IOS: x("ios"), WASM: x("wasm"), WindowsARM: x("windows-arm"),
		LinuxARM: x("linux-arm"), LinuxMIPS: x("linux-mips"), LinuxOdd: x("linux-odd"),
		AIX: x("aix"), Dragonfly: x("dragonfly"), Illumos: x("illumos"), FreeBSD: x("freebsd"),
		NetBSD: x("netbsd"), OpenBSD: x("openbsd"), Plan9: x("plan9"), Solaris: x("solaris"),
		Homebrew: x("homebrew"), NFPMS: x("nfpms"), Snapcraft: x("snapcraft"),
	}
}

type BuildOption struct {
	Key         string
	Title       string
	Description string
}

var AllBuildOptions = []*BuildOption{
	{Key: "publish", Title: "Publish", Description: "The release process will publish a full release"},

	{Key: "desktop", Title: "Desktop", Description: "Webview-based applications for the three major operating systems"},
	{Key: "notarize", Title: "Notarize", Description: "Sends build artifacts to Apple for notarization"},
	{Key: "signing", Title: "Signing", Description: "Signs the checksums using gpg"},

	{Key: "android", Title: "Android", Description: "Builds the application as an Android library"},
	{Key: "ios", Title: "iOS", Description: "Builds the application as an iOS framework "},
	{Key: "wasm", Title: "WASM", Description: "Builds the application for WebAssembly"},

	{Key: "windows-arm", Title: "Windows ARM", Description: "Builds the application for Windows on ARM and ARM64 architectures"},
	{Key: "linux-arm", Title: "Linux ARM", Description: "Builds the application for Linux on ARM and ARM64 architectures"},
	{Key: "linux-mips", Title: "Linux MIPS", Description: "Builds the application for Linux on MIPS architectures"},
	{Key: "linux-odd", Title: "Linux Odd", Description: "Builds the application for Linux using ppc64, ppc64le, riscv64, and s390x"},

	{Key: "aix", Title: "AIX", Description: "Builds the application for AIX"},
	{Key: "dragonfly", Title: "Dragonfly", Description: "Builds the application for Dragonfly"},
	{Key: "illumos", Title: "Illumos", Description: "Builds the application for Illumos"},
	{Key: "freebsd", Title: "FreeBSD", Description: "Builds the application for FreeBSD"},
	{Key: "netbsd", Title: "NetBSD", Description: "Builds the application for NetBSD"},
	{Key: "openbsd", Title: "OpenBSD", Description: "Builds the application for OpenBSD"},
	{Key: "plan9", Title: "Plan 9", Description: "Builds the application for Plan 9"},
	{Key: "solaris", Title: "Solaris", Description: "Builds the application for Solaris"},

	{Key: "homebrew", Title: "Homebrew", Description: "Publishes the builds to Homebrew"},
	{Key: "nfpms", Title: "NFPMS", Description: "Builds the application as RPMs, DEBs, and APKs for various Linux flavors "},
	{Key: "snapcraft", Title: "Snapcraft", Description: "Publishes the application as a Ubuntu Snap "},
}
