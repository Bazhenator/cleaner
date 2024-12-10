package delivery

import (
	"context"

	"github.com/Bazhenator/cleaner/configs"
	"github.com/Bazhenator/cleaner/internal/logic"
	"github.com/Bazhenator/cleaner/internal/logic/dto"
	cleaner "github.com/Bazhenator/cleaner/pkg/api/grpc"
	"github.com/Bazhenator/tools/src/logger"
)

type CleanerServer struct {
	cleaner.UnimplementedCleanerServiceServer

	c *configs.Config
	l *logger.Logger

	logic logic.CleanerService
}

func NewCleanerServer(c *configs.Config, l *logger.Logger, logic logic.CleanerService) *CleanerServer {
	return &CleanerServer{
		c: c,
		l: l,

		logic: logic,
	}
}

func (s *CleanerServer) ProceedCleaning(ctx context.Context, in *cleaner.ProceedCleaningIn) (*cleaner.ProceedCleaningOut, error) {
	s.l.DebugCtx(ctx, "ProceedCleaning data", logger.NewField("", in))
	req := in.GetReq()

	answer := s.logic.ProceedCleaningRequest(ctx, &dto.ProceedCleaningRequestIn{
		TeamId: in.GetTeamId(),
		Request: &dto.Request{
			Id:           req.GetId(),
			ClientId:     req.GetClientId(),
			CleaningType: uint(req.GetCleaningType()),
			Priority:     uint(req.GetPriority()),
		},
	},
	)

	return &cleaner.ProceedCleaningOut{Duration: answer.Duration}, nil
}

func (s *CleanerServer) GetAvailableTeams(ctx context.Context, _ *cleaner.GetAvailableTeamsIn) *cleaner.GetAvailableTeamsOut {
	s.l.DebugCtx(ctx, "GetAvailableTeams data", logger.NewField("in", nil))

	answer := s.logic.GetAvailableTeams(ctx)

	return &cleaner.GetAvailableTeamsOut{TeamsIds: answer.Teams}
}
