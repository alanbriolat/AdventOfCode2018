package day23

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
	"strconv"
	"strings"
)

type Nanobot struct {
	Position util.Vec3D
	Range int
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
		if bot.Position.Sub(largestRange.Position).Manhattan() <= largestRange.Range {
			count++
		}
	}

	return count
}

func init() {
	util.RegisterSolution("day23test1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day23/input_test.txt"))
	})
	util.RegisterSolution("day23part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day23/input.txt"))
	})
}
