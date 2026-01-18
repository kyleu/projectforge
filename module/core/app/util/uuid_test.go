//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"regexp"
	"testing"

	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestUUIDFromString(t *testing.T) {
	t.Parallel()

	t.Run("valid UUID string", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDFromString(testUUID)
		if result == nil {
			t.Fatal("expected non-nil UUID for valid string")
		}
		if result.String() != testUUID {
			t.Errorf("expected %s, got %s", testUUID, result.String())
		}
	})

	t.Run("empty string returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDFromString("")
		if result != nil {
			t.Error("expected nil for empty string")
		}
	})

	t.Run("invalid UUID string returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDFromString("not-a-uuid")
		if result != nil {
			t.Error("expected nil for invalid UUID string")
		}
	})
}

func TestUUIDFromStringOK(t *testing.T) {
	t.Parallel()

	t.Run("valid UUID string", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDFromStringOK(testUUID)
		if result.String() != testUUID {
			t.Errorf("expected %s, got %s", testUUID, result.String())
		}
	})

	t.Run("empty string returns default UUID", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDFromStringOK("")
		if result != util.UUIDDefault {
			t.Error("expected default UUID for empty string")
		}
	})

	t.Run("invalid UUID string returns default", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDFromStringOK("not-a-uuid")
		if result != util.UUIDDefault {
			t.Error("expected default UUID for invalid string")
		}
	})
}

func TestUUIDString(t *testing.T) {
	t.Parallel()

	t.Run("nil pointer returns empty string", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDString(nil)
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("valid pointer returns string", func(t *testing.T) {
		t.Parallel()
		u := uuid.MustParse(testUUID)
		result := util.UUIDString(&u)
		if result != testUUID {
			t.Errorf("expected UUID string, got %s", result)
		}
	})
}

func TestUUID(t *testing.T) {
	t.Parallel()

	t.Run("generates valid UUID", func(t *testing.T) {
		t.Parallel()
		result := util.UUID()
		if result == util.UUIDDefault {
			t.Error("expected non-default UUID")
		}
		matched, _ := regexp.MatchString(util.UUIDRegex, result.String())
		if !matched {
			t.Errorf("UUID %s doesn't match expected pattern", result.String())
		}
	})

	t.Run("generates unique UUIDs", func(t *testing.T) {
		t.Parallel()
		u1 := util.UUID()
		u2 := util.UUID()
		if u1 == u2 {
			t.Error("expected different UUIDs")
		}
	})
}

func TestUUIDP(t *testing.T) {
	t.Parallel()

	t.Run("returns non-nil pointer", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDP()
		if result == nil {
			t.Error("expected non-nil pointer")
		}
	})

	t.Run("returns valid UUID", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDP()
		matched, _ := regexp.MatchString(util.UUIDRegex, result.String())
		if !matched {
			t.Errorf("UUID %s doesn't match expected pattern", result.String())
		}
	})
}

func TestUUIDV7(t *testing.T) {
	t.Parallel()

	t.Run("generates valid UUIDv7", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDV7()
		matched, _ := regexp.MatchString(util.UUIDRegex, result.String())
		if !matched {
			t.Errorf("UUIDv7 %s doesn't match expected pattern", result.String())
		}
	})

	t.Run("generates unique UUIDv7s", func(t *testing.T) {
		t.Parallel()
		u1 := util.UUIDV7()
		u2 := util.UUIDV7()
		if u1 == u2 {
			t.Error("expected different UUIDv7s")
		}
	})
}

func TestUUIDV7P(t *testing.T) {
	t.Parallel()

	t.Run("returns non-nil pointer", func(t *testing.T) {
		t.Parallel()
		result := util.UUIDV7P()
		if result == nil {
			t.Error("expected non-nil pointer")
		}
	})
}

func TestUUIDRegex(t *testing.T) {
	t.Parallel()

	t.Run("matches valid UUIDs", func(t *testing.T) {
		t.Parallel()
		re := regexp.MustCompile(util.UUIDRegex)
		validUUIDs := []string{
			testUUID,
			"00000000-0000-0000-0000-000000000000",
			"ffffffff-ffff-ffff-ffff-ffffffffffff",
		}
		for _, u := range validUUIDs {
			if !re.MatchString(u) {
				t.Errorf("expected %s to match UUID regex", u)
			}
		}
	})

	t.Run("does not match invalid UUIDs", func(t *testing.T) {
		t.Parallel()
		re := regexp.MustCompile("^" + util.UUIDRegex + "$")
		invalidUUIDs := []string{
			"not-a-uuid",
			"550e8400-e29b-41d4-a716",
			"550e8400e29b41d4a716446655440000",
		}
		for _, u := range invalidUUIDs {
			if re.MatchString(u) {
				t.Errorf("expected %s to not match UUID regex", u)
			}
		}
	})
}
