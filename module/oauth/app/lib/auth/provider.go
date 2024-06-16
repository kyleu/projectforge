package auth

import (
	"fmt"
	"path"
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

func (p *Provider) Goth(proto string, host string, redirPathOverride ...string) (goth.Provider, error) {
	u := GothURL(p.ID, proto, host, redirPathOverride...)
	gothPrv, err := toGoth(p.ID, p.Key, p.Secret, u, p.Scopes...)
	if err != nil {
		return nil, err
	}
	goth.UseProviders(gothPrv)
	return gothPrv, nil
}

func GothURL(key string, proto string, host string, redirPathOverride ...string) string {
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
	if len(redirPathOverride) == 0 {
		u = fmt.Sprintf("%s/auth/callback/%s", u, key)
	} else {
		if strings.HasPrefix(redirPathOverride[0], proto) {
			return redirPathOverride[0]
		}
		u = fmt.Sprintf("%s%s", u, path.Join(redirPathOverride...))
	}
	return u
}
