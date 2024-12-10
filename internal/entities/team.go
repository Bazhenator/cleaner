package entities

import (
	"math/rand/v2"
	"time"

	"github.com/Bazhenator/cleaner/internal/logic/dto"
)

type Speed byte // Speed is a special type wich describes cleaning speed

const (
	Fast Speed = iota + 1
	Mid
	Slow
)

type Status byte // Status is a special type wich describes cleaning team's busyness

const (
	Available Status = iota
	Busy
)

type CleaningTeam struct {
	Id         uint64
	Request    *dto.Request
	Status     Status
	Speed      Speed
	StartedAt  time.Time
	FinishedAt time.Time
}

// AssignRequest assigns a cleaning request to the team
func (ct *CleaningTeam) AssignRequest(req *dto.Request) {
	ct.Request = req
	ct.Status = Busy
	ct.StartedAt = time.Now()
}

// CompleteCleaning marks the cleaning as completed
func (ct *CleaningTeam) CompleteCleaning() {
	ct.Status = Available
	ct.Request = nil
	ct.StartedAt = time.Time{}
	ct.FinishedAt = time.Time{}
}

// GetCleaningTime calculates the cleaning duration based on team speed and exponential distribution
func (ct *CleaningTeam) GetCleaningTime() time.Duration {
	baseTime := time.Duration(10) * time.Minute // Base cleaning time for Slow speed
	
	switch ct.Speed {
	case Mid:
		baseTime /= 2
	case Fast:
		baseTime /= 4
	}
	// Exponential distribution simulation
	lambda := 1 / float64(baseTime.Seconds())
	expTime := time.Duration(rand.ExpFloat64() / lambda) * time.Second
	return expTime
}
