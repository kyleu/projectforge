//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

var (
	errRoot = errors.New("root")
	errA    = errors.New("a")
	errB    = errors.New("b")
)

func TestErrorDetailAndMerge(t *testing.T) {
	t.Parallel()

	wrapped := fmt.Errorf("wrap: %w", errRoot)

	detail := util.GetErrorDetail(wrapped, true)
	if detail == nil || detail.Message == "" {
		t.Fatalf("expected error detail")
	}
	if detail.Cause == nil || detail.Cause.Message != "root" {
		t.Fatalf("expected wrapped cause, got %+v", detail.Cause)
	}

	if util.ErrorMerge(nil, nil) != nil {
		t.Fatalf("expected nil error merge to return nil")
	}

	merged := util.ErrorMerge(nil, errRoot, nil)
	if !errors.Is(merged, errRoot) {
		t.Fatalf("expected single error merge to return original")
	}

	err := util.ErrorMerge(errA, errB)
	if err == nil || !strings.Contains(err.Error(), "a") || !strings.Contains(err.Error(), "b") {
		t.Fatalf("unexpected merged error: %v", err)
	}
}
