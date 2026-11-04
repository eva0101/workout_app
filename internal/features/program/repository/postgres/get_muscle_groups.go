package program_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (r *ProgramRepository) GetMuscleGroups(
	ctx context.Context,
) ([]core_domain.MuscleGroup, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT id, name, title, image_path
		FROM workoutapp.muscle_group
		ORDER BY id;
		`,
	)
	if err != nil {
		r.log.Error("failed to select muscle groups", zap.Error(err))
		return nil, fmt.Errorf("select muscle groups: %w", err)
	}
	defer rows.Close()

	groups := make([]core_domain.MuscleGroup, 0)
	for rows.Next() {
		var (
			group     core_domain.MuscleGroup
			imagePath pgtype.Text
		)
		if err := rows.Scan(&group.ID, &group.Name, &group.Title, &imagePath); err != nil {
			r.log.Error("failed to scan muscle group", zap.Error(err))
			return nil, fmt.Errorf("scan muscle group: %w", err)
		}
		if imagePath.Valid {
			group.ImagePath = &imagePath.String
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate muscle groups", zap.Error(err))
		return nil, fmt.Errorf("iterate muscle groups: %w", err)
	}

	return groups, nil
}
