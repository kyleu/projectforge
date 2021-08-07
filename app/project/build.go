package project

type Build struct {
	SkipDesktop  bool `json:"skipDesktop,omitempty"`
	SkipNotarize bool `json:"skipNotarize,omitempty"`
	SkipHomebrew bool `json:"skipHomebrew,omitempty"`

	SkipWASM    bool `json:"skipWASM,omitempty"`
	SkipIOS     bool `json:"skipIOS,omitempty"`
	SkipAndroid bool `json:"skipAndroid,omitempty"`

	SkipLinuxArm  bool `json:"skipLinuxArm,omitempty"`
	SkipLinuxMips bool `json:"skipLinuxMips,omitempty"`
	SkipLinuxOdd  bool `json:"skipLinuxOdd,omitempty"`

	SkipAIX       bool `json:"skipAIX,omitempty"`
	SkipDragonfly bool `json:"skipDragonfly,omitempty"`
	SkipIllumos   bool `json:"skipIllumos,omitempty"`
	SkipFreeBSD   bool `json:"skipFreeBSD,omitempty"`
	SkipNetBSD    bool `json:"skipNetBSD,omitempty"`
	SkipOpenBSD   bool `json:"skipOpenBSD,omitempty"`
	SkipPlan9     bool `json:"skipPlan9,omitempty"`
	SkipSolaris   bool `json:"skipSolaris,omitempty"`

	SkipSigning   bool `json:"skipSigning,omitempty"`
	SkipNFPMS     bool `json:"skipNFPMS,omitempty"`
	SkipScoop     bool `json:"skipScoop,omitempty"`
	SkipSnapcraft bool `json:"skipSnapcraft,omitempty"`
}
