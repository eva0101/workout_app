package authorization_repository_postgres

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *AuthorizationRepository) RegisterUser(
	ctx context.Context,
	login string,
	password string,
) (core_domain.User, error) {
	var user core_domain.User

	err := r.pool.QueryRow(
		ctx,
		`
		INSERT INTO workoutapp.users (
			login,
			password_hash,
			created_at
		)
		VALUES ($1, $2, NOW())
		RETURNING id, login, created_at
		`,
		login,
		password,
	).Scan(
		&user.ID,
		&user.Login,
		&user.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return core_domain.User{}, core_errors.ErrLoginExists
		}

		return core_domain.User{}, err
	}

	return user, nil
}
