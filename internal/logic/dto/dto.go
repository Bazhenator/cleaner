package dto

import (
	"time"
)

type Request struct {
	Id            uint64
	ClientId      uint64
	TeamId        uint64
	CleaningType  uint
	Priority      uint
	TimeInCleaner time.Duration
}

type ProceedCleaningRequestIn struct {
	TeamId  uint64
	Request *Request
}

type ProceedCleaningRequestOut struct {
	Req *Request
}

type GetAvailableTeamsOut struct {
	Teams []uint64
}

type TeamStats struct {
	Id                uint64
	Speed             uint32
	ProcessedRequests uint64
	TotalBusyTime     time.Duration
}

type GetTeamsStatsOut struct {
	Stats []*TeamStats
}