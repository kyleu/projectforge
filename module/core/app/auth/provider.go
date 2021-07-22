package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"$PF_PACKAGE$/app/util"
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
	cb := fmt.Sprintf("%s://%s/auth/%s/callback", proto, host, p.ID)
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
		s.baseURL = "http://localhost:$PF_PORT$"
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
		s.logger.Debugf("authentication enabled for [%s]", util.OxfordComma(ret.Titles(), "and"))
	}

	return nil
}
