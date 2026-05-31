package pb_service

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/pb_cache"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
)

type pbCache interface {
	Get(ctx context.Context, key pb_cache.Key) (models.WPM, bool)
	Save(ctx context.Context, key pb_cache.Key, pb models.WPM) error
}

type statisticsRepository interface {
	Get(ctx context.Context, in *statistics_repository.GetIn) ([]models.Statistics, error)
}

type Service struct {
	pbCache              pbCache
	statisticsRepository statisticsRepository
	logger               internal.Logger
}

func (s *Service) IsPB(ctx context.Context, stats *models.Statistics) (isPB bool, wpmShift float64) {
	if stats == nil {
		return false, 0
	}

	pb := s.getPB(ctx, stats.UserID, stats.Language, stats.Mode, stats.SubMode)
	if stats.WPM <= pb {
		return false, 0
	}

	err := s.pbCache.Save(ctx, pb_cache.Key{
		UserID:   stats.UserID,
		Language: stats.Language,
		Mode:     stats.Mode,
		Submode:  stats.SubMode,
	}, stats.WPM)
	if err != nil {
		s.logger.Warning(s.logger.WithError(ctx, err))
		return false, 0
	}

	wpmShift = float64(stats.WPM - pb)

	return true, wpmShift
}

func (s *Service) getPB(
	ctx context.Context,
	userID models.ID,
	language models.Language, mode models.Mode, submode models.Submode,
) models.WPM {
	pb, exists := s.pbCache.Get(ctx, pb_cache.Key{
		UserID:   userID,
		Language: language,
		Mode:     mode,
		Submode:  submode,
	})
	if !exists {
		stats, err := s.statisticsRepository.Get(ctx, &statistics_repository.GetIn{
			UserID:   userID,
			Language: &language,
			Mode:     &mode,
			Submode:  &submode,
			Sorting: &statistics_repository.Sorting{
				SortBy:    statistics_repository.SortByWPM,
				Ascending: false, // По убыванию. PB будет на нулевом индексе
			},
		})
		if err != nil {
			s.logger.Warning(s.logger.WithError(ctx, err))
			return 0
		}
		if len(stats) == 0 {
			return 0
		}

		return stats[0].WPM
	}

	return pb
}

func New(pbCache pbCache, statisticsRepository statisticsRepository, logger internal.Logger) *Service {
	return &Service{
		pbCache:              pbCache,
		statisticsRepository: statisticsRepository,
		logger:               logger,
	}
}
