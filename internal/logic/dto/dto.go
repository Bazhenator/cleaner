package dto

type Request struct {
	Id           uint64
	ClientId     uint64
	CleaningType uint
	Priority     uint
}

type ProceedCleaningRequestIn struct {
	TeamId  uint64
	Request *Request
}

type ProceedCleaningRequestOut struct {
	Duration string
}

type GetAvailableTeamsOut struct {
	Teams []uint64
}
