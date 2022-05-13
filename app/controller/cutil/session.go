// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

const (
	WebAuthKey  = "auth"
	WebFlashKey = "flash"
	ReferKey    = "refer"
)

func NewCookie(v string) *fasthttp.Cookie {
	ret := &fasthttp.Cookie{}
	ret.SetPath("/")
	ret.SetHTTPOnly(true)
	ret.SetMaxAge(365 * 24 * 60 * 60)
	ret.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	ret.SetKey(util.AppKey)
	ret.SetValue(v)
	return ret
}

func StoreInSession(k string, v string, rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger) error {
	websess[k] = v
	return SaveSession(rc, websess, logger)
}

func RemoveFromSession(k string, rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger) error {
	delete(websess, k)
	return SaveSession(rc, websess, logger)
}

func SaveSession(rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger) error {
	enc, err := util.EncryptMessage(nil, util.ToJSONCompact(websess), logger)
	if err != nil {
		return err
	}
	c := NewCookie(enc)
	rc.Response.Header.SetCookie(c)
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

func SaveProfile(n *user.Profile, rc *fasthttp.RequestCtx, sess util.ValueMap, logger util.Logger) error {
	if n == nil || n.Equals(user.DefaultProfile) {
		err := RemoveFromSession("profile", rc, sess, logger)
		if err != nil {
			return errors.Wrap(err, "unable to remove profile from session")
		}
		return nil
	}
	err := StoreInSession("profile", util.ToJSON(n), rc, sess, logger)
	if err != nil {
		return errors.Wrap(err, "unable to save profile in session")
	}
	return nil
}
