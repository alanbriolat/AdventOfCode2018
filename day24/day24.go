package day24

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"github.com/alecthomas/participle"
	"log"
	"os"
	"sort"
)

type AttackType string

type AttackSet map[AttackType]bool

type ParsedAttribute struct {
	Kind  string   `@("weak" | "immune") "to"`
	Types []string `@Ident ("," @Ident)*`
}

type ParsedGroup struct {
	Units      int               `@Int "units"`
	UnitHP     int               `"each" "with" @Int "hit" "points"`
	Attributes []ParsedAttribute `("(" @@ (";" @@)* ")")?`
	Power      int               `"with" "an" "attack" "that" "does" @Int`
	Attack     string            `@Ident "damage"`
	Initiative int               `"at" "initiative" @Int`
}

type ParsedBattle struct {
	ImmuneSystem []ParsedGroup `"Immune" "System" ":" (@@)*`
	Infection    []ParsedGroup `"Infection" ":" (@@)*`
}

type Group struct {
	IsInfection bool
	Units       int
	UnitHP      int
	Multiplier  map[string]int
	Power       int
	Attack      string
	Initiative  int
}

func (g *ParsedGroup) ToGroup(isInfection bool) *Group {
	group := &Group{
		IsInfection: isInfection,
		Units:       g.Units,
		UnitHP:      g.UnitHP,
		Multiplier:  make(map[string]int),
		Power:       g.Power,
		Attack:      g.Attack,
		Initiative:  g.Initiative,
	}
	for _, attribute := range g.Attributes {
		m := 1
		switch attribute.Kind {
		case "weak":
			m = 2
		case "immune":
			m = 0
		}
		for _, damageType := range attribute.Types {
			group.Multiplier[damageType] = m
		}
	}
	return group
}

func (g *Group) Copy() *Group {
	// The copy has a reference to the same Multiplier map, but that shouldn't ever be changed
	result := *g
	return &result
}

func (g *Group) EffectiveDamage(damage int, damageType string) int {
	if m, ok := g.Multiplier[damageType]; ok {
		return damage * m
	} else {
		return damage
	}
}

func (g *Group) EffectivePower() int {
	return g.Units * g.Power
}

func (g *Group) Damaged(damage int) {
	g.Units -= damage / g.UnitHP
}

func (g *Group) IsDead() bool {
	return g.Units <= 0
}

type Battle struct {
	Groups []*Group
}

func (b *ParsedBattle) ToBattle() *Battle {
	size := len(b.ImmuneSystem) + len(b.Infection)
	battle := &Battle{
		Groups:      make([]*Group, 0, size),
	}
	for _, parsedGroup := range b.ImmuneSystem {
		g := parsedGroup.ToGroup(false)
		battle.Groups = append(battle.Groups, g)
	}
	for _, parsedGroup := range b.Infection {
		g := parsedGroup.ToGroup(true)
		battle.Groups = append(battle.Groups, g)
	}
	return battle
}

func (b *Battle) Copy() *Battle {
	result := &Battle{
		Groups: make([]*Group, len(b.Groups)),
	}
	for i, g := range b.Groups {
		result.Groups[i] = g.Copy()
	}
	return result
}

func (b *Battle) BoostImmuneSystem(amount int) {
	for _, g := range b.Groups {
		if !g.IsInfection {
			g.Power += amount
		}
	}
}

func (b *Battle) FilterDeadGroups() {
	filter := func(slice *[]*Group) {
		next := 0
		for i := 0; i < len(*slice); i++ {
			g := (*slice)[i]
			if !g.IsDead() {
				(*slice)[next] = g
				next++
			}
		}
		*slice = (*slice)[0:next]
	}
	filter(&b.Groups)
}

func (b *Battle) CountUnits(isInfection bool) int {
	result := 0
	for _, g := range b.Groups {
		if g.IsInfection == isInfection && g.Units > 0 {
			result += g.Units
		}
	}
	return result
}

func (b *Battle) Fight() (immuneCount, infectionCount int) {
	// Target selection phase:
	// Groups choose targets in decreasing effective power order, tie-broken by decreasing initiative order
	selectionOrder := make([]*Group, len(b.Groups))
	copy(selectionOrder, b.Groups)
	sort.Slice(selectionOrder, func(i, j int) bool {
		a, b := selectionOrder[i], selectionOrder[j]
		pa, pb := a.EffectivePower(), b.EffectivePower()
		return pa > pb || (pa == pb && a.Initiative > b.Initiative)
	})
	targeting := make(map[*Group]*Group)
	targeted := make(map[*Group]*Group)
	for _, g := range selectionOrder {
		p := g.EffectivePower()
		var target *Target
		for _, t := range b.Groups {
			if t.IsInfection == g.IsInfection {
				// Not a valid target
				continue
			}
			if _, ok := targeted[t]; ok {
				// Each group can only be targeted once
				continue
			}
			newTarget := &Target{
				Group: t,
				PotentialDamage: t.EffectiveDamage(p, g.Attack),
			}
			if newTarget.PotentialDamage > 0 && newTarget.BetterThan(target) {
				target = newTarget
			}
		}
		if target != nil {
			//fmt.Printf("%+v will attack %+v\n", g, target.Group)
			targeting[g] = target.Group
			targeted[target.Group] = g
		}
	}

	// Attacking phase:
	attackOrder := make([]*Group, len(b.Groups))
	copy(attackOrder, b.Groups)
	sort.Slice(attackOrder, func(i, j int) bool {
		a, b := attackOrder[i], attackOrder[j]
		return a.Initiative > b.Initiative
	})
	for _, g := range attackOrder {
		// If we have a target, and neither the attacker nor the target are dead yet
		if t, ok := targeting[g]; ok && !g.IsDead() && !t.IsDead() {
			t.Damaged(t.EffectiveDamage(g.EffectivePower(), g.Attack))
		}
	}

	// Clean up dead groups
	b.FilterDeadGroups()

	// Count remaining forces
	return b.CountUnits(false), b.CountUnits(true)
}

func (b *Battle) Run() (immuneCount, infectionCount int) {
	immuneCount, infectionCount = b.CountUnits(false), b.CountUnits(true)
	//fmt.Printf("start: immune = %d, infection = %d\n", immuneCount, infectionCount)
	for i := 0; immuneCount > 0 && infectionCount > 0; i++ {
		newImmuneCount, newInfectionCount := b.Fight()
		if newImmuneCount == immuneCount && newInfectionCount == infectionCount {
			// Reached a stalemate!
			return
		} else {
			immuneCount, infectionCount = newImmuneCount, newInfectionCount
		}
		//fmt.Printf("fight %d: immune = %d, infection = %d\n", i+1, immuneCount, infectionCount)
	}
	return
}

type Target struct {
	Group *Group
	PotentialDamage int
}

func (t *Target) BetterThan(o *Target) bool {
	if o == nil {
		return true
	}
	switch {
	case t.PotentialDamage > o.PotentialDamage:
		return true
	case t.PotentialDamage < o.PotentialDamage:
		return false
	}
	pt, po := t.Group.EffectivePower(), o.Group.EffectivePower()
	switch {
	case pt > po:
		return true
	case pt < po:
		return false
	}
	switch {
	case t.Group.Initiative > o.Group.Initiative:
		return true
	default:
		return false
	}
}

func parseBattle(filename string) *Battle {
	var err error

	reader, err := os.Open(filename)
	util.Check(err)

	parser := participle.MustBuild(&ParsedBattle{})
	parsedBattle := &ParsedBattle{}
	err = parser.Parse(reader, parsedBattle)
	util.Check(err)

	return parsedBattle.ToBattle()
}

func part1impl(logger *log.Logger, filename string) int {
	battle := parseBattle(filename)
	immuneCount, infectionCount := battle.Run()
	// One of these should be 0
	return immuneCount + infectionCount
}

/*
Do a binary search on immune system boost amounts to find the smallest amount where the immune system wins.
 */
func part2impl(logger *log.Logger, filename string) int {
	prototype := parseBattle(filename)

	// Evaluate if `boost` is sufficient to win
	evaluationFunc := func(boost int) (int, bool) {
		battle := prototype.Copy()
		battle.BoostImmuneSystem(boost)
		immuneCount, infectionCount := battle.Run()
		return immuneCount, infectionCount == 0
	}

	start := 0
	end := 0
	result := 0
	initialStep := 1024

	// Find an interval which contains the optimum
	if _, win := evaluationFunc(start); win {
		panic("winning without a boost")
	}
	for {
		end += initialStep
		_, win := evaluationFunc(end)
		if win {
			break
		} else {
			start = end
		}
	}
	//logger.Printf("minimum boost is between %d and %d\n", start, end)

	// Optimise `start` to the last value that loses and `end` to the first value that wins
	for {
		boost := (start + end) / 2
		immuneCount, win := evaluationFunc(boost)
		//logger.Printf("evaluated %d, immune = %d, win = %v", boost, immuneCount, win)
		if win {
			end = boost
			result = immuneCount
		} else {
			start = boost
		}
		//logger.Printf("minimum boost is between %d and %d\n", start, end)
		// Once the values are adjacent, `end` should be the lowest boost that wins
		if end - start == 1 {
			logger.Printf("binary search found boost = %d, immune = %d\n", end, result)
			break
		}
	}

	// We may have made a bad assumption that this is solvable by binary search, so let's scan
	// from 1 to `end` and find the first win
	for boost := 1; boost <= end; boost++ {
		if immuneCount, win := evaluationFunc(boost); win {
			logger.Printf("linear search found boost = %d, immune = %d\n", boost, immuneCount)
			end = boost
			result = immuneCount
			break
		}
	}

	return result
}

func init() {
	//util.RegisterSolution("day24test1", func(logger *log.Logger) string {
	//	return fmt.Sprint(part1impl(logger, "day24/input_test.txt"))
	//})
	util.RegisterSolution("day24part1", func(logger *log.Logger) string {
		return fmt.Sprint(part1impl(logger, "day24/input.txt"))
	})
	//util.RegisterSolution("day24test2", func(logger *log.Logger) string {
	//	return fmt.Sprint(part2impl(logger, "day24/input_test.txt"))
	//})
	util.RegisterSolution("day24part2", func(logger *log.Logger) string {
		return fmt.Sprint(part2impl(logger, "day24/input.txt"))
	})
}
