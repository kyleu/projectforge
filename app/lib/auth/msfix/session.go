// Content managed by Project Forge, see [projectforge.md] for details.
package msfix

import (
	"encoding/json"
	"time"

	"github.com/markbates/goth"
	"github.com/pkg/errors"
)

type Session struct {
	AuthURL     string
	AccessToken string
	ExpiresAt   time.Time
}

func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}

	return s.AuthURL, nil
}

func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p, ok := provider.(*Provider)
	if !ok {
		return "", errors.Errorf("invalid provider of type [%T]", provider)
	}
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.ExpiresAt = token.Expiry

	return token.AccessToken, err
}

func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}
