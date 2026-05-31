package statistics_service

import (
	"context"
	"time"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type SaveIn struct {
	UserID                     models.ID
	WPM                        float64
	CPM                        float64
	Accuracy                   float64
	Duration                   time.Duration
	Language                   string
	Mode                       string
	SubMode                    string
	IsPunctuation              bool
	UncompletedTestsDurationMs time.Duration
	UncompletedTestsCount      uint64
	UID                        string
	Sign                       string
	CreatedAt                  time.Time
	StartedAt                  time.Time
	FinishedAt                 time.Time
}

type SaveOut struct {
	StatisticsID uint
	IsPB         bool
	WPMShift     float64
	RankChange   *models.RankChange
}

func (s *Service) Save(ctx context.Context, in *SaveIn) (*SaveOut, error) {
	var (
		user  *models.User
		stats *models.Statistics
		err   error
	)

	// 1. Поиск пользователя
	user, err = s.userRepository.GetOne(ctx, &user_repository.GetOneIn{
		UserID: &in.UserID,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, UserNotFoundError{id: in.UserID}
	}

	// 2. Идемпотентность статистики
	if exists := s.statisticsRepository.Exists(ctx, in.UID); exists {
		return nil, AlreadyHandledError{}
	}

	// 3. Проверка запроса на фрод
	if s.antifroad.IsFroad(newAntifroadPayload(in), in.Sign) {
		return nil, FroadError{}
	}

	// 4. Получение языков
	languages := s.languageRepository.Get()

	stats, err = models.NewStatistics(models.StatisticsOptions{
		UserID:                        user.ID,
		WPM:                           in.WPM,
		CPM:                           in.CPM,
		Accuracy:                      in.Accuracy,
		Duration:                      in.Duration,
		Language:                      in.Language,
		Mode:                          in.Mode,
		Submode:                       in.SubMode,
		IsPunctuation:                 in.IsPunctuation,
		UncompletedTestsCount:         in.UncompletedTestsCount,
		UncompletedTestsTotalDuration: in.UncompletedTestsDurationMs,
		IdempotencyKey:                in.UID,
	}, languages)
	if err != nil {
		return nil, err
	}

	isPB, wpmShift := s.pbService.IsPB(ctx, stats) // Должно быть перед сохранением. Иначе ПБ никогда не будет

	stats, err = s.statisticsRepository.Save(ctx, stats)
	if err != nil {
		return nil, err
	}

	var rankChange *models.RankChange
	if isPB {
		rankChange, err = s.leaderboardService.UpdateRank(
			ctx,
			stats.UserID,
			stats.Language,
			stats.Mode,
			stats.SubMode,
			float64(stats.WPM),
			float64(stats.Accuracy),
			stats.PlayedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return s.newSaveOut(stats, isPB, wpmShift, rankChange), nil
}

func (s *Service) newSaveOut(stats *models.Statistics, isPB bool, wpmShift float64, rankChange *models.RankChange) *SaveOut {
	if stats == nil {
		return nil
	}

	return &SaveOut{
		StatisticsID: uint(stats.ID),
		IsPB:         isPB,
		WPMShift:     wpmShift,
		RankChange:   rankChange,
	}
}

func newAntifroadPayload(in *SaveIn) models.SignPayload {
	return models.SignPayload{
		UID:                      in.UID,
		WPM:                      in.WPM,
		CPM:                      in.CPM,
		Accuracy:                 in.Accuracy,
		Duration:                 in.Duration,
		Language:                 in.Language,
		Mode:                     in.Mode,
		SubMode:                  in.SubMode,
		IsPunctuation:            in.IsPunctuation,
		UncompletedTestsDuration: in.UncompletedTestsDurationMs,
		UncompletedTestsCount:    in.UncompletedTestsCount,
		CreatedAt:                in.CreatedAt,
		StartedAt:                in.StartedAt,
		FinishedAt:               in.FinishedAt,
	}
}
