package day14

import (
	"log"
	"os"
	"testing"
)

func TestPart1Impl(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	tables := []struct{
		previous int
		slice int
		result []byte
	}{
		{9, 10, []byte{5, 1, 5, 8, 9, 1, 6, 7, 7, 9}},
		{5, 10, []byte{0, 1, 2, 4, 5, 1, 5, 8, 9, 1}},
		{18, 10, []byte{9, 2, 5, 1, 0, 7, 1, 0, 8, 5}},
		{2018, 10, []byte{5, 9, 4, 1, 4, 2, 9, 8, 8, 2}},
	}
	
	for _, table := range tables {
		result := part1impl(logger, table.previous, table.slice)
		if string(result) != string(table.result) {
			t.Errorf("expected %v, got %v", table.result, result)
		}
	}
}

func TestPart2Impl(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	tables := []struct{
		match []byte
		result int
	}{
		{[]byte{5, 1, 5, 8, 9}, 9},
		{[]byte{0, 1, 2, 4, 5}, 5},
		{[]byte{9, 2, 5, 1, 0}, 18},
		{[]byte{5, 9, 4, 1, 4}, 2018},
	}

	for _, table := range tables {
		result := part2impl(logger, table.match)
		if result != table.result {
			t.Errorf("expected %v, got %v", table.result, result)
		}
	}
}


