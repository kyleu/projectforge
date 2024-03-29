package model

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
)

func AsEnum(t types.Type) (*types.Enum, error) {
	w, ok := t.(*types.Wrapped)
	if ok {
		t = w.T
	}
	ref, ok := t.(*types.Enum)
	if !ok {
		return nil, errors.Errorf("InvalidType(%T)", w.T)
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
	w, ok := t.(*types.Wrapped)
	if ok {
		t = w.T
	}
	ref, ok := t.(*types.Reference)
	if !ok {
		return nil, errors.Errorf("InvalidType(%T)", w.T)
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
