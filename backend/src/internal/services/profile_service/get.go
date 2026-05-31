package profile_service

import (
	"context"
	"fmt"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type UserNotFoundError struct {
	Username string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user with nickname '%s' not found", e.Username)
}

func IsUserNotFoundError(err error) bool {
	return models.IsCustomError[UserNotFoundError](err)
}

type GetIn struct {
	Username string
}

func (s *Service) Get(ctx context.Context, in *GetIn) (*models.Profile, error) {
	if s.useRedis {
		profile, ok := s.profileRepository.Get(ctx, in.Username)
		if ok {
			return profile, nil
		}
	}

	user, err := s.userRepository.GetOne(ctx, &user_repository.GetOneIn{
		Username: &in.Username,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, UserNotFoundError{Username: in.Username}
	}

	stats, err := s.statisticsRepository.Get(ctx, &statistics_repository.GetIn{
		UserID: user.ID,
	})
	if err != nil {
		return nil, err
	}

	leaderboardPositions, err := s.leaderboardService.GetUserPositions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	profile := models.NewProfile(user, stats, leaderboardPositions)

	if s.useRedis {
		if err = s.profileRepository.Save(ctx, profile); err != nil {
			return nil, err
		}
	}

	return profile, nil
}
