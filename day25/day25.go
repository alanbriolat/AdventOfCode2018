package day25

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Constellation []util.Vec4D

func parseConstellation(filename string) Constellation {
	result := Constellation{}
	rawReader, err := os.Open(filename)
	util.Check(err)
	reader := bufio.NewReader(rawReader)
	for {
		point := util.Vec4D{}
		s, err := reader.ReadString(',')
		if len(s) == 0 || err == io.EOF {
			// end of input
			break
		} else {
			util.Check(err)
		}
		point.X, err = strconv.Atoi(strings.TrimRight(s, ","))
		util.Check(err)
		s, err = reader.ReadString(',')
		point.Y, err = strconv.Atoi(strings.TrimRight(s, ","))
		util.Check(err)
		s, err = reader.ReadString(',')
		point.Z, err = strconv.Atoi(strings.TrimRight(s, ","))
		util.Check(err)
		s, err = reader.ReadString('\n')
		point.T, err = strconv.Atoi(strings.TrimRight(s, "\r\n"))
		util.Check(err)
		result = append(result, point)
	}
	return result
}

func part1impl(logger *log.Logger, filename string) int {
	allPoints := parseConstellation(filename)
	//logger.Println(allPoints)

	membership := make(map[util.Vec4D]*Constellation)

	for _, p := range allPoints {
		var firstCandidate *Constellation
		candidateSet := make(map[*Constellation]bool)
		for o, c := range membership {
			if o.Sub(p).Manhattan() <= 3 {
				if firstCandidate == nil {
					firstCandidate = c
				}
				candidateSet[c] = true
			}
		}
		if firstCandidate == nil {
			// No candidates (yet), start a new constellation
			c := &Constellation{p}
			membership[p] = c
		} else {
			// Join the first neighbour's constellation
			*firstCandidate = append(*firstCandidate, p)
			membership[p] = firstCandidate
			// Merge any other candidates into this constellation (this point bridged them)
			for c := range candidateSet {
				if c == firstCandidate {
					continue
				}
				for _, op := range *c {
					*firstCandidate = append(*firstCandidate, op)
					membership[op] = firstCandidate
				}
			}
		}
	}

	constellationSet := make(map[*Constellation]bool)
	for _, c := range membership {
		constellationSet[c] = true
	}

	return len(constellationSet)
}

func init() {
	util.RegisterSolution("day25part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day25/input.txt"))
	})
}
