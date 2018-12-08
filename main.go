package main

import (
	"flag"
	_ "github.com/alanbriolat/AdventOfCode2018/day01"
	_ "github.com/alanbriolat/AdventOfCode2018/day02"
	_ "github.com/alanbriolat/AdventOfCode2018/day03"
	_ "github.com/alanbriolat/AdventOfCode2018/day04"
	_ "github.com/alanbriolat/AdventOfCode2018/day05"
	_ "github.com/alanbriolat/AdventOfCode2018/day06"
	_ "github.com/alanbriolat/AdventOfCode2018/day07"
	_ "github.com/alanbriolat/AdventOfCode2018/day08"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io/ioutil"
	"log"
	"os"
)

var verbose = flag.Bool("v", false, "verbose logging")

func main() {
	mainLog := log.New(os.Stdout, "main: ", 0)
	t := util.NewTimer(mainLog, "")
	defer t.LogCheckpoint("ran all solutions")

	flag.Parse()

	for _, s := range util.GetSolutions() {
		var logger *log.Logger
		if *verbose {
			logger = log.New(os.Stdout, s.Name + ": ", 0)
		} else {
			logger = log.New(ioutil.Discard, "", 0)
		}
		logger.Println("----------------")
		result := s.Run(logger)
		mainLog.Printf("%v answer: %v", s.Name, result)
		logger.Println("----------------")
	}
}
