package util

import "github.com/google/uuid"

const UUIDRegex = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

var UUIDDefault = uuid.UUID{}

func UUIDFromString(s string) *uuid.UUID {
	if s != "" {
		s, err := uuid.Parse(s)
		if err == nil {
			return &s
		}
	}
	return nil
}

func UUIDFromStringOK(s string) uuid.UUID {
	if s != "" {
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

func UUIDV7() uuid.UUID {
	ret, _ := uuid.NewV7()
	return ret
}

func UUIDV7P() *uuid.UUID {
	ret := UUIDV7()
	return &ret
}
