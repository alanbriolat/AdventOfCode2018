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

type FuelGrid struct {
	SerialNo int
}

func (fg *FuelGrid) CellPower(x, y int) int {
	rackId := x + 10
	result := ((rackId * y) + fg.SerialNo) * rackId
	result =((result % 1000) / 100) - 5
	return result
}

func (fg *FuelGrid) GroupPower(x, y int) int {
	initX, initY := x, y
	maxX, maxY := x + 2, y + 2
	result := 0
	for x = initX; x <= maxX; x++ {
		for y = initY; y <= maxY; y++ {
			power := fg.CellPower(x, y)
			result += power
		}
	}
	return result
}

func (fg *FuelGrid) FindBestGroup() (bestX, bestY, bestPower int) {
	bestX, bestY, bestPower = 0, 0, math.MinInt32
	for x := 1; x <= (300 - 2); x++ {
		for y := 1; y <= (300 - 2); y++ {
			if power := fg.GroupPower(x, y); power > bestPower {
				bestX, bestY, bestPower = x, y, power
			}
		}
	}
	return
}

func impl(logger *log.Logger, serialNo int) (x, y int) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	fg := FuelGrid{SerialNo: serialNo}
	var power int
	x, y, power = fg.FindBestGroup()
	t.Printf("found best fuel cell group at %v,%v power %v", x, y, power)
	return x, y
}

func part1(logger *log.Logger) string {
	x, y := impl(logger, 7315)
	return fmt.Sprint(x, ",", y)
}

func init() {
	util.RegisterSolution("day11part1", part1)
}
