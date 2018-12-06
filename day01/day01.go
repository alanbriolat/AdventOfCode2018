package main

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
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


func part1and2() {
	changes, err := util.ReadIntsFromFile("input1.txt")
	util.Check(err)
	state := State()

	for _, x := range changes {
		state.Update(x)
	}
	fmt.Println("Resulting Frequency:", state.Frequency)
	for !state.FoundRepeat {
		for _, x := range changes {
			state.Update(x)
			if state.FoundRepeat {
				break
			}
		}
	}
	fmt.Println("First repeated Frequency:", state.Repeat)
}

func main() {
	part1and2()
}