package user

import (
	"strings"

	"{{{ .Package }}}/app/util"
)

type Account struct {
	Provider string `json:"provider"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Token    string `json:"-"`
}

func (a *Account) String() string {
	ret := a.Provider + ":" + a.Email
	if a.Token != "" {
		msg := a.Token
		if a.Picture != "" {
			msg += "@@" + a.Picture
		}
		if enc, err := util.EncryptMessage(nil, msg, nil); err == nil {
			ret += "|" + enc
		} else {
			ret += "|" + msg
		}
	}
	return ret
}

func (a *Account) TitleString() string {
	return a.Provider + ":" + a.Email
}

func (a *Account) Domain() string {
	if a.Email == "" || !strings.Contains(a.Email, "@") {
		return ""
	}
	_, r := util.StringSplitLast(a.Email, '@', true)
	return r
}

func accountFromString(s string) *Account {
	p, e := util.StringSplit(s, ':', true)
	var t, pic string
	if strings.Contains(e, "|") {
		e, t = util.StringSplit(e, '|', true)
		if decr, err := util.DecryptMessage(nil, t, nil); err == nil {
			t = decr
			if idx := strings.LastIndex(t, "@@"); idx > -1 {
				pic = t[idx+2:]
				t = t[:idx]
			}
		}
	}
	return &Account{Provider: p, Email: e, Picture: pic, Token: t}
}
