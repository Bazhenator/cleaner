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
	CleaningServiceSize = 20
)

type Service struct {
	c  *configs.Config
	l  *logger.Logger
	mu sync.Mutex

	teams []*entities.CleaningTeam
	stats      map[uint64]*entities.TeamStats
}

func NewService(c *configs.Config, l *logger.Logger) *Service {
	// Cleaning teams' initializing
	teams := initTeams()
	stats := make(map[uint64]*entities.TeamStats, len(teams))

	for _, team := range teams {
		stats[team.Id] = &entities.TeamStats{}
	}

	return &Service{
		c: c,
		l: l,

		stats: stats,
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

		// Update statistics
		stats := s.stats[team.Id]
		stats.ProcessedRequests++
		stats.TotalCleaningTime += duration

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

// GenerateReport generates a report about cleaning teams' workloads
func (s *Service) GenerateReport() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("Cleaning Teams Workload Report:")
	fmt.Println("================================")
	for _, team := range s.teams {
		stats := s.stats[team.Id]
		fmt.Printf("Team ID: %d | Speed: %d | Processed Requests: %d | Total Cleaning Time: %s\n",
			team.Id, team.Speed, stats.ProcessedRequests, stats.TotalCleaningTime)
	}
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
