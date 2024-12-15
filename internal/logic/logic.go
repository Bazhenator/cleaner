package logic

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"

	"github.com/Bazhenator/cleaner/configs"
	"github.com/Bazhenator/cleaner/internal/entities"
	"github.com/Bazhenator/cleaner/internal/logic/dto"
	"github.com/Bazhenator/tools/src/logger"
)

type Service struct {
	c  *configs.Config
	l  *logger.Logger
	mu sync.Mutex

	teams []*entities.CleaningTeam
}

func NewService(c *configs.Config, l *logger.Logger) *Service {
	// Cleaning teams' initializing
	teams := initTeams(c.TeamsAmount)

	return &Service{
		c: c,
		l: l,

		teams: teams,
	}
}

// ProceedCleaningRequest proceeds request from user, assigns it to cleaning team and processes it.
// Returns cleaning duration
func (s *Service) ProceedCleaningRequest(ctx context.Context, in *dto.ProceedCleaningRequestIn) (*dto.ProceedCleaningRequestOut, error) {
	team := s.teams[in.TeamId]
	duration := team.GetCleaningTime(s.c.BaseSpeed)
	team.AssignRequest(in.Request)

	go func(team *entities.CleaningTeam, duration time.Duration) {
		time.Sleep(duration)
		s.mu.Lock()
		defer s.mu.Unlock()

		team.CompleteCleaning(team.StartedAt)

		s.l.Info(fmt.Sprintf("Team %d completed cleaning.", team.Id))
	}(team, duration)

	team.Request.TimeInCleaner += duration

	processedReq := team.Request
	if processedReq == nil {
		s.l.ErrorCtx(ctx, "Request came nil after cleaning", logger.NewErrorField(errors.New("nil req")))
		return nil, errors.New("nil req")
	}

	return &dto.ProceedCleaningRequestOut{Req: processedReq}, nil
}

// GetAvailableTeams checks available teams in cleaning service.
// Returns available cleaning teams' IDs
func (s *Service) GetAvailableTeams(ctx context.Context) (*dto.GetAvailableTeamsOut, error) {
	availables := make([]uint64, 0, s.c.TeamsAmount)

	if len(s.teams) == 0 {
		s.l.Error("teams are not initialized")
		return nil, errors.New("empty cleaning service")
	}

	for _, team := range s.teams {
		if team.Status == entities.Available {
			availables = append(availables, team.Id)
		}
	}

	return &dto.GetAvailableTeamsOut{Teams: availables}, nil
}

// GetTeamsStats gets statistics of each team in cleaning service, while working to build statistic table for dispatcher.
// Returns all cleaning teams' statistics.
func (s *Service) GetTeamsStats(ctx context.Context) (*dto.GetTeamsStatsOut, error) {
	stats := s.teams
	if stats == nil {
		s.l.Error("teams array is nil")
		return nil, errors.New("no cleanning teams in service")
	}

	answer := make([]*dto.TeamStats, 0, len(stats))
	for _, stat := range stats {
		answer = append(answer, &dto.TeamStats{
			Id: stat.Id,
			Speed: uint32(stat.Speed),
			ProcessedRequests: stat.ProcessedRequests,
			TotalBusyTime: stat.TotalBusyTime,
		})
	}

	return &dto.GetTeamsStatsOut{Stats: answer}, nil
}

// initTeams - private func for initializing cleaner service's teams during the first connection to service
func initTeams(size uint64) []*entities.CleaningTeam {
	teams := make([]*entities.CleaningTeam, 0, size)

	for i := uint64(0); i < size; i++ {
		teams = append(teams, &entities.CleaningTeam{
			Id:        uint64(i),
			Request:   nil,
			Status:    entities.Available,
			Speed:     entities.Speed(rand.Intn(3) + 1),
			StartedAt: time.Time{},
		})
	}

	return teams
}
