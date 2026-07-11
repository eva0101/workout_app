package authorization_service

import (
	"context"
	"strings"
	"unicode/utf8"
	core_domain "workout_app/internal/core/domain"
	core_errors "workout_app/internal/core/errors"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthorizationService) RegisterUser(
	ctx context.Context,
	login string,
	password string,
) (core_domain.User, error) {
	login = strings.TrimSpace(login)
	loginLen := utf8.RuneCountInString(login)

	if loginLen < 3 {
		return core_domain.User{}, core_errors.ErrLoginFewSymbols
	}
	if loginLen > 32 {
		return core_domain.User{}, core_errors.ErrLoginManySymbols
	}

	passwordLen := utf8.RuneCountInString(password)

	if passwordLen < 8 {
		return core_domain.User{}, core_errors.ErrPasswordFewSymbols
	}
	if passwordLen > 64 {
		return core_domain.User{}, core_errors.ErrPasswordManySymbols
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return core_domain.User{}, err
	}

	user, err := s.authorizationRepository.RegisterUser(
		ctx,
		login,
		string(hash),
	)
	if err != nil {
		return core_domain.User{}, err
	}

	return user, nil
}
