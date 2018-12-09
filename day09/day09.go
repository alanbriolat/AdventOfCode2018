package day09

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"strings"
)

type Marble struct {
	Value int
	Left, Right *Marble
}

func (m *Marble) Forward(n int) *Marble {
	result := m
	for ; n > 0; n-- {
		result = result.Right
	}
	return result
}

func (m *Marble) Back(n int) *Marble {
	result := m
	for ; n > 0; n-- {
		result = result.Left
	}
	return result
}

func (m *Marble) AddRight(x *Marble) {
	right := m.Right
	m.Right = x
	x.Right = right
	right.Left = x
	x.Left = m
}

func (m *Marble) RemoveLeft() *Marble {
	left := m.Left
	m.Left = left.Left
	m.Left.Right = m
	left.Left = nil
	left.Right = nil
	return left
}

type GameState struct {
	Start *Marble
	Current *Marble
	Marbles []Marble
}

func NewGameState(size int) GameState {
	result := GameState{}
	result.Marbles = make([]Marble, size + 1)
	m := &result.Marbles[0]
	result.Start, result.Current, m.Left, m.Right = m, m, m, m
	return result
}

func (g *GameState) Play(x int) (score int) {
	score = 0
	if x % 23 == 0 {
		score += x
		g.Current = g.Current.Back(6)
		removed := g.Current.RemoveLeft()
		score += removed.Value
	} else {
		g.Current = g.Current.Forward(1)
		g.Current.AddRight(&g.Marbles[x])
		g.Current = g.Current.Forward(1)
	}

	return score
}

func (g *GameState) String() string {
	b := strings.Builder{}
	m := g.Start
	for {
		if m == g.Current {
			b.WriteString(fmt.Sprint("[", m.Value, "] "))
		} else {
			b.WriteString(fmt.Sprint(m.Value, " "))
		}

		m = m.Right
		if m == g.Start {
			break
		}
	}
	return b.String()
}

func part1impl(logger *log.Logger, players, max int) (highScore int) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	state := NewGameState(max)
	scores := make([]int, players)
	player := 0
	for i := 1; i <= max; i++ {
		scores[player] += state.Play(i)
		//logger.Println(state.String())
		player = (player + 1) % players
	}

	winner, score := 0, scores[0]
	for i, x := range scores {
		if x > score {
			winner, score = i, x
		}
	}
	logger.Println("elf", winner+1, "won with score", score)

	return score
}

func part1(logger *log.Logger) string {
	//highScore := part1impl(logger, 7, 25)
	highScore := part1impl(logger, 441, 71032)
	return fmt.Sprint(highScore)
}

func part2(logger *log.Logger) string {
	highScore := part1impl(logger, 441, 7103200)
	return fmt.Sprint(highScore)
}

func init() {
	util.RegisterSolution("day09part1", part1)
	util.RegisterSolution("day09part2", part2)
}
