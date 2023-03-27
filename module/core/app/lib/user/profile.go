package user
{{{ if .HasModule "user" }}}
import "github.com/google/uuid"
{{{ end }}}
type Profile struct {
{{{ if .HasModule "user" }}}	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Mode  string    `json:"mode,omitempty"`
	Theme string    `json:"theme,omitempty"`
{{{ else }}}	Name  string `json:"name"`
	Mode  string `json:"mode,omitempty"`
	Theme string `json:"theme,omitempty"`
{{{ end }}}}

var DefaultProfile = &Profile{Name: "Guest"}

func (p *Profile) String() string {
	return p.Name
}

func (p *Profile) Clone() *Profile {
	{{{ if .HasModule "user" }}}return &Profile{ID: p.ID, Name: p.Name, Mode: p.Mode, Theme: p.Theme}{{{ else }}}return &Profile{Name: p.Name, Mode: p.Mode, Theme: p.Theme}{{{ end }}}
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
