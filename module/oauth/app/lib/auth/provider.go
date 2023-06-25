package auth

import (
	"fmt"
	"strings"

	"github.com/markbates/goth"

	"{{{ .Package }}}/app/util"
)

type Provider struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Key    string   `json:"-"`
	Secret string   `json:"-"`
	Scopes []string `json:"-"`
}

func (p *Provider) Goth(proto string, host string) (goth.Provider, error) {
	if p := util.GetEnv("oauth_protocol"); p != "" {
		proto = p
	}
	if proto == "" {
		proto = "http"
	}
	u := fmt.Sprintf("%s://%s", proto, host)

	if env := util.GetEnv(util.AppKey + "_oauth_redirect"); env != "" {
		u = env
	}
	if env := util.GetEnv("oauth_redirect"); env != "" {
		u = env
	}
	u = strings.TrimSuffix(u, "/")
	cb := fmt.Sprintf("%s/auth/callback/%s", u, p.ID)
	gothPrv, err := toGoth(p.ID, p.Key, p.Secret, cb, p.Scopes...)
	if err != nil {
		return nil, err
	}
	goth.UseProviders(gothPrv)
	return gothPrv, nil
}
