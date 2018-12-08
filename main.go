package main

import (
	_ "github.com/alanbriolat/AdventOfCode2018/day01"
	_ "github.com/alanbriolat/AdventOfCode2018/day02"
	_ "github.com/alanbriolat/AdventOfCode2018/day03"
	_ "github.com/alanbriolat/AdventOfCode2018/day04"
	_ "github.com/alanbriolat/AdventOfCode2018/day05"
	_ "github.com/alanbriolat/AdventOfCode2018/day06"
	_ "github.com/alanbriolat/AdventOfCode2018/day07"
	_ "github.com/alanbriolat/AdventOfCode2018/day08"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"log"
	"os"
)

func main() {
	t := util.NewTimer(log.New(os.Stdout, "main: ", 0), "")
	defer t.LogCheckpoint("ran all solutions")

	for _, s := range util.GetSolutions() {
		logger := log.New(os.Stdout, s.Name + ": ", 0)
		logger.Println("----------------")
		s.Run(logger)
		logger.Println("----------------")
	}
}
