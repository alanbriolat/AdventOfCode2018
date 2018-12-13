package day13

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"sort"
)

const (
	TrackH    = '-'
	TrackV    = '|'
	Intersect = '+'
	CornerL   = '\\'
	CornerR   = '/'
	CartU     = '^'
	CartD     = 'v'
	CartL     = '<'
	CartR     = '>'
	TurnLeft  = 0
	TurnNone  = 1
	TurnRight = 2
)

type Cart struct {
	Position        util.Vec2D
	Velocity        util.Vec2D
	IntersectAction int
	Crashed         bool
}

func NewCart(x, y int, direction byte) (cart Cart, track byte) {
	cart = Cart{
		Position:        util.Vec2D{x, y},
		IntersectAction: TurnLeft,
		Crashed:         false,
	}
	switch direction {
	case CartU:
		cart.Velocity, track = util.Vec2D{0, -1}, '|'
	case CartD:
		cart.Velocity, track = util.Vec2D{0, 1}, '|'
	case CartL:
		cart.Velocity, track = util.Vec2D{-1, 0}, '-'
	case CartR:
		cart.Velocity, track = util.Vec2D{1, 0}, '-'
	default:
		panic(fmt.Sprint("invalid direction: ", direction))
	}
	return
}

func (c *Cart) RotateCW() {
	result := util.Vec2D{
		c.Velocity.Y,
		c.Velocity.X,
	}
	if result.X != 0 {
		result.X = -result.X
	}
	c.Velocity = result
}

func (c *Cart) RotateCCW() {
	result := util.Vec2D{
		c.Velocity.Y,
		c.Velocity.X,
	}
	if result.Y != 0 {
		result.Y = -result.Y
	}
	c.Velocity = result
}

func (c *Cart) HandleIntersection() {
	switch c.IntersectAction {
	case TurnLeft:
		c.RotateCCW()
	case TurnRight:
		c.RotateCW()
	}
	c.IntersectAction = (c.IntersectAction + 1) % 3
}

func (c *Cart) HandleCorner(corner byte) {
	switch corner {
	case CornerR:
		if c.Velocity.Y != 0 {
			// Approaching from top or bottom
			c.RotateCW()
		} else {
			// Approaching from left or right
			c.RotateCCW()
		}
	case CornerL:
		if c.Velocity.Y != 0 {
			// Approaching from top or bottom
			c.RotateCCW()
		} else {
			// Approaching from left or right
			c.RotateCW()
		}
	}
}

type CartSystem struct {
	Width, Height int
	Carts         []Cart
	Crashes       []util.Vec2D
	Tracks        [][]byte
	Time          int
}

func NewCartSystem(input []string) CartSystem {
	cs := CartSystem{}
	cs.Height = len(input)
	cs.Width = len(input[0])
	cs.Carts = make([]Cart, 0)
	cs.Crashes = make([]util.Vec2D, 0)
	cs.Tracks = util.NewByteGrid(cs.Width, cs.Height)
	for x := 0; x < cs.Width; x++ {
		for y := 0; y < cs.Height; y++ {
			track := input[y][x]
			// If this is a cart, record it and replace it with the correct track
			switch track {
			case CartU, CartD, CartL, CartR:
				var cart Cart
				cart, track = NewCart(x, y, track)
				cs.Carts = append(cs.Carts, cart)
			}
			cs.Tracks[x][y] = track
		}
	}
	cs.Time = 0
	return cs
}

func (cs *CartSystem) Tick() {
	cs.Time++
	// Sort carts by position, to process them in the correct order
	sort.Slice(cs.Carts, func(i, j int) bool {
		a := &cs.Carts[i]
		b := &cs.Carts[j]
		return a.Position.Y < b.Position.Y || (a.Position.Y == b.Position.Y && a.Position.X < b.Position.X)
	})

	// Process each cart: move, turn, collide
	for i := range cs.Carts {
		cart := &cs.Carts[i]
		if cart.Crashed {
			// Skip crashed carts, because they can't move
			continue
		}

		// Move the cart
		cart.Position.AddInPlace(cart.Velocity)
		// Update cart state
		track := cs.Tracks[cart.Position.X][cart.Position.Y]
		switch track {
		case Intersect:
			cart.HandleIntersection()
		case CornerL, CornerR:
			cart.HandleCorner(track)
		case TrackH, TrackV:
			// No other state change required
		default:
			panic("cart off the track!")
		}

		// Check for collision
		for j := range cs.Carts {
			if j == i {
				continue
			}
			other := &cs.Carts[j]
			if other.Crashed {
				// Skip crashed carts, because they get removed instantly
				continue
			}
			if cart.Position == other.Position {
				cs.Crashes = append(cs.Crashes, cart.Position)
				cart.Crashed = true
				other.Crashed = true
				// can't collide more than once, because only one cart moved
				break
			}
		}
	}

	// Clean up crashed carts
	clean := cs.Carts[:0]
	for _, cart := range cs.Carts {
		if !cart.Crashed {
			clean = append(clean, cart)
		}
	}
	cs.Carts = clean
}

func part1(logger *log.Logger, filename string) util.Vec2D {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	lines, err := util.ReadLinesFromFile(filename)
	util.Check(err)
	cs := NewCartSystem(lines)
	t.Printf("read %vx%v cart system with %v carts", cs.Width, cs.Height, len(cs.Carts))

	for len(cs.Crashes) < 1 {
		cs.Tick()
	}
	logger.Printf("%d crash(es) at tick %d: %v\n", len(cs.Crashes), cs.Time, cs.Crashes)

	return cs.Crashes[0]
}

func part2(logger *log.Logger, filename string) util.Vec2D {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	lines, err := util.ReadLinesFromFile(filename)
	util.Check(err)
	cs := NewCartSystem(lines)
	t.Printf("read %vx%v cart system with %v carts", cs.Width, cs.Height, len(cs.Carts))

	for len(cs.Carts) > 1 {
		cs.Tick()
	}
	logger.Printf("%d cart(s) remaining at tick %d: %+v\n", len(cs.Carts), cs.Time, cs.Carts[0])

	return cs.Carts[0].Position
}

func init() {
	util.RegisterSolution("day13part1example", func(logger *log.Logger) string {
		p := part1(logger, "day13/input_test1.txt")
		return fmt.Sprint(p.X, ",", p.Y)
	})
	
	util.RegisterSolution("day13part1", func(logger *log.Logger) string {
		p := part1(logger, "day13/input.txt")
		return fmt.Sprint(p.X, ",", p.Y)
	})

	util.RegisterSolution("day13part2example", func(logger *log.Logger) string {
		p := part2(logger, "day13/input_test2.txt")
		return fmt.Sprint(p.X, ",", p.Y)
	})

	util.RegisterSolution("day13part2", func(logger *log.Logger) string {
		p := part2(logger, "day13/input.txt")
		return fmt.Sprint(p.X, ",", p.Y)
	})
}
