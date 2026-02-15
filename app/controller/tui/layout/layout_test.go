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
	if r.Main.H != 35 {
		t.Fatalf("unexpected main height: %d", r.Main.H)
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
