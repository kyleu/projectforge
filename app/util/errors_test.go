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

func TestErrorDetailAndMerge(t *testing.T) {
	root := errors.New("root")
	wrapped := fmt.Errorf("wrap: %w", root)

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

	merged := util.ErrorMerge(nil, root, nil)
	if merged != root {
		t.Fatalf("expected single error merge to return original")
	}

	err := util.ErrorMerge(errors.New("a"), errors.New("b"))
	if err == nil || !strings.Contains(err.Error(), "a") || !strings.Contains(err.Error(), "b") {
		t.Fatalf("unexpected merged error: %v", err)
	}
}
