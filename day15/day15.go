package day15

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	InputWall   = '#'
	InputFloor  = '.'
	InputElf    = 'E'
	InputGoblin = 'G'
	MapWall     = 254
	MapFloor    = 255
)

func tieBreak(p1, p2 util.Vec2D) bool {
	return p1.Y < p2.Y || (p1.Y == p2.Y && p1.X < p2.X)
}

type BitMap [][]bool

func (b BitMap) String() string {
	sb := strings.Builder{}
	sb.Grow(len(b) * (len((b)[0]) + 1))
	for y := 0; y < len(b); y++ {
		for x := 0; x < len(b[0]); x++ {
			switch b[x][y] {
			case true:
				sb.WriteByte('*')
			case false:
				sb.WriteByte('_')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (b BitMap) Get(p util.Vec2D) bool {
	return b[p.X][p.Y]
}

func (b BitMap) Set(p util.Vec2D) {
	b[p.X][p.Y] = true
}

func (b BitMap) Unset(p util.Vec2D) {
	b[p.X][p.Y] = false
}

type Unit struct {
	Battle        *Battle
	Id            byte
	Position      util.Vec2D
	IsGoblin      bool
	HitPoints     int
	AttackPower   int
	reachable     BitMap
	reachableFrom util.Vec2D
}

func (u *Unit) String() string {
	unitType := InputElf
	if u.IsGoblin {
		unitType = InputGoblin
	}
	return fmt.Sprintf("%s%d(%d)@%d,%d", string(unitType), u.Id, u.HitPoints, u.Position.X, u.Position.Y)
}

func (u *Unit) IsAlive() bool {
	return u.HitPoints > 0
}

func (u *Unit) IsEnemy(o *Unit) bool {
	// Don't need to check for same ID, because the IsGoblin value would be the same
	return u.IsGoblin != o.IsGoblin
}

/*
Reachable creates a map of reachable squares from a Unit's position using
a flood fill algorithm.
*/
func (u *Unit) Reachable() BitMap {
	// Next squares to process - upperbound length is all non-wall squares
	queue := make([]util.Vec2D, 0, u.Battle.NonWallCount)
	// Squares already in the processing queue
	queued := BitMap(util.NewBoolGrid(u.Battle.MapSize.X, u.Battle.MapSize.Y))
	// Squares that are reachable from u.Position
	reachable := BitMap(util.NewBoolGrid(u.Battle.MapSize.X, u.Battle.MapSize.Y))

	// Start off with u.Position in the queue - can reach own position by not moving
	queue = append(queue, u.Position)
	queued[u.Position.X][u.Position.Y] = true

	for len(queue) > 0 {
		// Pop from the front of the queue
		var next util.Vec2D
		next, queue = queue[0], queue[1:]
		// Mark as reachable
		reachable.Set(next)
		// Find candidates for other reachable squares
		adjacent := u.Battle.Adjacent(next, true)
		for _, p := range adjacent {
			// Skip squares already processed or queued
			if !queued.Get(p) {
				queue = append(queue, p)
				queued.Set(p)
			}
		}
	}
	return reachable
}

/*
FindPath finds the shortest path to reach a destination, implemented as A*
search.

The following details are specific to the requirements of this problem:

- The heuristic function is distance to any destination, not to a specific
  destination.
- Specifically, distance is measured as Manhattan distance.
- If there are multiple candidates with the same score, they are sorted by
  "reading order", i.e. top-to-bottom, left-to-right.

Mostly just following the pseudocode at https://en.wikipedia.org/wiki/A*_search_algorithm
*/
func (u *Unit) FindPath(destinations []util.Vec2D) ([]util.Vec2D, error) {
	if len(destinations) == 0 {
		return nil, fmt.Errorf("no destinations")
	}

	// f(n) = g(n) + h(n)
	// g(n): cost to get to n from start
	// h(n): estimate of cost from n to goal
	gScore := make(map[util.Vec2D]int)
	g := func(n util.Vec2D) int {
		if result, ok := gScore[n]; ok {
			return result
		} else {
			return math.MaxInt32
		}
	}
	fScore := make(map[util.Vec2D]int)
	f := func(n util.Vec2D) int {
		if result, ok := fScore[n]; ok {
			return result
		} else {
			return math.MaxInt32
		}
	}
	h := func(n util.Vec2D) int {
		min := math.MaxInt32
		for _, d := range destinations {
			if distance := d.Sub(n).Manhattan(); distance < min {
				min = distance
			}
		}
		return min
	}

	// Position processing queue
	openSet := make([]util.Vec2D, 0, u.Battle.NonWallCount)
	// Positions already processed
	closedSet := BitMap(util.NewBoolGrid(u.Battle.MapSize.X, u.Battle.MapSize.Y))
	// Positions queued or processed
	visited := BitMap(util.NewBoolGrid(u.Battle.MapSize.X, u.Battle.MapSize.Y))

	start := u.Position
	gScore[start] = 0
	fScore[start] = h(start)
	openSet = append(openSet, start)
	visited.Set(start)

	// Keep track of most efficient path to each position
	cameFrom := make(map[util.Vec2D]util.Vec2D)
	path := func(destination util.Vec2D) []util.Vec2D {
		result := make([]util.Vec2D, 0, g(destination)+1)
		next := destination
		for next != start {
			result = append(result, next)
			next = cameFrom[next]
		}
		// Reverse the path, so it's from start to goal
		for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
			result[left], result[right] = result[right], result[left]
		}
		return result
	}

	for len(openSet) > 0 {
		// Sort by f(n), tie-break on reading order
		sort.Slice(openSet, func(i, j int) bool {
			p1, p2 := openSet[i], openSet[j]
			f1, f2 := f(p1), f(p2)
			return f1 < f2 || (f1 == f2 && tieBreak(p1, p2))
		})
		// Get most promising next position
		var current util.Vec2D
		current, openSet = openSet[0], openSet[1:]
		// Did we find a goal?
		for _, d := range destinations {
			if current == d {
				return path(current), nil
			}
		}
		// Don't visit this position again
		closedSet.Set(current)

		// Score potential next positions
		nScore := gScore[current] + 1  // Distance from start to neighbour (path cost is always 1)
		for _, neighbour := range u.Battle.Adjacent(current, true) {
			if closedSet.Get(neighbour) {
				// Already have a shortest path to this neighbour
				continue
			}
			if !visited.Get(neighbour) {
				// Position we've never seen before, add to the queue
				openSet = append(openSet, neighbour)
				visited.Set(neighbour)
			} else if nScore >= g(neighbour) {
				continue
			}
			// This position is already in the queue, and we've found a quicker path to it
			cameFrom[neighbour] = current
			gScore[neighbour] = nScore
			fScore[neighbour] = nScore + h(neighbour)
		}
	}

	return nil, fmt.Errorf("no path found to any destination")
}

type Battle struct {
	Units        []*Unit
	Map          [][]byte
	MapSize      util.Vec2D
	WallCount    int
	NonWallCount int
}

func NewBattle(input []string) Battle {
	b := Battle{}
	b.Units = make([]*Unit, 0)
	b.MapSize = util.Vec2D{len(input[0]), len(input)}
	b.WallCount = 0
	b.Map = util.NewByteGrid(b.MapSize.X, b.MapSize.Y)

	for y, line := range input {
		for x := range line {
			switch line[x] {
			case InputWall:
				b.Map[x][y] = MapWall
				b.WallCount += 1
			case InputFloor:
				b.Map[x][y] = MapFloor
			case InputElf:
				b.CreateUnit(false, x, y)
			case InputGoblin:
				b.CreateUnit(true, x, y)
			}
		}
	}

	b.NonWallCount = b.MapSize.Area() - b.WallCount

	return b
}

func (b *Battle) String() string {
	sb := strings.Builder{}
	sb.Grow(b.MapSize.X*b.MapSize.Y + len(b.Units)*12)

	for y := 0; y < b.MapSize.Y; y++ {
		units := make([]byte, 0)
		for x := 0; x < b.MapSize.X; x++ {
			switch i := b.Map[x][y]; i {
			case MapWall:
				sb.WriteByte(InputWall)
			case MapFloor:
				sb.WriteByte(InputFloor)
			default:
				units = append(units, i)
				switch b.Units[i].IsGoblin {
				case false:
					sb.WriteByte(InputElf)
				case true:
					sb.WriteByte(InputGoblin)
				}
			}
		}
		sb.WriteString("   ")
		for _, i := range units {
			sb.WriteByte(' ')
			sb.WriteString(b.Units[i].String())
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func (b *Battle) SortUnits() {
	// Sort by the "tie break" criteria
	sort.Slice(b.Units, func(i, j int) bool {
		return tieBreak(b.Units[i].Position, b.Units[j].Position)
	})
	// Update all the unit IDs
	for i, u := range b.Units {
		u.Id = byte(i)
		if u.IsAlive() {
			*b.At(u.Position) = u.Id
		}
	}
}

func (b *Battle) At(p util.Vec2D) *byte {
	return &b.Map[p.X][p.Y]
}

func (b *Battle) ValidPosition(p util.Vec2D) bool {
	switch {
	case p.X < 0, p.X >= b.MapSize.X, p.Y < 0, p.Y >= b.MapSize.Y:
		return false
	default:
		return true
	}
}

func (b *Battle) Adjacent(p util.Vec2D, floorOnly bool) []util.Vec2D {
	// Possible adjacent squares, in "reading order"
	candidates := []util.Vec2D{
		p.Add(util.Vec2D{0, -1}),
		p.Add(util.Vec2D{-1, 0}),
		p.Add(util.Vec2D{1, 0}),
		p.Add(util.Vec2D{0, 1}),
	}
	result := make([]util.Vec2D, 0, 4)
	for _, c := range candidates {
		if b.ValidPosition(c) && (!floorOnly || *b.At(c) == MapFloor) {
			result = append(result, c)
		}
	}
	return result
}

func (b *Battle) CreateUnit(isGoblin bool, x, y int) {
	u := Unit{
		b,
		byte(len(b.Units)),
		util.Vec2D{x, y},
		isGoblin,
		200,
		3,
		nil,
		util.Vec2D{-1, -1},
	}
	b.Units = append(b.Units, &u)
	b.Map[x][y] = u.Id
}

func (b *Battle) MoveUnit(u *Unit, p util.Vec2D) {
	*b.At(p) = u.Id
	*b.At(u.Position) = MapFloor
	u.Position = p
}

func (b *Battle) AttackUnit(u *Unit, damage int) {
	u.HitPoints -= damage
	if u.HitPoints <= 0 {
		*b.At(u.Position) = MapFloor
	}
}

func (b *Battle) FindTargets(u *Unit) []*Unit {
	result := make([]*Unit, 0)
	for _, t := range b.Units {
		if u.IsEnemy(t) && t.IsAlive() {
			result = append(result, t)
		}
	}
	return result
}

func (b *Battle) FindDestinations(u *Unit, targets []*Unit) []util.Vec2D {
	reachable := u.Reachable()
	resultSet := make(map[util.Vec2D]struct{})

	// Find all the reachable locations in range of a target, de-duplicated
	for _, t := range targets {
		// Don't need floorOnly because non-floor squares are not in reachable set
		for _, p := range b.Adjacent(t.Position, false) {
			if reachable.Get(p) {
				resultSet[p] = struct{}{}
			}
		}
	}

	// Convert to slice
	result := make([]util.Vec2D, 0, len(resultSet))
	for p := range resultSet {
		result = append(result, p)
	}

	// Sort by distance
	sort.Slice(result, func(i, j int) bool {
		return (&result[i]).Sub(u.Position).Length() < (&result[j]).Sub(u.Position).Length()
	})
	return result
}

func (b *Battle) RemainingHitPoints() int {
	result := 0
	for _, u := range b.Units {
		if u.IsAlive() {
			result += u.HitPoints
		}
	}
	return result
}

func (b *Battle) NextRound() (combatEnded bool) {
	b.SortUnits()
	fmt.Print("start of round:\n", b.String())

	for _, u := range b.Units {
		if !u.IsAlive() {
			// Dead units don't move!
			continue
		}
		fmt.Printf("new turn: %s\n", u.String())
		// Find targets
		targets := b.FindTargets(u)
		if len(targets) == 0 {
			// Combat ended, one side has no remaining units
			fmt.Println("  no targets, combat ended")
			return true
		}
		// Find path to nearest position in range of a target, and move towards it
		destinations := b.FindDestinations(u, targets)
		fmt.Println("  destinations:", destinations)
		path, err := u.FindPath(destinations)
		if err != nil {
			// Can't find any targets, so end turn
			fmt.Println("  no path found")
			continue
		}
		fmt.Println("  path found:", path)
		if len(path) > 0 {
			// Not already in position to attack, so move 1 step
			fmt.Println("  moving from", u.Position, "to", path[0])
			b.MoveUnit(u, path[0])
		}
		// Find best adjacent enemy
		var target *Unit
		for _, a := range b.Adjacent(u.Position, false) {
			at := *b.At(a)
			// Not a unit, skip this square
			if int(at) >= len(b.Units) {
				continue
			}
			newTarget := b.Units[at]
			switch {
			case !u.IsEnemy(newTarget):
				// Not an enemy, so not a target
				continue
			case target == nil:
				// First enemy found
				target = newTarget
			case newTarget.HitPoints < target.HitPoints:
				// Weaker enemy than current target. Equal HP tie-break is
				// already handled by Adjacent being in tie-break order (i.e.
				// first target with a specific HP remains the target).
				target = newTarget
			}
		}
		// If we have an enemy, attack it
		if target != nil {
			fmt.Println("  attacking target:", target.String())
			b.AttackUnit(target, u.AttackPower)
			if !target.IsAlive() {
				fmt.Println("  killed target:", target.String())
			}
		}
	}
	return false
}

func part1impl(logger *log.Logger, input []string, maxRounds int, interactive bool) (rounds, remainingHP int) {
	battle := NewBattle(input)
	combatEnded := false
	reader := bufio.NewReader(os.Stdin)
	var i int
	for i = 0; !combatEnded && i < maxRounds; i++ {
		if interactive {
			fmt.Print("Hit enter to continue...")
			reader.ReadString('\n')
		}
		combatEnded = battle.NextRound()
	}
	return i - 1, battle.RemainingHitPoints()
}

func part1(logger *log.Logger, filename string, maxRounds int, interactive bool) string {
	input, _ := util.ReadLinesFromFile(filename)
	rounds, remainingHP := part1impl(logger, input, maxRounds, interactive)
	return fmt.Sprintf("%dx%d = %d", rounds, remainingHP, rounds*remainingHP)
}

func init() {
	util.RegisterSolution("day15test1", func(logger *log.Logger) string {
		return part1(logger, "day15/input_test1.txt", 3, false)
	})
	util.RegisterSolution("day15test2", func(logger *log.Logger) string {
		return part1(logger, "day15/input_test2.txt", 50, false)
	})
	util.RegisterSolution("day15part1", func(logger *log.Logger) string {
		return part1(logger, "day15/input.txt", math.MaxInt32, true)
	})
}
