package domain

import (
	"time"
)

type Message struct {
	Text     string
	Priority Priority
	Delay    time.Duration
	Tags     []string
	Title    string
}

type Priority uint8

const (
	PriorityUnknown Priority = iota
	PriorityMin
	PriorityLow
	PriorityDefault
	PriorityHigh
	PriorityMax
)
