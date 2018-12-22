package util

import (
	"container/heap"
	"fmt"
	"github.com/cheekybits/genny/generic"
	"math"
)

type SearchNode generic.Type

// Priority queue based on https://golang.org/pkg/container/heap/#example__priorityQueue
type SearchQueueItem struct {
	value    SearchNode
	priority func() int
	index    int
}

type SearchQueue struct {
	ctx *AStarSearchContext
	data []*SearchQueueItem
}

func NewSearchQueue(ctx *AStarSearchContext) SearchQueue {
	return SearchQueue{
		ctx,
		make([]*SearchQueueItem, 0, ctx.NodeCountMax),
	}
}

func (pq SearchQueue) Len() int {
	return len(pq.data)
}

func (pq SearchQueue) Less(i, j int) bool {
	x1, x2 := pq.data[i], pq.data[j]
	p1, p2 := x1.priority(), x2.priority()
	return p1 < p2 || p1 == p2 && pq.ctx.TieBreak(x1.value, x2.value)
}

func (pq SearchQueue) Swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
	pq.data[i].index = i
	pq.data[j].index = j
}

func (pq *SearchQueue) Push(x interface{}) {
	n := len(pq.data)
	item := x.(*SearchQueueItem)
	item.index = n
	pq.data = append(pq.data, item)
}

func (pq *SearchQueue) Pop() interface{} {
	old := pq.data
	n := len(old)
	item := old[n-1]
	item.index = -1
	pq.data = old[0:n-1]
	return item
}

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

	// Priority queue of the "open set"
	openSetQueue := NewSearchQueue(ctx)
	heap.Init(&openSetQueue)
	// Find "open set" queue item by search node value
	openSetMap := make(map[SearchNode]*SearchQueueItem)
	// Nodes already processed
	closedSet := make(map[SearchNode]bool)
	// Nodes queued or processed
	visited := make(map[SearchNode]bool)

	pushFunc := func(n SearchNode) {
		item := &SearchQueueItem{
			value: n,
			priority: func() int {
				return f(n)
			},
		}
		heap.Push(&openSetQueue, item)
		openSetMap[n] = item
		visited[n] = true
	}

	popFunc := func() SearchNode {
		item := heap.Pop(&openSetQueue).(*SearchQueueItem)
		delete(openSetMap, item.value)
		closedSet[item.value] = true	// Don't visit this node again
		return item.value
	}

	updateFunc := func(n SearchNode) {
		if item, ok := openSetMap[n]; ok {
			heap.Fix(&openSetQueue, item.index)
		} else {
			panic("AStarSearch: could not find openSetMap item")
		}
	}

	start := ctx.Start
	gScore[start] = 0
	fScore[start] = h(start)
	// First node to process is starting node
	pushFunc(start)

	// Keep track of most efficient path to each node
	cameFrom := make(map[SearchNode]SearchNode)
	path := func(destination SearchNode) []SearchNode {
		result := make([]SearchNode, 0)
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

	for openSetQueue.Len() > 0 {
		// Get most promising next node
		current := popFunc()
		// Did we find a goal?
		for _, d := range ctx.Destinations {
			if current == d {
				return path(current), nil
			}
		}

		// Score potential next nodes
		for _, neighbour := range ctx.Adjacent(current) {
			// Calculate new path cost
			nScore := gScore[current] + ctx.Cost(current, neighbour)

			if closedSet[neighbour] {
				// Already have a shortest path to this neighbour
				continue
			}
			if !visited[neighbour] {
				// Position we've never seen before, add to the queue
				pushFunc(neighbour)
			} else if nScore > g(neighbour) {
				// Already a better path to neighbour
				continue
			} else if nScore == g(neighbour) && !ctx.TieBreak(current, cameFrom[neighbour]) {
				// Already an equal path to neighbour which came from a "better" source
				// (according to the tie break function)
				continue
			}
			// This node is already in the queue, and we've found a better path to it
			cameFrom[neighbour] = current
			gScore[neighbour] = nScore
			fScore[neighbour] = nScore + h(neighbour)
			// Force the priority queue to update
			updateFunc(neighbour)
		}
	}

	return nil, fmt.Errorf("no path found to any destination")
}
