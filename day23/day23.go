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

/*
Applies an "evolutionary strategy" to discover the optimum location, along the following lines:

- The fitness comparison optimises for number of nanobots in range, then for distance from (0, 0, 0)
- The initial population is the location of every nanobot
- Each generation:
	- Preserve `keep` best locations from previous generation
	- Generate `multiply` new (unique) locations for each preserved location
	- New locations are perturbed by a random Manhattan distance, which is randomly partitioned into (X, Y, Z)
	- If the best location from this generation is fitter than from last generation, record it
- The perturbation is limited by `energy`, decreasing exponentially from the size of the entire search volume
- When the best location hasn't been surpassed in `threshold` generations, terminate the algorithm

(This isn't a genetic algorithm, because it has mutation and selection but no crossover.)
 */
func part2impl(logger *log.Logger, filename string) int {
	nanobots := readNanobots(filename)
	min, max := util.MaxVec3D(), util.MinVec3D()
	population := make([]Location, 0, len(nanobots))

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

	// Populate candidate locations from nanobot locations
	for _, b := range nanobots {
		population = append(population, Location{
			Position: b.Position,
			InRangeOf: evaluateFunc(b.Position),
		})
	}
	// Sort by number of bots in range
	sort.Slice(population, locationSort(population))
	//logger.Printf("most connected positions: %v...", population[:util.MinInt(10, len(population)-1)])

	energy := max.Sub(min).Manhattan()
	threshold := 10
	keep := len(nanobots)
	multiply := 5
	generate := keep * multiply

	best := Location{}
	bestSurvival := 0	// How long the best location has remained the best location

	generation := 0
	for ; bestSurvival < threshold; generation++ {
		//logger.Printf("new generation with energy=%d threshold=%d keep=%d generate=%d", energy, threshold, keep, generate)
		newPopulation := make([]Location, 0, keep + generate)
		newPopulationSet := make(map[Location]bool)
		for _, loc := range population[0:keep] {
			// Keep the existing location
			newPopulation = append(newPopulation, loc)
			newPopulationSet[loc] = true
			for i, generated := 0, 0; generated < multiply && i < multiply*multiply; i++ {
				// Pick a random manhattan distance to perturb by (at least 1)
				distance := rand.Intn(energy)+1
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
				newLoc := Location{
					Position: newPos,
					InRangeOf: evaluateFunc(newPos),
				}
				if !newPopulationSet[newLoc] {
					newPopulation = append(newPopulation, newLoc)
					newPopulationSet[newLoc] = true
					generated++
				}
			}
		}

		// Sort the new population
		sort.Slice(newPopulation, locationSort(newPopulation))
		// See what the best location is so far
		newBest := newPopulation[0]

		// Did we get better?
		if lessFunc(&newBest, &best) {
			best = newBest
			bestSurvival = 0
		} else {
			bestSurvival++
		}
		//logger.Printf("best location so far: %+v, distance %d, survived %d generations",
		//	best, best.Position.Manhattan(), bestSurvival)
		population = newPopulation

		// Exponentially reduce the randomness
		if energy > 1 {
			energy /= 2
		}
	}

	logger.Printf("best location after %d generations: %+v, distance=%d", generation, best, best.Position.Manhattan())

	bestCount := 1
	for i := 1; i < len(population); i++ {
		b := &population[i]
		if lessFunc(&best, b) {
			break
		}
		bestCount++
	}
	logger.Printf("found %d best locations, %+v", bestCount, population[bestCount-1])

	return best.Position.Manhattan()
}

func init() {
	//util.RegisterSolution("day23test1", func(logger *log.Logger) string {
	//	return fmt.Sprint(part1impl(logger, "day23/input_test.txt"))
	//})
	util.RegisterSolution("day23part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day23/input.txt"))
	})
	//util.RegisterSolution("day23test2", func(logger *log.Logger) string {
	//	return fmt.Sprint(part2impl(logger, "day23/input_test2.txt"))
	//})
	util.RegisterSolution("day23part2", func(logger *log.Logger) string {
		return fmt.Sprint(part2impl(logger, "day23/input.txt"))
	})
}
