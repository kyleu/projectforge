package msfix

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/markbates/going/defaults"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"{{{ .Package }}}/app/util"
)

const (
	authURL         string = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize"
	tokenURL        string = "https://login.microsoftonline.com/%s/oauth2/v2.0/token" //nolint:gosec
	endpointProfile string = "https://graph.microsoft.com/v1.0/me"
)

var defaultScopes = []string{"openid", "offline_access", "user.read"}

// Note that this is a copy of the `microsoftonline` provider, but accepts a tenant.
func New(clientKey, secret, callbackURL string, tenant string, scopes ...string) *Provider {
	if tenant == "" {
		tenant = "common"
	}
	p := &Provider{ClientKey: clientKey, Secret: secret, CallbackURL: callbackURL, Tenant: tenant, providerName: "microsoft"}
	p.config = newConfig(p, scopes)
	return p
}

type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	Tenant       string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

func (p *Provider) Name() string {
	return p.providerName
}

func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *Provider) Debug(_ bool) {}

func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	au := p.config.AuthCodeURL(state)
	return &Session{AuthURL: au}, nil
}

func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	msSession, ok := session.(*Session)
	if !ok {
		return goth.User{}, errors.Errorf("invalid session of type [%T]", session)
	}
	user := goth.User{
		AccessToken: msSession.AccessToken,
		Provider:    p.Name(),
		ExpiresAt:   msSession.ExpiresAt,
	}

	if user.AccessToken == "" {
		return user, errors.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpointProfile, http.NoBody)
	if err != nil {
		return user, err
	}

	req.Header.Set(authorizationHeader(msSession))

	response, err := p.Client().Do(req)
	if err != nil {
		return user, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return user, errors.Errorf("%s responded with a %d trying to fetch user information", p.providerName, response.StatusCode)
	}

	user.AccessToken = msSession.AccessToken
	if len(user.AccessToken) > 1024 {
		user.AccessToken = ""
	}

	err = userFromReader(response.Body, &user)
	return user, err
}

func (p *Provider) RefreshTokenAvailable() bool {
	return false
}

func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, errors.Errorf("no refresh token provided")
	}

	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(goth.ContextForClient(p.Client()), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf(authURL, provider.Tenant),
			TokenURL: fmt.Sprintf(tokenURL, provider.Tenant),
		},
		Scopes: []string{},
	}

	c.Scopes = append(c.Scopes, scopes...)
	if len(scopes) == 0 {
		c.Scopes = append(c.Scopes, defaultScopes...)
	}

	return c
}

func userFromReader(r io.Reader, user *goth.User) error {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)

	u := struct {
		ID                string `json:"id"`
		Name              string `json:"displayName"`
		Email             string `json:"mail"`
		FirstName         string `json:"givenName"`
		LastName          string `json:"surname"`
		UserPrincipalName string `json:"userPrincipalName"`
	}{}

	if err := util.FromJSONReader(tee, &u); err != nil {
		return err
	}
	raw := map[string]any{}
	if err := util.FromJSONReader(buf, &raw); err != nil {
		return err
	}

	user.UserID = u.ID
	user.Email = defaults.String(u.Email, u.UserPrincipalName)
	user.Name = u.Name
	user.NickName = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.RawData = raw

	return nil
}

func authorizationHeader(session *Session) (string, string) {
	return "Authorization", fmt.Sprintf("Bearer %s", session.AccessToken)
}

func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	return util.FromJSONObj[*Session]([]byte(data))
}
