package authorization_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
)

type AuthorizationService struct {
	authorizationRepository AuthorizationRepository
}

type AuthorizationRepository interface {
	RegisterUser(
		ctx context.Context,
		login string,
		password string,
	) (core_domain.User, error)
}

func NewAuthorizationService(
	authorizationRepository AuthorizationRepository,
) *AuthorizationService {
	return &AuthorizationService{
		authorizationRepository: authorizationRepository,
	}
}
