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
	}

	for _, table := range tables {
		riskTotal := part1impl(logger, table.depth, table.target)
		if riskTotal != table.riskTotal {
			t.Errorf("expected %d, got %d", table.riskTotal, riskTotal)
		}
	}
}
