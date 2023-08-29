// $PF_GENERATE_ONCE$
package user

import (
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Users []*User

func (u Users) Get(id uuid.UUID) *User {
	return lo.FindOrElse(u, nil, func(x *User) bool {
		return x.ID == id
	})
}

func (u Users) GetByIDs(ids ...uuid.UUID) Users {
	return lo.Filter(u, func(x *User, _ int) bool {
		return lo.Contains(ids, x.ID)
	})
}

func (u Users) IDs() []uuid.UUID {
	return lo.Map(u, func(x *User, _ int) uuid.UUID {
		return x.ID
	})
}

func (u Users) IDStrings(includeNil bool) []string {
	ret := make([]string, 0, len(u)+1)
	if includeNil {
		ret = append(ret, "")
	}
	lo.ForEach(u, func(x *User, _ int) {
		ret = append(ret, x.ID.String())
	})
	return ret
}

func (u Users) TitleStrings(nilTitle string) []string {
	ret := make([]string, 0, len(u)+1)
	if nilTitle != "" {
		ret = append(ret, nilTitle)
	}
	lo.ForEach(u, func(x *User, _ int) {
		ret = append(ret, x.TitleString())
	})
	return ret
}

func (u Users) Clone() Users {
	return slices.Clone(u)
}
