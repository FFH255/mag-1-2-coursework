package statistics_service

import (
	"context"

	"github.com/AlekSi/pointer"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

/*
DeleteAllForUser - удаляет всю статистику пользователя
1. Получается пользователя
2. Удаляет всю его статистику
*/
func (s *Service) DeleteAllForUser(ctx context.Context, userID uint64) error {
	user, err := s.userRepository.GetOne(ctx, &user_repository.GetOneIn{
		UserID: pointer.To(models.ID(userID)),
	})
	if err != nil {
		return err
	}
	if user == nil {
		return models.UserNotFoundError{UserID: models.ID(userID)}
	}

	if err := s.statisticsRepository.DeleteAllByUserID(ctx, user.ID); err != nil {
		return err
	}

	return nil
}
