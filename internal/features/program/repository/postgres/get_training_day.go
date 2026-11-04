package program_repository_postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *ProgramRepository) GetTrainingDay(
	ctx context.Context,
	userID uuid.UUID,
	dayID int,
) (core_domain.TrainingDay, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var trainingDay core_domain.TrainingDay

	trainingDay.Exercises = make([]core_domain.TrainingExercise, 0)

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT
			td.id,
			td.day_number,
			td.name,
			tde.id,
			tde.exercise_id,
			tde.sets,
			COALESCE(
				(
					SELECT json_agg(
						json_build_object(
							'set_number', tes.set_number,
							'repetitions', tes.repetitions,
							'weight', tes.weight
						)
						ORDER BY tes.set_number
					)
					FROM workoutapp.training_day_exercise_sets tes
					WHERE tes.training_day_exercise_id = tde.id
				),
				'[]'::json
			)
		FROM workoutapp.training_days td
		JOIN workoutapp.program p
		ON td.program_id = p.id
		LEFT JOIN workoutapp.training_day_exercises tde
		ON td.id = tde.training_day_id
		WHERE td.id = $1
		AND p.user_id = $2;
		`,
		dayID,
		userID,
	)
	if err != nil {
		r.log.Error(
			"failed to select training day",
			zap.Error(err),
		)

		return core_domain.TrainingDay{}, fmt.Errorf("select training day: %w", err)
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		found = true

		var (
			exerciseTableID *int
			exerciseID      *int
			sets            *int
			plannedSetsJSON []byte
		)

		err := rows.Scan(
			&trainingDay.ID,
			&trainingDay.DayNumber,
			&trainingDay.Name,
			&exerciseTableID,
			&exerciseID,
			&sets,
			&plannedSetsJSON,
		)
		if err != nil {
			r.log.Error(
				"failed to scan training day",
				zap.Error(err),
			)

			return core_domain.TrainingDay{}, fmt.Errorf("scan training day: %w", err)
		}

		if exerciseID != nil {
			var plannedSets []struct {
				SetNumber   int      `json:"set_number"`
				Repetitions *int     `json:"repetitions"`
				Weight      *float64 `json:"weight"`
			}
			if err := json.Unmarshal(plannedSetsJSON, &plannedSets); err != nil {
				r.log.Error("failed to decode planned exercise sets", zap.Error(err))
				return core_domain.TrainingDay{}, fmt.Errorf("decode planned exercise sets: %w", err)
			}

			domainSets := make([]core_domain.PlannedSet, 0, len(plannedSets))
			for _, plannedSet := range plannedSets {
				domainSets = append(domainSets, core_domain.PlannedSet{
					SetNumber:   plannedSet.SetNumber,
					Repetitions: plannedSet.Repetitions,
					Weight:      plannedSet.Weight,
				})
			}

			trainingDay.Exercises = append(trainingDay.Exercises, core_domain.TrainingExercise{
				ID:          *exerciseTableID,
				ExerciseID:  *exerciseID,
				Sets:        *sets,
				PlannedSets: domainSets,
			})
		}
	}

	if err := rows.Err(); err != nil {
		r.log.Error(
			"failed to iterate rows",
			zap.Error(err),
		)

		return core_domain.TrainingDay{}, fmt.Errorf("next rows: %w", err)
	}

	if !found {
		return core_domain.TrainingDay{}, core_errors.ErrTrainingDayNotFound
	}

	return trainingDay, nil
}
