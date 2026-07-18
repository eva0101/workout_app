package core_middleware

import (
	"context"
	"net/http"
	"strings"
	"workout_app/internal/core/pkg/core_pkg_jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	jwtService core_pkg_jwt.JWTService
}

func NewAuthMiddleware(jwtService core_pkg_jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := m.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
