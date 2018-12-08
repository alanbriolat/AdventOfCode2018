package util

import (
	"log"
	"sort"
)

type Implementation func(logger *log.Logger) string

type Solution struct {
	Name string
	Run  Implementation
}

var solutions = make([]Solution, 0)

func RegisterSolution(name string, run Implementation) {
	solutions = append(solutions, Solution{name, run})
}

func GetSolutions() []Solution {
	result := make([]Solution, len(solutions))
	copy(result, solutions)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}
