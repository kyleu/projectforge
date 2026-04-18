package template

func (t *Context) BuildAndroid() bool {
	ret := t.HasModules("android") && t.Build.Android
	return ret
}

func (t *Context) BuildIOS() bool {
	return t.HasModules("ios") && t.Build.IOS
}

func (t *Context) BuildDesktop() bool {
	return t.HasModules("desktop") && t.Build.Desktop
}

func (t *Context) BuildMobile() bool {
	return t.BuildIOS() || t.BuildAndroid()
}

func (t *Context) BuildWASM() bool {
	return t.HasModules("wasmserver") && t.Build.WASM
}

func (t *Context) BuildNotarize() bool {
	return t.HasModules("notarize") && t.Build.Notarize
}

func (t *Context) IsNotarized() bool {
	return t.HasModule("notarize") && t.Build != nil && t.Build.Notarize
}

func (t *Context) IsArmAndMips() bool {
	return t.Build.HasArm() && t.Build.LinuxMIPS
}

func (t *Context) HasDocker() bool {
	return !t.Build.SkipDocker
}
