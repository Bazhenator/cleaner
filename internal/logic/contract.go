package logic

import (
	"context"

	"github.com/Bazhenator/cleaner/internal/logic/dto"
)

type CleanerService interface {
	ProceedCleaningRequest(context.Context, *dto.ProceedCleaningRequestIn) (*dto.ProceedCleaningRequestOut, error)
	GetAvailableTeams(context.Context) (*dto.GetAvailableTeamsOut, error)
	GetTeamsStats(context.Context) (*dto.GetTeamsStatsOut, error)
}
