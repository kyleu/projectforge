package auth

import (
	"fmt"
	"strings"
	"sync"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

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

type Providers []*Provider

func (p Providers) Get(id string) *Provider {
	for _, x := range p {
		if x.ID == id {
			return x
		}
	}
	return nil
}

func (p Providers) Contains(id string) bool {
	return p.Get(id) != nil
}

func (p Providers) IDs() []string {
	ret := make([]string, 0, len(p))
	for _, x := range p {
		ret = append(ret, x.ID)
	}
	return ret
}

func (p Providers) Titles() []string {
	ret := make([]string, 0, len(p))
	for _, x := range p {
		ret = append(ret, x.Title)
	}
	return ret
}

func (s *Service) Providers(logger util.Logger) (Providers, error) {
	if s.providers == nil {
		err := s.load(logger)
		if err != nil {
			return nil, err
		}
	}
	return s.providers, nil
}

var initMu sync.Mutex

func (s *Service) load(logger util.Logger) error {
	initMu.Lock()
	defer initMu.Unlock()

	if s.providers != nil {
		return errors.New("called [load] twice")
	}
	if s.baseURL == "" {
		s.baseURL = util.GetEnv(util.AppKey + "_oauth_redirect")
	}
	if s.baseURL == "" {
		s.baseURL = fmt.Sprintf("http://localhost:%d", s.port)
	}
	s.baseURL = strings.TrimSuffix(s.baseURL, "/")

	initAvailable()

	ret := Providers{}
	for _, k := range AvailableProviderKeys {
		envKey := util.GetEnv(k + "_key")
		envSecret := util.GetEnv(k + "_secret")
		envScopes := util.StringSplitAndTrim(util.GetEnv(k+"_scopes"), ",")
		if envKey != "" {
			ret = append(ret, &Provider{ID: k, Title: AvailableProviderNames[k], Key: envKey, Secret: envSecret, Scopes: envScopes})
		}
	}

	s.providers = ret

	if len(ret) == 0 {
		logger.Debug("authentication disabled, no providers configured in environment")
	} else {
		const msg = "authentication enabled for [%s], using [%s] as a base URL"
		logger.Infof(msg, util.StringArrayOxfordComma(ret.Titles(), "and"), s.baseURL)
	}

	return nil
}
