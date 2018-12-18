/*
Day 12
======

This problem looks familiar! It's a 1D cellular automaton with a neighbourhood
of 5, which means 2^5 = 32 possible patterns.
 */
package day12

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strings"
)

const (
	True = '#'
	False = '.'
	PatternSize = 5						// How wide the pattern is
	PatternCount = 1 << PatternSize		// Number of possible patterns
	Growth = PatternSize - 1			// How much the state grows at each end, each generation
)

type CellularAutomaton struct {
	Patterns [PatternCount]byte
	State string
	Generation int
	GrowthString string
}

func NewCellularAutomaton(initialState string, patterns []string) CellularAutomaton {
	ca := CellularAutomaton{}
	ca.GrowthString = strings.Repeat(string(False), Growth)
	ca.State = initialState
	// Build pattern table for rapid matching
	for i := range ca.Patterns {
		// Fill the table with false values, because example input only includes true values
		ca.Patterns[i] = False
	}
	for _, p := range patterns {
		i := patternIndex(p[0:5])
		if p[9] == False {
			ca.Patterns[i] = p[9]
		} else {
			// give the patterns "names" to make patterns easier to spot
			ca.Patterns[i] = byte(';' + i)
		}
	}
	return ca
}

func (ca *CellularAutomaton) String() string {
	return fmt.Sprintf("%02d: %s", ca.Generation, ca.State)
}

func (ca *CellularAutomaton) Apply(group string) byte {
	return ca.Patterns[patternIndex(group)]
}

func (ca *CellularAutomaton) NextGeneration() (int, string) {
	gen := ca.Generation + 1
	old := strings.Join([]string{ca.GrowthString, ca.State, ca.GrowthString}, "")
	b := strings.Builder{}
	b.Grow(len(ca.State) + Growth)
	for i := 0; i < len(old) - Growth; i++ {
		group := old[i:i+PatternSize]
		b.WriteByte(ca.Apply(group))
	}
	return gen, b.String()
}

func (ca *CellularAutomaton) Advance() {
	ca.Generation, ca.State = ca.NextGeneration()
}

func (ca *CellularAutomaton) IndexSum() int {
	offset := ca.Generation * Growth / 2
	sum := 0
	for i, v := range ca.State {
		if v != False {
			sum += i - offset
		}
	}
	return sum
}

/*
patternIndex treats a string of length PatternSize as a binary string, turning
it into an integer.
 */
func patternIndex(pattern string) int {
	index := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != False {
			index |= 1 << uint8(len(pattern)-i-1)
		}
	}
	return index
}

func readInput(filename string) CellularAutomaton {
	lines, err := util.ReadLinesFromFile(filename)
	util.Check(err)
	ca := NewCellularAutomaton(lines[0][15:], lines[2:])
	return ca
}

func part1(logger *log.Logger, filename string, generations int) int {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	ca := readInput(filename)
	t.LogCheckpoint("read input")

	sum := ca.IndexSum()
	sumDiff := sum
	sumDiffDiff := sum
	var i int
	for i = 0; i < generations && sumDiffDiff != 0; i++ {
		ca.Advance()
		newSum := ca.IndexSum()
		newSumDiff := newSum - sum
		sumDiffDiff = newSumDiff - sumDiff
		sum = newSum
		sumDiff = newSumDiff
	}
	remaining := generations - i
	t.Printf("ran %d generations, fast-forwarding by %d", i, remaining)
	sum += remaining * sumDiff

	return sum
}

func init() {
	//util.RegisterSolution("day12part0", func(logger *log.Logger) string {
	//	return fmt.Sprint(part1(logger,"day12/input_test.txt", 20))
	//})
	util.RegisterSolution("day12part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1(logger,"day12/input.txt", 20))
	})

	util.RegisterSolution("day12part2", func(logger *log.Logger) string {
		return fmt.Sprint(part1(logger,"day12/input.txt", 50000000000))
	})
}
