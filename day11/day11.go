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

type FuelGrid struct {
	SerialNo int
	Grid [][]int
}

func NewFuelGrid(serialNo int) FuelGrid {
	fg := FuelGrid{SerialNo: serialNo}
	raw := make([]int, SIZE*SIZE)
	fg.Grid = make([][]int, SIZE)
	for x := 0; x < SIZE; x++ {
		var row []int
		row, raw = raw[:SIZE], raw[SIZE:]
		fg.Grid[x] = row
		for y := 0; y < SIZE; y++ {
			row[y] = fg.CalcCellPower(x + 1, y + 1)
		}
	}
	return fg
}

func (fg *FuelGrid) CalcCellPower(x, y int) int {
	rackId := x + 10
	result := ((rackId * y) + fg.SerialNo) * rackId
	result =((result % 1000) / 100) - 5
	return result
}

func (fg *FuelGrid) CellPower(x, y int) int {
	return fg.Grid[x-1][y-1]
}

func (fg *FuelGrid) GroupPower(x, y, size int) int {
	initX, initY := x, y
	maxX, maxY := x + size - 1, y + size - 1
	result := 0
	for x = initX; x <= maxX; x++ {
		for y = initY; y <= maxY; y++ {
			power := fg.CellPower(x, y)
			result += power
		}
	}
	return result
}

/*
EdgePower finds the sum of power values along the right and bottom edges of the
area defined by (minX, minY) to (maxX, maxY).
 */
func (fg *FuelGrid) EdgePower(minX, maxX, minY, maxY int) int {
	result := 0
	for x := minX; x <= maxX; x++ {
		result += fg.CellPower(x, maxY)
	}
	// Avoid double-counting (maxX, maxY)!
	for y := minY; y < maxY; y++ {
		result += fg.CellPower(maxX, y)
	}
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
	var power int
	x, y, power = fg.FindBestGroup(3)
	t.Printf("found best fuel cell group at %v,%v power %v", x, y, power)
	return x, y
}

func part2impl(logger *log.Logger, serialNo int) (x, y, size int) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	fg := NewFuelGrid(serialNo)
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
