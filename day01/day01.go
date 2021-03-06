package day01

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
)

type state struct {
	Frequency   int
	seen        map[int]bool
	FoundRepeat bool
	Repeat      int
}

func State() state {
	return state{0,map[int]bool{0: true},false,0}
}

func (s *state) Update(change int) {
	s.Frequency += change
	if !s.FoundRepeat && s.seen[s.Frequency] {
		s.FoundRepeat = true
		s.Repeat = s.Frequency
	}
	s.seen[s.Frequency] = true
}


func part1and2(logger *log.Logger) string {
	changes, err := util.ReadIntsFromFile("day01/input1.txt")
	util.Check(err)
	state := State()

	for _, x := range changes {
		state.Update(x)
	}
	finalFrequency := state.Frequency
	logger.Println("Resulting Frequency:", finalFrequency)
	for !state.FoundRepeat {
		for _, x := range changes {
			state.Update(x)
			if state.FoundRepeat {
				break
			}
		}
	}
	logger.Println("First repeated Frequency:", state.Repeat)

	return fmt.Sprintf("part1 = %v , part2 = %v", finalFrequency, state.Repeat)
}

func init() {
	util.RegisterSolution("day01", part1and2)
}