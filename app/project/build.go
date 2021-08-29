package project

import (
	"fmt"

	"github.com/kyleu/projectforge/app/util"
)

type Build struct {
	SkipPublish bool `json:"skipPublish,omitempty"`

	SkipDesktop  bool `json:"skipDesktop,omitempty"`
	SkipNotarize bool `json:"skipNotarize,omitempty"`
	SkipSigning  bool `json:"skipSigning,omitempty"`

	SkipAndroid bool `json:"skipAndroid,omitempty"`
	SkipIOS     bool `json:"skipIOS,omitempty"`
	SkipWASM    bool `json:"skipWASM,omitempty"`

	SkipLinuxARM  bool `json:"skipLinuxARM,omitempty"`
	SkipLinuxMIPS bool `json:"skipLinuxMIPS,omitempty"`
	SkipLinuxOdd  bool `json:"skipLinuxOdd,omitempty"`

	SkipAIX       bool `json:"skipAIX,omitempty"`
	SkipDragonfly bool `json:"skipDragonfly,omitempty"`
	SkipIllumos   bool `json:"skipIllumos,omitempty"`
	SkipFreeBSD   bool `json:"skipFreeBSD,omitempty"`
	SkipNetBSD    bool `json:"skipNetBSD,omitempty"`
	SkipOpenBSD   bool `json:"skipOpenBSD,omitempty"`
	SkipPlan9     bool `json:"skipPlan9,omitempty"`
	SkipSolaris   bool `json:"skipSolaris,omitempty"`

	SkipHomebrew  bool `json:"skipHomebrew,omitempty"`
	SkipNFPMS     bool `json:"skipNFPMS,omitempty"`
	SkipScoop     bool `json:"skipScoop,omitempty"`
	SkipSnapcraft bool `json:"skipSnapcraft,omitempty"`
}

func (b *Build) Empty() bool {
	return !(b.SkipPublish || b.SkipDesktop || b.SkipNotarize || b.SkipSigning || b.SkipAndroid || b.SkipIOS || b.SkipWASM || b.SkipLinuxARM ||
		b.SkipLinuxMIPS || b.SkipLinuxOdd || b.SkipAIX || b.SkipDragonfly || b.SkipIllumos || b.SkipFreeBSD || b.SkipNetBSD ||
		b.SkipOpenBSD || b.SkipPlan9 || b.SkipSolaris || b.SkipHomebrew || b.SkipNFPMS || b.SkipScoop || b.SkipSnapcraft)
}

func (b *Build) ToMap() map[string]bool {
	return map[string]bool{
		"publish": !b.SkipPublish, "desktop": !b.SkipDesktop, "notarize": !b.SkipNotarize, "signing": !b.SkipSigning,
		"android": !b.SkipAndroid, "ios": !b.SkipIOS, "wasm": !b.SkipWASM,
		"linux-arm": !b.SkipLinuxARM, "linux-mips": !b.SkipLinuxMIPS, "linux-odd": !b.SkipLinuxOdd,
		"aix": !b.SkipAIX, "dragonfly": !b.SkipDragonfly, "illumos": !b.SkipIllumos, "freebsd": !b.SkipFreeBSD,
		"netbsd": !b.SkipNetBSD, "openbsd": !b.SkipOpenBSD, "plan9": !b.SkipPlan9, "solaris": !b.SkipSolaris,
		"homebrew": !b.SkipHomebrew, "nfpms": !b.SkipNFPMS, "scoop": !b.SkipScoop, "snapcraft": !b.SkipSnapcraft,
	}
}

func BuildFromMap(frm util.ValueMap) *Build {
	x := func(k string) bool {
		v := fmt.Sprint(frm["build-"+k])
		return v != "true"
	}
	return &Build{
		SkipPublish: x("publish"), SkipDesktop: x("desktop"), SkipNotarize: x("notarize"), SkipSigning: x("signing"),
		SkipAndroid: x("android"), SkipIOS: x("ios"), SkipWASM: x("wasm"),
		SkipLinuxARM: x("linux-arm"), SkipLinuxMIPS: x("linux-mips"), SkipLinuxOdd: x("linux-odd"),
		SkipAIX: x("aix"), SkipDragonfly: x("dragonfly"), SkipIllumos: x("illumos"), SkipFreeBSD: x("freebsd"),
		SkipNetBSD: x("netbsd"), SkipOpenBSD: x("openbsd"), SkipPlan9: x("plan9"), SkipSolaris: x("solaris"),
		SkipHomebrew: x("homebrew"), SkipNFPMS: x("nfpms"), SkipScoop: x("scoop"), SkipSnapcraft: x("snapcraft"),
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

	{Key: "linux-arm", Title: "Linux ARM", Description: "Builds the application for Linux on ARM architectures"},
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
	{Key: "scoop", Title: "Scoop", Description: "Builds the application as a Windows Scoop "},
	{Key: "snapcraft", Title: "Snapcraft", Description: "Publishes the application as a Ubuntu Snap "},
}
