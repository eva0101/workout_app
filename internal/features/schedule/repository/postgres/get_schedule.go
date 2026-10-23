package schedule_repository_postgres

import (
	"context"
	"fmt"
	"time"
	core_domain "workout_app/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *ScheduleRepository) GetSchedule(
	ctx context.Context,
	userID uuid.UUID,
	from time.Time,
	to time.Time,
) ([]core_domain.ScheduledWorkout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			sw.id,
			sw.training_day_id,
			sw.scheduled_date,
			td.day_number,
			p.name
		FROM workoutapp.scheduled_workouts sw
		JOIN workoutapp.training_days td ON td.id = sw.training_day_id
		JOIN workoutapp.program p ON p.id = td.program_id
		WHERE sw.user_id = $1
			AND sw.scheduled_date BETWEEN $2 AND $3
		ORDER BY sw.scheduled_date, td.day_number, sw.id;
	`, userID, from, to)
	if err != nil {
		r.log.Error("failed to query scheduled workouts", zap.Error(err))
		return nil, fmt.Errorf("query scheduled workouts: %w", err)
	}
	defer rows.Close()

	scheduled := make([]core_domain.ScheduledWorkout, 0)
	for rows.Next() {
		var item core_domain.ScheduledWorkout
		if err := rows.Scan(
			&item.ID,
			&item.TrainingDayID,
			&item.ScheduledDate,
			&item.DayNumber,
			&item.ProgramName,
		); err != nil {
			r.log.Error("failed to scan scheduled workout", zap.Error(err))
			return nil, fmt.Errorf("scan scheduled workout: %w", err)
		}
		scheduled = append(scheduled, item)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("failed to iterate scheduled workouts", zap.Error(err))
		return nil, fmt.Errorf("iterate scheduled workouts: %w", err)
	}
	return scheduled, nil
}
