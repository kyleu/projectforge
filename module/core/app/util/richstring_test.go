//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestRichString(t *testing.T) {
	s := util.Str(" Hello ")
	if !s.Contains("ell") {
		t.Fatalf("expected Contains to match")
	}
	if s.TrimSpace().String() != "Hello" {
		t.Fatalf("unexpected TrimSpace result: %q", s.TrimSpace())
	}
	left, right := s.TrimSpace().Cut('/', true)
	if left.String() != "Hello" || right.String() != "" {
		t.Fatalf("unexpected Cut result: %q %q", left, right)
	}

	path := util.Str("root").Path("child", "leaf")
	if path.String() != "root/child/leaf" {
		t.Fatalf("unexpected Path result: %q", path)
	}

	prefixed := util.Str("a").WithPrefix(util.Str("b"), util.Str("-"))
	if prefixed.String() != "b-a" {
		t.Fatalf("unexpected WithPrefix result: %q", prefixed)
	}

	suffixed := util.Str("a").WithSuffix(util.Str("-"), util.Str("c"))
	if suffixed.String() != "a-c" {
		t.Fatalf("unexpected WithSuffix result: %q", suffixed)
	}
}
