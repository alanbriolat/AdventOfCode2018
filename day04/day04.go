/*
Day 04

Part 1
======

- Read input as stream of event objects
- Sort the input by the time
- Record the events as statistics for guards
- Search the statistics for the most sleepy guard's most sleepy minute
 */
package day04

import (
	"bufio"
	"fmt"
	"github.com/alanbriolat/AdventOfCode2018/util"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type EventType int

const (
	Begin EventType = 1
	Sleep EventType = 2
	Wake EventType = 3
)

type Event struct {
	Timestamp time.Time
	Type EventType
	GuardId int
}

type EventStream []Event

// Implement sort.Interface for slice of events
func (es EventStream) Len() int 			{ return len(es) }
func (es EventStream) Swap(i, j int)		{ es[i], es[j] = es[j], es[i] }
func (es EventStream) Less(i, j int) bool	{ return es[i].Timestamp.Before(es[j].Timestamp) }

// Propagate guard IDs to subsequent events
func (es EventStream) Propagate() {
	id := 0
	for i := range es {
		e := &es[i]
		if e.GuardId != 0 {
			id = e.GuardId
		} else {
			e.GuardId = id
		}
	}
}

type Guard struct {
	Id int
	Asleep bool
	AsleepSince time.Time
	SleepStats [60]int
	TotalSleep int
}

func (g *Guard) FindSleepyMinute() (minute, total int) {
	minute, total = 0, 0
	for i, x := range g.SleepStats {
		if x > total {
			minute, total = i, x
		}
	}
	return
}

type Roster map[int]*Guard

func (r Roster) FindGuard(id int) *Guard {
	g, ok := r[id]
	if !ok {
		g = &Guard{Id: id}
		r[id] = g
	}
	return g
}

func (r Roster) ApplyEvent(e Event) {
	g := r.FindGuard(e.GuardId)
	switch e.Type {
	case Sleep:
		g.Asleep = true
		g.AsleepSince = e.Timestamp
	case Wake:
		g.Asleep = false
		startMinute := g.AsleepSince.Minute()
		endMinute := e.Timestamp.Minute()
		g.TotalSleep += endMinute - startMinute
		for i := startMinute; i < endMinute; i++ {
			g.SleepStats[i]++
		}
	}
}

func (r Roster) FindSleepyGuard() *Guard {
	var result *Guard = nil
	for _, g := range r {
		if result == nil || g.TotalSleep > result.TotalSleep {
			result = g
		}
	}
	return result
}

func (r Roster) FindConsistentlySleepyGuard() (guard *Guard, minute int) {
	guard = nil
	timesAsleep := 0
	for _, g := range r {
		for i, x := range g.SleepStats {
			if guard == nil || x > timesAsleep {
				guard = g
				minute = i
				timesAsleep = x
			}
		}
	}
	return
}

func ReadEventsFromFile(name string) ([]Event, error) {
	result := make([]Event, 0, 100)
	var err error
	var reader io.Reader
	if reader, err = os.Open(name); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		event := Event{}
		line := scanner.Text()
		timestamp := string(line[1:17])
		data := string(line[19:])
		event.Timestamp, err = time.Parse("2006-01-02 15:04", timestamp)
		if err != nil {
			return nil, err
		}
		switch data {
		case "falls asleep":
			event.Type = Sleep
		case "wakes up":
			event.Type = Wake
		default:
			buf := strings.NewReader(data)
			_, err := fmt.Fscanf(buf, "Guard #%d begins shift", &event.GuardId)
			if err != nil {
				return nil, fmt.Errorf("failed to read event type: %v", err)
			}
			event.Type = Begin
		}
		result = append(result, event)
	}
	return result, nil
}

func part1and2(logger *log.Logger) string {
	t := util.NewTimer(logger, "")
	defer t.LogCheckpoint("end")

	events, err := ReadEventsFromFile("day04/input1.txt")
	util.Check(err)
	eventStream := EventStream(events)
	t.LogCheckpoint(fmt.Sprint("read ", len(events), " events"))

	sort.Sort(eventStream)
	t.LogCheckpoint(fmt.Sprint("sorted events"))

	eventStream.Propagate()
	t.LogCheckpoint("propagated guard IDs")

	roster := Roster{}
	for _, e := range eventStream {
		roster.ApplyEvent(e)
	}
	t.LogCheckpoint("applied events to roster")

	sleepyGuard := roster.FindSleepyGuard()
	logger.Printf("sleepiest guard: #%v, %v minutes\n", sleepyGuard.Id, sleepyGuard.TotalSleep)
	sleepyMinute, sleepyMinuteTotal := sleepyGuard.FindSleepyMinute()
	logger.Printf("sleepiest minute: 00:%02d, %d times\n", sleepyMinute, sleepyMinuteTotal)
	logger.Println("part1 answer:", sleepyGuard.Id, "*", sleepyMinute, "=", sleepyGuard.Id * sleepyMinute)
	t.LogCheckpoint("found sleepiest")

	consistentGuard, consistentMinute := roster.FindConsistentlySleepyGuard()
	logger.Printf("consistent guard: #%v, minute %02d, %v times\n",
		consistentGuard.Id, consistentMinute, consistentGuard.SleepStats[consistentMinute])
	logger.Printf("part2 answer: %v * %v = %v\n",
		consistentGuard.Id, consistentMinute, consistentGuard.Id * consistentMinute)
	t.LogCheckpoint("found consistent")

	return fmt.Sprintf("part1 = %v , part2 = %v",
		sleepyGuard.Id * sleepyMinute, consistentGuard.Id * consistentMinute)
}

func init() {
	util.RegisterSolution("day04", part1and2)
}
