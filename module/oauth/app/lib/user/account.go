package user

import (
	"strings"

	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/util"
)

type Account struct {
	Provider string `json:"provider"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Token    string `json:"-"`
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

func (a Account) String() string {
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

func (a Account) TitleString() string {
	return a.Provider + ":" + a.Email
}

func (a Account) Domain() string {
	if a.Email == "" || !strings.Contains(a.Email, "@") {
		return ""
	}
	_, r := util.StringSplitLast(a.Email, '@', true)
	return r
}

type Accounts []*Account

func (a Accounts) String() string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
		ret = append(ret, x.String())
	}
	return strings.Join(ret, ",")
}

func (a Accounts) TitleString() string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
		ret = append(ret, x.TitleString())
	}
	return strings.Join(ret, ",")
}

func (a Accounts) Images() []string {
	ret := make(util.KeyVals[string], 0, len(a))
	for _, x := range a {
		if x.Picture != "" {
			ret = append(ret, &util.KeyVal[string]{Key: x.Provider, Val: x.Picture})
		}
	}
	return ret.Values()
}

func (a Accounts) Image() string {
	if is := a.Images(); len(is) > 0 {
		return is[0]
	}
	return ""
}

func (a Accounts) Sort() {
	slices.SortFunc(a, func(l *Account, r *Account) bool {
		if l.Provider == r.Provider {
			return strings.ToLower(l.Email) < strings.ToLower(r.Email)
		}
		return l.Provider < r.Provider
	})
}

func (a Accounts) GetByProvider(p string) Accounts {
	var ret Accounts
	for _, x := range a {
		if x.Provider == p {
			ret = append(ret, x)
		}
	}
	return ret
}

func (a Accounts) GetByProviderDomain(p string, d string) *Account {
	for _, x := range a {
		if x.Provider == p && x.Domain() == d {
			return x
		}
	}
	return nil
}

func (a Accounts) Matches(match string) bool {
	if match == "" || match == "*" {
		return true
	}
	if strings.Contains(match, ",") {
		xs := util.StringSplitAndTrim(match, ",")
		for _, x := range xs {
			if a.Matches(x) {
				return true
			}
		}
		return false
	}
	prv, acct := util.StringSplit(match, ':', true)
	for _, x := range a {
		if x.Provider == prv {
			if acct == "" {
				return true
			}
			return strings.HasSuffix(x.Email, acct)
		}
	}
	return false
}

func (a Accounts) Purge(keys ...string) Accounts {
	ret := make(Accounts, 0, len(a))
	for _, ss := range a {
		hit := false
		for _, key := range keys {
			if ss.Provider == key {
				hit = true
			}
		}
		if !hit {
			ret = append(ret, ss)
		}
	}
	return ret
}

func AccountsFromString(s string) Accounts {
	split := util.StringSplitAndTrim(s, ",")
	ret := make(Accounts, 0, len(split))
	for _, x := range split {
		ret = append(ret, accountFromString(x))
	}
	return ret
}
