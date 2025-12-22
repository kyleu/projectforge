// $PF_GENERATE_ONCE$
package user

import (
	"time"

	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

type User struct {
	ID      uuid.UUID  `json:"id,omitzero"`
	Name    string     `json:"name,omitzero"`
	Picture string     `json:"picture,omitzero"`
	Created time.Time  `json:"created,omitzero"`
	Updated *time.Time `json:"updated,omitzero"`
}

func New(id uuid.UUID) *User {
	return &User{ID: id}
}

func Random() *User {
	return &User{
		ID:      util.UUID(),
		Name:    util.RandomID(),
		Picture: "https://" + util.RandomString(6) + ".com/" + util.RandomString(6),
		Created: util.TimeCurrent(),
		Updated: util.TimeCurrentP(),
	}
}

func FromMap(m util.ValueMap, setPK bool) (*User, error) {
	ret := &User{}
	var err error
	if setPK {
		retID, e := m.ParseUUID("id", true, true)
		if e != nil {
			return nil, e
		}
		if retID != nil {
			ret.ID = *retID
		}
		// $PF_SECTION_START(pkchecks)$
		// $PF_SECTION_END(pkchecks)$
	}
	ret.Name, err = m.ParseString("name", true, true)
	if err != nil {
		return nil, err
	}
	ret.Picture, err = m.ParseString("picture", true, true)
	if err != nil {
		return nil, err
	}
	// $PF_SECTION_START(extrachecks)$
	// $PF_SECTION_END(extrachecks)$
	return ret, nil
}

func (u *User) Clone() *User {
	return &User{u.ID, u.Name, u.Picture, u.Created, u.Updated}
}

func (u *User) String() string {
	return u.ID.String()
}

func (u *User) TitleString() string {
	return u.Name
}

func (u *User) WebPath() string {
	return "/admin/db/user/" + u.ID.String()
}

func (u *User) Diff(ux *User) util.Diffs {
	var diffs util.Diffs
	if u.ID != ux.ID {
		diffs = append(diffs, util.NewDiff("id", u.ID.String(), ux.ID.String()))
	}
	if u.Name != ux.Name {
		diffs = append(diffs, util.NewDiff("name", u.Name, ux.Name))
	}
	if u.Picture != ux.Picture {
		diffs = append(diffs, util.NewDiff("picture", u.Picture, ux.Picture))
	}
	if u.Created != ux.Created {
		diffs = append(diffs, util.NewDiff("created", u.Created.String(), ux.Created.String()))
	}
	return diffs
}

func (u *User) ToData() []any {
	return []any{u.ID, u.Name, u.Picture, u.Created, u.Updated}
}
