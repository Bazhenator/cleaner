package logic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"

	"github.com/Bazhenator/cleaner/configs"
	"github.com/Bazhenator/cleaner/internal/entities"
	"github.com/Bazhenator/cleaner/internal/logic/dto"
	"github.com/Bazhenator/tools/src/logger"
)

const (
	CleaningServiceSize = 10
)

type Service struct {
	c  *configs.Config
	l  *logger.Logger
	mu sync.Mutex

	teams []*entities.CleaningTeam
}

func NewService(c *configs.Config, l *logger.Logger) *Service {
	// Cleaning teams' initializing
	teams := initTeams()

	return &Service{
		c: c,
		l: l,

		teams: teams,
	}
}

// ProceedCleaningRequest proceeds request from user, assigns it to cleaning team and processes it.
// Returns cleaning duration
func (s *Service) ProceedCleaningRequest(ctx context.Context, in *dto.ProceedCleaningRequestIn) *dto.ProceedCleaningRequestOut {
	team := s.teams[in.TeamId]
	team.AssignRequest(in.Request)

	duration := team.GetCleaningTime()
	team.FinishedAt = team.StartedAt.Add(duration)

	go func(team *entities.CleaningTeam, duration time.Duration) {
		time.Sleep(duration)
		s.mu.Lock()
		defer s.mu.Unlock()

		team.CompleteCleaning()

		s.l.Info(fmt.Sprintf("Team %d completed cleaning.", team.Id))
	}(team, duration)

	return &dto.ProceedCleaningRequestOut{Duration: duration.String()}
}

// GetAvailableTeams checks available teams in cleaning service.
// Returns available cleaning teams' IDs
func (s *Service) GetAvailableTeams(ctx context.Context) *dto.GetAvailableTeamsOut {
	availables := make([]uint64, 0, CleaningServiceSize)

	for _, team := range s.teams {
		if team.Status == entities.Available {
			availables = append(availables, team.Id)
		}
	}

	return &dto.GetAvailableTeamsOut{Teams: availables}
}

// initTeams - private func for initializing cleaner service's teams during the first connection to service
func initTeams() []*entities.CleaningTeam {
	teams := make([]*entities.CleaningTeam, 0, CleaningServiceSize)

	for i := 0; i < CleaningServiceSize; i++ {
		teams = append(teams, &entities.CleaningTeam{
			Id:         uint64(i),
			Request:    nil,
			Status:     entities.Available,
			Speed:      entities.Speed(rand.Intn(3) + 1),
			StartedAt:  time.Time{},
			FinishedAt: time.Time{},
		})
	}

	return teams
}
