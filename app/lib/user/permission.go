// Content managed by Project Forge, see [projectforge.md] for details.
package user

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/util"
)

const permPrefix = "perm: "

var (
	PermissionsLogger util.Logger
	perms             Permissions
)

func SetPermissions(allowDefault bool, ps ...*Permission) {
	perms = make(Permissions, 0, len(ps)+4)
	perms = append(perms, Perm("/auth", "*", true), Perm("/profile", "*", true))
	perms = append(perms, ps...)
	perms = append(perms, Perm("/admin", "*", false), Perm("/about", "*", true), Perm("/", "*", allowDefault))
}

func GetPermissions() Permissions {
	ret := make(Permissions, 0, len(perms))
	return append(ret, perms...)
}

type Permission struct {
	Path  string `json:"path"`
	Match string `json:"match"`
	Allow bool   `json:"allow"`
}

func Check(path string, accounts Accounts) (bool, string) {
	return perms.Check(path, accounts)
}

func Perm(p string, m string, a bool) *Permission {
	return &Permission{Path: p, Match: m, Allow: a}
}

func (p Permission) Matches(path string) bool {
	return strings.HasPrefix(path, p.Path)
}

func (p Permission) String() string {
	return fmt.Sprintf("%s [%s::%t]", p.Path, p.Match, p.Allow)
}

type Permissions []*Permission

func (p Permissions) Sort() {
	slices.SortFunc(p, func(l *Permission, r *Permission) bool {
		if l.Path == r.Path {
			return l.Match < r.Match
		}
		return l.Path < r.Path
	})
}

func (p Permissions) Check(path string, accounts Accounts) (bool, string) {
	if PermissionsLogger != nil {
		PermissionsLogger.Debugf(permPrefix+"checking [%d] permissions for [%s]", len(p), accounts.String())
	}
	if len(p) == 0 {
		const msg = "no permissions configured"
		if PermissionsLogger != nil {
			PermissionsLogger.Debug(permPrefix + msg)
		}
		return true, msg
	}
	for _, perm := range p {
		if perm.Matches(path) {
			if accounts.Matches(perm.Match) {
				msg := fmt.Sprintf("matched [%s], result [%t]", perm.Match, perm.Allow)
				if PermissionsLogger != nil {
					PermissionsLogger.Debug(permPrefix + msg)
				}
				return perm.Allow, msg
			}
		}
	}
	msg := fmt.Sprintf("no matches among [%d] permissions", len(p))
	if PermissionsLogger != nil {
		PermissionsLogger.Debug(permPrefix + msg)
	}
	return false, msg
}
