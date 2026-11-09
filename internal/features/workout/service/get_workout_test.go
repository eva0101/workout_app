package workout_service

import (
	"context"
	"errors"
	"testing"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
)

type workoutReadRepositoryStub struct {
	WorkoutRepository
	activeWorkout core_domain.WorkoutDetails
	workout       core_domain.WorkoutDetails
	activeErr     error
	workoutErr    error
	activeUserID  uuid.UUID
	requestedUser uuid.UUID
	requestedID   int
}

func (r *workoutReadRepositoryStub) GetActiveWorkout(
	_ context.Context,
	userID uuid.UUID,
) (core_domain.WorkoutDetails, error) {
	r.activeUserID = userID
	return r.activeWorkout, r.activeErr
}

func (r *workoutReadRepositoryStub) GetWorkout(
	_ context.Context,
	userID uuid.UUID,
	workoutID int,
) (core_domain.WorkoutDetails, error) {
	r.requestedUser = userID
	r.requestedID = workoutID
	return r.workout, r.workoutErr
}

func TestGetActiveWorkoutReturnsCurrentUsersWorkout(t *testing.T) {
	userID := uuid.New()
	expected := core_domain.WorkoutDetails{ID: 17, Status: "in_progress"}
	repository := &workoutReadRepositoryStub{activeWorkout: expected}
	service := &WorkoutService{workoutRepository: repository}

	got, err := service.GetActiveWorkout(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetActiveWorkout() error = %v", err)
	}
	if got.ID != expected.ID || repository.activeUserID != userID {
		t.Fatalf("GetActiveWorkout() = %+v, user ID passed = %s", got, repository.activeUserID)
	}
}

func TestGetActiveWorkoutReturnsNotFound(t *testing.T) {
	repository := &workoutReadRepositoryStub{activeErr: errors.New("active workout not found")}
	service := &WorkoutService{workoutRepository: repository}

	if _, err := service.GetActiveWorkout(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected missing active workout error")
	}
}

func TestGetWorkoutReturnsCompletedWorkoutWithSnapshotSets(t *testing.T) {
	repetitions, weight, repsDone, actualWeight := 8, 40.0, 7, 37.5
	expected := core_domain.WorkoutDetails{
		ID:     23,
		Status: "completed",
		Exercises: []core_domain.WorkoutExerciseSnapshot{
			{
				ID: 101, ExerciseID: 2, Name: "row", Title: "Row", Order: 1,
				PlannedSets: []core_domain.WorkoutPlannedSet{
					{ID: 201, SetNumber: 1, Repetitions: &repetitions, Weight: &weight},
				},
				ActualSets: []core_domain.WorkoutActualSet{
					{ID: 301, SetNumber: 1, RepsDone: &repsDone, Weight: &actualWeight},
				},
			},
			{
				ID: 102, ExerciseID: 3, Name: "press", Title: "Press", Order: 2,
				PlannedSets: []core_domain.WorkoutPlannedSet{},
				ActualSets:  []core_domain.WorkoutActualSet{},
			},
		},
	}
	userID := uuid.New()
	repository := &workoutReadRepositoryStub{workout: expected}
	service := &WorkoutService{workoutRepository: repository}

	got, err := service.GetWorkout(context.Background(), userID, expected.ID)
	if err != nil {
		t.Fatalf("GetWorkout() error = %v", err)
	}
	if got.Status != "completed" || len(got.Exercises) != 2 {
		t.Fatalf("GetWorkout() returned unexpected workout: %+v", got)
	}
	if got.Exercises[0].ID != 101 || len(got.Exercises[0].PlannedSets) != 1 || len(got.Exercises[0].ActualSets) != 1 {
		t.Fatalf("first snapshot exercise or its sets are missing: %+v", got.Exercises[0])
	}
	if got.Exercises[1].PlannedSets == nil || got.Exercises[1].ActualSets == nil {
		t.Fatalf("empty set lists must be non-nil arrays: %+v", got.Exercises[1])
	}
	if repository.requestedUser != userID || repository.requestedID != expected.ID {
		t.Fatalf("repository received user=%s workout=%d", repository.requestedUser, repository.requestedID)
	}
}

func TestGetWorkoutPropagatesNotFoundForForeignWorkout(t *testing.T) {
	repository := &workoutReadRepositoryStub{workoutErr: errors.New("workout not found")}
	service := &WorkoutService{workoutRepository: repository}

	if _, err := service.GetWorkout(context.Background(), uuid.New(), 99); err == nil {
		t.Fatal("expected not-found error for a workout not owned by the user")
	}
}
