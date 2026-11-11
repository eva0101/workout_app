package workout_repository_postgres

import (
	"context"
	"os"
	"testing"
	"time"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestGetStatisticsDashboardDataIsUserScopedAndUsesCompletedActualSets(t *testing.T) {
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

	ownerID, otherUserID := uuid.New(), uuid.New()
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.users (id, login, password_hash, created_at)
		VALUES
			($1, $2, 'test', NOW()),
			($3, $4, 'test', NOW())
	`, ownerID, "statistics-"+ownerID.String(), otherUserID, "statistics-"+otherUserID.String()); err != nil {
		t.Fatalf("create test users: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `
			DELETE FROM workoutapp.workout_set
			WHERE workout_id IN (SELECT id FROM workoutapp.workout WHERE user_id = ANY($1))
		`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.workout WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.program WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.users WHERE id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
	})

	createDay := func(userID uuid.UUID, name string) int {
		t.Helper()
		var programID, dayID int
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.program (user_id, name, started_at)
			VALUES ($1, $2, NOW())
			RETURNING id
		`, userID, "Statistics "+name).Scan(&programID); err != nil {
			t.Fatalf("create program: %v", err)
		}
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.training_days (program_id, day_number, name)
			VALUES ($1, 1, $2)
			RETURNING id
		`, programID, name).Scan(&dayID); err != nil {
			t.Fatalf("create training day: %v", err)
		}
		return dayID
	}
	createWorkout := func(userID uuid.UUID, dayID int, status string, completedAt *time.Time) int {
		t.Helper()
		var workoutID int
		beginAt := time.Now().UTC().Add(-30 * time.Minute)
		if completedAt != nil {
			beginAt = completedAt.Add(-30 * time.Minute)
		}
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.workout (
				user_id, training_day_id, training_day_name, status, begin_at, completed_at
			)
			VALUES ($1, $2, 'Statistics test', $3, $4, $5)
			RETURNING id
		`, userID, dayID, status, beginAt, completedAt).Scan(&workoutID); err != nil {
			t.Fatalf("create workout: %v", err)
		}
		return workoutID
	}
	createActualSet := func(workoutID, sourceID, setNumber, reps int, weight float64, group string) {
		t.Helper()
		var workoutExerciseID int
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.workout_exercise (
				workout_id, source_training_day_exercise_id, exercise_id,
				exercise_name, exercise_title, muscle_group_name, sort_order
			)
			VALUES ($1, $2, 1, 'statistics exercise', 'Упражнение статистики', $3, $4)
			RETURNING id
		`, workoutID, sourceID, group, sourceID).Scan(&workoutExerciseID); err != nil {
			t.Fatalf("create workout exercise snapshot: %v", err)
		}
		if _, err := pool.Exec(ctx, `
			INSERT INTO workoutapp.workout_set (
				workout_id, workout_exercise_id, set_number, reps_done, weight
			)
			VALUES ($1, $2, $3, $4, $5)
		`, workoutID, workoutExerciseID, setNumber, reps, weight); err != nil {
			t.Fatalf("create actual workout set: %v", err)
		}
	}

	ownerDayID := createDay(ownerID, "Owner")
	otherDayID := createDay(otherUserID, "Other")
	now := time.Now().UTC().Truncate(time.Microsecond)
	completedAt := now.Add(-time.Minute)
	ownerWorkoutID := createWorkout(ownerID, ownerDayID, "completed", &completedAt)
	createActualSet(ownerWorkoutID, 1, 1, 10, 20, "chest")
	createActualSet(ownerWorkoutID, 2, 1, 8, 30, "back")

	otherWorkoutID := createWorkout(otherUserID, otherDayID, "completed", &completedAt)
	createActualSet(otherWorkoutID, 3, 1, 100, 100, "biceps")

	activeWorkoutID := createWorkout(ownerID, ownerDayID, "in_progress", nil)
	createActualSet(activeWorkoutID, 4, 1, 100, 100, "chest")

	repository := &WorkoutRepository{pool: &core_pool_pgx.Pool{Pool: pool}}
	from := now.AddDate(0, 0, -1)
	workouts, groups, err := repository.GetStatisticsDashboardData(ctx, ownerID, &from, now)
	if err != nil {
		t.Fatalf("GetStatisticsDashboardData() error = %v", err)
	}
	if len(workouts) != 1 {
		t.Fatalf("got %d workouts, want only owner's completed workout", len(workouts))
	}
	workout := workouts[0]
	if workout.CompletedSets != 2 || workout.TotalWeight != 440 || workout.ActualReps != 18 ||
		workout.DurationSeconds != 1800 {
		t.Fatalf("unexpected workout totals: %+v", workout)
	}
	if len(groups) != 7 {
		t.Fatalf("got %d muscle groups, want 7", len(groups))
	}
	counts := make(map[string]int64, len(groups))
	for _, group := range groups {
		counts[group.Name] = group.CompletedSets
	}
	if counts["chest"] != 1 || counts["back"] != 1 || counts["biceps"] != 0 {
		t.Fatalf("muscle group totals are not user-scoped or include planned/active sets: %v", counts)
	}
}
