package day20

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strings"
)

const (
	Indent = "  "
)

type Visitor func(s string)

type Expr interface {
	fmt.Stringer
	/*
	Visit every prefix with f, return final strings.

	e.g. ABCD -> f(A), f(AB), f(ABC), f(ABCD), return {ABCD,}
	 */
	Enumerate(prefix string, f Visitor) []string
	/*
	Count how many final strings will be generated
	 */
	StringCount() int
	/*
	Count how many syntax tree nodes there are
	 */
	SyntaxNodeCount() int

	BuildTreeString(builder *strings.Builder, depth int)
}

type Literal string
type Choice []Expr
type Sequence []Expr

func (e Literal) Enumerate(prefix string, f Visitor) []string {
	if f == nil {
		return []string{prefix + string(e)}
	}
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

func (e Literal) StringCount() int {
	return 1
}

func (e Literal) SyntaxNodeCount() int {
	return 1
}

func (e Literal) BuildTreeString(builder *strings.Builder, depth int) {
	builder.WriteString(strings.Repeat(Indent, depth))
	builder.WriteString("Literal(\"")
	builder.WriteString(string(e))
	builder.WriteString("\")")
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

func (e Choice) StringCount() int {
	sum := 0
	for _, next := range e {
		sum += next.StringCount()
	}
	return sum
}

func (e Choice) SyntaxNodeCount() int {
	sum := 1
	for _, next := range e {
		sum += next.SyntaxNodeCount()
	}
	return sum
}

func (e Choice) BuildTreeString(builder *strings.Builder, depth int) {
	builder.WriteString(strings.Repeat(Indent, depth))
	builder.WriteString(fmt.Sprintf("Choice{  // depth=%d\n", depth))
	for _, next := range e {
		next.BuildTreeString(builder, depth+1)
		builder.WriteString(",\n")
	}
	builder.WriteString(strings.Repeat(Indent, depth))
	builder.WriteString("}")
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

func (e Sequence) StringCount() int {
	prod := 1
	for _, next := range e {
		prod *= next.StringCount()
	}
	return prod
}

func (e Sequence) SyntaxNodeCount() int {
	sum := 1
	for _, next := range e {
		sum += next.SyntaxNodeCount()
	}
	return sum
}

func (e Sequence) BuildTreeString(builder *strings.Builder, depth int) {
	builder.WriteString(strings.Repeat(Indent, depth))
	builder.WriteString(fmt.Sprintf("Sequence{  // depth=%d\n", depth))
	for _, next := range e {
		next.BuildTreeString(builder, depth+1)
		builder.WriteString(",\n")
	}
	builder.WriteString(strings.Repeat(Indent, depth))
	builder.WriteString("}")
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

func directionToVec2D(d byte) util.Vec2D {
	switch d {
	case 'N':
		return util.Vec2D{0, -1}
	case 'E':
		return util.Vec2D{1, 0}
	case 'S':
		return util.Vec2D{0, 1}
	case 'W':
		return util.Vec2D{-1, 0}
	default:
		panic("invalid direction")
	}
}

func traverse(distances map[util.Vec2D]int, position util.Vec2D, distance int, expr Expr) (util.Vec2D, int) {
	switch expr := expr.(type) {
	case Literal:
		for i := range expr {
			position.AddInPlace(directionToVec2D(expr[i]))
			distance++
			if prevDistance, ok := distances[position]; ok {
				distance = util.MinInt(distance, prevDistance)
			}
			distances[position] = distance
		}
		return position, distance
	case Sequence:
		for _, e := range expr {
			position, distance = traverse(distances, position, distance, e)
		}
		return position, distance
	case Choice:
		for _, e := range expr {
			// We assume choices are either detours (return to starting point) or terminal, rather than true branches,
			// so no need to figure out what the resulting position and distance of a choice is.
			traverse(distances, position, distance, e)
		}
		return position, distance
	default:
		panic("can't traverse() unknown type")
	}
}

func RoomStats(regex string, threshold int) (int, int) {
	expr := ReadExpression(strings.NewReader(regex))

	distances := make(map[util.Vec2D]int)
	traverse(distances, util.Vec2D{0, 0}, 0, expr)

	maxDistance := -1
	thresholdCount := 0
	for _, distance := range distances {
		maxDistance = util.MaxInt(maxDistance, distance)
		if distance >= threshold {
			thresholdCount++
		}
	}

	return maxDistance, thresholdCount
}

func init() {
	util.RegisterSolution("day20", func(logger *log.Logger) string {
		lines, err := util.ReadLinesFromFile("day20/input.txt")
		util.Check(err)
		maxDistance, thresholdCount := RoomStats(lines[0], 1000)
		return fmt.Sprintf("part1 = %d, part2 = %d", maxDistance, thresholdCount)
	})
}
