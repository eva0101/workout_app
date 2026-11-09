package workout_repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
	core_errors "workout_app/internal/core/errors"
	core_logger "workout_app/internal/core/logger"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func TestScheduledWorkoutAndExerciseCompletionIntegration(t *testing.T) {
	databaseURL := os.Getenv("WORKOUT_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("WORKOUT_TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	t.Cleanup(pool.Close)

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
	`, ownerID, "sw-"+ownerID.String(), otherUserID, "sw-"+otherUserID.String()); err != nil {
		t.Fatalf("create test users: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.scheduled_workouts WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.workout_set WHERE workout_id IN (
			SELECT id FROM workoutapp.workout WHERE user_id = ANY($1)
		)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.workout WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.program WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.users WHERE id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
	})

	ownerDayID := createWorkoutTestDay(t, ctx, pool, ownerID, "Scheduled day")
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.training_day_exercises (training_day_id, exercise_id, sets)
		VALUES ($1, $2, 1)
	`, ownerDayID, catalogExerciseID); err != nil {
		t.Fatalf("create planned exercise: %v", err)
	}
	var sourceExerciseID int
	if err := pool.QueryRow(ctx, `
		SELECT id FROM workoutapp.training_day_exercises WHERE training_day_id = $1
	`, ownerDayID).Scan(&sourceExerciseID); err != nil {
		t.Fatalf("read planned exercise: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.training_day_exercise_sets (
			training_day_exercise_id, set_number, repetitions, weight
		)
		VALUES ($1, 1, 8, 42.5)
	`, sourceExerciseID); err != nil {
		t.Fatalf("create planned set: %v", err)
	}

	otherDayID := createWorkoutTestDay(t, ctx, pool, otherUserID, "Rollback day")
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.training_day_exercises (training_day_id, exercise_id, sets)
		VALUES ($1, $2, 1)
	`, otherDayID, catalogExerciseID); err != nil {
		t.Fatalf("create rollback exercise: %v", err)
	}
	scheduleDate := time.Now().UTC().Format("2006-01-02")
	var ownerScheduleID, otherScheduleID int
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.scheduled_workouts (user_id, training_day_id, scheduled_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, ownerID, ownerDayID, scheduleDate).Scan(&ownerScheduleID); err != nil {
		t.Fatalf("schedule owner workout: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.scheduled_workouts (user_id, training_day_id, scheduled_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, otherUserID, otherDayID, scheduleDate).Scan(&otherScheduleID); err != nil {
		t.Fatalf("schedule foreign workout: %v", err)
	}

	repository := &WorkoutRepository{
		pool: &core_pool_pgx.Pool{Pool: pool},
		log:  &core_logger.Logger{Logger: zap.NewNop()},
	}
	workout, err := repository.StartScheduledTraining(ctx, ownerID, ownerScheduleID)
	if err != nil {
		t.Fatalf("StartScheduledTraining() error = %v", err)
	}
	if workout.Status != "in_progress" || workout.TrainingDayID != ownerDayID {
		t.Fatalf("unexpected started workout: %+v", workout)
	}

	var remainingSchedules int
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM workoutapp.scheduled_workouts WHERE id = $1
	`, ownerScheduleID).Scan(&remainingSchedules); err != nil || remainingSchedules != 0 {
		t.Fatalf("successful start must consume schedule: count=%d err=%v", remainingSchedules, err)
	}
	if _, err := repository.StartScheduledTraining(ctx, ownerID, otherScheduleID); !errors.Is(err, core_errors.ErrScheduledWorkoutNotFound) {
		t.Fatalf("foreign schedule should be hidden, got %v", err)
	}

	var secondOwnerScheduleID int
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.scheduled_workouts (user_id, training_day_id, scheduled_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, ownerID, ownerDayID, time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")).Scan(&secondOwnerScheduleID); err != nil {
		t.Fatalf("create schedule for active workout check: %v", err)
	}
	if _, err := repository.StartScheduledTraining(ctx, ownerID, secondOwnerScheduleID); !errors.Is(err, core_errors.ErrWorkoutAlreadyInProgress) {
		t.Fatalf("expected active workout conflict, got %v", err)
	}
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM workoutapp.scheduled_workouts WHERE id = $1
	`, secondOwnerScheduleID).Scan(&remainingSchedules); err != nil || remainingSchedules != 1 {
		t.Fatalf("failed start must preserve schedule: count=%d err=%v", remainingSchedules, err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM workoutapp.scheduled_workouts WHERE id = $1`, secondOwnerScheduleID); err != nil {
		t.Fatalf("delete schedule for missing-record check: %v", err)
	}
	if _, err := repository.StartScheduledTraining(ctx, ownerID, secondOwnerScheduleID); !errors.Is(err, core_errors.ErrScheduledWorkoutNotFound) {
		t.Fatalf("deleted schedule should return not found, got %v", err)
	}

	triggerSuffix := uuid.NewString()
	triggerName := "test_fail_snapshot_" + triggerSuffix[:8]
	functionName := triggerName + "_fn"
	if _, err := pool.Exec(ctx, fmt.Sprintf(`
		CREATE FUNCTION workoutapp.%s() RETURNS trigger AS $$
		BEGIN
			RAISE EXCEPTION 'injected snapshot failure';
		END;
		$$ LANGUAGE plpgsql;
	`, functionName)); err != nil {
		t.Fatalf("create snapshot failure function: %v", err)
	}
	if _, err := pool.Exec(ctx, fmt.Sprintf(`
		CREATE TRIGGER %s
		BEFORE INSERT ON workoutapp.workout_exercise
		FOR EACH ROW EXECUTE FUNCTION workoutapp.%s();
	`, triggerName, functionName)); err != nil {
		t.Fatalf("create snapshot failure trigger: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON workoutapp.workout_exercise`, triggerName))
		_, _ = pool.Exec(ctx, fmt.Sprintf(`DROP FUNCTION IF EXISTS workoutapp.%s()`, functionName))
	})
	if _, err := repository.StartScheduledTraining(ctx, otherUserID, otherScheduleID); err == nil {
		t.Fatal("expected snapshot creation failure")
	}
	if _, err := pool.Exec(ctx, fmt.Sprintf(`DROP TRIGGER %s ON workoutapp.workout_exercise`, triggerName)); err != nil {
		t.Fatalf("drop snapshot failure trigger: %v", err)
	}
	if _, err := pool.Exec(ctx, fmt.Sprintf(`DROP FUNCTION workoutapp.%s()`, functionName)); err != nil {
		t.Fatalf("drop snapshot failure function: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM workoutapp.scheduled_workouts WHERE id = $1
	`, otherScheduleID).Scan(&remainingSchedules); err != nil || remainingSchedules != 1 {
		t.Fatalf("snapshot failure must preserve schedule: count=%d err=%v", remainingSchedules, err)
	}
	var otherWorkoutCount int
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM workoutapp.workout WHERE user_id = $1
	`, otherUserID).Scan(&otherWorkoutCount); err != nil || otherWorkoutCount != 0 {
		t.Fatalf("snapshot failure must rollback workout: count=%d err=%v", otherWorkoutCount, err)
	}

	var workoutExerciseID int
	if err := pool.QueryRow(ctx, `
		SELECT id FROM workoutapp.workout_exercise WHERE workout_id = $1
	`, workout.ID).Scan(&workoutExerciseID); err != nil {
		t.Fatalf("read snapshot exercise: %v", err)
	}
	completion, err := repository.CompleteExercise(ctx, ownerID, workoutExerciseID)
	if err != nil {
		t.Fatalf("CompleteExercise() error = %v", err)
	}
	repeatedCompletion, err := repository.CompleteExercise(ctx, ownerID, workoutExerciseID)
	if err != nil {
		t.Fatalf("repeated CompleteExercise() error = %v", err)
	}
	if completion.ID != workoutExerciseID || !completion.CompletedAt.Equal(repeatedCompletion.CompletedAt) {
		t.Fatalf("completion must be stable and idempotent: first=%+v repeated=%+v", completion, repeatedCompletion)
	}
	if _, err := repository.CompleteExercise(ctx, otherUserID, workoutExerciseID); !errors.Is(err, core_errors.ErrWorkoutExerciseNotFound) {
		t.Fatalf("foreign exercise should return not found, got %v", err)
	}

	var plannedSetCount, actualSetCount int
	if err := pool.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM workoutapp.workout_exercise_sets WHERE workout_exercise_id = $1),
			(SELECT COUNT(*) FROM workoutapp.workout_set WHERE workout_exercise_id = $1)
	`, workoutExerciseID).Scan(&plannedSetCount, &actualSetCount); err != nil {
		t.Fatalf("count snapshot and actual sets: %v", err)
	}
	if plannedSetCount != 1 || actualSetCount != 0 {
		t.Fatalf("expected one planned and zero actual sets, got planned=%d actual=%d", plannedSetCount, actualSetCount)
	}

	if _, err := pool.Exec(ctx, `UPDATE workoutapp.training_days SET name = 'Renamed day' WHERE id = $1`, ownerDayID); err != nil {
		t.Fatalf("rename source day: %v", err)
	}
	details, err := repository.GetWorkout(ctx, ownerID, workout.ID)
	if err != nil {
		t.Fatalf("GetWorkout() error = %v", err)
	}
	if details.TrainingDayName == nil || *details.TrainingDayName != "Scheduled day" ||
		len(details.Exercises) != 1 || details.Exercises[0].CompletedAt == nil {
		t.Fatalf("workout snapshot did not preserve name and completion: %+v", details)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM workoutapp.training_days WHERE id = $1`, ownerDayID); err != nil {
		t.Fatalf("delete source day: %v", err)
	}
	details, err = repository.GetWorkout(ctx, ownerID, workout.ID)
	if err != nil || details.TrainingDayName == nil || *details.TrainingDayName != "Scheduled day" {
		t.Fatalf("source day deletion changed workout name: details=%+v err=%v", details, err)
	}

	if _, err := pool.Exec(ctx, `
		UPDATE workoutapp.workout SET status = 'completed', completed_at = NOW()
		WHERE id = $1
	`, workout.ID); err != nil {
		t.Fatalf("complete workout for exercise status check: %v", err)
	}
	if _, err := repository.CompleteExercise(ctx, ownerID, workoutExerciseID); !errors.Is(err, core_errors.ErrWorkoutExerciseAlreadyCompleted) {
		t.Fatalf("exercise in completed workout should return conflict, got %v", err)
	}

	directDayID := createWorkoutTestDay(t, ctx, pool, ownerID, "Direct start day")
	directWorkout, err := repository.StartTraining(ctx, ownerID, directDayID)
	if err != nil {
		t.Fatalf("existing direct StartTraining() error = %v", err)
	}
	if directWorkout.TrainingDayID != directDayID || directWorkout.Status != "in_progress" {
		t.Fatalf("unexpected workout from direct start: %+v", directWorkout)
	}
}

func createWorkoutTestDay(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	userID uuid.UUID,
	name string,
) int {
	t.Helper()
	var programID, dayID int
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.program (user_id, name, started_at)
		VALUES ($1, 'Workout start test', NOW())
		RETURNING id
	`, userID).Scan(&programID); err != nil {
		t.Fatalf("create test program: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		INSERT INTO workoutapp.training_days (program_id, day_number, name)
		VALUES ($1, 1, $2)
		RETURNING id
	`, programID, name).Scan(&dayID); err != nil {
		t.Fatalf("create test training day: %v", err)
	}
	return dayID
}
