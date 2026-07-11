package authorization_repository_postgres

import (
	core_pool_pgx "workout_app/internal/core/repository/postgres/pool/pgx"
)

type AuthorizationRepository struct {
	pool *core_pool_pgx.Pool
}

func NewAuthorizationRepository(pool *core_pool_pgx.Pool) *AuthorizationRepository {
	return &AuthorizationRepository{
		pool: pool,
	}
}
