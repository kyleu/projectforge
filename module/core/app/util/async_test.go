//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"{{{ .Package }}}/app/util"
)

func TestAsyncCollect(t *testing.T) {
	items := []int{1, 2, 3, 4}
	res, errs := util.AsyncCollect(items, func(x int) (int, error) {
		if x == 3 {
			return 0, errors.New("boom")
		}
		return x * x, nil
	})

	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if !strings.Contains(errs[0].Error(), "item [3]") {
		t.Fatalf("expected error to mention item [3], got %q", errs[0].Error())
	}

	sort.Ints(res)
	if res[0] != 1 || res[1] != 4 || res[2] != 16 {
		t.Fatalf("unexpected results: %v", res)
	}
}

func TestAsyncCollectMap(t *testing.T) {
	items := []string{"a", "b", "c"}
	res, errs := util.AsyncCollectMap(items, strings.ToUpper, func(x string) (int, error) {
		if x == "b" {
			return 0, errors.New("boom")
		}
		return len(x), nil
	})

	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if _, ok := res["A"]; !ok {
		t.Fatalf("expected result for key A")
	}
	if _, ok := res["C"]; !ok {
		t.Fatalf("expected result for key C")
	}
	if err, ok := errs["B"]; !ok || !strings.Contains(err.Error(), "item [B]") {
		t.Fatalf("expected error for key B, got %v", errs)
	}
}

func TestAsyncRateLimit(t *testing.T) {
	items := []int{1, 2, 3}
	res, errs := util.AsyncRateLimit("test", items, func(x int) (int, error) {
		time.Sleep(5 * time.Millisecond)
		return x * x, nil
	}, 2, 250*time.Millisecond)

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
}

func TestAsyncRateLimitTimeout(t *testing.T) {
	items := []int{1}
	res, errs := util.AsyncRateLimit("timeout", items, func(x int) (int, error) {
		return x, nil
	}, 0, 10*time.Millisecond)

	if len(res) != 0 {
		t.Fatalf("expected no results, got %v", res)
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %v", errs)
	}
	if !strings.Contains(errs[0].Error(), "timed out") {
		t.Fatalf("expected timeout error, got %q", errs[0].Error())
	}
}
