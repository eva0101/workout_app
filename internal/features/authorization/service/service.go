package authorization_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_dto "workout_app/internal/core/dto"
	"workout_app/internal/core/pkg/core_pkg_jwt"
)

type AuthorizationService struct {
	authorizationRepository AuthorizationRepository
	jwtService              *core_pkg_jwt.JWTService
}

type AuthorizationRepository interface {
	RegisterUser(
		ctx context.Context,
		login string,
		password string,
	) (core_domain.User, error)

	GetPasswordHashByLogin(
		ctx context.Context,
		login string,
	) (core_dto.DTOUserPasswordHash, error)
}

func NewAuthorizationService(
	authorizationRepository AuthorizationRepository,
	jwtService *core_pkg_jwt.JWTService,
) *AuthorizationService {
	return &AuthorizationService{
		authorizationRepository: authorizationRepository,
		jwtService:              jwtService,
	}
}
