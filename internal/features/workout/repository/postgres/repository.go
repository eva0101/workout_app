package workout_repository_postgres

import core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"

type WorkoutRepository struct {
	pool *core_pool_pgx.Pool
}

func NewWorkoutRepository(pool *core_pool_pgx.Pool) *WorkoutRepository {

	return &WorkoutRepository{
		pool: pool,
	}
}
