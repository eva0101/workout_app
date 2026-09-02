package workout_repository_postgres

import (
	core_logger "workout_app/internal/core/logger"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
)

type WorkoutRepository struct {
	pool *core_pool_pgx.Pool
	log  *core_logger.Logger
}

func NewWorkoutRepository(
	pool *core_pool_pgx.Pool,
	log *core_logger.Logger,
) *WorkoutRepository {
	return &WorkoutRepository{
		pool: pool,
		log:  log,
	}
}
