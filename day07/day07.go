package day07

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
	"sort"
)

type Dependency struct {
	Step byte
	DependsOn byte
}

type DependencyTree struct {
	controls map[byte]map[byte]bool
	depends map[byte]map[byte]bool
}

func NewDependencyTree() DependencyTree {
	return DependencyTree{
		make(map[byte]map[byte]bool),
		make(map[byte]map[byte]bool),
	}
}

func (t *DependencyTree) Controls(s byte) map[byte]bool {
	result, ok := t.controls[s]
	if !ok {
		result = make(map[byte]bool)
		t.controls[s] = result
	}
	return result
}

func (t *DependencyTree) Depends(s byte) map[byte]bool {
	result, ok := t.depends[s]
	if !ok {
		result = make(map[byte]bool)
		t.depends[s] = result
	}
	return result
}

/*
Record that a depends on b
 */
func (t *DependencyTree) AddDependency(s, d byte) {
	// mark the relationship between the two
	t.Depends(s)[d] = true
	t.Controls(d)[s] = true
	// make sure the dependency exists in depends list too
	t.Depends(d)
}

func (t *DependencyTree) Resolve(d byte) []byte {
	result := make([]byte, 0)
	// find steps "controlled by" d
	c := t.Controls(d)
	// resolving d, so it doesn't control anything in the future
	delete(t.controls, d)
	// for each step that was "controlled by" d...
	for s := range c {
		// remove d from the list steps it "depends on"
		sDeps := t.Depends(s)
		delete(sDeps, d)
		// if there are no more dependencies, we're free to resolve s now
		if len(sDeps) == 0 {
			result = append(result, s)
		}
	}
	// return list of steps that can now be resolved
	return result
}

func (t *DependencyTree) ResolveAll() []byte {
	steps := make([]byte, 0, len(t.depends))
	next := steps[:]

	// Find steps that have no dependencies to start with
	for s, deps := range t.depends {
		if len(deps) == 0 {
			next = append(next, s)
		}
	}

	for len(next) > 0 {
		// Sort alphabetically, because that's how we tiebreak
		sort.Slice(next, func(i, j int) bool {
			return next[i] < next[j]
		})
		// Resolve the step, find more steps that can be resolved
		more := t.Resolve(next[0])
		//logger.Println("steps:", steps, "next:", next, "more:", more)
		// Drop the step we just resolved, include new ones
		next = append(next[1:], more...)
		// Make sure the step we just resolved is kept in steps
		steps = steps[:len(steps)+1]
	}

	return steps
}

func (t *DependencyTree) ResolveAllParallel(workers int) int {
	type Job struct {
		Step byte
		FinishedAt int
	}

	jobs := make([]Job, 0, len(t.depends))
	complete := jobs[:]
	running := jobs[:]
	next := jobs[:]
	now := 0

	for s, deps := range t.depends {
		if len(deps) == 0 {
			next = append(next, Job{Step: s})
		}
	}

	for len(next) > 0 || len(running) > 0 {
		// Sort next complete alphabetically, because that's how we tiebreak
		sort.Slice(next, func(i, j int) bool {
			return next[i].Step < next[j].Step
		})
		// Allocate work to free workers
		for len(next) > 0 && len(running) <= workers {
			j := &next[0]
			j.FinishedAt = now + cost(j.Step)
			// Move job from next to running
			running, next = running[:len(running)+1], next[1:]
		}
		// Sort running by finish order, to free up workers in the right order
		sort.Slice(running, func(i, j int) bool {
			return running[i].FinishedAt < running[j].FinishedAt
		})
		// Resolve the step from the next running job
		j := running[0]
		now = j.FinishedAt
		more := t.Resolve(j.Step)
		// Move job from running to complete
		complete, running = complete[:len(complete)+1], running[1:]
		// Include new jobs
		for _, s := range more {
			next = append(next, Job{Step: s})
		}
	}

	return now
}

func cost(s byte) int {
	return 60 + int(s - 'A' + 1)
}

func ReadDependencies(filename string) []Dependency {
	result := make([]Dependency, 0)
	reader, err := os.Open(filename)
	util.Check(err)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		d := Dependency{line[36], line[5]}
		result = append(result, d)
	}
	return result
}

func part1(logger *log.Logger) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	dependencies := ReadDependencies("day07/input.txt")
	t.LogCheckpoint(fmt.Sprintf("read %v dependencies", len(dependencies)))

	depTree := NewDependencyTree()
	for _, d := range dependencies {
		depTree.AddDependency(d.Step, d.DependsOn)
	}
	t.LogCheckpoint("build dependency tree")

	steps := depTree.ResolveAll()
	logger.Println("steps:", string(steps))
	t.LogCheckpoint("resolved dependency graph")

	return string(steps)
}

func part2(logger *log.Logger) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	dependencies := ReadDependencies("day07/input.txt")
	t.LogCheckpoint(fmt.Sprintf("read %v dependencies", len(dependencies)))

	depTree := NewDependencyTree()
	for _, d := range dependencies {
		depTree.AddDependency(d.Step, d.DependsOn)
	}
	t.LogCheckpoint("build dependency tree")

	duration := depTree.ResolveAllParallel(5)
	logger.Println("completed steps in", duration, "seconds")
	t.LogCheckpoint("resolved dependency graph")

	return fmt.Sprint(duration)
}

func init() {
	util.RegisterSolution("day07part1", part1)
	util.RegisterSolution("day07part2", part2)
}
