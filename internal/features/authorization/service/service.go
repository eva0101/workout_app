package authorization_service

import (
	"context"
	core_domain "workout_app/internal/core/domain"
	core_dto "workout_app/internal/core/dto"
	core_logger "workout_app/internal/core/logger"
	"workout_app/internal/core/pkg/core_pkg_jwt"
)

type AuthorizationService struct {
	authorizationRepository AuthorizationRepository
	jwtService              core_pkg_jwt.JWTService
	log                     *core_logger.Logger
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
	jwtService core_pkg_jwt.JWTService,
	log *core_logger.Logger,
) *AuthorizationService {
	return &AuthorizationService{
		authorizationRepository: authorizationRepository,
		jwtService:              jwtService,
		log:                     log,
	}
}
