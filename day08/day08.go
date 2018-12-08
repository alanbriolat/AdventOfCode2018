package main

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
)

type Node struct {
	UnreadChildren int
	UnreadMetadata int
	ChildValues []int
	Value int
}

func NewNode(children, metadata int) Node {
	return Node{
		children,
		metadata,
		make([]int, 0, children),
		0,
	}
}

type NodeStack struct {
	Data []Node
}

func NewNodeStack(size int) NodeStack {
	return NodeStack{make([]Node, 0, size) }
}

func (s *NodeStack) Push(x Node) {
	s.Data = append(s.Data, x)
}

func (s *NodeStack) Peek() (*Node, bool) {
	last := len(s.Data) - 1
	if last < 0 {
		return nil, false
	} else {
		return &s.Data[last], true
	}
}

func (s *NodeStack) Pop() (*Node, bool) {
	if result, ok := s.Peek(); !ok {
		return nil, ok
	} else {
		s.Data = s.Data[:len(s.Data)-1]
		return result, true
	}
}

func (s *NodeStack) Count() int {
	return len(s.Data)
}

type Input struct {
	Data      []int
	Remaining []int
}

func (i *Input) Next() (result int) {
	result, i.Remaining = i.Remaining[0], i.Remaining[1:]
	return
}

func readInput(filename string) Input {
	result, err := util.ReadIntsFromFile(filename)
	util.Check(err)
	return Input{result[:], result[:]}
}

func part1and2() {
	t := util.NewTimer("day08part1and2")
	defer t.PrintCheckpoint("end")

	input := readInput("input.txt")
	t.PrintCheckpoint(fmt.Sprintf("read %v numbers", len(input.Data)))


	stack := NewNodeStack(0)
	stack.Push(NewNode(input.Next(), input.Next()))
	sum := 0
	var top *Node
	for stack.Count() > 0 {
		top, _ = stack.Peek()
		switch {
		case top.UnreadChildren > 0:
			// Still have child nodes to read - read one and it'll get processed on next loop
			top.UnreadChildren--
			stack.Push(NewNode(input.Next(), input.Next()))
		case top.UnreadMetadata > 0:
			// No child nodes, but have metadata - read it all and sum it
			for top.UnreadMetadata > 0 {
				top.UnreadMetadata--
				i := input.Next()
				sum += i
				if len(top.ChildValues) > 0 {
					if i > 0 && i <= len(top.ChildValues) {
						top.Value += top.ChildValues[i-1]
					}
				} else {
					top.Value += i
				}
			}
			fallthrough
		default:
			// If there are no child nodes, and all metadata has been processed, finished with this node
			stack.Pop()
			if next, ok := stack.Peek(); ok {
				next.ChildValues = append(next.ChildValues, top.Value)
			}
		}
	}
	fmt.Println("sum of metadata entries:", sum)
	fmt.Println("value of root node:", top.Value)
	t.PrintCheckpoint(fmt.Sprintf("results"))
}

func main() {
	part1and2()
}
