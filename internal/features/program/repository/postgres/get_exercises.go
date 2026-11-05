package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (r *ProgramRepository) GetExercises(
	ctx context.Context,
	muscleGroupID *int,
) ([]core_domain.CatalogExercise, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT id, name, title, muscle_group_id, description, image_path
		FROM workoutapp.exercise
		WHERE $1::INT IS NULL OR muscle_group_id = $1
		ORDER BY id;
		`,
		muscleGroupID,
	)
	if err != nil {
		r.log.Error("failed to select exercises", zap.Error(err))
		return nil, fmt.Errorf("select exercises: %w", err)
	}
	defer rows.Close()

	exercises := make([]core_domain.CatalogExercise, 0)
	for rows.Next() {
		var (
			exercise    core_domain.CatalogExercise
			description pgtype.Text
			imagePath   pgtype.Text
		)
		if err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Title,
			&exercise.MuscleGroupID,
			&description,
			&imagePath,
		); err != nil {
			r.log.Error("failed to scan exercise", zap.Error(err))
			return nil, fmt.Errorf("scan exercise: %w", err)
		}
		if description.Valid {
			exercise.Description = &description.String
		}
		if imagePath.Valid {
			exercise.ImagePath = &imagePath.String
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate exercises", zap.Error(err))
		return nil, fmt.Errorf("iterate exercises: %w", err)
	}

	return exercises, nil
}
