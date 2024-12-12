package logic

import (
	"context"

	"github.com/Bazhenator/cleaner/internal/logic/dto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CleanerService interface {
	ProceedCleaningRequest(context.Context, *dto.ProceedCleaningRequestIn) *dto.ProceedCleaningRequestOut
	GetAvailableTeams(context.Context) *dto.GetAvailableTeamsOut
	GenerateReport()
}
