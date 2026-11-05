package workout_repository_postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type workoutExerciseSnapshotJSON struct {
	ID          int    `json:"id"`
	ExerciseID  int    `json:"exercise_id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Order       int    `json:"order"`
	PlannedSets []struct {
		ID          int      `json:"id"`
		SetNumber   int      `json:"set_number"`
		Repetitions *int     `json:"repetitions"`
		Weight      *float64 `json:"weight"`
	} `json:"planned_sets"`
	ActualSets []struct {
		ID        int      `json:"id"`
		SetNumber int      `json:"set_number"`
		RepsDone  *int     `json:"reps_done"`
		Weight    *float64 `json:"weight"`
	} `json:"actual_sets"`
}

func (r *WorkoutRepository) GetActiveWorkout(
	ctx context.Context,
	userID uuid.UUID,
) (core_domain.WorkoutDetails, error) {
	workout, err := r.getWorkout(ctx, userID, nil, true)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_domain.WorkoutDetails{}, core_errors.ErrActiveWorkoutNotFound
	}
	return workout, err
}

func (r *WorkoutRepository) GetWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID int,
) (core_domain.WorkoutDetails, error) {
	workout, err := r.getWorkout(ctx, userID, &workoutID, false)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_domain.WorkoutDetails{}, core_errors.ErrWorkoutNotFound
	}
	return workout, err
}

func (r *WorkoutRepository) getWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID *int,
	activeOnly bool,
) (core_domain.WorkoutDetails, error) {
	var workout core_domain.WorkoutDetails
	var exercisesJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT
			w.id,
			w.status,
			w.begin_at,
			w.completed_at,
			COALESCE(
				jsonb_agg(
					jsonb_build_object(
						'id', we.id,
						'exercise_id', we.exercise_id,
						'name', we.exercise_name,
						'title', we.exercise_title,
						'order', we.sort_order,
						'planned_sets', COALESCE((
							SELECT jsonb_agg(
								jsonb_build_object(
									'id', wes.id,
									'set_number', wes.set_number,
									'repetitions', wes.repetitions,
									'weight', wes.weight
								)
								ORDER BY wes.set_number
							)
							FROM workoutapp.workout_exercise_sets wes
							WHERE wes.workout_exercise_id = we.id
						), '[]'::jsonb),
						'actual_sets', COALESCE((
							SELECT jsonb_agg(
								jsonb_build_object(
									'id', ws.id,
									'set_number', ws.set_number,
									'reps_done', ws.reps_done,
									'weight', ws.weight
								)
								ORDER BY ws.set_number
							)
							FROM workoutapp.workout_set ws
							WHERE ws.workout_id = w.id
								AND ws.workout_exercise_id = we.id
						), '[]'::jsonb)
					)
					ORDER BY we.sort_order
				) FILTER (WHERE we.id IS NOT NULL),
				'[]'::jsonb
			)
		FROM workoutapp.workout w
		LEFT JOIN workoutapp.workout_exercise we
			ON we.workout_id = w.id
		WHERE w.user_id = $1
			AND ($2::INT IS NULL OR w.id = $2)
			AND (NOT $3 OR w.status = 'in_progress')
		GROUP BY w.id;
	`, userID, workoutID, activeOnly).Scan(
		&workout.ID,
		&workout.Status,
		&workout.StartedAt,
		&workout.CompletedAt,
		&exercisesJSON,
	)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			r.log.Error("failed to query workout details", zap.Error(err))
			return core_domain.WorkoutDetails{}, fmt.Errorf("query workout details: %w", err)
		}
		return core_domain.WorkoutDetails{}, err
	}

	workout.Exercises = make([]core_domain.WorkoutExerciseSnapshot, 0)
	var exercises []workoutExerciseSnapshotJSON
	if err := json.Unmarshal(exercisesJSON, &exercises); err != nil {
		r.log.Error("failed to decode workout exercises", zap.Error(err))
		return core_domain.WorkoutDetails{}, fmt.Errorf("decode workout exercises: %w", err)
	}

	for _, exercise := range exercises {
		item := core_domain.WorkoutExerciseSnapshot{
			ID:          exercise.ID,
			ExerciseID:  exercise.ExerciseID,
			Name:        exercise.Name,
			Title:       exercise.Title,
			Order:       exercise.Order,
			PlannedSets: make([]core_domain.WorkoutPlannedSet, 0, len(exercise.PlannedSets)),
			ActualSets:  make([]core_domain.WorkoutActualSet, 0, len(exercise.ActualSets)),
		}
		for _, set := range exercise.PlannedSets {
			item.PlannedSets = append(item.PlannedSets, core_domain.WorkoutPlannedSet{
				ID:          set.ID,
				SetNumber:   set.SetNumber,
				Repetitions: set.Repetitions,
				Weight:      set.Weight,
			})
		}
		for _, set := range exercise.ActualSets {
			item.ActualSets = append(item.ActualSets, core_domain.WorkoutActualSet{
				ID:        set.ID,
				SetNumber: set.SetNumber,
				RepsDone:  set.RepsDone,
				Weight:    set.Weight,
			})
		}
		workout.Exercises = append(workout.Exercises, item)
	}

	return workout, nil
}
