package user

import (
	"fmt"

	"github.com/go-gem/sessions"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/auth"
	"{{{ .Package }}}/app/util"
)

type Profile struct {
	Name  string `json:"name"`
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

func (p *Profile) AuthString(a auth.Sessions) string {
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

func SaveProfile(n *Profile, ctx *fasthttp.RequestCtx, sess *sessions.Session, logger *zap.SugaredLogger) error {
	if n == nil || n.Equals(DefaultProfile) {
		err := auth.RemoveFromSession("profile", ctx, sess, logger)
		if err != nil {
			return errors.Wrap(err, "unable to remove profile from session")
		}
		return nil
	}
	err := auth.StoreInSession("profile", util.ToJSON(n), ctx, sess, logger)
	if err != nil {
		return errors.Wrap(err, "unable to save profile in session")
	}
	return nil
}
