package day22

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
)

const (
	// Terrain types
	Rocky = 0
	Wet = 1
	Narrow = 2
	// Equipment types, numbered so that terrain==equipment means incompatible
	Neither = 0
	Torch = 1
	ClimbingGear = 2
)

func Compatible(terrain, equipment byte) bool {
	return terrain != equipment
}

type State struct {
	Position util.Vec2D
	Equipment byte
}

func (s1 State) LessThan(s2 State) bool {
	switch {
	case s1.Position.Y < s2.Position.Y:
		return true
	case s1.Position.Y > s2.Position.Y:
		return false
	case s1.Position.X < s2.Position.X:
		return true
	case s1.Position.X > s2.Position.X:
		return false
	case s1.Equipment < s2.Equipment:
		return true
	default:
		return false
	}
}

func MakeErosionMap(depth int, target util.Vec2D, scale util.Vec2D) util.IntGrid {
	erosion := util.NewIntGrid(target.X*scale.X+1, target.Y*scale.Y+1)

	geologicIndex := func(p util.Vec2D) int {
		switch {
		case p == util.Vec2D{0, 0}, p == target:
			return 0
		case p.Y == 0:
			return p.X * 16807
		case p.X == 0:
			return p.Y * 48271
		default:
			// Assumes erosion grid is filled from top-left corner,
			// so these values will already exist
			left := p.Add(util.Vec2D{-1, 0})
			up := p.Add(util.Vec2D{0, -1})
			return *erosion.At(left) * *erosion.At(up)
		}
	}

	erosionLevel := func(p util.Vec2D) int {
		return (geologicIndex(p) + depth) % 20183
	}

	erosion.Traverse(func(p util.Vec2D, data *int) {
		*data = erosionLevel(p)
	})

	return erosion
}

func MakeTerrainMap(erosion util.IntGrid) util.ByteGrid {
	terrain := util.NewByteGrid(erosion.Width(), erosion.Height())
	erosion.Traverse(func(p util.Vec2D, data *int) {
		*terrain.At(p) = byte(*data % 3)
	})
	return terrain
}

/*
Sum the terrain/risk value of every square, in a map from (0, 0) to target (inclusive).
 */
func part1impl(logger *log.Logger, depth int, target util.Vec2D) int {
	erosion := MakeErosionMap(depth, target, util.Vec2D{1, 1})
	terrain := MakeTerrainMap(erosion)

	sum := 0
	terrain.Traverse(func(p util.Vec2D, data *byte) {
		sum += int(*data)
	})

	return sum
}

//go:generate genny -in=../util/astar.go -out=gen-astar.go -pkg=day22 gen "SearchNode=State"
/*
Find the shortest path from (0, 0) to target, taking equipment into account.

TODO: combine equipment change and navigation? equipment needs to be valid in both new and old location!
 */
func part2impl(logger *log.Logger, depth int, target util.Vec2D) int {
	scale := util.Vec2D{1, 1}
	var erosion util.IntGrid
	var terrain util.ByteGrid
	regenerate := func() {
		logger.Printf("generating map: %d x %d", target.X*scale.X+1, target.Y*scale.Y+1)
		erosion = MakeErosionMap(depth, target, scale)
		terrain = MakeTerrainMap(erosion)
	}
	regenerate()

	valid := func(p util.Vec2D) bool {
		// This is definitely valid
		if terrain.Valid(p) {
			return true
		}
		// This will never be valid
		if p.X < 0 || p.Y < 0 {
			return false
		}
		// If p isn't valid because it's beyond the right and/or bottom,
		// expand the grid until it will be enclosed
		logger.Printf("expanding to include %v", p)
		for ; p.X >= target.X*scale.X+1; scale.X++ {}
		for ; p.Y >= target.Y*scale.Y+1; scale.Y++ {}
		regenerate()
		// Should be valid now
		return terrain.Valid(p)
	}

	search := AStarSearchContext{
		Start: State{
			Position:  util.Vec2D{0, 0},
			Equipment: Torch,
		},
		Destinations: []State{
			{
				Position:  target,
				Equipment: Torch,
			},
		},
		NodeCountMax: terrain.Width() * terrain.Height() * 2,
		Adjacent: func(n State) []State {
			result := make([]State, 0)
			offsets := []util.Vec2D{
				{1, 0},
				{0, 1},
				{-1, 0},
				{0, -1},
			}
			equipment := [3]byte{Neither, Torch, ClimbingGear}
			// Moving to adjacent locations with the current equipment
			for _, offset := range offsets {
				next := n
				next.Position.AddInPlace(offset)
				if valid(next.Position) && Compatible(*terrain.At(next.Position), next.Equipment) {
					result = append(result, next)
				}
			}
			// Changing equipment at the current location
			for _, equip := range equipment {
				if equip != n.Equipment && Compatible(*terrain.At(n.Position), equip) {
					next := n
					next.Equipment = equip
					result = append(result, next)
				}
			}
			return result
		},
		Heuristic: func(n1, n2 State) int {
			return n2.Position.Sub(n1.Position).Manhattan()
		},
		Cost: func(n1, n2 State) int {
			cost := n2.Position.Sub(n1.Position).Manhattan()
			if n1.Equipment != n2.Equipment {
				cost += 7
			}
			return cost
		},
		TieBreak: State.LessThan,
	}

	path, err := AStarSearch(&search)
	util.Check(err)

	prev := search.Start
	cost := 0
	for _, next := range path {
		cost += search.Cost(prev, next)
		prev = next
	}
	return cost
}

func init() {
	util.RegisterSolution("day22part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, 3879, util.Vec2D{8, 713}))
	})
	util.RegisterSolution("day22part2", func(logger *log.Logger) string {
		return fmt.Sprint(part2impl(logger, 3879, util.Vec2D{8, 713}))
	})
}
