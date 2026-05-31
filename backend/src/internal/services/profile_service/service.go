package profile_service

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type statisticsRepository interface {
	Get(ctx context.Context, in *statistics_repository.GetIn) ([]models.Statistics, error)
}

type userRepository interface {
	GetOne(ctx context.Context, in *user_repository.GetOneIn) (*models.User, error)
}

type profileRepository interface {
	Get(ctx context.Context, username string) (*models.Profile, bool)
	Save(ctx context.Context, in *models.Profile) error
}

type leaderboardService interface {
	GetUserPositions(ctx context.Context, userID models.ID) ([]models.LeaderboardPosition, error)
}

type Service struct {
	statisticsRepository statisticsRepository
	userRepository       userRepository
	profileRepository    profileRepository
	leaderboardService   leaderboardService
	useRedis             bool
}

func New(
	statisticsRepository statisticsRepository,
	userRepository userRepository,
	profileRepository profileRepository,
	leaderboardService leaderboardService,
	useRedis bool,
) *Service {
	return &Service{
		statisticsRepository: statisticsRepository,
		userRepository:       userRepository,
		profileRepository:    profileRepository,
		leaderboardService:   leaderboardService,
		useRedis:             useRedis,
	}
}
