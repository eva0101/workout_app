package core_dto

import "github.com/google/uuid"

type DTOAuthorizationUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DTOUserPasswordHash struct {
	ID           uuid.UUID
	PasswordHash string
}

type DTOToken struct {
	Token string `json:"token"`
}
