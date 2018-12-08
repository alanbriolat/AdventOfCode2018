package day02

import (
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strings"
)

type Counter struct {
	Counts map[rune]int
	HasDouble bool
	HasTriple bool
}

func NewCounter() Counter {
	return Counter{Counts: make(map[rune]int)}
}

func (c *Counter) Count(s string) {
	for _, x := range s {
		c.Counts[x]++
	}
	c.HasDouble, c.HasTriple = false, false
	for _, v := range c.Counts {
		switch v {
		case 2:
			c.HasDouble = true
		case 3:
			c.HasTriple = true
		}
	}
}

type Comparison struct {
	s1, s2 string
	distance int
}

func NewComparison(s1, s2 string) Comparison {
	return Comparison{s1, s2, SubstitutionDistance(s1, s2)}
}

/*
SubstitutionDistance finds the substitution-only edit distance between two
strings of the same length.
 */
func SubstitutionDistance(s1, s2 string) int {
	distance := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			distance++
		}
	}
	return distance
}

func SharedString(s1, s2 string) string {
	s1r, s2r := []rune(s1), []rune(s2)
	builder := strings.Builder{}
	maxBytes := len(s1)
	if len(s2) < maxBytes {
		maxBytes = len(s2)
	}
	builder.Grow(maxBytes)
	for i := 0; i < len(s1r) && i < len(s2r); i++ {
		if s1r[i] == s2r[i] {
			builder.WriteRune(s1r[i])
		}
	}
	return builder.String()
}

func part1and2(logger *log.Logger) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")
	lines, err := util.ReadLinesFromFile("day02/input1.txt")
	util.Check(err)
	t.LogCheckpoint("readInput")

	doubles, triples := 0, 0
	for _, s := range lines {
		counter := NewCounter()
		counter.Count(s)
		if counter.HasDouble {
			doubles++
		}
		if counter.HasTriple {
			triples++
		}
		//logger.Println(s, "HasDouble?", counter.HasDouble, "HasTriple?", counter.HasTriple, counter.Counts)
	}
	t.LogCheckpoint("checksum")
	logger.Println("Checksum:", doubles * triples)


	var closest *Comparison = nil
	for i, s1 := range lines {
		for _, s2 := range lines[i+1:] {
			comparison := NewComparison(s1, s2)
			if closest == nil || comparison.distance < closest.distance {
				closest = &comparison
			}
		}
	}
	t.LogCheckpoint("closestComparison")
	logger.Println("Closest IDs:", closest, "shared string:", SharedString(closest.s1, closest.s2))
}

func init() {
	util.RegisterSolution("day02", part1and2)
}
