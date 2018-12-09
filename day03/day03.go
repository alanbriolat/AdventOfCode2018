package day03

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io"
	"log"
	"os"
	"strconv"
)

type Square struct {
	X, Y int
}

type ClaimError string

func (e ClaimError) Error() string {
	return string(e)
}

type Claim struct {
	Id int
	X, Y, W, H int
}

/*
Overlap finds the Claim that is the intersection of c and o, or gives an error
if no such intersection exists.
*/
func (c Claim) Overlap(o Claim) (Claim, error) {
	overlap := Claim{}
	overlap.X, overlap.W = LinearOverlap(c.X, c.W, o.X, o.W)
	if overlap.W == 0 {
		return Claim{}, ClaimError("no overlap in X direction")
	}
	overlap.Y, overlap.H = LinearOverlap(c.Y, c.H, o.Y, o.H)
	if overlap.H == 0 {
		return Claim{}, ClaimError("no overlap in Y direction")
	}
	return overlap, nil
}

func (c Claim) Area() int {
	return c.W * c.H
}

func (c Claim) Squares() []Square {
	result := make([]Square, 0, c.W * c.H)
	for x := c.X; x < c.X + c.W; x++ {
		for y := c.Y; y < c.Y + c.H; y++ {
			result = append(result, Square{x, y})
		}
	}
	return result
}

/*
LinearOverlap calculates the starting point p and distance d of overlap between
two ranges of integers defined by their own starting points and distances.
*/
func LinearOverlap(p1, d1, p2, d2 int) (p, d int) {
	// Simplify by making second range always start after first range
	if p2 < p1 {
		p1, d1, p2, d2 = p2, d2, p1, d1
	}
	// Take a shortcut if there is no overlap
	if d1 == 0 || d2 == 0 || p2 >= p1+d1 {
		return 0, 0
	}
	return p2, util.MinInt(p1+d1-p2, d2)
}

func ReadClaimsFromFile(name string) ([]Claim, error) {
	result := make([]Claim, 0, 1500)
	var rawReader io.Reader
	var err error
	if rawReader, err = os.Open(name); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(rawReader)
	for {
		claim := Claim{}
		// read start of line
		bytes, err := reader.ReadBytes('#')
		if len(bytes) == 0 || err == io.EOF {
			// end of input
			break
		} else {
			util.Check(err)
		}
		// read the claim definition
		bytes, err = reader.ReadBytes('@')
		claim.Id, err = strconv.Atoi(string(bytes[:len(bytes)-2]))
		util.Check(err)
		bytes, err = reader.ReadBytes(',')
		claim.X, err = strconv.Atoi(string(bytes[1:len(bytes)-1]))
		util.Check(err)
		bytes, err = reader.ReadBytes(':')
		claim.Y, err = strconv.Atoi(string(bytes[:len(bytes)-1]))
		util.Check(err)
		bytes, err = reader.ReadBytes('x')
		claim.W, err = strconv.Atoi(string(bytes[1:len(bytes)-1]))
		util.Check(err)
		bytes, err = reader.ReadBytes('\n')
		claim.H, err = strconv.Atoi(string(bytes[:len(bytes)-1]))
		util.Check(err)
		
		result = append(result, claim)
	}
	return result, nil
}

func part1and2(logger *log.Logger) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	claims, err := ReadClaimsFromFile("day03/input1.txt")
	util.Check(err)
	t.LogCheckpoint(fmt.Sprint("read ", len(claims), " claims"))

	// Create sparse grid of claim counts per square
	counts := make(map[Square]int)
	for _, claim := range claims {
		for _, square := range claim.Squares() {
			counts[square]++
		}
	}
	t.LogCheckpoint(fmt.Sprint("counted ", len(counts), " claimed squares"))

	// Count squares with more than one claim
	contested := 0
	for _, count := range counts {
		if count > 1 {
			contested += 1
		}
	}
	t.LogCheckpoint(fmt.Sprint("found ", contested, " contested squares"))

	// Find claim where every square was only claimed once
	var intact *Claim = nil
	for _, claim := range claims {
		overlapped := false
		for _, square := range claim.Squares() {
			if counts[square] > 1 {
				overlapped = true
				break
			}
		}
		if !overlapped {
			intact = &claim
			break
		}
	}
	t.LogCheckpoint(fmt.Sprintf("found intact claim #%v at %v,%v %vx%v",
		intact.Id, intact.X, intact.Y, intact.W, intact.H))

	return fmt.Sprintf("part1 = %v , part2 = %v", contested, intact.Id)
}

func init() {
	util.RegisterSolution("day03", part1and2)
}
