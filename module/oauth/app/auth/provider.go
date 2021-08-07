package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

type Provider struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Key    string `json:"-"`
	Secret string `json:"-"`
	cache  map[string]goth.Provider
}

func (p *Provider) Goth(proto string, host string) (goth.Provider, error) {
	if proto == "" {
		proto = "http"
	}
	u := fmt.Sprintf("%s://%s", proto, host)

	env := os.Getenv(util.AppKey + "_oauth_redirect")
	if env != "" {
		u = env
	}
	cb := fmt.Sprintf("%s/auth/callback/%s", u, p.ID)
	if g, ok := p.cache[cb]; ok {
		return g, nil
	}
	gothPrv, err := toGoth(p.ID, p.Key, p.Secret, cb)
	if err != nil {
		return nil, err
	}
	goth.UseProviders(gothPrv)
	p.cache[cb] = gothPrv
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

func (s *Service) Providers() (Providers, error) {
	if s.providers == nil {
		err := s.load()
		if err != nil {
			return nil, err
		}
	}
	return s.providers, nil
}

func (s *Service) load() error {
	if s.providers != nil {
		return errors.New("called [load] twice")
	}
	if s.baseURL == "" {
		s.baseURL = os.Getenv(util.AppKey + "_oauth_redirect")
	}
	if s.baseURL == "" {
		s.baseURL = fmt.Sprintf("http://localhost:%d", util.AppPort)
	}
	s.baseURL = strings.TrimSuffix(s.baseURL, "/")

	initAvailable()

	ret := Providers{}
	for _, k := range AvailableProviderKeys {
		envKey := os.Getenv(k + "_key")
		envSecret := os.Getenv(k + "_secret")
		if envKey != "" {
			ret = append(ret, &Provider{ID: k, Title: AvailableProviderNames[k], Key: envKey, Secret: envSecret, cache: map[string]goth.Provider{}})
		}
	}

	s.providers = ret

	if len(ret) == 0 {
		s.logger.Debug("authentication disabled, no providers configured in environment")
	} else {
		msg := "authentication enabled for [%s], using [%s] as a base URL"
		s.logger.Debugf(msg, util.OxfordComma(ret.Titles(), "and"), s.baseURL)
	}

	return nil
}
