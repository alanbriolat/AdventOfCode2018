/*
Day 6
=====

Part 1
------

The area around a Location is defined by points that are closer to it than any
other location.

The extent of the Map can at least be defined relative to the bounding box that
contains all Locations.

If a Location is on the edge of the Map, a "cone" of points that are closest to
that Location extends out infinitely, giving that Location an infinite area.

Therefore we should be able to define an extent of the Map such that for each
point that lies on the perimeter, it is closest to a location that has infinite
area.

Defining how much larger the extent of the Map should be, than the bounding box
of all Locations, is the difficult/mathematical problems. It must be large
enough to encompass the area of all Locations with non-infinite areas.

Manhattan distance messes up intuition this. With Euclidian geometry, it's
trivial to construct a set of points where a cell is bounded by the bounding
box, but should include a finite area outside of it, with converging edges.
However, with Manhattan distance, edges never converge. (As you get further
from the boundary, the distance normal to the boundary overwhelms and reduces
the influence of the distance parallel to the boundary.)
*/
package day06

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io"
	"log"
	"math"
	"os"
)

type Point struct {
	X, Y int
}

func (p Point) ManhattanDistance(o Point) int {
	return util.AbsInt(p.X-o.X) + util.AbsInt(p.Y-o.Y)
}

func MinPoint(p1, p2 Point) Point {
	return Point{util.MinInt(p1.X, p2.X), util.MinInt(p1.Y, p2.Y)}
}

func MaxPoint(p1, p2 Point) Point {
	return Point{util.MaxInt(p1.X, p2.X), util.MaxInt(p1.Y, p2.Y)}
}

type Location struct {
	Coordinates Point
	Infinite    bool
	Area        int
}

type Map struct {
	Min, Max  Point
	Locations []Location
}

func NewMap() Map {
	return Map{
		Point{math.MaxInt32, math.MaxInt32},
		Point{math.MinInt32, math.MinInt32},
		make([]Location, 0),
	}
}

func (m *Map) CreateLocation(p Point) {
	m.Locations = append(m.Locations, Location{p, false, 0})
	m.Min = MinPoint(m.Min, p)
	m.Max = MaxPoint(m.Max, p)
}

func (m *Map) ClosestLocation(p Point) (result *Location, ok bool) {
	result = nil
	resultDistance := 0
	for i := range m.Locations {
		l := &m.Locations[i]
		distance := p.ManhattanDistance(l.Coordinates)
		switch {
		case distance == resultDistance:
			result = nil
		case result == nil, distance < resultDistance:
			result = l
			resultDistance = distance
		}
	}
	return result, result != nil
}

/*
Traverses a perimeter that is 1 unit of distance outside of the bounding box
containing all locations, so that minimum distance is 1. Marks locations closest
to a perimeter point as infinite.
*/
func (m *Map) MarkInfinite() {
	x, y := m.Min.X-1, m.Min.Y-1
	for ; x <= m.Max.X+1; x++ {
		if l, ok := m.ClosestLocation(Point{x, y}); ok {
			l.Infinite = true
		}
	}
	for y++; y <= m.Max.Y+1; y++ {
		if l, ok := m.ClosestLocation(Point{x, y}); ok {
			l.Infinite = true
		}
	}
	for x--; x >= -1; x-- {
		if l, ok := m.ClosestLocation(Point{x, y}); ok {
			l.Infinite = true
		}
	}
	for y--; y > -1; y-- {
		if l, ok := m.ClosestLocation(Point{x, y}); ok {
			l.Infinite = true
		}
	}
}

func (m *Map) CalculateAreas() {
	for x := m.Min.X - 1; x <= m.Max.X+1; x++ {
		for y := m.Min.Y - 1; y <= m.Max.Y+1; y++ {
			if l, ok := m.ClosestLocation(Point{x, y}); ok {
				l.Area++
			}
		}
	}
}

func (m *Map) FindMostRemoteLocation() (result *Location) {
	result = nil
	for i := range m.Locations {
		l := &m.Locations[i]
		if !l.Infinite && (result == nil || l.Area > result.Area) {
			result = l
		}
	}
	return result
}

func (m *Map) LocationDistanceSum(p Point) int {
	result := 0
	for _, l := range m.Locations {
		result += p.ManhattanDistance(l.Coordinates)
	}
	return result
}

func (m *Map) CountPointsWithinRange(r int) int {
	result := 0
	for x := m.Min.X - 1; x <= m.Max.X+1; x++ {
		for y := m.Min.Y - 1; y <= m.Max.Y+1; y++ {
			distanceSum := m.LocationDistanceSum(Point{x, y})
			if distanceSum < r {
				result++
			}
		}
	}
	return result
}

func ReadPoints(name string) []Point {
	reader, err := os.Open(name)
	util.Check(err)
	result := make([]Point, 0)
	for {
		p := Point{}
		n, err := fmt.Fscanf(reader, "%d, %d", &p.X, &p.Y)
		if n == 0 || err == io.EOF {
			// end of input
			break
		}
		util.Check(err)
		result = append(result, p)
	}
	return result
}

func part1(logger *log.Logger) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	points := ReadPoints("day06/input.txt")
	t.LogCheckpoint(fmt.Sprint("read ", len(points), " points"))

	worldMap := NewMap()
	for _, p := range points {
		worldMap.CreateLocation(p)
	}
	t.LogCheckpoint(fmt.Sprintf("populated map %v to %v", worldMap.Min, worldMap.Max))

	worldMap.MarkInfinite()
	t.LogCheckpoint(fmt.Sprintf("marked locations as infinite"))

	worldMap.CalculateAreas()
	t.LogCheckpoint(fmt.Sprintf("calculated areas"))

	bestLocation := worldMap.FindMostRemoteLocation()
	t.LogCheckpoint(fmt.Sprintf("found destination: %+v", bestLocation))

	return fmt.Sprint(bestLocation.Area)
}

func part2(logger *log.Logger) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	points := ReadPoints("day06/input.txt")
	t.LogCheckpoint(fmt.Sprint("read ", len(points), " points"))

	worldMap := NewMap()
	for _, p := range points {
		worldMap.CreateLocation(p)
	}
	t.LogCheckpoint(fmt.Sprintf("populated map %v to %v", worldMap.Min, worldMap.Max))

	area := worldMap.CountPointsWithinRange(10000)
	t.LogCheckpoint(fmt.Sprintf("found %v points with distance sum < 10000", area))

	return fmt.Sprint(area)
}

func init() {
	util.RegisterSolution("day06part1", part1)
	util.RegisterSolution("day06part2", part2)
}
