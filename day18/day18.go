package day18

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strings"
)

const (
	Open       = '.'
	Trees      = '|'
	Lumberyard = '#'
)

type Forest struct {
	Map           util.ByteGrid
	Width, Height int
	Time int
}

func NewForest(input []string) Forest {
	f := Forest{}
	f.Width = len(input[0])
	f.Height = len(input)
	f.Map = util.NewByteGrid(f.Width, f.Height)
	for y, line := range input {
		for x := range line {
			f.Map[x][y] = line[x]
		}
	}
	f.Time = 0
	return f
}

func (f *Forest) String() string {
	sb := strings.Builder{}
	sb.Grow((f.Width + 1) * f.Height)
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			sb.WriteByte(f.Map[x][y])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (f *Forest) Valid(p util.Vec2D) bool {
	return f.Map.Valid(p)
}

func (f *Forest) At(p util.Vec2D) *byte {
	return f.Map.At(p)
}

func (f *Forest) CountAdjacent(p util.Vec2D) map[byte]int {
	result := make(map[byte]int)
	adjacentPoints := [8]util.Vec2D{
		p.Add(util.Vec2D{-1, -1}),
		p.Add(util.Vec2D{0, -1}),
		p.Add(util.Vec2D{1, -1}),
		p.Add(util.Vec2D{-1, 0}),
		p.Add(util.Vec2D{1, 0}),
		p.Add(util.Vec2D{-1, 1}),
		p.Add(util.Vec2D{0, 1}),
		p.Add(util.Vec2D{1, 1}),
	}
	for _, a := range adjacentPoints {
		if f.Valid(a) {
			result[*f.At(a)]++
		}
	}
	return result
}

func (f *Forest) CountAll() map[byte]int {
	result := make(map[byte]int)
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			result[f.Map[x][y]]++
		}
	}
	return result
}

func (f *Forest) AdvanceTime() {
	f.Time++
	newMap := util.NewByteGrid(f.Width, f.Height)
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			p := util.Vec2D{x, y}
			counts := f.CountAdjacent(p)
			switch *f.At(p) {
			case Open:
				if counts[Trees] >= 3 {
					*newMap.At(p) = Trees
				} else {
					*newMap.At(p) = Open
				}
			case Trees:
				if counts[Lumberyard] >= 3 {
					*newMap.At(p) = Lumberyard
				} else {
					*newMap.At(p) = Trees
				}
			case Lumberyard:
				if counts[Lumberyard] >= 1 && counts[Trees] >= 1 {
					*newMap.At(p) = Lumberyard
				} else {
					*newMap.At(p) = Open
				}
			}
		}
	}
	f.Map = newMap
}

func part1impl(logger *log.Logger, filename string, duration int) (trees, lumberyards int) {
	input, err := util.ReadLinesFromFile(filename)
	util.Check(err)
	forest := NewForest(input)
	//logger.Print("start:\n", forest.String())
	for i := 0; i < duration; i++ {
		forest.AdvanceTime()
		//logger.Print("t = ", forest.Time, ":\n", forest.String())
	}
	counts := forest.CountAll()
	return counts[Trees], counts[Lumberyard]
}

func init() {
	util.RegisterSolution("day18test1", func(logger *log.Logger) string {
		trees, lumberyards := part1impl(logger, "day18/input_test.txt", 10)
		return fmt.Sprintf("%d x %d = %d", trees, lumberyards, trees*lumberyards)
	})
	util.RegisterSolution("day18part1", func(logger *log.Logger) string {
		trees, lumberyards := part1impl(logger, "day18/input.txt", 10)
		return fmt.Sprintf("%d x %d = %d", trees, lumberyards, trees*lumberyards)
	})
}
