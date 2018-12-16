package day16

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strconv"
	"strings"
)

type Registers [4]int
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

type Op [4]int

func (o Op) Opcode() int { return o[0] }
func (o Op) A() int { return o[1] }
func (o Op) B() int { return o[2] }
func (o Op) C() int { return o[3] }

type TestCase struct {
	In, Out Registers
	Op Op
}

type Program []Op

func parseRegister(input string) Registers {
	result := Registers{}
	result[0], _ = strconv.Atoi(input[9:10])
	result[1], _ = strconv.Atoi(input[12:13])
	result[2], _ = strconv.Atoi(input[15:16])
	result[3], _ = strconv.Atoi(input[18:19])
	return result
}

func parseOp(input string) Op {
	result := Op{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < 4; i++ {
		scanner.Scan()
		result[i], _ = strconv.Atoi(scanner.Text())
	}
	return result
}

func readInput(filename string) ([]TestCase, Program) {
	lines, err := util.ReadLinesFromFile(filename)
	util.Check(err)

	tests := make([]TestCase, 0)
	program := make(Program, 0)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch {
		case line == "":
			// Empty line, do nothing
		case line[0] == 'B':
			// Read a test case
			test := TestCase{
				In: parseRegister(line),
				Op: parseOp(lines[i+1]),
				Out: parseRegister(lines[i+2]),
			}
			i += 2
			tests = append(tests, test)
		default:
			// Read a bit of the program
			op := parseOp(line)
			program = append(program, op)
		}
	}

	return tests, program
}

func part1(logger *log.Logger) string {
	tests, _ := readInput("day16/input.txt")

	veryAmbiguousCount := 0
	for _, t := range tests {
		opcodeCount := 0
		for _, f := range Operations {
			out := t.In
			f(&t.In, &out, t.Op.A(), t.Op.B(), t.Op.C())
			if out == t.Out {
				opcodeCount++
			}
		}
		if opcodeCount >= 3 {
			veryAmbiguousCount++
		}
	}

	return fmt.Sprint(veryAmbiguousCount)
}

func part2(logger *log.Logger) string {
	tests, program := readInput("day16/input.txt")

	opcodeToFunc := [16]OpFunc{}
	funcToOpcode := make(map[string]int)

	// Repeat test cases until we've mapped all opcodes
	for len(funcToOpcode) < len(opcodeToFunc) {
		for _, t := range tests {
			// Already mapped this opcode, so skip the test case
			if opcodeToFunc[t.Op.Opcode()] != nil {
				continue
			}
			// Find operations that give the right result
			candidates := 0
			var lastFuncName string
			var lastOpFunc OpFunc
			for funcName, f := range Operations {
				// If we already know the opcode for this function, skip it
				if _, ok := funcToOpcode[funcName]; ok {
					continue
				}
				// Compare operation result to test result
				out := t.In
				f(&t.In, &out, t.Op.A(), t.Op.B(), t.Op.C())
				if out == t.Out {
					candidates++
					lastFuncName = funcName
					lastOpFunc = f
				}
			}
			// If there's only one possible operation, we can map it to an opcode
			if candidates == 1 {
				opcodeToFunc[t.Op.Opcode()] = lastOpFunc
				funcToOpcode[lastFuncName] = t.Op.Opcode()
			}
		}
	}

	// Run the program
	registers := Registers{}
	for _, op := range program {
		f := opcodeToFunc[op.Opcode()]
		f(&registers, &registers, op.A(), op.B(), op.C())
	}

	return fmt.Sprint(registers[0])
}

func init() {
	util.RegisterSolution("day16part1", part1)
	util.RegisterSolution("day16part2", part2)
}
