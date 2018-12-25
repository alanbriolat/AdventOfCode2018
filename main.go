package main

import (
	"flag"
	"fmt"
	"runtime/pprof"

	_ "github.com/alanbriolat/AdventOfCode2018/day01"
	_ "github.com/alanbriolat/AdventOfCode2018/day02"
	_ "github.com/alanbriolat/AdventOfCode2018/day03"
	_ "github.com/alanbriolat/AdventOfCode2018/day04"
	_ "github.com/alanbriolat/AdventOfCode2018/day05"
	_ "github.com/alanbriolat/AdventOfCode2018/day06"
	_ "github.com/alanbriolat/AdventOfCode2018/day07"
	_ "github.com/alanbriolat/AdventOfCode2018/day08"
	_ "github.com/alanbriolat/AdventOfCode2018/day09"
	_ "github.com/alanbriolat/AdventOfCode2018/day10"
	_ "github.com/alanbriolat/AdventOfCode2018/day11"
	_ "github.com/alanbriolat/AdventOfCode2018/day12"
	_ "github.com/alanbriolat/AdventOfCode2018/day13"
	_ "github.com/alanbriolat/AdventOfCode2018/day14"
	_ "github.com/alanbriolat/AdventOfCode2018/day15"
	_ "github.com/alanbriolat/AdventOfCode2018/day16"
	_ "github.com/alanbriolat/AdventOfCode2018/day17"
	_ "github.com/alanbriolat/AdventOfCode2018/day18"
	_ "github.com/alanbriolat/AdventOfCode2018/day19"
	_ "github.com/alanbriolat/AdventOfCode2018/day20"
	_ "github.com/alanbriolat/AdventOfCode2018/day22"
	_ "github.com/alanbriolat/AdventOfCode2018/day23"
	_ "github.com/alanbriolat/AdventOfCode2018/day24"
	_ "github.com/alanbriolat/AdventOfCode2018/day25"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io/ioutil"
	"log"
	"os"
)

var verbose = flag.Bool("v", false, "verbose logging")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	mainLog := log.New(os.Stdout, "main: ", 0)
	t := util.NewTimer(mainLog, "")
	defer t.LogCheckpoint("ran all solutions")

	flag.Parse()

	only := make(map[string]struct{})
	for _, name := range flag.Args() {
		only[name] = struct{}{}
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		util.Check(err)
		err = pprof.StartCPUProfile(f)
		util.Check(err)
		defer pprof.StopCPUProfile()
	}

	for _, s := range util.GetSolutions() {
		if _, ok := only[s.Name]; len(only) > 0 && !ok {
			continue
		}
		var logger *log.Logger
		if *verbose {
			logger = log.New(os.Stdout, s.Name + ": ", 0)
		} else {
			logger = log.New(ioutil.Discard, "", 0)
		}
		logger.Println("----------------")
		result := s.Run(logger)
		t.LogCheckpoint(fmt.Sprintf("%v answer: %v", s.Name, result))
		logger.Println("----------------")
	}
}
