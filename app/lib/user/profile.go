package user

type Profile struct {
	Name  string `json:"name,omitzero"`
	Mode  string `json:"mode,omitzero"`
	Theme string `json:"theme,omitzero"`
}

var DefaultProfile = &Profile{Name: "Guest"}

func (p *Profile) String() string {
	return p.Name
}

func (p *Profile) Clone() *Profile {
	return &Profile{Name: p.Name, Mode: p.Mode, Theme: p.Theme}
}

func (p *Profile) ModeClass() string {
	if p.Mode == "" {
		return ""
	}
	return "mode-" + p.Mode
}

func (p *Profile) SetName(n string) bool {
	if p.Name == n {
		return false
	}
	if p.Name != "" && p.Name != DefaultProfile.Name {
		return false
	}
	p.Name = n
	return true
}

func (p *Profile) Equals(x *Profile) bool {
	return p.Name == x.Name && p.Mode == x.Mode && p.Theme == x.Theme
}
