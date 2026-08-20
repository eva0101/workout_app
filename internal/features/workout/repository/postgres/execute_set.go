package workout_repository_postgres

import (
	"context"
	"errors"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *WorkoutRepository) ExecuteSet(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayExerciseID int,
	repsDone *int,
	weight *float64,
) (core_domain.WorkoutSet, error) {
	var workoutSet core_domain.WorkoutSet

	query := `
	INSERT INTO workoutapp.workout_set (
		workout_id,
		training_day_exercises_id,
		set_number,
		reps_done,
		weight
	)
	SELECT
		w.id,
		tde.id,
		COALESCE(MAX(ws.set_number), 0) + 1,
		$3,
		$4
	FROM workoutapp.workout w
	JOIN workoutapp.training_day_exercises tde 
		ON tde.training_day_id = w.training_day_id
	LEFT JOIN workoutapp.workout_set ws
		ON ws.workout_id = w.id
		AND ws.training_day_exercises_id = tde.id
	WHERE w.user_id = $1
		AND w.status = 'in_progress'
		AND tde.id = $2
	GROUP BY w.id, tde.id
	RETURNING
		training_day_exercises_id,
		set_number,
		reps_done,
		weight;
	`
	err := r.pool.QueryRow(
		ctx,
		query,
		userID,
		trainingDayExerciseID,
		repsDone,
		weight,
	).Scan(
		&workoutSet.TrainingDayExercisesID,
		&workoutSet.SetNumber,
		&workoutSet.RepsDone,
		&workoutSet.Weight,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return core_domain.WorkoutSet{}, core_errors.ErrActiveWorkoutNotFound
		}

		return core_domain.WorkoutSet{}, err
	}

	return workoutSet, nil
}
