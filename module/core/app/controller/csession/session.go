package csession

import (
	"net/http"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

const ({{{ if .HasAccount }}}
	WebAuthKey  = "auth"{{{ end }}}
	WebFlashKey = "flash"
	ReferKey    = "refer"
)

func NewCookie(v string) *http.Cookie {
	return &http.Cookie{Name: util.AppKey, Value: v, Path: "/", MaxAge: 365 * 24 * 60 * 60, HttpOnly: true, SameSite: http.SameSiteLaxMode}
}

func StoreInSession(k string, v string, w http.ResponseWriter, websess util.ValueMap, logger util.Logger) error {
	websess[k] = v
	return SaveSession(w, websess, logger)
}

func RemoveFromSession(k string, w http.ResponseWriter, websess util.ValueMap, logger util.Logger) error {
	delete(websess, k)
	return SaveSession(w, websess, logger)
}

func SaveSession(w http.ResponseWriter, websess util.ValueMap, logger util.Logger) error {
	js := util.ToJSONCompact(websess)
	enc, err := util.EncryptMessage(nil, js, logger)
	if err != nil {
		return err
	}
	http.SetCookie(w, NewCookie(enc))
	return nil
}

func GetFromSession(key string, websess util.ValueMap) (string, error) {
	value, ok := websess[key]
	if !ok {
		return "", errors.Errorf("could not find a matching session value with key [%s] for this request", key)
	}
	s, ok := value.(string)
	if !ok {
		return "", errors.Errorf("session value with key [%s] is of type [%T], not [string]", key, value)
	}
	return s, nil
}

func SaveProfile(n *user.Profile, w http.ResponseWriter, sess util.ValueMap, logger util.Logger) error {
	if n != nil && n.Name == "" {
		n.Name = user.DefaultProfile.Name
	}
	if n == nil || n.Equals(user.DefaultProfile) {
		return errors.Wrap(RemoveFromSession("profile", w, sess, logger), "unable to remove profile from session")
	}
	if n.Name == user.DefaultProfile.Name {
		n.Name = ""
	}
	err := StoreInSession("profile", util.ToJSON(n), w, sess, logger)
	if err != nil {
		return errors.Wrap(err, "unable to save profile in session")
	}
	return nil
}
