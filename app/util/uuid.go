// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"github.com/google/uuid"
)

var UUIDDefault = uuid.UUID{}

func UUIDFromString(s string) *uuid.UUID {
	var retID *uuid.UUID

	if len(s) > 0 {
		s, err := uuid.Parse(s)
		if err == nil {
			retID = &s
		}
	}

	return retID
}

func UUIDFromStringOK(s string) uuid.UUID {
	if len(s) > 0 {
		s, err := uuid.Parse(s)
		if err == nil {
			return s
		}
	}
	return UUIDDefault
}

func UUIDString(u *uuid.UUID) string {
	if u == nil {
		return ""
	}
	return u.String()
}

func UUID() uuid.UUID {
	return uuid.New()
}

func UUIDP() *uuid.UUID {
	ret := UUID()
	return &ret
}
