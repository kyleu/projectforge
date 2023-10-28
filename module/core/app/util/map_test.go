//go:build test_all || !func_test
package util_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func assert(t *testing.T, name string, a any, b any, err error, messages ...string) {
	t.Helper()
	if a == b {
		return
	}
	msg := fmt.Sprintf("%s [%v != %v]", name, a, b)
	lo.ForEach(messages, func(m string, _ int) {
		if len(m) > 0 {
			msg += "  " + m
		}
	})
	t.Fatal(msg, err)
}

var tm = util.ValueMap{
	"s":  "str",
	"i":  42,
	"b":  true,
	"u":  util.UUIDFromString("00000000-0000-0000-0000-000000000042"),
	"a":  []any{"a", "b"},
	"ae": []any{},
}

func TestValueMap(t *testing.T) {
	t.Parallel()

	s, err := tm.ParseString("s", true, true)
	assert(t, "string-t-t", s, "str", err)
	s, err = tm.ParseString("s", false, false)
	assert(t, "string-f-f", s, "str", err)
	s, err = tm.ParseString("sx", false, false)
	assert(t, "string-f-f-x", s, "", err)

	i, err := tm.ParseInt("i", true, true)
	assert(t, "int-t-t", i, 42, err)
	i, err = tm.ParseInt("i", true, false)
	assert(t, "int-t-f", i, 42, err)
	i, err = tm.ParseInt("i", false, false)
	assert(t, "int-f-f", i, 42, err)
	i, err = tm.ParseInt("ix", true, true)
	assert(t, "int-f-f-x", i, 0, err)

	ut := tm["u"]
	u, err := tm.ParseUUID("u", true, true)
	assert(t, "uuid-t-t", u, ut, err)
	u, err = tm.ParseUUID("u", true, false)
	assert(t, "uuid-t-f", u, ut, err)
	u, err = tm.ParseUUID("u", false, false)
	assert(t, "uuid-f-f", u, ut, err)
	u, err = tm.ParseUUID("ux", true, true)
	var un *uuid.UUID
	assert(t, "uuid-f-f-x", u, un, err)

	a, err := tm.ParseArray("a", true, true, false)
	assert(t, "array-t-t", a[0], "a", err)
	assert(t, "array-t-t", a[1], "b", err)

	_, err = tm.ParseArray("ae", false, false, false)
	if err == nil {
		t.Fatal(err)
	}
}
