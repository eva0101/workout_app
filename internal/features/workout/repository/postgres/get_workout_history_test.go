package workout_repository_postgres

import (
	"context"
	"os"
	"testing"
	"time"
	core_logger "workout_app/internal/core/logger"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func TestGetWorkoutHistoryIntegration(t *testing.T) {
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
	`, ownerID, "history-"+ownerID.String(), otherUserID, "history-"+otherUserID.String()); err != nil {
		t.Fatalf("create test users: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.workout WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.program WHERE user_id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
		_, _ = pool.Exec(ctx, `DELETE FROM workoutapp.users WHERE id = ANY($1)`, []uuid.UUID{ownerID, otherUserID})
	})

	type actualSet struct {
		reps   *int
		weight *float64
	}
	type scenario struct {
		name         string
		plannedReps  []*int
		actualSets   []actualSet
		wantSets     int64
		wantWeight   float64
		wantActual   int64
		wantPlanned  int64
		wantPercent  *int
		fatigueScore *int
	}
	intPointer := func(value int) *int { return &value }
	floatPointer := func(value float64) *float64 { return &value }
	eighty, oneHundred, zero, goodFeeling := 80, 100, 0, 7
	scenarios := []scenario{
		{
			name:        "ordinary",
			plannedReps: []*int{intPointer(10), intPointer(5), nil},
			actualSets:  []actualSet{{reps: intPointer(8), weight: floatPointer(20)}, {reps: intPointer(4)}},
			wantSets:    2, wantWeight: 160, wantActual: 12, wantPlanned: 15, wantPercent: &eighty,
			fatigueScore: &goodFeeling,
		},
		{
			name:        "over-completion",
			plannedReps: []*int{intPointer(10)},
			actualSets:  []actualSet{{reps: intPointer(12), weight: floatPointer(5)}},
			wantSets:    1, wantWeight: 60, wantActual: 12, wantPlanned: 10, wantPercent: &oneHundred,
		},
		{
			name:        "no planned repetitions",
			plannedReps: []*int{nil},
			actualSets:  []actualSet{{reps: intPointer(6), weight: floatPointer(3)}},
			wantSets:    1, wantWeight: 18, wantActual: 6, wantPlanned: 0,
		},
		{
			name:        "no actual sets",
			plannedReps: []*int{intPointer(9)},
			wantSets:    0, wantWeight: 0, wantActual: 0, wantPlanned: 9, wantPercent: &zero,
		},
	}

	base := time.Date(2026, 10, 15, 8, 0, 0, 0, time.UTC)
	wantByID := make(map[int]scenario, len(scenarios))
	var snapshotIDs []int
	for index, item := range scenarios {
		dayID := createWorkoutTestDay(t, ctx, pool, ownerID, "Current "+item.name)
		startedAt := base.Add(time.Duration(index) * time.Hour)
		completedAt := startedAt.Add(100 * time.Second)
		var workoutID, snapshotExerciseID int
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.workout (
				user_id, training_day_id, training_day_name, status, begin_at, completed_at, fatigue_score
			)
			VALUES ($1, $2, $3, 'completed', $4, $5, $6)
			RETURNING id
		`, ownerID, dayID, "Snapshot "+item.name, startedAt, completedAt, item.fatigueScore).Scan(&workoutID); err != nil {
			t.Fatalf("create %s workout: %v", item.name, err)
		}
		if err := pool.QueryRow(ctx, `
			INSERT INTO workoutapp.workout_exercise (
				workout_id, source_training_day_exercise_id, exercise_id,
				exercise_name, exercise_title, muscle_group_name, sort_order
			)
			VALUES ($1, $2, $3, $4, $4, 'chest', 1)
			RETURNING id
		`, workoutID, index+1, catalogExerciseID, item.name).Scan(&snapshotExerciseID); err != nil {
			t.Fatalf("create snapshot exercise for %s: %v", item.name, err)
		}
		for setIndex, reps := range item.plannedReps {
			if _, err := pool.Exec(ctx, `
				INSERT INTO workoutapp.workout_exercise_sets (
					workout_exercise_id, set_number, repetitions
				)
				VALUES ($1, $2, $3)
			`, snapshotExerciseID, setIndex+1, reps); err != nil {
				t.Fatalf("create planned set for %s: %v", item.name, err)
			}
		}
		for setIndex, set := range item.actualSets {
			if _, err := pool.Exec(ctx, `
				INSERT INTO workoutapp.workout_set (
					workout_id, workout_exercise_id, set_number, reps_done, weight
				)
				VALUES ($1, $2, $3, $4, $5)
			`, workoutID, snapshotExerciseID, setIndex+1, set.reps, set.weight); err != nil {
				t.Fatalf("create actual set for %s: %v", item.name, err)
			}
		}
		wantByID[workoutID] = item
		snapshotIDs = append(snapshotIDs, workoutID)
		if _, err := pool.Exec(ctx, `
			UPDATE workoutapp.training_days SET name = $2 WHERE id = $1
		`, dayID, "Changed "+item.name); err != nil {
			t.Fatalf("change current day name for %s: %v", item.name, err)
		}
	}

	foreignDayID := createWorkoutTestDay(t, ctx, pool, otherUserID, "Foreign day")
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.workout (
			user_id, training_day_id, training_day_name, status, begin_at, completed_at, fatigue_score
		)
		VALUES ($1, $2, 'Foreign snapshot', 'completed', $3, $4, 5)
	`, otherUserID, foreignDayID, base, base.Add(time.Minute)); err != nil {
		t.Fatalf("create foreign workout: %v", err)
	}
	activeDayID := createWorkoutTestDay(t, ctx, pool, ownerID, "Active day")
	if _, err := pool.Exec(ctx, `
		INSERT INTO workoutapp.workout (user_id, training_day_id, status, begin_at)
		VALUES ($1, $2, 'in_progress', $3)
	`, ownerID, activeDayID, base); err != nil {
		t.Fatalf("create in-progress workout: %v", err)
	}

	repository := &WorkoutRepository{
		pool: &core_pool_pgx.Pool{Pool: pool},
		log:  &core_logger.Logger{Logger: zap.NewNop()},
	}
	history, err := repository.GetWorkoutHistory(ctx, ownerID, 20, 0)
	if err != nil {
		t.Fatalf("GetWorkoutHistory() error = %v", err)
	}
	if len(history) != len(scenarios) {
		t.Fatalf("GetWorkoutHistory() returned %d items, want %d", len(history), len(scenarios))
	}
	for index, item := range history {
		if index > 0 && item.CompletedAt.After(history[index-1].CompletedAt) {
			t.Fatalf("history is not sorted newest first: %+v", history)
		}
		want, ok := wantByID[item.ID]
		if !ok {
			t.Fatalf("history contains foreign or unfinished workout id %d", item.ID)
		}
		if item.TrainingDayName == nil || *item.TrainingDayName != "Snapshot "+want.name {
			t.Errorf("workout %d training day name = %v, want snapshot %q",
				item.ID, item.TrainingDayName, "Snapshot "+want.name)
		}
		expectedStartedAt := base.Add(time.Duration(3-index) * time.Hour)
		if !item.StartedAt.Equal(expectedStartedAt) ||
			!item.CompletedAt.Equal(expectedStartedAt.Add(100*time.Second)) {
			t.Errorf("workout %q timestamps = %s/%s",
				want.name, item.StartedAt, item.CompletedAt)
		}
		if item.DurationSeconds != 100 || item.CompletedSets != want.wantSets ||
			item.TotalWeight != want.wantWeight || item.ActualReps != want.wantActual ||
			item.PlannedReps != want.wantPlanned {
			t.Errorf("workout %q metrics = %+v", want.name, item)
		}
		if (item.CompletionPercent == nil) != (want.wantPercent == nil) {
			t.Errorf("workout %q completion percent = %v, want %v",
				want.name, item.CompletionPercent, want.wantPercent)
		} else if want.wantPercent != nil && *item.CompletionPercent != *want.wantPercent {
			t.Errorf("workout %q completion percent = %d, want %d",
				want.name, *item.CompletionPercent, *want.wantPercent)
		}
		if (item.FatigueScore == nil) != (want.fatigueScore == nil) {
			t.Errorf("workout %q fatigue score = %v, want %v",
				want.name, item.FatigueScore, want.fatigueScore)
		} else if want.fatigueScore != nil && *item.FatigueScore != *want.fatigueScore {
			t.Errorf("workout %q fatigue score = %d, want %d",
				want.name, *item.FatigueScore, *want.fatigueScore)
		}
	}

	page, err := repository.GetWorkoutHistory(ctx, ownerID, 2, 1)
	if err != nil {
		t.Fatalf("GetWorkoutHistory() with pagination error = %v", err)
	}
	if len(page) != 2 || page[0].ID != snapshotIDs[2] || page[1].ID != snapshotIDs[1] {
		t.Fatalf("unexpected paginated history IDs: %+v", page)
	}

	foreignHistory, err := repository.GetWorkoutHistory(ctx, otherUserID, 20, 0)
	if err != nil {
		t.Fatalf("GetWorkoutHistory() for other user error = %v", err)
	}
	if len(foreignHistory) != 1 {
		t.Fatalf("other user should only see their workout, got %d items", len(foreignHistory))
	}
	if foreignHistory[0].TrainingDayName == nil ||
		*foreignHistory[0].TrainingDayName != "Foreign snapshot" ||
		foreignHistory[0].FatigueScore == nil ||
		*foreignHistory[0].FatigueScore != 5 {
		t.Fatalf("other user received unexpected workout history: %+v", foreignHistory[0])
	}
}
