package layout_test

import (
	"testing"

	"projectforge.dev/projectforge/app/controller/tui/layout"
)

func TestSolveNonCompact(t *testing.T) {
	t.Parallel()
	r := layout.Solve(140, 40)
	if r.Compact {
		t.Fatal("expected non-compact layout")
	}
	if r.Main.W <= 0 || r.Sidebar.W <= 0 {
		t.Fatalf("invalid widths: main=%d sidebar=%d", r.Main.W, r.Sidebar.W)
	}
	if r.Header.H != 1 {
		t.Fatalf("unexpected header height: %d", r.Header.H)
	}
	if r.Main.H != 36 {
		t.Fatalf("unexpected main height: %d", r.Main.H)
	}
	if r.Main.W+r.Sidebar.W != 142 {
		t.Fatalf("unexpected content width: main=%d sidebar=%d", r.Main.W, r.Sidebar.W)
	}
	if r.Sidebar.X != r.Main.W {
		t.Fatalf("unexpected sidebar x: %d", r.Sidebar.X)
	}
}

func TestSolveCompact(t *testing.T) {
	t.Parallel()
	r := layout.Solve(80, 20)
	if !r.Compact {
		t.Fatal("expected compact layout")
	}
	if r.Header.H != 3 {
		t.Fatalf("unexpected header height: %d", r.Header.H)
	}
	if r.Main.H <= 0 {
		t.Fatalf("unexpected main height: %d", r.Main.H)
	}
}

func TestSolveNonCompactNoSidebar(t *testing.T) {
	t.Parallel()
	r := layout.SolveWithSidebar(140, 40, false)
	if r.Compact {
		t.Fatal("expected non-compact layout")
	}
	if r.Sidebar.W != 0 {
		t.Fatalf("expected hidden sidebar, got width=%d", r.Sidebar.W)
	}
	if r.Main.W != 140 {
		t.Fatalf("expected full-width main pane, got %d", r.Main.W)
	}
}

func TestSolveBreakpointTransition(t *testing.T) {
	t.Parallel()
	c := layout.Solve(99, 24)
	if !c.Compact {
		t.Fatal("expected compact layout below width breakpoint")
	}
	n := layout.Solve(100, 24)
	if n.Compact {
		t.Fatal("expected non-compact layout at width breakpoint")
	}
}
