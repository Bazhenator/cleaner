package delivery

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

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
	s.l.DebugCtx(ctx, "ProceedCleaning started with", logger.NewField("data", in))
	req := in.GetReq()

	answer, err := s.logic.ProceedCleaningRequest(ctx, &dto.ProceedCleaningRequestIn{
		TeamId: in.GetTeamId(),
		Request: &dto.Request{
			Id:           req.GetId(),
			ClientId:     req.GetClientId(),
			CleaningType: uint(req.GetCleaningType()),
			Priority:     uint(req.GetPriority()),
		},
	})
	if err != nil {
		s.l.ErrorCtx(ctx, "error occurred:", logger.NewErrorField(err))
		return nil, err
	}

	totalTime := answer.Req.TimeInCleaner.Seconds()

	return &cleaner.ProceedCleaningOut{Req: &cleaner.Request{
		Id:            answer.Req.Id,
		ClientId:      answer.Req.ClientId,
		Priority:      uint32(answer.Req.Priority),
		CleaningType:  uint32(answer.Req.CleaningType),
		TeamId:        &answer.Req.TeamId,
		TimeInCleaner: &totalTime,
	}}, nil
}

func (s *CleanerServer) GetAvailableTeams(ctx context.Context, _ *emptypb.Empty) (*cleaner.GetAvailableTeamsOut, error) {
	s.l.Debug("GetAvailableTeams requested teams")

	answer, err := s.logic.GetAvailableTeams(ctx)
	if err != nil {
		s.l.ErrorCtx(ctx, "error occurred:", logger.NewErrorField(err))
		return nil, err
	}

	return &cleaner.GetAvailableTeamsOut{TeamsIds: answer.Teams}, nil
}

func (s *CleanerServer) GetTeamsStats(ctx context.Context, _ *emptypb.Empty) (*cleaner.GetTeamsStatsOut, error) {
	s.l.Debug("GetTeamsStats requested stats")

	stats, err := s.logic.GetTeamsStats(ctx)
	if err != nil {
		s.l.ErrorCtx(ctx, "error occurred:", logger.NewErrorField(err))
		return nil, err
	}

	answer := make([]*cleaner.Team, 0, len(stats.Stats)) 
	for _, stat := range stats.Stats {
		totalTime := stat.TotalBusyTime.Seconds()

		answer = append(answer, &cleaner.Team{
			Id: stat.Id,
			Speed: stat.Speed,
			ProcessedRequests: stat.ProcessedRequests,
			TotalBusyTime: totalTime,
		})
	}

	return &cleaner.GetTeamsStatsOut{Teams: answer}, nil
}
