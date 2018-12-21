package day19

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/elfcode"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
)

const (
	RegisterCount = 6
)

func readInput(filename string) elfcode.Program {
	reader, err := os.Open(filename)
	util.Check(err)
	return elfcode.ParseProgram(reader, RegisterCount)
}

/*
Run the program, emulating the instructions, and return the final state.
 */
func emulated(logger *log.Logger, filename string, initialState elfcode.Registers) elfcode.Registers {
	program := readInput(filename)
	state, _ := program.Run(initialState)
	return state
}

/*
A re-implementation of what the instructions in input.txt do: sum the factors of a number.
 */
func translated(logger *log.Logger, seed int) int {
	var c int
	if seed == 0 {
		c = 877
	} else {
		c = 10551277
	}
	a := 0
	for d := 1; d <= c; d++ {
		for b := 1; b <= c; b++ {
			if d * b == c {
				a += d
			}
		}
	}
	return a
}

/*
A faster re-implementation that takes O(n) time instead of O(n^2).
 */
func translatedOptimised(logger *log.Logger, seed int) int {
	var c int
	if seed == 0 {
		c = 877
	} else {
		c = 10551277
	}
	// 1 and c are always going to be factors
	a := 1 + c
	// The second-largest factor cannot be larger than c/2
	for d := 2; d <= c/2; d++ {
		// No remainder means it's a factor
		if c % d == 0 {
			a += d
		}
	}
	return a
}

func init() {
	//util.RegisterSolution("day19test1emu", func(logger *log.Logger) string {
	//	return fmt.Sprint(emulated(logger, "day19/input_test.txt", Registers{}))
	//})

	//util.RegisterSolution("day19part1emu", func(logger *log.Logger) string {
	//	return fmt.Sprint(emulated(logger, "day19/input.txt", elfcode.Registers{})[0])
	//})
	//util.RegisterSolution("day19part1trans", func(logger *log.Logger) string {
	//	return fmt.Sprint(translated(logger, 0))
	//})
	util.RegisterSolution("day19part1opt", func(logger *log.Logger) string {
		return fmt.Sprint(translatedOptimised(logger, 0))
	})

	//util.RegisterSolution("day19part2emu", func(logger *log.Logger) string {
	//	return fmt.Sprint(emulated(logger, "day19/input.txt", Registers{1})[0])
	//})
	//util.RegisterSolution("day19part2trans", func(logger *log.Logger) string {
	//	return fmt.Sprint(translated(logger, 1))
	//})
	util.RegisterSolution("day19part2opt", func(logger *log.Logger) string {
		return fmt.Sprint(translatedOptimised(logger, 1))
	})
}
