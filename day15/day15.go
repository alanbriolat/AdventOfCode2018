package day15

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
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

type Unit struct {
	Battle      *Battle
	Id          byte
	Position    util.Vec2D
	IsGoblin    bool
	HitPoints   int
	AttackPower int
}

func (u *Unit) String() string {
	unitType := InputElf
	if u.IsGoblin {
		unitType = InputGoblin
	}
	return fmt.Sprintf("%s%d(%d)@%d,%d", string(unitType), u.Id, u.HitPoints, u.Position.X, u.Position.Y)
}

type Battle struct {
	Units   []*Unit
	Map     [][]byte
	MapSize util.Vec2D
}

func NewBattle(input []string) Battle {
	b := Battle{}
	b.Units = make([]*Unit, 0)
	b.MapSize = util.Vec2D{len(input[0]), len(input)}
	b.Map = util.NewByteGrid(b.MapSize.X, b.MapSize.Y)

	for y, line := range input {
		for x := range line {
			switch line[x] {
			case InputWall:
				b.Map[x][y] = MapWall
			case InputFloor:
				b.Map[x][y] = MapFloor
			case InputElf:
				b.CreateUnit(false, x, y)
			case InputGoblin:
				b.CreateUnit(true, x, y)
			}
		}
	}

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

func (b *Battle) CreateUnit(isGoblin bool, x, y int) {
	u := Unit{
		b,
		byte(len(b.Units)),
		util.Vec2D{x, y},
		isGoblin,
		200,
		3,
	}
	b.Units = append(b.Units, &u)
	b.Map[x][y] = u.Id
}

func part1(logger *log.Logger, filename string) string {
	input, _ := util.ReadLinesFromFile(filename)
	battle := NewBattle(input)
	logger.Print("\n", battle.String())

	return ""
}

func init() {
	util.RegisterSolution("day15test1", func(logger *log.Logger) string {
		return part1(logger, "day15/input_test1.txt")
	})
	util.RegisterSolution("day15test2", func(logger *log.Logger) string {
		return part1(logger, "day15/input_test2.txt")
	})
}
