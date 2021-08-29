package user

import (
	"strings"
)

type Permission struct {
	Path  string `json:"path"`
	Match string `json:"match"`
	Allow bool   `json:"allow"`
}

func (p Permission) Matches(path string) bool {
	return strings.HasPrefix(path, p.Path)
}

type Permissions []*Permission

func (p Permissions) Check(path string, accounts Accounts) bool {
	if p == nil {
		return true
	}
	for _, perm := range p {
		if accounts.Matches(perm.Match) {
			if perm.Matches(path) {
				println("####", perm.Path, "::", perm.Allow)
				return true
			}
		}
	}
	return false
}
