package day14

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
)

func part1impl(logger *log.Logger, previous int, slice int) []byte {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")
	
	// How many recipes we need in total = sample length + number of preceding recipes
	recipeCount := slice + previous
	// Pre-allocate upper-bound of how many recipes there can be - account for potential overshoot
	// of 1 if the last round is a double-digit score
	recipes := make([]byte, 0, recipeCount + 1)
	recipes = append(recipes, 3, 7)
	elves := [2]int{0, 1}
	
	// Run the first round differently, because it's special
	recipes = append(recipes, 1, 0)
	
	// Run rounds until we have enough recipes
	for len(recipes) < recipeCount {
		sum := byte(0)
		// Find the recipe score total
		for _, i := range elves {
			sum += recipes[i]
		}
		// Turn the score total into more recipes
		for _, x := range []byte(fmt.Sprint(sum)) {
			recipes = append(recipes, x - '0') // Convert ASCII digit to integer
		}
		// Move the elves
		for i := range elves {
			elves[i] = (elves[i] + 1 + int(recipes[elves[i]])) % len(recipes)
		}
	}
	
	return recipes[previous:previous+slice]
}

func part1(logger *log.Logger, previous int, slice int) string {
	result := part1impl(logger, previous, slice)
	for i := range result {
		result[i] += '0'
	}
	return string(result)
}

func init() {
	util.RegisterSolution("day14part1", func(logger *log.Logger) string {
		return part1(logger, 635041, 10)
	})
}
