package workout_repository_postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestGetWorkoutSnapshotIntegration(t *testing.T) {
	databaseURL := os.Getenv("WORKOUT_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("WORKOUT_TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	defer pool.Close()

	var catalogExerciseID int
	if err := pool.QueryRow(ctx, `
		SELECT id FROM workoutapp.exercise ORDER BY id LIMIT 1
	`).Scan(&catalogExerciseID); err != nil {
		t.Fatalf("select catalog exercise: %v", err)
	}

	ownerID, otherUserID := uuid.New(), uuid.New()
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.users (id, login, password_hash, created_at)
		VALUES
			($1, $2, 'test', NOW()),
			($3, $4, 'test', NOW())
	`, ownerID, "get-workout-"+ownerID.String(), otherUserID, "get-workout-"+otherUserID.String()); err != nil {
		t.Fatalf("create test users: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.workout_set WHERE workout_id IN (
			SELECT id FROM workoutapp.workout WHERE user_id = ANY($1)
		)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.workout WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.program WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.users WHERE id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
	})

	var programID, dayID, sourceExerciseOne, sourceExerciseTwo int
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.program (user_id, name, started_at)
		VALUES ($1, 'Workout read test', NOW())
		RETURNING id
	`, ownerID).Scan(&programID); err != nil {
		t.Fatalf("create test program: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.training_days (program_id, day_number, name)
		VALUES ($1, 1, 'Snapshot day')
		RETURNING id
	`, programID).Scan(&dayID); err != nil {
		t.Fatalf("create test day: %v", err)
	}
	for _, destination := range []*int{&sourceExerciseOne, &sourceExerciseTwo} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.training_day_exercises (training_day_id, exercise_id, sets)
			VALUES ($1, $2, 2)
			RETURNING id
		`, dayID, catalogExerciseID).Scan(destination); err != nil {
			t.Fatalf("create source exercise: %v", err)
		}
	}

	var activeWorkoutID, completedWorkoutID int
	completedAt := time.Now().UTC().Truncate(time.Microsecond)
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.workout (user_id, training_day_id, status, begin_at)
		VALUES ($1, $2, 'in_progress', NOW())
		RETURNING id
	`, ownerID, dayID).Scan(&activeWorkoutID); err != nil {
		t.Fatalf("create active workout: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.workout (
			user_id, training_day_id, status, begin_at, completed_at
		)
		VALUES ($1, $2, 'completed', NOW() - INTERVAL '1 hour', $3)
		RETURNING id
	`, ownerID, dayID, completedAt).Scan(&completedWorkoutID); err != nil {
		t.Fatalf("create completed workout: %v", err)
	}

	var snapshotExerciseOne, snapshotExerciseTwo int
	for index, sourceID := range []int{sourceExerciseOne, sourceExerciseTwo} {
		var destination *int
		if index == 0 {
			destination = &snapshotExerciseOne
		} else {
			destination = &snapshotExerciseTwo
		}
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.workout_exercise (
				workout_id, source_training_day_exercise_id, exercise_id,
				exercise_name, exercise_title, sort_order
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, completedWorkoutID, sourceID, catalogExerciseID, "snapshot-name", "Snapshot title", index+1).Scan(destination); err != nil {
			t.Fatalf("create workout exercise snapshot: %v", err)
		}
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.workout_exercise_sets (
			workout_exercise_id, set_number, repetitions, weight
		)
		VALUES
			($1, 1, 10, 40),
			($1, 2, 8, 45)
	`, snapshotExerciseOne); err != nil {
		t.Fatalf("create planned snapshot sets: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.workout_set (
			workout_id, workout_exercise_id, set_number, reps_done, weight
		)
		VALUES
			($1, $2, 1, 9, 37.5),
			($1, $2, 2, 8, 40)
	`, completedWorkoutID, snapshotExerciseOne); err != nil {
		t.Fatalf("create actual workout sets: %v", err)
	}

	repository := &WorkoutRepository{pool: &core_pool_pgx.Pool{Pool: pool}}
	activeWorkout, err := repository.GetActiveWorkout(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetActiveWorkout() error = %v", err)
	}
	if activeWorkout.ID != activeWorkoutID || activeWorkout.Status != "in_progress" ||
		activeWorkout.StartedAt.IsZero() || activeWorkout.Exercises == nil || len(activeWorkout.Exercises) != 0 {
		t.Fatalf("unexpected active workout: %+v", activeWorkout)
	}

	if _, err := repository.GetActiveWorkout(ctx, otherUserID); !errors.Is(err, core_errors.ErrActiveWorkoutNotFound) {
		t.Fatalf("expected active-not-found for user without active workout, got %v", err)
	}
	if _, err := repository.GetWorkout(ctx, otherUserID, completedWorkoutID); !errors.Is(err, core_errors.ErrWorkoutNotFound) {
		t.Fatalf("expected workout-not-found for foreign workout, got %v", err)
	}

	workout, err := repository.GetWorkout(ctx, ownerID, completedWorkoutID)
	if err != nil {
		t.Fatalf("GetWorkout() error = %v", err)
	}
	if workout.Status != "completed" || workout.CompletedAt == nil || !workout.CompletedAt.Equal(completedAt) {
		t.Fatalf("completed workout metadata is incorrect: %+v", workout)
	}
	if len(workout.Exercises) != 2 {
		t.Fatalf("expected 2 snapshot exercises, got %d", len(workout.Exercises))
	}
	first := workout.Exercises[0]
	if first.ID != snapshotExerciseOne || first.Name != "snapshot-name" || first.Title != "Snapshot title" || first.Order != 1 {
		t.Fatalf("snapshot exercise metadata is incorrect: %+v", first)
	}
	if len(first.PlannedSets) != 2 || len(first.ActualSets) != 2 {
		t.Fatalf("planned or actual sets missing: %+v", first)
	}
	if *first.PlannedSets[1].Repetitions != 8 || *first.PlannedSets[1].Weight != 45 ||
		*first.ActualSets[0].RepsDone != 9 || *first.ActualSets[0].Weight != 37.5 {
		t.Fatalf("planned or actual set values are incorrect: %+v", first)
	}
	second := workout.Exercises[1]
	if second.ID != snapshotExerciseTwo || second.PlannedSets == nil || second.ActualSets == nil ||
		len(second.PlannedSets) != 0 || len(second.ActualSets) != 0 {
		t.Fatalf("empty set lists must be empty arrays: %+v", second)
	}

	var _ core_domain.WorkoutDetails = workout
}
