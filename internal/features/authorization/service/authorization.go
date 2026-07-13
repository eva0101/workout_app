package authorization_service

import (
	"context"
	core_dto "workout_app/internal/core/dto"
	core_errors "workout_app/internal/core/errors"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthorizationService) AuthorizationUser(
	ctx context.Context,
	login string,
	password string,
) (core_dto.DTOToken, error) {
	user, err := s.authorizationRepository.GetPasswordHashByLogin(ctx, login)
	if err != nil {
		return core_dto.DTOToken{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return core_dto.DTOToken{}, core_errors.ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return core_dto.DTOToken{}, err
	}

	return core_dto.DTOToken{
		Token: token,
	}, nil
}
