package main

import "testing"

func TestLinearOverlap(t *testing.T) {
	tables := []struct {
		name string
		p1, d1 int	// First range
		p2, d2 int	// Second range
		p, d int	// Expected overlap range
	}{
		{"a entirely before b",
			2, 4, 8, 3, 0, 0},
		{"a ends at start of b",
			4, 4, 8, 3, 0, 0},
		{"a overlaps start of b",
			5, 4, 8, 3, 8, 1},
		{"a contained within b, same start",
			8, 2, 8, 3, 8, 2},
		{"a contained within b",
			9, 2, 8, 4, 9, 2},
		{"b contained within a",
			3, 10, 5, 3, 5, 3},
		{"a contained within b, same end",
			10, 2, 8, 4, 10, 2},
		{"a overlaps end of b",
			10, 3, 8, 4, 10, 2},
		{"a starts at end of b",
			12, 3, 8, 4, 0, 0},
		{"a entirely after b",
			14, 3, 4, 3, 0, 0},
		{"zero-length a",
			5, 0, 3, 4, 0, 0},
		{"zero-length b",
			5, 3, 6, 0, 0, 0},
		{"zero-length a and b",
			5, 0, 3, 0, 0, 0},
	}
	for _, table := range tables {
		p, d := LinearOverlap(table.p1, table.d1, table.p2, table.d2)
		if p != table.p || d != table.d {
			t.Errorf("%v: expected %v + %v, got %v + %v",
				table.name, table.p, table.d, p, d)
		}
	}
}

func TestClaim_Overlap(t *testing.T) {
	a := Claim{1, 3, 4, 4}
	b := Claim{3, 1, 4, 4}
	expected := Claim{3, 3, 2, 2}
	overlap, err := a.Overlap(b)
	switch {
	case err != nil:
		t.Error("unexpected error", err)
	case overlap != expected:
		t.Error("expected", expected, "got", overlap)
	}
}