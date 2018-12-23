package day23

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Nanobot struct {
	Position   util.Vec3D
	Range      int
}

func (b *Nanobot) IsInRange(p util.Vec3D) bool {
	return b.Position.Sub(p).Manhattan() <= b.Range
}

func readNanobots(filename string) []Nanobot {
	file, err := os.Open(filename)
	util.Check(err)
	reader := bufio.NewReader(file)
	result := make([]Nanobot, 0)
	for {
		nanobot := Nanobot{}
		var s string
		var err error
		if _, err := reader.ReadString('<'); err != nil {
			break
		}
		s, err = reader.ReadString(',')
		nanobot.Position.X, err = strconv.Atoi(strings.TrimRight(s, ","))
		util.Check(err)
		s, err = reader.ReadString(',')
		nanobot.Position.Y, err = strconv.Atoi(strings.TrimRight(s, ","))
		util.Check(err)
		s, err = reader.ReadString('>')
		nanobot.Position.Z, err = strconv.Atoi(strings.TrimRight(s, ">"))
		util.Check(err)
		s, err = reader.ReadString('=')
		util.Check(err)
		s, err = reader.ReadString('\n')
		nanobot.Range, err = strconv.Atoi(strings.TrimRight(s, "\r\n"))
		util.Check(err)
		result = append(result, nanobot)
	}
	return result
}

type Location struct {
	Position util.Vec3D
	InRangeOf int
}

func part1impl(logger *log.Logger, filename string) int {
	nanobots := readNanobots(filename)

	// Find nanobot with largest range
	var largestRange *Nanobot
	for i := range nanobots {
		bot := &nanobots[i]
		if largestRange == nil || bot.Range > largestRange.Range {
			largestRange = bot
		}
	}

	// Find nanobots in range
	count := 0
	for i := range nanobots {
		bot := &nanobots[i]
		if largestRange.IsInRange(bot.Position) {
			count++
		}
	}

	return count
}

func part2impl(logger *log.Logger, filename string) int {
	nanobots := readNanobots(filename)
	min, max := util.MaxVec3D(), util.MinVec3D()
	candidates := make([]Location, 0, len(nanobots))

	// Find extend of coordinate space
	for _, bot := range nanobots {
		min.MinInPlace(bot.Position)
		max.MaxInPlace(bot.Position)
	}
	size := max.Sub(min).Add(util.Vec3D{1, 1, 1})
	logger.Printf("search space: %v to %v, size = %v, volume = %d", min, max, size, size.Volume())

	evaluateFunc := func(p util.Vec3D) int {
		count := 0
		for _, b := range nanobots {
			if b.IsInRange(p) {
				count++
			}
		}
		return count
	}

	lessFunc := func(a, b *Location) bool {
		return a.InRangeOf > b.InRangeOf || (a.InRangeOf == b.InRangeOf && a.Position.Manhattan() < b.Position.Manhattan())
	}

	locationSort := func(locations []Location) func(i, j int) bool {
		return func(i, j int) bool {
			return lessFunc(&locations[i], &locations[j])
		}
	}

	//validFunc := func(p util.Vec3D) bool {
	//	return p.X >= min.X && p.X <= max.X && p.Y >= min.Y && p.Y <= max.Y && p.Z >= min.Z && p.Z <= max.Z
	//}

	// Populate candidate locations from nanobot locations
	for _, b := range nanobots {
		candidates = append(candidates, Location{
			Position: b.Position,
			InRangeOf: evaluateFunc(b.Position),
		})
	}
	// Sort by number of bots in range
	sort.Slice(candidates, locationSort(candidates))
	logger.Printf("most connected positions: %v...", candidates[:util.MinInt(10, len(candidates)-1)])

	amount := max.Sub(min).Manhattan()
	threshold := 100
	keep := len(nanobots)
	multiply := 10
	generate := keep * multiply

	var best Location
	bestSurvival := 0	// How long the best location has remained the best location

	for bestSurvival < threshold {
		logger.Printf("new generation with amount=%d threshold=%d keep=%d generate=%d", amount, threshold, keep, generate)
		newCandidates := make([]Location, 0, keep + generate)
		for _, loc := range candidates[0:keep] {
			// Keep the existing location
			newCandidates = append(newCandidates, loc)
			for i := 0; i < multiply; i++ {
				// Pick a random manhattan distance to perturb by (at least 1)
				distance := rand.Intn(amount)+1
				// Partition the distance into random X, Y and Z amounts
				firstPartition := rand.Intn(distance+1)
				secondPartition := rand.Intn(distance+1)
				if firstPartition > secondPartition {
					firstPartition, secondPartition = secondPartition, firstPartition
				}
				perturb := util.Vec3D{
					firstPartition,
					secondPartition - firstPartition,
					distance - secondPartition,
				}
				// Randomise directions
				if rand.Float32() >= 0.5 { perturb.X *= -1 }
				if rand.Float32() >= 0.5 { perturb.Y *= -1 }
				if rand.Float32() >= 0.5 { perturb.Z *= -1 }
				// Create new position
				newPos := loc.Position.Add(perturb)
				// Clamp it inside the bounds of the known universe
				newPos.MaxInPlace(min)
				newPos.MinInPlace(max)
				// Create and add the new location
				//logger.Printf("adding new candidate: original=%v, perturb=%v, new=%v", loc.Position, perturb, newPos)
				newCandidates = append(newCandidates, Location{
					Position: newPos,
					InRangeOf: evaluateFunc(newPos),
				})
			}
		}

		// Sort the new candidates
		sort.Slice(newCandidates, locationSort(newCandidates))
		// See what the best candidate is so far
		newBest := newCandidates[0]
		if newBest == best {
			bestSurvival++
		} else {
			best = newBest
			bestSurvival = 0
		}
		logger.Printf("best location so far: %+v, distance %d, survived %d generations",
			best, best.Position.Manhattan(), bestSurvival)
		candidates = newCandidates

		if amount > 1 {
			amount /= 2
		}
	}

	logger.Printf("final best location: %+v, distance %d", best, best.Position.Manhattan())
	return best.Position.Manhattan()
}

func init() {
	util.RegisterSolution("day23test1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day23/input_test.txt"))
	})
	util.RegisterSolution("day23part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day23/input.txt"))
	})
	util.RegisterSolution("day23test2", func(logger *log.Logger) string {
		return fmt.Sprint(part2impl(logger, "day23/input_test2.txt"))
	})
	util.RegisterSolution("day23part2", func(logger *log.Logger) string {
		return fmt.Sprint(part2impl(logger, "day23/input.txt"))
	})
}
