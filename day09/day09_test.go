package day09

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestDay09(t *testing.T) {
	logger := log.New(ioutil.Discard, "", 0)

	tables := []struct {
		players, max, score int
	}{
		{9, 25, 32},
		{441, 71032, 393229},
		{441, 7103200, 3273405195},
	}

	for _, table := range tables {
		result := part1impl(logger, table.players, table.max)
		if result != table.score {
			t.Errorf("%+v != %v", table, result)
		}
	}
}