package elfcode

import (
	"bufio"
	"io"
	"strconv"
)

type Registers []int

type Operation func(in, out Registers, a, b, c int)

var Operations = map[string]Operation{
	"addr": func(in, out Registers, a, b, c int) { out[c] = in[a] + in[b] },
	"addi": func(in, out Registers, a, b, c int) { out[c] = in[a] + b },
	"mulr": func(in, out Registers, a, b, c int) { out[c] = in[a] * in[b] },
	"muli": func(in, out Registers, a, b, c int) { out[c] = in[a] * b },
	"banr": func(in, out Registers, a, b, c int) { out[c] = in[a] & in[b] },
	"bani": func(in, out Registers, a, b, c int) { out[c] = in[a] & b },
	"borr": func(in, out Registers, a, b, c int) { out[c] = in[a] | in[b] },
	"bori": func(in, out Registers, a, b, c int) { out[c] = in[a] | b },
	"setr": func(in, out Registers, a, b, c int) { out[c] = in[a] },
	"seti": func(in, out Registers, a, b, c int) { out[c] = a },
	"gtir": func(in, out Registers, a, b, c int) { if a > in[b] { out[c] = 1 } else { out[c] = 0 } },
	"gtri": func(in, out Registers, a, b, c int) { if in[a] > b { out[c] = 1 } else { out[c] = 0 } },
	"gtrr": func(in, out Registers, a, b, c int) { if in[a] > in[b] { out[c] = 1 } else { out[c] = 0 } },
	"eqir": func(in, out Registers, a, b, c int) { if a == in[b] { out[c] = 1 } else { out[c] = 0 } },
	"eqri": func(in, out Registers, a, b, c int) { if in[a] == b { out[c] = 1 } else { out[c] = 0 } },
	"eqrr": func(in, out Registers, a, b, c int) { if in[a] == in[b] { out[c] = 1 } else { out[c] = 0 } },
}

type Instruction struct {
	Op string
	A, B, C int
}

func (i *Instruction) Execute(in, out Registers) {
	Operations[i.Op](in, out, i.A, i.B, i.C)
}

type Program struct {
	RegisterCount int
	IP int
	Code []Instruction
}

func ParseProgram(reader io.Reader, registerCount int) Program {
	p := Program{}
	p.RegisterCount = registerCount
	p.Code = make([]Instruction, 0)
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

func (p *Program) Run(initialState Registers) (Registers, int) {
	instructionCount := 0
	state := make(Registers, p.RegisterCount)
	for i := 0; i < len(state) && i < len(initialState); i++ {
		state[i] = initialState[i]
	}
	ip := &state[p.IP]
	for ; *ip < len(p.Code); *ip++ {
		instructionCount++
		p.Code[*ip].Execute(state, state)
	}
	return state, instructionCount
}
