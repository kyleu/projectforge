package user

import (
	"fmt"{{{ if .HasModule "user" }}}

	"github.com/google/uuid"{{{ end }}}
)

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
	return &Profile{Name: p.Name, Mode: p.Mode, Theme: p.Theme}
}

func (p *Profile) ModeClass() string {
	if p.Mode == "" {
		return ""
	}
	return "mode-" + p.Mode
}

func (p *Profile) AuthString(a Accounts) string {
	msg := fmt.Sprintf("signed in as %s", p.String())
	if len(a) == 0 {
		if p.Name == DefaultProfile.Name {
			return "click to sign in"
		}
		return msg
	}
	return fmt.Sprintf("%s using [%s]", msg, a.String())
}

func (p *Profile) Equals(x *Profile) bool {
	return p.Name == x.Name && p.Mode == x.Mode && p.Theme == x.Theme
}
