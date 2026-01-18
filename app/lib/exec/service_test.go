package exec_test

import (
	"sync"
	"testing"

	"projectforge.dev/projectforge/app/lib/exec"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	s := exec.NewService()
	if s == nil {
		t.Fatal("expected non-nil service")
	}
	if len(s.Execs) != 0 {
		t.Errorf("expected empty execs, got %d", len(s.Execs))
	}
}

func TestServiceNewExec_CreatesAndTracks(t *testing.T) {
	t.Parallel()
	s := exec.NewService()
	e := s.NewExec("build", "make all", "/project", false)

	if e == nil {
		t.Fatal("expected non-nil exec")
	}
	if e.Key != "build" {
		t.Errorf("expected key 'build', got %q", e.Key)
	}
	if e.Idx != 1 {
		t.Errorf("expected idx 1 (first exec), got %d", e.Idx)
	}
	if e.Cmd != "make all" {
		t.Errorf("expected cmd 'make all', got %q", e.Cmd)
	}
	if len(s.Execs) != 1 {
		t.Errorf("expected 1 tracked exec, got %d", len(s.Execs))
	}
}

func TestServiceNewExec_IncrementsIdxPerKey(t *testing.T) {
	t.Parallel()
	s := exec.NewService()

	e1 := s.NewExec("build", "make", ".", false)
	e2 := s.NewExec("build", "make test", ".", false)
	e3 := s.NewExec("deploy", "deploy.sh", ".", false)
	e4 := s.NewExec("build", "make clean", ".", false)

	if e1.Idx != 1 {
		t.Errorf("e1: expected idx 1, got %d", e1.Idx)
	}
	if e2.Idx != 2 {
		t.Errorf("e2: expected idx 2, got %d", e2.Idx)
	}
	if e3.Idx != 1 {
		t.Errorf("e3: expected idx 1 (different key), got %d", e3.Idx)
	}
	if e4.Idx != 3 {
		t.Errorf("e4: expected idx 3, got %d", e4.Idx)
	}
}

func TestServiceNewExec_WithEnvVars(t *testing.T) {
	t.Parallel()
	s := exec.NewService()
	e := s.NewExec("build", "make", ".", true, "DEBUG=1", "CI=true")

	if !e.Debug {
		t.Error("expected debug true")
	}
	if len(e.Env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(e.Env))
	}
}

func TestServiceNewExec_MaintainsSortedOrder(t *testing.T) {
	t.Parallel()
	s := exec.NewService()

	s.NewExec("zebra", "cmd", ".", false)
	s.NewExec("alpha", "cmd", ".", false)
	s.NewExec("beta", "cmd", ".", false)

	// After adding, should be sorted
	if s.Execs[0].Key != "alpha" {
		t.Errorf("expected first key 'alpha', got %q", s.Execs[0].Key)
	}
	if s.Execs[1].Key != "beta" {
		t.Errorf("expected second key 'beta', got %q", s.Execs[1].Key)
	}
	if s.Execs[2].Key != "zebra" {
		t.Errorf("expected third key 'zebra', got %q", s.Execs[2].Key)
	}
}

func TestServiceNewExec_ConcurrentSafety(t *testing.T) {
	t.Parallel()
	s := exec.NewService()
	var wg sync.WaitGroup
	const numGoroutines = 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			s.NewExec("concurrent", "cmd", ".", false)
		}()
	}
	wg.Wait()

	if len(s.Execs) != numGoroutines {
		t.Errorf("expected %d execs, got %d", numGoroutines, len(s.Execs))
	}

	// Verify all have unique idx values
	idxSeen := make(map[int]bool)
	for _, e := range s.Execs {
		if idxSeen[e.Idx] {
			t.Errorf("duplicate idx %d found", e.Idx)
		}
		idxSeen[e.Idx] = true
	}
}
