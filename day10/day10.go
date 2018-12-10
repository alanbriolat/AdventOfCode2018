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
assumption that the stars are converging.

Convergence on a solution may not be direct, maybe need a bit of randomised
hill climbing.
 */
package day10

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"math"
	"strconv"
	"strings"
)

type Vec2D struct {
	X, Y int
}

func MaxVec2D() Vec2D {
	return Vec2D{math.MaxInt32, math.MaxInt32}
}

func MinVec2D() Vec2D {
	return Vec2D{math.MinInt32, math.MinInt32}
}

func (v *Vec2D) Add(o Vec2D) Vec2D {
	result := *v
	result.AddInPlace(o)
	return result
}

func (v *Vec2D) Sub(o Vec2D) Vec2D {
	result := *v
	result.SubInPlace(o)
	return result
}

func (v *Vec2D) Scale(s int) Vec2D {
	result := *v
	result.X *= s
	result.Y *= s
	return result
}

func (v *Vec2D) AddInPlace(o Vec2D) {
	v.X += o.X
	v.Y += o.Y
}

func (v *Vec2D) SubInPlace(o Vec2D) {
	v.X -= o.X
	v.Y -= o.Y
}

func (v *Vec2D) MinInPlace(o Vec2D) {
	if o.X < v.X { v.X = o.X }
	if o.Y < v.Y { v.Y = o.Y }
}

func (v *Vec2D) MaxInPlace(o Vec2D) {
	if o.X > v.X { v.X = o.X }
	if o.Y > v.Y { v.Y = o.Y }
}

type Star struct {
	Position, Velocity Vec2D
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
	Lookup map[Vec2D]struct{}
	Min, Max Vec2D
}

func NewStarField(lines []string) (result StarField, err error) {
	result = StarField{}
	result.Stars = make([]Star, len(lines))
	for i, line := range lines {
		if err = ParseStar(&result.Stars[i], line); err != nil {
			return result, err
		}
	}
	result.GenerateLookup()
	return result, nil
}

func (sf *StarField) GenerateLookup() {
	sf.Min = MaxVec2D()
	sf.Max = MinVec2D()
	sf.Lookup = make(map[Vec2D]struct{})
	for i := range sf.Stars {
		s := &sf.Stars[i]
		sf.Lookup[s.Position] = struct{}{}
		sf.Max.MaxInPlace(s.Position)
		sf.Min.MinInPlace(s.Position)
	}
}

// TODO: instead of manipulating the state of the whole starfield, just update
//		the lookups and area to look into the future
func (sf *StarField) AdvanceTime(n int) {
	// Update star positions by 1 second
	for i := range sf.Stars {
		s := &sf.Stars[i]
		s.Position.AddInPlace(s.Velocity.Scale(n))
	}
	// Re-generate lookup table
	sf.GenerateLookup()
}

func (sf *StarField) Show(star string, space string) string {
	b := strings.Builder{}
	for y := sf.Min.Y; y <= sf.Max.Y; y++ {
		for x := sf.Min.X; x <= sf.Max.X; x++ {
			if _, ok := sf.Lookup[Vec2D{x, y}]; ok {
				b.WriteString(star)
			} else {
				b.WriteString(space)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (sf *StarField) HasAdjacent4(p Vec2D) bool {
	search := []Vec2D{
		p.Add(Vec2D{1, 0}),
		p.Add(Vec2D{0, 1}),
		p.Add(Vec2D{-1, 0}),
		p.Add(Vec2D{0, -1}),
	}
	for _, s := range search {
		if _, ok := sf.Lookup[s]; ok {
			return true
		}
	}
	return false
}

func (sf *StarField) CountOutliers(hasAdjacent func(p Vec2D) bool) int {
	result := 0
	for i := range sf.Stars {
		s := &sf.Stars[i]
		if !hasAdjacent(s.Position) {
			result++
		}
	}
	return result
}

func (sf *StarField) Area() int {
	return (sf.Max.X - sf.Min.X) * (sf.Max.Y - sf.Min.Y)
}

func findOptimum(sf *StarField, heuristic func(sf *StarField)int, compare func(a, b int)int) int {
	value := heuristic(sf)
	newValue := value

	time := 0
	for ; ; time++ {
		sf.AdvanceTime(1)
		newValue = heuristic(sf)
		check := compare(value, newValue)
		if check < 0 {
			sf.AdvanceTime(-1)
			return time - 1
		}
		value = newValue
	}
}

// TODO: try starting with an increment value, every time heuristic worsens
// 		reduce the increment and try again, until increment = 1. Effectively
//		a binary search with a first stage to find the "end". Will need to be
//		a bit different to a classical binary search, since it's not one list
//		of values sorted, it's 2 lists sorted in opposite directions and we're
//		trying to find the "join".
func part1impl(logger *log.Logger, filename string) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	lines, err := util.ReadLinesFromFile(filename)
	util.Check(err)

	starField, err := NewStarField(lines)
	util.Check(err)
	//logger.Print("initial:\n", starField.Show("#", " "))
	//logger.Print("outliers: ", starField.CountOutliers(starField.HasAdjacent4))
	//for i := 0; i < 3; i++ {
	//count := starField.CountOutliers(starField.HasAdjacent4)
	//time := 0
	//logger.Print("outlier count after ", time, " seconds: ", count)
	//for ; count > 0; time++ {
	//	starField.AdvanceTime()
	//	count = starField.CountOutliers(starField.HasAdjacent4)
	//	logger.Print("outlier count after ", time, " seconds: ", count)
	//	//logger.Print("After ", i+1, " seconds:\n", starField.Show("#", " "))
	//	//logger.Print("outliers: ", starField.CountOutliers(starField.HasAdjacent4))
	//}

	time := findOptimum(
		&starField,
		//func(sf *StarField) int { return sf.CountOutliers(sf.HasAdjacent4) },
		(*StarField).Area,
		func(a, b int) int { return a - b },
	)

	return fmt.Sprint("after ", time, " seconds:\n", starField.Show("#", " "))
}

func part0(logger *log.Logger) string {
	return part1impl(logger, "day10/input_test.txt")
}

func part1(logger *log.Logger) string {
	return part1impl(logger, "day10/input.txt")
}

func init() {
	util.RegisterSolution("day10part0", part0)
	util.RegisterSolution("day10part1", part1)
}
