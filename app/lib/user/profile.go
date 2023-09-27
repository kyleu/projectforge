// Package user - Content managed by Project Forge, see [projectforge.md] for details.
package user

type Profile struct {
	Name  string `json:"name,omitempty"`
	Mode  string `json:"mode,omitempty"`
	Theme string `json:"theme,omitempty"`
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

func (p *Profile) Equals(x *Profile) bool {
	return p.Name == x.Name && p.Mode == x.Mode && p.Theme == x.Theme
}
