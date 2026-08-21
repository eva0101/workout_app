package workout_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

func (s *WorkoutService) GetStatistic(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.Statistic, error) {
	statistic, err := s.workoutRepository.GetStatistic(
		ctx,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return statistic, nil
}
