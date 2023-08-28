package user

import (
	"cmp"
	"slices"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Accounts []*Account

func (a Accounts) String() string {
	return strings.Join(lo.Map(a, func(x *Account, _ int) string {
		return x.String()
	}), ",")
}

func (a Accounts) TitleString() string {
	return strings.Join(lo.Map(a, func(x *Account, _ int) string {
		return x.TitleString()
	}), ",")
}

func (a Accounts) Images() []string {
	ret := make(util.KeyVals[string], 0, len(a))
	lo.ForEach(a, func(x *Account, _ int) {
		if x.Picture != "" {
			ret = append(ret, &util.KeyVal[string]{Key: x.Provider, Val: x.Picture})
		}
	})
	return ret.Values()
}

func (a Accounts) Image() string {
	if is := a.Images(); len(is) > 0 {
		return is[0]
	}
	return ""
}

func (a Accounts) Sort() {
	slices.SortFunc(a, func(l *Account, r *Account) int {
		if l.Provider == r.Provider {
			return cmp.Compare(strings.ToLower(l.Email), strings.ToLower(r.Email))
		}
		return cmp.Compare(l.Provider, r.Provider)
	})
}

func (a Accounts) GetByProvider(p string) Accounts {
	return lo.Filter(a, func(x *Account, _ int) bool {
		return x.Provider == p
	})
}

func (a Accounts) GetByProviderDomain(p string, d string) *Account {
	return lo.FindOrElse(a, nil, func(x *Account) bool {
		return x.Provider == p && x.Domain() == d
	})
}

func (a Accounts) Matches(match string) bool {
	if match == "" || match == "*" {
		return true
	}
	if strings.Contains(match, ",") {
		return lo.ContainsBy(util.StringSplitAndTrim(match, ","), func(x string) bool {
			return a.Matches(x)
		})
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
	return lo.Map(util.StringSplitAndTrim(s, ","), func(x string, _ int) *Account {
		return accountFromString(x)
	})
}
