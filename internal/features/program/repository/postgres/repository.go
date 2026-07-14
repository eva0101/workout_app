package program_repository_postgres

import core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"

type ProgramRepository struct {
	pool *core_pool_pgx.Pool
}

func NewProgramRepository(pool *core_pool_pgx.Pool) *ProgramRepository {
	return &ProgramRepository{
		pool: pool,
	}
}
