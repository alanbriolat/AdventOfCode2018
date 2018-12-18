/*
Day 10
======

Part 1
------

Using character recognition wouldn't work, because we don't know what all the
characters are supposed to look like. Instead we need a heuristic of when to
stop.

One such heuristic might be: when every star has at least one adjacent star.
This would hopefully detect when the stars have all been arranged into lines.
Might get away with just 4-direction adjacency, but some character
representations might use diagonal lines, which would force 8-direction
adjacency instead.

Another possible heuristic: minimise the bounding box area instead, on the
assumption that the message appears when the stars are at their most converged.

Convergence on a solution may not be direct, maybe need a bit of randomised
hill climbing? Will find out by trying a naive solution first.
 */
package day10

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strconv"
	"strings"
)

type Star struct {
	Position, Velocity util.Vec2D
}

func ParseStar(star *Star, input string) (err error) {
	if star.Position.X, err = strconv.Atoi(strings.TrimSpace(string(input[10:16]))); err != nil { return err }
	if star.Position.Y, err = strconv.Atoi(strings.TrimSpace(string(input[18:24]))); err != nil { return err }
	if star.Velocity.X, err = strconv.Atoi(strings.TrimSpace(string(input[36:38]))); err != nil { return err }
	if star.Velocity.Y, err = strconv.Atoi(strings.TrimSpace(string(input[40:42]))); err != nil { return err }
	return nil
}

type StarField struct {
	Stars []Star
	Time int
	Lookup map[util.Vec2D]struct{}
	Min, Max util.Vec2D
}

func NewStarField(lines []string) (result StarField, err error) {
	result = StarField{}
	result.Stars = make([]Star, len(lines))
	for i, line := range lines {
		if err = ParseStar(&result.Stars[i], line); err != nil {
			return result, err
		}
	}
	// Make sure lookup has been generated
	result.TimeTravel(0)
	return result, nil
}

func (sf *StarField) TimeTravel(time int) {
	sf.Time = time
	sf.Min = util.MaxVec2D()
	sf.Max = util.MinVec2D()
	sf.Lookup = make(map[util.Vec2D]struct{})
	for i := range sf.Stars {
		s := &sf.Stars[i]
		p := s.Position.Add(s.Velocity.Scale(time))
		sf.Lookup[p] = struct{}{}
		sf.Min.MinInPlace(p)
		sf.Max.MaxInPlace(p)
	}
}

func (sf *StarField) Show(star string, space string) string {
	b := strings.Builder{}
	for y := sf.Min.Y; y <= sf.Max.Y; y++ {
		for x := sf.Min.X; x <= sf.Max.X; x++ {
			if _, ok := sf.Lookup[util.Vec2D{x, y}]; ok {
				b.WriteString(star)
			} else {
				b.WriteString(space)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (sf *StarField) Area() int {
	return (sf.Max.X - sf.Min.X) * (sf.Max.Y - sf.Min.Y)
}

/*
simpleHillClimbing searches the solution space linearly as the fitness improves,
until the fitness worsens instead, returning the index of the first best
solution found.

size is the length of the problem space.

fitness(i) is a function that calculates the fitness for index i.

compare(a,b) is a function that evaluates if the fitness has improved or
	worsened; a positive result means improved, negative means worsened, zero
	means no change.
 */
func simpleHillClimbing(size int, fitness func(i int)int, compare func(a, b int)int) (int, error) {
	bestIndex := 0
	value := fitness(bestIndex)
	newValue := value

	for i := 0; i < size; i++ {
		newValue = fitness(i)
		check := compare(value, newValue)
		switch {
		case check < 0:
			return bestIndex, nil
		case check > 0:
			bestIndex = i
			value = newValue
		}
	}

	return fitness(size - 1), fmt.Errorf("did not find optimum within first %v solutions", size)
}

/*
oscillate searches forwards and backwards with a resolution that decays
exponentially, switching direction each time the fitness gets worse between
two consecutive points. Once the resolution is at its minimum and the
inflection point is encountered again, the solution has been found.
 */
func oscillate(resolution int, fitness func(i int)int, compare func(a, b int)int) (int, error) {
	index := 0
	value := fitness(index)
	newValue := value

	for {
		for compare(value, newValue) >= 0 {
			value = newValue
			index += resolution
			newValue = fitness(index)
		}
		if util.AbsInt(resolution) > 1 {
			resolution = -(resolution / 2)
			value = newValue
		} else {
			return index - resolution, nil
		}
	}
}

func part1impl(logger *log.Logger, filename string) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	lines, err := util.ReadLinesFromFile(filename)
	util.Check(err)
	t.Printf("read %v lines", len(lines))

	starField, err := NewStarField(lines)
	util.Check(err)
	t.Printf("read %v stars", len(starField.Stars))

	//time, err := simpleHillClimbing(
	//	math.MaxInt32,
	//	func(i int)int { starField.TimeTravel(i); return starField.Area() },
	//	func(a, b int)int { return a - b },
	//)
	time, err := oscillate(
		1000,
		func(i int)int { starField.TimeTravel(i); return starField.Area() },
		func(a, b int)int { return a - b },
	)
	util.Check(err)

	starField.TimeTravel(time)
	return fmt.Sprint("after ", time, " time steps:\n", starField.Show("#", " "))
}

func part0(logger *log.Logger) string {
	return part1impl(logger, "day10/input_test.txt")
}

func part1and2(logger *log.Logger) string {
	return part1impl(logger, "day10/input.txt")
}

func init() {
	//util.RegisterSolution("day10part0", part0)
	util.RegisterSolution("day10", part1and2)
}
