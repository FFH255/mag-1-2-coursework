package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type UserNotFoundError struct {
	id models.ID
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user with id %d not found", e.id)
}

func IsUserNotFoundError(err error) bool {
	return models.IsCustomError[UserNotFoundError](err)
}

type FroadError struct{}

func (e FroadError) Error() string {
	return fmt.Sprintf("froad")
}

func IsFroadError(err error) bool {
	return models.IsCustomError[FroadError](err)
}

type AlreadyHandledError struct{}

func (e AlreadyHandledError) Error() string {
	return fmt.Sprintf("already handled")
}

func IsAlreadyHandledError(err error) bool {
	return models.IsCustomError[AlreadyHandledError](err)
}

type statisticsRepository interface {
	Save(ctx context.Context, stats *models.Statistics) (*models.Statistics, error)
	Exists(ctx context.Context, idempotencyKey string) bool
	Get(ctx context.Context, in *statistics_repository.GetIn) ([]models.Statistics, error)
	DeleteAllByUserID(ctx context.Context, userID models.ID) error
}

type userRepository interface {
	GetOne(ctx context.Context, in *user_repository.GetOneIn) (*models.User, error)
}

type languageRepository interface {
	Get() models.Languages
}

type antifroad interface {
	IsFroad(payload models.SignPayload, sign string) bool
}

type pbService interface {
	IsPB(ctx context.Context, stats *models.Statistics) (isPB bool, wpmShift float64)
}

type leaderboardService interface {
	UpdateRank(ctx context.Context, userID models.ID, language models.Language, mode models.Mode, submode models.Submode, wpm float64, accuracy float64, playedAt time.Time) (*models.RankChange, error)
}

type Service struct {
	statisticsRepository statisticsRepository
	userRepository       userRepository
	languageRepository   languageRepository
	antifroad            antifroad
	pbService            pbService
	leaderboardService   leaderboardService
}

func New(
	statisticsRepository statisticsRepository,
	userRepository userRepository,
	languageRepository languageRepository,
	antifroad antifroad,
	pbService pbService,
	leaderboardService leaderboardService,
) *Service {
	return &Service{
		statisticsRepository: statisticsRepository,
		userRepository:       userRepository,
		languageRepository:   languageRepository,
		antifroad:            antifroad,
		pbService:            pbService,
		leaderboardService:   leaderboardService,
	}
}
