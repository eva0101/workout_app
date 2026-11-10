package workout_service

import (
	"context"
	"errors"
	"testing"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
)

type workoutHistoryRepositoryStub struct {
	WorkoutRepository
	history         []core_domain.WorkoutHistoryItem
	err             error
	requestedUser   uuid.UUID
	requestedLimit  int
	requestedOffset int
}

func (r *workoutHistoryRepositoryStub) GetWorkoutHistory(
	_ context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]core_domain.WorkoutHistoryItem, error) {
	r.requestedUser = userID
	r.requestedLimit = limit
	r.requestedOffset = offset
	return r.history, r.err
}

func TestGetWorkoutHistoryPassesUserAndPagination(t *testing.T) {
	userID := uuid.New()
	fatigueScore := 8
	expected := []core_domain.WorkoutHistoryItem{{ID: 17, FatigueScore: &fatigueScore}}
	repository := &workoutHistoryRepositoryStub{history: expected}
	service := &WorkoutService{workoutRepository: repository}

	got, err := service.GetWorkoutHistory(context.Background(), userID, 15, 30)
	if err != nil {
		t.Fatalf("GetWorkoutHistory() error = %v", err)
	}
	if len(got) != 1 || got[0].ID != expected[0].ID ||
		got[0].FatigueScore == nil || *got[0].FatigueScore != fatigueScore {
		t.Fatalf("GetWorkoutHistory() = %+v, want %+v", got, expected)
	}
	if repository.requestedUser != userID || repository.requestedLimit != 15 || repository.requestedOffset != 30 {
		t.Fatalf("repository received user=%s limit=%d offset=%d",
			repository.requestedUser, repository.requestedLimit, repository.requestedOffset)
	}
}

func TestGetWorkoutHistoryRejectsInvalidPagination(t *testing.T) {
	repository := &workoutHistoryRepositoryStub{}
	service := &WorkoutService{workoutRepository: repository}

	for _, pagination := range [][2]int{
		{0, 0},
		{core_domain.WorkoutHistoryMaxLimit + 1, 0},
		{10, -1},
	} {
		if _, err := service.GetWorkoutHistory(
			context.Background(),
			uuid.New(),
			pagination[0],
			pagination[1],
		); !errors.Is(err, core_errors.ErrInvalidArgument) {
			t.Errorf("GetWorkoutHistory(limit=%d, offset=%d) error = %v, want invalid argument",
				pagination[0], pagination[1], err)
		}
	}
}
