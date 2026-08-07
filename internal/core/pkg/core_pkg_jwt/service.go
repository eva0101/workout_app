package core_pkg_jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	secretKey string
}

type JWTService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (s *jwtService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(s.secretKey), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	userIDstring, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id missing")
	}

	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
