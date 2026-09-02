package authorization_repository_postgres

import (
	core_logger "workout_app/internal/core/logger"
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
)

type AuthorizationRepository struct {
	pool *core_pool_pgx.Pool
	log  *core_logger.Logger
}

func NewAuthorizationRepository(
	pool *core_pool_pgx.Pool,
	log *core_logger.Logger,
) *AuthorizationRepository {
	return &AuthorizationRepository{
		pool: pool,
		log:  log,
	}
}
