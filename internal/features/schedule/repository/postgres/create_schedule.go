package schedule_repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *ScheduleRepository) CreateSchedule(
	ctx context.Context,
	userID uuid.UUID,
	trainingDayID int,
	dates []time.Time,
) ([]core_domain.ScheduledWorkout, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin schedule transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var ownedDayID int
	err = tx.QueryRow(ctx, `
		SELECT td.id
		FROM workoutapp.training_days td
		JOIN workoutapp.program p ON p.id = td.program_id
		WHERE td.id = $1 AND p.user_id = $2
		FOR UPDATE OF td;
	`, trainingDayID, userID).Scan(&ownedDayID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core_errors.ErrTrainingDayNotFound
		}
		r.log.Error("failed to verify training day ownership", zap.Error(err))
		return nil, fmt.Errorf("verify training day ownership: %w", err)
	}

	result := make([]core_domain.ScheduledWorkout, 0, len(dates))
	for _, date := range dates {
		var scheduled core_domain.ScheduledWorkout
		err = tx.QueryRow(ctx, `
			INSERT INTO workoutapp.scheduled_workouts (
				user_id, training_day_id, scheduled_date
			)
			VALUES ($1, $2, $3)
			RETURNING id, user_id, training_day_id, scheduled_date, created_at;
		`, userID, ownedDayID, date).Scan(
			&scheduled.ID,
			&scheduled.UserID,
			&scheduled.TrainingDayID,
			&scheduled.ScheduledDate,
			&scheduled.CreatedAt,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return nil, core_errors.ErrScheduleConflict
			}
			r.log.Error("failed to insert scheduled workout", zap.Error(err))
			return nil, fmt.Errorf("insert scheduled workout: %w", err)
		}
		result = append(result, scheduled)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit schedule transaction", zap.Error(err))
		return nil, fmt.Errorf("commit schedule transaction: %w", err)
	}
	return result, nil
}
