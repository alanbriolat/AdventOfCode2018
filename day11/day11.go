package day11

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"math"
)

type Point struct {
	X, Y int
}

const (
	SIZE = 300
)

type Grid [][]int

func makeGrid(w, h int) Grid {
	raw := make([]int, w*h)
	grid := make([][]int, w)
	for x := 0; x < w; x++ {
		grid[x], raw = raw[:w], raw[w:]
	}
	return grid
}

type FuelGrid struct {
	SerialNo int
	Grid Grid
	Area Grid
}

func NewFuelGrid(serialNo int) FuelGrid {
	fg := FuelGrid{SerialNo: serialNo}
	fg.Grid = makeGrid(SIZE, SIZE)
	fg.Area = makeGrid(SIZE, SIZE)
	for x := 0; x < SIZE; x++ {
		for y, colSum := 0, 0; y < SIZE; y++ {
			power := fg.CalcCellPower(x + 1, y + 1)
			fg.Grid[x][y] = power
			colSum += power
			fg.Area[x][y] = colSum
			if x > 0 {
				fg.Area[x][y] += fg.Area[x-1][y]
			}
		}
	}
	return fg
}

func (fg *FuelGrid) CalcCellPower(x, y int) int {
	rackId := x + 10
	result := ((rackId * y) + fg.SerialNo) * rackId
	result = ((result % 1000) / 100) - 5
	return result
}

func (fg *FuelGrid) CellPower(x, y int) int {
	return fg.Grid[x-1][y-1]
}

func (fg *FuelGrid) GroupPower(x, y, size int) int {
	// Adjust to grid coordinates
	x -= 1
	y -= 1
	// Calculate bottom-right corner
	maxX := x+size-1
	maxY := y+size-1
	// Find sum at bottom right corner
	result := fg.Area[maxX][maxY]
	// Find sum before bottom left
	if x > 0 { result -= fg.Area[x-1][maxY] }
	// Find sum before top right
	if y > 0 { result -= fg.Area[maxX][y-1] }
	// Find sum before top left
	if x > 0 && y > 0 { result += fg.Area[x-1][y-1] }
	return result
}

func (fg *FuelGrid) FindBestGroup(size int) (bestX, bestY, bestPower int) {
	bestX, bestY, bestPower = 0, 0, math.MinInt32
	maxX := 300 - size + 1
	maxY := maxX
	for x := 1; x <= maxX; x++ {
		for y := 1; y <= maxY; y++ {
			if power := fg.GroupPower(x, y, size); power > bestPower {
				bestX, bestY, bestPower = x, y, power
			}
		}
	}
	return
}

func (fg *FuelGrid) FindBestGroupAnySize() (bestX, bestY, bestSize, bestPower int) {
	x, y, size, bestPower := 0, 0, 0, math.MinInt32
	for i := 1; i <= 300; i++ {
		groupX, groupY, groupPower := fg.FindBestGroup(i)
		if groupPower > bestPower {
			x, y, size, bestPower = groupX, groupY, i, groupPower
		}
	}
	return x, y, size, bestPower
}

func part1impl(logger *log.Logger, serialNo int) (x, y int) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	fg := NewFuelGrid(serialNo)
	t.Printf("generated grid")
	var power int
	x, y, power = fg.FindBestGroup(3)
	t.Printf("found best fuel cell group at %v,%v power %v", x, y, power)
	return x, y
}

func part2impl(logger *log.Logger, serialNo int) (x, y, size int) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	fg := NewFuelGrid(serialNo)
	t.Printf("generated grid")
	x, y, size, bestPower := fg.FindBestGroupAnySize()
	t.Printf("found best fuel cell group at %v,%v,%v power %v", x, y, size, bestPower)
	return x, y, size
}

func part1(logger *log.Logger) string {
	x, y := part1impl(logger, 7315)
	return fmt.Sprint(x, ",", y)
}

func part2(logger *log.Logger) string {
	x, y, size := part2impl(logger, 7315)
	return fmt.Sprint(x, ",", y, ",", size)
}

func init() {
	util.RegisterSolution("day11part1", part1)
	util.RegisterSolution("day11part2", part2)
}
