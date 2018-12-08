package util

import (
	"log"
	"time"
)

type Timer struct {
	Log            *log.Logger
	Prefix         string
	StartedAt      time.Time
	LastCheckpoint time.Time
}

func NewTimer(log *log.Logger, prefix string) Timer {
	now := time.Now()
	return Timer{log, prefix, now, now}
}

func (t *Timer) Checkpoint() (sinceStart time.Duration, sinceLast time.Duration) {
	now := time.Now()
	sinceStart = now.Sub(t.StartedAt)
	sinceLast = now.Sub(t.LastCheckpoint)
	t.LastCheckpoint = now
	return
}

func (t *Timer) LogCheckpoint(checkpointName string) {
	sinceStart, sinceLast := t.Checkpoint()
	t.Log.Printf("%v%v (elapsed %v, total %v)\n",
		t.Prefix, checkpointName, sinceLast, sinceStart)
}
