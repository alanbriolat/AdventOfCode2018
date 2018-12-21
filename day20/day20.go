/*
Part 1
======

Naive approach:
- Expand the regex to every prefix of every string it can match
- For each of those strings, reduce it by removing "redundant" operations
- Update the end point with the path length, if it's shorter than a previous path to that location
- Find the room with the highest path length
 */
package day20

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"math"
	"strings"
)

type Visitor func(s string)

type Expr interface {
	fmt.Stringer
	/*
	Visit every prefix with f, return final strings.

	e.g. ABCD -> f(A), f(AB), f(ABC), f(ABCD), return {ABCD,}
	  */
	Enumerate(prefix string, f Visitor) []string
}

type Literal string
type Choice []Expr
type Sequence []Expr

func (e Literal) Enumerate(prefix string, f Visitor) []string {
	sb := strings.Builder{}
	sb.WriteString(prefix)
	for i := range e {
		sb.WriteByte(e[i])
		f(sb.String())
	}
	return []string{sb.String()}
}

func (e Literal) String() string {
	return string(e)
}

func (e Choice) Enumerate(prefix string, f Visitor) []string {
	result := make([]string, 0)
	for _, next := range e {
		result = append(result, next.Enumerate(prefix, f)...)
	}
	return result
}

func (e Choice) String() string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	for i, next := range e {
		if i > 0 {
			sb.WriteByte('|')
		}
		sb.WriteString(next.String())
	}
	sb.WriteByte(')')
	return sb.String()
}

func (e Sequence) Enumerate(prefix string, f Visitor) []string {
	prefixes := []string{prefix}
	for _, next := range e {
		newPrefixes := make([]string, 0)
		for _, prefix := range prefixes {
			newPrefixes = append(newPrefixes, next.Enumerate(prefix, f)...)
		}
		prefixes = newPrefixes
	}
	return prefixes
}

func (e Sequence) String() string {
	sb := strings.Builder{}
	for _, next := range e {
		sb.WriteString(next.String())
	}
	return sb.String()
}

func ReadLiteral(reader *strings.Reader) Literal {
	sb := strings.Builder{}
readLoop:
	for reader.Len() > 0 {
		b, _ := reader.ReadByte()
		switch b {
		case '(', '|', ')', '$':
			reader.UnreadByte()
			break readLoop
		default:
			sb.WriteByte(b)
		}
	}
	return Literal(sb.String())
}

func ReadChoice(reader *strings.Reader) Choice {
	choice := Choice{}
	for reader.Len() > 0 {
		b, _  := reader.ReadByte()
		switch b {
		case '(', '|':
			// Drop the separator, read an option
			choice = append(choice, ReadExpression(reader))
		case ')':
			// Finished all options for this Choice
			return choice
		default:
			// Read another option, including this byte
			reader.UnreadByte()
			choice = append(choice, ReadExpression(reader))
		}
	}
	return nil
}

func ReadSequence(reader *strings.Reader) Sequence {
	seq := Sequence{}
readLoop:
	for reader.Len() > 0 {
		b, _ := reader.ReadByte()
		switch b {
		case '^', '$':
			// Drop start and end characters
			break
		case '(':
			// Found start of a Choice
			reader.UnreadByte()
			seq = append(seq, ReadChoice(reader))
		case '|', ')':
			// Found the end of a sequence
			reader.UnreadByte()
			break readLoop
		default:
			// Otherwise, read a literal, including this byte (may be zero-length if b terminates it)
			reader.UnreadByte()
			seq = append(seq, ReadLiteral(reader))
		}
	}
	if len(seq) == 0 {
		// Special case: a sequence that has no items is an empty string.
		// E.g. in the case of (A|), the second option is an empty sequence.
		seq = append(seq, Literal(""))
	}
	return seq
}

func ReadExpression(reader *strings.Reader) Expr {
	seq := ReadSequence(reader)
	if len(seq) == 1 {
		// Strip unnecessary nesting
		return seq[0]
	} else {
		return seq
	}
}

func EnumerateExpression(expr Expr, visitPartial Visitor, visitFinal Visitor) {
	finals := expr.Enumerate("", visitPartial)
	for _, final := range finals {
		visitFinal(final)
	}
}

func SimplifyPath(directions string) string {
	stack := util.NewByteStack(len(directions))
	for i := range directions {
		stack.Push(directions[i])
		// Remove clockwise/anti-clockwise cycles
		if top, ok := stack.PeekMany(4); ok {
			pattern := string(top)
			switch pattern {
			case "NESW", "ESWN", "SWNE", "WNES",	// Clockwise
				 "WSEN", "NWSE", "ENWS", "SENW":	// Anti-clockwise
				stack.PopMany(4)
				continue
			}
		}
		// Remove up/down and left/right oscillations
		if top, ok := stack.PeekMany(2); ok {
			pattern := string(top)
			switch pattern {
			case "NS", "SN", "WE", "EW":
				stack.PopMany(2)
				continue
			}
		}
	}
	return string(stack.Data)
}

func ResolvePath(directions string) util.Vec2D {
	result := util.Vec2D{0, 0}
	for i := range directions {
		switch directions[i] {
		case 'N':
			result.Y--
		case 'E':
			result.X++
		case 'S':
			result.Y++
		case 'W':
			result.X--
		}
	}
	return result
}

func FurthestRoom(regex string) int {
	// Shortest path to each room
	rooms := make(map[util.Vec2D]string)

	visitPartial := func (path string) {
		// Remove redundancy from path, to find shortest version of the path
		simplified := SimplifyPath(path)
		// Find the room the path ends at
		room := ResolvePath(simplified)
		// Record first/shorter path to the room
		if oldPath, ok := rooms[room]; !ok || len(simplified) < len(oldPath) {
			rooms[room] = simplified
		}
	}
	visitFinal := func (path string) {
	}

	// Run every path
	expr := ReadExpression(strings.NewReader(regex))
	EnumerateExpression(expr, visitPartial, visitFinal)

	// Find the room where the longest shortest path
	furthestDistance := math.MinInt32
	//var furthestRoom util.Vec2D
	for _, path := range rooms {
		if len(path) > furthestDistance {
			furthestDistance = len(path)
			//furthestRoom = room
		}
	}
	return furthestDistance
}

func init() {
	util.RegisterSolution("day20part1", func(logger *log.Logger) string {
		lines, err := util.ReadLinesFromFile("day20/input.txt")
		util.Check(err)
		return fmt.Sprint(FurthestRoom(lines[0]))
	})
}
