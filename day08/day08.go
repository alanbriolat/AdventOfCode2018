package day08

import (
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
)

type Node struct {
	UnreadChildren int
	UnreadMetadata int
	ChildValues []int
	Value int
}

func NewNode(children, metadata int) *Node {
	return &Node{
		children,
		metadata,
		make([]int, 0, children),
		0,
	}
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

func part1and2(logger *log.Logger) {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	input := readInput("day08/input.txt")
	t.LogCheckpoint(fmt.Sprintf("read %v numbers", len(input.Data)))


	stack := util.NewGenericStack(0)
	stack.Push(NewNode(input.Next(), input.Next()))
	sum := 0
	var top *Node
	for stack.Count() > 0 {
		top = stack.Top().(*Node)
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
			if stack.Count() > 0 {
				next := stack.Top().(*Node)
				next.ChildValues = append(next.ChildValues, top.Value)
			}
		}
	}
	logger.Println("sum of metadata entries:", sum)
	logger.Println("value of root node:", top.Value)
	t.LogCheckpoint(fmt.Sprintf("results"))
}

func init() {
	util.RegisterSolution("day08", part1and2)
}
