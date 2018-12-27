package day21

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/elfcode"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
)

func translatedOptimised(seed int) {
	a, c, d, e := seed, 0, 0, 0

	d = 123
	for d != 72 {
		d &= 456
	}

	d = 0
	for {
		c = d | 65536
		d = 1505483

		// For each byte of c, from last to first (should be 3 bytes)
		// e.g. for first round, c = 65536, then 256, then 1
		for {
			// e.g. for first round, e = 0, 0, 1
			e = c & 255
			// d = (d + e) * 65899, in modulo 2^24
			d = (((d + e) & 16777215) * 65899) & 16777215

			// Ran out of bytes of c
			if c < 256 {
				break
			}
			// Otherwise, shift right by 8 bits
			c /= 256
		}

		// Or, to flatten the above out a bit:
		//d = (((d + (c & 255)) & 16777215) * 65899) & 16777215
		//d = (((d + ((c >> 8) & 255)) & 16777215) * 65899) & 16777215
		//d = (((d + ((c >> 16) & 255)) & 16777215) * 65899) & 16777215

		fmt.Println("d=", d)

		if d == a {
			return
		}
	}
}

/*
The input code exits when register 0 == register 3, so calculate what the value of register 3 will
be at the first point where the equality check is made - this ensures the fewest instructions
executed.
 */
func reverseEngineered() int {
	d := 0
	//c := d | 65536
	d = 1505483
	d = (((d + 0) & 16777215) * 65899) & 16777215
	d = (((d + 0) & 16777215) * 65899) & 16777215
	d = (((d + 1) & 16777215) * 65899) & 16777215
	return d
}

func part1impl(logger *log.Logger, filename string) int {
	// Reverse-engineer the value
	value := reverseEngineered()
	logger.Printf("reverse engineered value: %d\n", value)

	// Verify it terminates
	reader, err := os.Open(filename)
	util.Check(err)
	program := elfcode.ParseProgram(reader, 6)
	processor := elfcode.Processor{Program: &program}
	processor.Init(elfcode.Registers{value})
	for {
		if halted := processor.Step(); halted {
			break
		}
	}
	logger.Printf("register 0 = %d, executed %d instructions\n", value, processor.InstructionCount)

	return value
}

func init() {
	util.RegisterSolution("day21part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day21/input.txt"))
	})
}
