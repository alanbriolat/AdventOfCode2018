package day12

import (
	"testing"
)

func TestPatternIndex(t *testing.T) {
	tables := []struct{
		pattern string
		index int
	}{
		{".....", 0},
		{"....#", 1},
		{"...#.", 2},
		{"..#..", 4},
		{".#...", 8},
		{"#....", 16},
		{"#####", 31},
	}

	for _, table := range tables {
		index := patternIndex(table.pattern)
		if index != table.index {
			t.Errorf("expected %v, got %v", table.index, index)
		}
	}
}
