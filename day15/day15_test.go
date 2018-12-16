package day15

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestPart1Impl(t *testing.T) {
	logger := log.New(ioutil.Discard, "", 0)
	tables := []struct{
		input []string
		rounds int
		remainingHP int
	}{
		{
			[]string{
				"#######",
				"#.G...#",
				"#...EG#",
				"#.#.#G#",
				"#..G#E#",
				"#.....#",
				"#######",
			},
			47, 590,
		},
		{
			[]string{
				"#######",
				"#G..#E#",
				"#E#E.E#",
				"#G.##.#",
				"#...#E#",
				"#...E.#",
				"#######",
			},
			37, 982,
		},
		{
			[]string{
				"#######",
				"#E..EG#",
				"#.#G.E#",
				"#E.##E#",
				"#G..#.#",
				"#..E#.#",
				"#######",
			},
			46, 859,
		},
		{
			[]string{
				"#######",
				"#E.G#.#",
				"#.#G..#",
				"#G.#.G#",
				"#G..#.#",
				"#...E.#",
				"#######",
			},
			35, 793,
		},
		{
			[]string{
				"#######",
				"#.E...#",
				"#.#..G#",
				"#.###.#",
				"#E#G#G#",
				"#...#G#",
				"#######",
			},
			54, 536,
		},
		{
			[]string{
				"#########",
				"#G......#",
				"#.E.#...#",
				"#..##..G#",
				"#...##..#",
				"#...#...#",
				"#.G...G.#",
				"#.....G.#",
				"#########",
			},
			20, 937,
		},
	}

	for _, table := range tables {
		rounds, remainingHP := part1impl(logger, table.input, table.rounds + 10, false)
		if rounds != table.rounds || remainingHP != table.remainingHP {
			t.Errorf("expected %dx%d, got %dx%d", table.rounds, table.remainingHP, rounds, remainingHP)
		}
	}
}
