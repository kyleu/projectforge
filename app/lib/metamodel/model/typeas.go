package model

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func AsEnum(t types.Type) (*types.Enum, error) {
	w, err := util.Cast[*types.Wrapped](t)
	if err == nil {
		t = w.T
	}
	ref, err := util.Cast[*types.Enum](t)
	if err != nil {
		return nil, errors.Wrapf(err, "InvalidType(%T)", w.T)
	}
	return ref, nil
}

func AsEnumInstance(t types.Type, enums enum.Enums) (*enum.Enum, error) {
	e, err := AsEnum(t)
	if err != nil {
		return nil, err
	}
	ret := enums.Get(e.Ref)
	if ret == nil {
		return nil, errors.Errorf("no enum found with name [%s]", e.Ref)
	}
	return ret, nil
}

func AsRef(t types.Type) (*types.Reference, error) {
	w, err := util.Cast[*types.Wrapped](t)
	if err == nil {
		t = w.T
	}
	ref, err := util.Cast[*types.Reference](t)
	if err != nil {
		return nil, errors.Wrapf(err, "InvalidType(%T)", w.T)
	}
	if ref.K == "" {
		return nil, errors.New("invalid empty reference")
	}
	return ref, nil
}

func asRefK(t types.Type) string {
	ref, err := AsRef(t)
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return ref.K
}
