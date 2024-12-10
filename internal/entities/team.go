package entities

import "time"

type Speed byte // Speed is a special type wich describes cleaning speed

const (
	Fast Speed = iota + 1
	Mid
	Slow
)

type Status byte // Status is a special type wich describes cleaning team's busyness

const (
	Available Speed = iota
	Busy
)

type CleaningTeam struct {
	Id         uint64
	Status     Status
	Speed      Speed
	StartedAt  time.Time
	FinishedAt time.Time
}