package statistics_service

import (
	"context"
	"time"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type GetByUserIn struct {
	UserID   models.ID
	DateFrom *time.Time
}

func (s *Service) GetByUser(ctx context.Context, in *GetByUserIn) ([]models.Statistics, error) {
	user, err := s.userRepository.GetOne(ctx, &user_repository.GetOneIn{
		UserID: &in.UserID,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, UserNotFoundError{id: in.UserID}
	}

	stats, err := s.statisticsRepository.Get(ctx, &statistics_repository.GetIn{
		UserID:               user.ID,
		DateFromNotInclusive: in.DateFrom,
		Sorting: &statistics_repository.Sorting{
			SortBy:    statistics_repository.SortByPlayedAt,
			Ascending: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return stats, nil
}
