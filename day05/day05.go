/*
Looks like accumulating characters as a stack, and checking if the top two can annihilate
 */
package main

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io/ioutil"
	"unicode"
)

func CanReact(a, b rune) bool {
	// Swap the case of one byte
	if unicode.IsUpper(b) {
		b = unicode.ToLower(b)
	} else {
		b = unicode.ToUpper(b)
	}
	// If they match now, they were opposite to start with
	return a == b
}

func React(bytes []byte) []byte {
	stack := util.NewStack(len(bytes))
	for _, b := range bytes {
		top, ok := stack.Peek()
		if ok && CanReact(rune(top.(byte)), rune(b)) {
			stack.Pop()
		} else {
			stack.Push(b)
		}
	}
	result := make([]byte, len(stack.Data))
	for i, b := range stack.Data {
		result[i] = b.(byte)
	}
	return result
}

func StripUnit(bytes []byte, unit byte) []byte {
	upper := unit &^ (1 << 5)
	result := make([]byte, 0, len(bytes))
	for _, b := range bytes {
		if upper != b &^ (1 << 5) {
			result = append(result, b)
		}
	}
	return result
}

func ReadInput(name string) []byte {
	bytes, err := ioutil.ReadFile(name)
	util.Check(err)
	// strip newline
	bytes = bytes[:len(bytes)-1]
	return bytes
}

func part1() {
	t := util.NewTimer("day05part1")
	defer t.PrintCheckpoint("end")

	bytes := ReadInput("input1.txt")
	t.PrintCheckpoint(fmt.Sprint("read ", len(bytes), " bytes"))

	polymer := React(bytes)
	fmt.Printf("polymer is %v units long\n", len(polymer))
	t.PrintCheckpoint("reacted polymer")
}

func part2() {
	t := util.NewTimer("day05part2")
	defer t.PrintCheckpoint("end")

	bytes := ReadInput("input1.txt")
	t.PrintCheckpoint(fmt.Sprint("read ", len(bytes), " bytes"))

	shortest := len(bytes)
	best := byte(' ')
	for b := byte('A'); b <= byte('Z'); b++ {
		polymer := React(StripUnit(bytes, b))
		if len(polymer) < shortest {
			shortest = len(polymer)
			best = b
		}
	}
	fmt.Printf("shortest polymer is %v units long (after removing %v)\n", shortest, best)
	t.PrintCheckpoint("found shortest possible polymer")
}

func main() {
	part1()
	part2()
}
