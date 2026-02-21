package layout

import "testing"

func TestSolveNonCompact(t *testing.T) {
	r := Solve(140, 40)
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
	if r.Main.W+r.Sidebar.W != 140+nonCompactWidthCompensation {
		t.Fatalf("unexpected content width: main=%d sidebar=%d", r.Main.W, r.Sidebar.W)
	}
	if r.Sidebar.X != r.Main.W {
		t.Fatalf("unexpected sidebar x: %d", r.Sidebar.X)
	}
}

func TestSolveCompact(t *testing.T) {
	r := Solve(80, 20)
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

func TestSolveBreakpointTransition(t *testing.T) {
	c := Solve(compactBreakpointWidth-1, compactBreakpointHeight)
	if !c.Compact {
		t.Fatal("expected compact layout below width breakpoint")
	}
	n := Solve(compactBreakpointWidth, compactBreakpointHeight)
	if n.Compact {
		t.Fatal("expected non-compact layout at width breakpoint")
	}
}
