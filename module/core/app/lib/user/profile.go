package user
{{{ if .HasUser }}}
import "github.com/google/uuid"
{{{ end }}}
type Profile struct {
{{{ if .HasUser }}}	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name,omitempty"`
	Mode  string    `json:"mode,omitempty"`
	Theme string    `json:"theme,omitempty"`
{{{ else }}}	Name  string `json:"name,omitempty"`
	Mode  string `json:"mode,omitempty"`
	Theme string `json:"theme,omitempty"`
{{{ end }}}}

var DefaultProfile = &Profile{Name: "Guest"}

func (p *Profile) String() string {
	return p.Name
}

func (p *Profile) Clone() *Profile {
	{{{ if .HasUser }}}return &Profile{ID: p.ID, Name: p.Name, Mode: p.Mode, Theme: p.Theme}{{{ else }}}return &Profile{Name: p.Name, Mode: p.Mode, Theme: p.Theme}{{{ end }}}
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
