package model

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

var (
	TypeInt       = newType("int")
	TypeMap       = newType("map", "json", "jsonb")
	TypeString    = newType("string", "text")
	TypeTimestamp = newType("timestamp", "datetime")
	TypeUUID      = newType("uuid")
	AllTypes      = []*Type{TypeInt, TypeMap, TypeString, TypeTimestamp, TypeUUID}
)

type Types []*Type

func TypeFromString(s string) (*Type, error) {
	for _, t := range AllTypes {
		for _, k := range t.Names {
			if k == s {
				return t, nil
			}
		}
	}
	return nil, errors.New("No export type available with key [" + s + "]")
}

func (t *Type) String() string {
	return t.Key
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(t.Key, false), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := util.FromJSON(data, &s); err != nil {
		return err
	}
	x, err := TypeFromString(s)
	if err != nil {
		return errors.Wrapf(err, "no export type available with key [%s]", s)
	}

	*t = *x
	return nil
}
