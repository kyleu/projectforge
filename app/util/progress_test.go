//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"
	"time"

	"projectforge.dev/projectforge/app/util"
)

func TestProgressAndTimer(t *testing.T) {
	var called bool
	p := util.NewProgress("work", 0, func(p *util.Progress, delta int) {
		if delta == 5 {
			called = true
		}
	})
	p.Increment(5, nil)
	if !called {
		t.Fatalf("expected progress callback to be called")
	}
	if p.Total != 100 || p.Completed != 5 {
		t.Fatalf("unexpected progress values: %+v", p)
	}
	if p.String() != "work (5 of 100)" {
		t.Fatalf("unexpected progress string: %q", p.String())
	}

	var nilProgress *util.Progress
	nilProgress.Increment(1, nil)

	tr := util.TimerStart()
	time.Sleep(2 * time.Millisecond)
	_ = tr.End()
	if tr.Completed == 0 || tr.String() == "" {
		t.Fatalf("unexpected timer state: %+v", tr)
	}
}
