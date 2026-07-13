package authorization_repository_postgres

import (
	"context"
	"errors"
	core_dto "workout_app/internal/core/dto"
	core_errors "workout_app/internal/core/errors"

	"github.com/jackc/pgx/v5"
)

func (r *AuthorizationRepository) GetPasswordHashByLogin(
	ctx context.Context,
	login string,
) (core_dto.DTOUserPasswordHash, error) {
	var hash core_dto.DTOUserPasswordHash

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT id, password_hash
		FROM workoutapp.users
		WHERE login = $1;
		`,
		login,
	).Scan(
		&hash.ID,
		&hash.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_dto.DTOUserPasswordHash{}, core_errors.ErrInvalidCredentials
		}

		return core_dto.DTOUserPasswordHash{}, err
	}

	return hash, nil
}
