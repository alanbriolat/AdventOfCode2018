package day15

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
	"testing"
)

func TestPart1Impl(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	tables := []struct{
		input []string
		rounds int
		remainingHP int
	}{
		// Puzzle description worked example
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
		// Puzzle description samples
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
		// https://www.reddit.com/r/adventofcode/comments/a6f100/day_15_details_easy_to_be_wrong_on/ebvkuxr/
		{
			[]string{
				"####",
				"##E#",
				"#GG#",
				"####",
			},
			67, 200,
		},
		{
			[]string{
				"#####",
				"#GG##",
				"#.###",
				"#..E#",
				"#.#G#",
				"#.E##",
				"#####",
			},
			71, 197,
		},
		{
			[]string{
				"################",
				"#.......G......#",
				"#G.............#",
				"#..............#",
				"#....###########",
				"#....###########",
				"#.......EG.....#",
				"################",
			},
			38, 486,
		},
	}

	for _, table := range tables {
		rounds, remainingHP := part1impl(logger, table.input, table.rounds + 10, false)
		if rounds != table.rounds || remainingHP != table.remainingHP {
			t.Errorf("expected %dx%d, got %dx%d", table.rounds, table.remainingHP, rounds, remainingHP)
		}
	}
}

// https://www.reddit.com/r/adventofcode/comments/a6f100/day_15_details_easy_to_be_wrong_on/ebvkuxr/
func TestEdgeCase(t *testing.T) {
	input := []string{
		"#######",
		"#.E..G#",
		"#.#####",
		"#G#####",
		"#######",
	}
	b := NewBattle(input)
	b.NextRound()
	fmt.Println(b.String())
	expected := util.Vec2D{3, 1}
	if b.Units[0].Position != expected {
		t.Errorf("expected %v, got %v", expected, b.Units[0].Position)
	}
}

func TestEqualPathReadingOrder(t *testing.T) {
	input := []string{
		"###########",
		"#...#.....#",
		"#G#.......#",
		"#...#.....#",
		"#####.....#",
		"#####....E#",
		"###########",
	}
	b := NewBattle(input)
	b.NextRound()
	fmt.Println(b.String())
	expected := util.Vec2D{1, 1}
	if b.Units[0].Position != expected {
		t.Errorf("expected %v, got %v", expected, b.Units[0].Position)
	}
}
