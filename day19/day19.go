package day19

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
	"strconv"
)

const (
	RegisterCount = 6
)

type Registers [RegisterCount]int

type OpFunc func(in, out *Registers, a, b, c int)

var Operations = map[string]OpFunc{
	"addr": func(in, out *Registers, a, b, c int) { out[c] = in[a] + in[b] },
	"addi": func(in, out *Registers, a, b, c int) { out[c] = in[a] + b },
	"mulr": func(in, out *Registers, a, b, c int) { out[c] = in[a] * in[b] },
	"muli": func(in, out *Registers, a, b, c int) { out[c] = in[a] * b },
	"banr": func(in, out *Registers, a, b, c int) { out[c] = in[a] & in[b] },
	"bani": func(in, out *Registers, a, b, c int) { out[c] = in[a] & b },
	"borr": func(in, out *Registers, a, b, c int) { out[c] = in[a] | in[b] },
	"bori": func(in, out *Registers, a, b, c int) { out[c] = in[a] | b },
	"setr": func(in, out *Registers, a, b, c int) { out[c] = in[a] },
	"seti": func(in, out *Registers, a, b, c int) { out[c] = a },
	"gtir": func(in, out *Registers, a, b, c int) { if a > in[b] { out[c] = 1 } else { out[c] = 0 } },
	"gtri": func(in, out *Registers, a, b, c int) { if in[a] > b { out[c] = 1 } else { out[c] = 0 } },
	"gtrr": func(in, out *Registers, a, b, c int) { if in[a] > in[b] { out[c] = 1 } else { out[c] = 0 } },
	"eqir": func(in, out *Registers, a, b, c int) { if a == in[b] { out[c] = 1 } else { out[c] = 0 } },
	"eqri": func(in, out *Registers, a, b, c int) { if in[a] == b { out[c] = 1 } else { out[c] = 0 } },
	"eqrr": func(in, out *Registers, a, b, c int) { if in[a] == in[b] { out[c] = 1 } else { out[c] = 0 } },
}

type Instruction struct {
	Op string
	A, B, C int
}

type Program struct {
	IP int
	Code []Instruction
}

func readInput(filename string) Program {
	p := Program{}
	p.Code = make([]Instruction, 0)
	reader, err := os.Open(filename)
	util.Check(err)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	scanner.Scan()
	p.IP, _ = strconv.Atoi(scanner.Text())
	for scanner.Scan() {
		inst := Instruction{}
		inst.Op = scanner.Text()
		scanner.Scan()
		inst.A, _ = strconv.Atoi(scanner.Text())
		scanner.Scan()
		inst.B, _ = strconv.Atoi(scanner.Text())
		scanner.Scan()
		inst.C, _ = strconv.Atoi(scanner.Text())
		p.Code = append(p.Code, inst)
	}
	return p
}

/*
Run the program, emulating the instructions, and return the final state.
 */
func emulated(logger *log.Logger, filename string, initialState Registers) Registers {
	program := readInput(filename)
	state := initialState

	for ; state[program.IP] < len(program.Code); state[program.IP]++ {
		inst := program.Code[state[program.IP]]
		f := Operations[inst.Op]
		f(&state, &state, inst.A, inst.B, inst.C)
		//newState := state
		//f(&state, &newState, inst.A, inst.B, inst.C)
		//logger.Printf("ip=%d %v %v %v", state[program.IP], state, inst, newState)
		//state = newState
	}

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
	//	return fmt.Sprint(emulated(logger, "day19/input.txt", Registers{})[0])
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
