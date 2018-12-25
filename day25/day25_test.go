package day25

import (
	"log"
	"os"
	"testing"
)

func TestPart1Impl(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	tables := []struct{
		filename string
		result int
	}{
		{"input_test1.txt", 2},
		{"input_test2.txt", 4},
		{"input_test3.txt", 3},
		{"input_test4.txt", 8},
	}

	for _, table := range tables {
		result := part1impl(logger, table.filename)
		if result != table.result {
			t.Errorf("%s: expected %d, got %d", table.filename, table.result, result)
		}
	}
}
