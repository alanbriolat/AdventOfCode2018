package day22

import (
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
	"testing"
)

func TestPart1Impl(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	tables := []struct {
		depth     int
		target    util.Vec2D
		riskTotal int
	}{
		{
			510,
			util.Vec2D{10, 10},
			114,
		},
		{
			3879,
			util.Vec2D{8, 713},
			6323,
		},
	}

	for _, table := range tables {
		riskTotal := part1impl(logger, table.depth, table.target)
		if riskTotal != table.riskTotal {
			t.Errorf("expected %d, got %d", table.riskTotal, riskTotal)
		}
	}
}

func TestPart2Impl(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	tables := []struct {
		depth     int
		target    util.Vec2D
		shortestPath int
	}{
		{
			510,
			util.Vec2D{10, 10},
			45,
		},
		{
			3879,
			util.Vec2D{8, 713},
			982,
		},
	}

	for _, table := range tables {
		shortestPath := part2impl(logger, table.depth, table.target)
		if shortestPath != table.shortestPath {
			t.Errorf("expected %d, got %d", table.shortestPath, shortestPath)
		}
	}
}
