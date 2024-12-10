package logic

import (
	"context"

	"github.com/Bazhenator/cleaner/internal/logic/dto"
)

type CleanerService interface {
	ProceedCleaningRequest(context.Context, *dto.ProceedCleaningRequestIn) *dto.ProceedCleaningRequestOut
	GetAvailableTeams(context.Context) *dto.GetAvailableTeamsOut
}
