package util

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
	"math"
	"sort"
)

type SearchNode generic.Type

type AStarSearchContext struct {
	// Initial state
	Start SearchNode
	// All allowable end states
	Destinations []SearchNode
	// Estimate of number of possible states
	NodeCountMax int
	// All valid states adjacent to n
	Adjacent func(n SearchNode) []SearchNode
	// Estimate of cost to move from n1 to n2 (not necessarily adjacent)
	Heuristic func(n1, n2 SearchNode) int
	// Actual cost to move from n1 to n2 (adjacent)
	Cost func(n1, n2 SearchNode) int
	// Is n1 chosen first if it has the same path cost as n2?
	TieBreak func(n1, n2 SearchNode) bool
}

/*
https://en.wikipedia.org/wiki/A*_search_algorithm
 */
func AStarSearch(ctx *AStarSearchContext) ([]SearchNode, error) {
	if len(ctx.Destinations) == 0 {
		return nil, fmt.Errorf("no destinations")
	}

	// f(n) = g(n) + h(n)
	// g(n): cost to get to n from start
	// h(n): estimate of cost from n to goal
	gScore := make(map[SearchNode]int)
	g := func(n SearchNode) int {
		if result, ok := gScore[n]; ok {
			return result
		} else {
			return math.MaxInt32
		}
	}
	fScore := make(map[SearchNode]int)
	f := func(n SearchNode) int {
		if result, ok := fScore[n]; ok {
			return result
		} else {
			return math.MaxInt32
		}
	}
	h := func(n SearchNode) int {
		min := math.MaxInt32
		for _, d := range ctx.Destinations {
			if distance := ctx.Heuristic(n, d); distance < min {
				min = distance
			}
		}
		return min
	}

	// Position processing queue
	openSet := make([]SearchNode, 0, ctx.NodeCountMax)
	// Positions already processed
	closedSet := make(map[SearchNode]bool)
	// Positions queued or processed
	visited := make(map[SearchNode]bool)

	start := ctx.Start
	gScore[start] = 0
	fScore[start] = h(start)
	openSet = append(openSet, start)
	visited[start] = true

	// Keep track of most efficient path to each position
	cameFrom := make(map[SearchNode]SearchNode)
	path := func(destination SearchNode) []SearchNode {
		result := make([]SearchNode, 0, g(destination))		// Path cost is a good guess of path size to preallocate
		next := destination
		for next != start {
			result = append(result, next)
			next = cameFrom[next]
		}
		// Reverse the path, so it's from start to goal
		for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
			result[left], result[right] = result[right], result[left]
		}
		return result
	}

	for len(openSet) > 0 {
		// Sort by f(n), tie-break on reading order
		sort.Slice(openSet, func(i, j int) bool {
			p1, p2 := openSet[i], openSet[j]
			f1, f2 := f(p1), f(p2)
			return f1 < f2 || (f1 == f2 && ctx.TieBreak(p1, p2))
		})
		// Get most promising next position
		var current SearchNode
		current, openSet = openSet[0], openSet[1:]
		// Did we find a goal?
		for _, d := range ctx.Destinations {
			if current == d {
				return path(current), nil
			}
		}
		// Don't visit this position again
		closedSet[current] = true

		// Score potential next positions
		for _, neighbour := range ctx.Adjacent(current) {
			// Calculate new path cost
			nScore := gScore[current] + ctx.Cost(current, neighbour)

			if closedSet[neighbour] {
				// Already have a shortest path to this neighbour
				continue
			}
			if !visited[neighbour] {
				// Position we've never seen before, add to the queue
				openSet = append(openSet, neighbour)
				visited[neighbour] = true
			} else if nScore > g(neighbour) {
				// Already a better path to neighbour
				continue
			} else if nScore == g(neighbour) && !ctx.TieBreak(current, cameFrom[neighbour]) {
				// Already an equal path to neighbour which came from higher in
				// the "reading order". The tendency to follow the reading order
				// also means a square will be visited from a higher reading
				// order square if possible.
				continue
			}
			// This position is already in the queue, and we've found a better path to it
			cameFrom[neighbour] = current
			gScore[neighbour] = nScore
			fScore[neighbour] = nScore + h(neighbour)
		}
	}

	return nil, fmt.Errorf("no path found to any destination")
}
