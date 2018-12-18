package day17

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Sand    = '.'
	Clay    = '#'
	Water   = '~'
	Flowing = '|'
)

type Line struct {
	Start, End util.Vec2D
}

type Aquifer struct {
	Min, Max      util.Vec2D
	Width, Height int
	Data          util.ByteGrid
	Springs       []util.Vec2D
}

func NewAquifer(input []Line) Aquifer {
	a := Aquifer{
		Min: util.Vec2D{500, 0},
		Max: util.Vec2D{500, 0},
	}
	// Establish the boundary of the map
	for _, line := range input {
		fmt.Println("line:", line)
		a.Min.MinInPlace(line.Start)
		a.Max.MaxInPlace(line.End)
	}
	// Expand by 1 more each way in X direction, to allow flowing around edge features
	a.Min.SubInPlace(util.Vec2D{1, 0})
	a.Max.AddInPlace(util.Vec2D{1, 0})
	// Width and height are inclusive of max
	a.Width = a.Max.X - a.Min.X + 1
	a.Height = a.Max.Y - a.Min.Y + 1
	fmt.Println("min", a.Min, "max", a.Max, "width", a.Width, "height", a.Height)
	// Create and initialise map data
	a.Data = util.NewByteGrid(a.Width, a.Height)
	for x := 0; x < a.Width; x++ {
		for y := 0; y < a.Height; y++ {
			a.Data[x][y] = Sand
		}
	}
	// Draw clay onto the map
	for _, line := range input {
		if line.Start.X == line.End.X {
			// Vertical line
			for p := line.Start; p.Y <= line.End.Y; p.Y++ {
				*a.At(p) = Clay
			}
		} else {
			// Horizontal line
			for p := line.Start; p.X <= line.End.X; p.X++ {
				*a.At(p) = Clay
			}
		}
	}
	a.Springs = make([]util.Vec2D, 1)
	a.Springs[0] = util.Vec2D{500, 0}
	return a
}

func (a *Aquifer) String() string {
	sb := strings.Builder{}
	sb.Grow(a.Width * a.Height)
	for y := 0; y < a.Height; y++ {
		for x := 0; x < a.Width; x++ {
			sb.WriteByte(a.Data[x][y])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (a *Aquifer) Valid(p util.Vec2D) bool {
	return p.X >= a.Min.X && p.X <= a.Max.X &&
		p.Y >= a.Min.Y && p.Y <= a.Max.Y
}

func (a *Aquifer) At(p util.Vec2D) *byte {
	p.SubInPlace(a.Min)
	return &a.Data[p.X][p.Y]
}

func (a *Aquifer) Flow() {
	position := a.Springs[0]
	a.Springs = a.Springs[1:]

	// Fall until finding something that isn't sand
	for a.Valid(position) && *a.At(position) == Sand {
		*a.At(position) = Flowing
		position.AddInPlace(util.Vec2D{0, 1})
	}

	if !a.Valid(position) {
		// Fell off the map without anything else to do
		return
	}
	if *a.At(position) == Flowing {
		// Collided with another stream, would now follow the same path
		return
	}

	scanToEnd := func(start, direction util.Vec2D) (p util.Vec2D, closed bool) {
		p = start
		closed = false
		p.AddInPlace(direction)
		fmt.Println("scanning", direction)
		scan: for a.Valid(p) {
			fmt.Println("checking", p)
			switch *a.At(p) {
			case Clay:
				// Hit a wall, so gone as far as we can in this direction
				p.SubInPlace(direction)
				closed = true
				break scan
			case Flowing:
				// Hit a falling stream of water
				p.SubInPlace(direction)
				break scan
			case Sand:
				// If we ended up "dangling" over more sand, let's stop but create a new spring
				below := p.Add(util.Vec2D{0, 1})
				if a.Valid(below) && *a.At(below) == Sand {
					a.Springs = append(a.Springs, p)
					p.SubInPlace(direction)
					break scan
				} else {
					// Nothing in the way, let's flow
					*a.At(p) = Flowing
					p.AddInPlace(direction)
				}
			default:
				panic("shouldn't be able to flow into water/unknown")
			}
		}
		return
	}

	// We hit something to spread across, go back up  and start generating water
	for position.SubInPlace(util.Vec2D{0, 1}); a.Valid(position) && *a.At(position) == Flowing; position.SubInPlace(util.Vec2D{0, 1}) {
		left, leftClosed := scanToEnd(position, util.Vec2D{-1, 0})
		right, rightClosed := scanToEnd(position, util.Vec2D{1, 0})
		fmt.Println(left, leftClosed, right, rightClosed)
		if leftClosed && rightClosed {
			// Line is capped by clay at both ends, so fill it with water
			for ; left.X <= right.X; left.X++ {
				fmt.Println("adding water to", left)
				*a.At(left) = Water
			}
		} else {
			// Line wasn't capped, this line won't fill with water, so stop processing those above it too
			break
		}
	}
}

func readInput(filename string) []Line {
	var err error
	file, err := os.Open(filename)
	util.Check(err)
	reader := bufio.NewReader(file)
	result := make([]Line, 0)
	for {
		directionStr, err := reader.ReadString('=')
		if err != nil {
			break
		}
		positionStr, err := reader.ReadString(',')
		util.Check(err)
		position, err := strconv.Atoi(strings.TrimRight(positionStr, ","))
		util.Check(err)
		reader.ReadString('=')
		startStr, err := reader.ReadString('.')
		util.Check(err)
		start, err := strconv.Atoi(strings.TrimRight(startStr, "."))
		util.Check(err)
		reader.ReadByte()
		endStr, err := reader.ReadString('\n')
		util.Check(err)
		end, err := strconv.Atoi(strings.TrimRight(endStr, "\r\n"))
		util.Check(err)
		var line Line
		if directionStr[0] == 'x' {
			line = Line{
				util.Vec2D{position, start},
				util.Vec2D{position, end},
			}
		} else {
			line = Line{
				util.Vec2D{start, position},
				util.Vec2D{end, position},
			}
		}
		result = append(result, line)
	}
	return result
}

func part1impl(logger *log.Logger, filename string) string {
	input := readInput(filename)
	aquifer := NewAquifer(input)
	logger.Print("start:\n", aquifer.String())
	for i := 0; len(aquifer.Springs) > 0; i++ {
		aquifer.Flow()
		logger.Printf("after %d flow steps:\n%s", i+1, aquifer.String())
		logger.Println("new springs:", aquifer.Springs)
	}
	return ""
}

func init() {
	util.RegisterSolution("day17test1", func(logger *log.Logger) string {
		return part1impl(logger, "day17/input_test.txt")
	})
	util.RegisterSolution("day17part1", func(logger *log.Logger) string {
		return part1impl(logger, "day17/input.txt")
	})
}
