package auth

import (
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func getState(r *http.Request) string {
	state := r.URL.Query().Get("state")
	if state != "" {
		return state
	}
	nonceBytes := util.RandomBytes(64)
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func validateState(w http.ResponseWriter, r *http.Request, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	originalState := authURL.Query().Get("state")
	qs := r.URL.Query().Get("state")
	if originalState != "" && (originalState != qs) {
		return errors.New("state token mismatch")
	}
	return nil
}
