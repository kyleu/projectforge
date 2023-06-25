package auth

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

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

	s.providers = lo.FilterMap(AvailableProviderKeys, func(k string, _ int) (*Provider, bool) {
		envKey := util.GetEnv(k + "_key")
		if envKey == "" {
			return nil, false
		}
		envSecret := util.GetEnv(k + "_secret")
		envScopes := util.StringSplitAndTrim(util.GetEnv(k+"_scopes"), ",")
		return &Provider{ID: k, Title: AvailableProviderNames[k], Key: envKey, Secret: envSecret, Scopes: envScopes}, true
	})

	if len(s.providers) == 0 {
		logger.Debug("authentication disabled, no providers configured in environment")
	} else {
		const msg = "authentication enabled for [%s], using [%s] as a base URL"
		logger.Infof(msg, util.StringArrayOxfordComma(s.providers.Titles(), "and"), s.baseURL)
	}

	return nil
}

