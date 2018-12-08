package util

import (
	"fmt"
	"time"
)

type Timer struct {
	Name           string
	StartedAt      time.Time
	LastCheckpoint time.Time
}

func NewTimer(name string) Timer {
	now := time.Now()
	return Timer{name, now, now}
}

func (t *Timer) Checkpoint() (sinceStart time.Duration, sinceLast time.Duration) {
	now := time.Now()
	sinceStart = now.Sub(t.StartedAt)
	sinceLast = now.Sub(t.LastCheckpoint)
	t.LastCheckpoint = now
	return
}

func (t *Timer) PrintCheckpoint(checkpointName string) {
	sinceStart, sinceLast := t.Checkpoint()
	fmt.Printf("%v - %v: elapsed %v, total %v\n",
		t.Name, checkpointName, sinceLast, sinceStart)
}
