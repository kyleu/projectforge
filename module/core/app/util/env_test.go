//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"
	"time"

	"{{{ .Package }}}/app/util"
)

func TestEnvHelpers(t *testing.T) {
	t.Setenv("PF_TEST_UPPER", "upper")
	t.Setenv("pf_test_lower", "lower")
	t.Setenv("PF_BOOL", "true")
	t.Setenv("PF_INT", "42")
	t.Setenv("PF_DUR", "2s")

	if got := util.GetEnv("PF_TEST_UPPER"); got != "upper" {
		t.Fatalf("expected upper, got %q", got)
	}
	if got := util.GetEnv("pf_test_lower"); got != "lower" {
		t.Fatalf("expected lower, got %q", got)
	}
	if got := util.GetEnv("missing", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}

	if !util.GetEnvBool("PF_BOOL", false) {
		t.Fatalf("expected bool true")
	}
	if util.GetEnvBoolAny(false, "missing", "PF_BOOL") != true {
		t.Fatalf("expected bool true from GetEnvBoolAny")
	}
	if got := util.GetEnvInt("PF_INT", 1); got != 42 {
		t.Fatalf("expected int 42, got %d", got)
	}
	if got := util.GetEnvDuration("PF_DUR", time.Second); got != 2*time.Second {
		t.Fatalf("expected 2s, got %v", got)
	}

	res := util.ReplaceEnvVars("x=${PF_TEST_UPPER|fallback} y=${MISSING|dflt}", nil)
	if res != "x=upper y=dflt" {
		t.Fatalf("unexpected ReplaceEnvVars result: %q", res)
	}
}
