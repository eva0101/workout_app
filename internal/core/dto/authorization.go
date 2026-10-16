package core_dto

import (
	"time"

	"github.com/google/uuid"
)

type DTOAuthorizationUserResponse struct {
	ID        uuid.UUID `json:"id"         example:"8f369456-2371-45b5-85b3-1bf4faae5c8d"`
	Login     string    `json:"login"      example:"egor"`
	CreatedAt time.Time `json:"created_at" example:"2026-09-23T11:12:44.513975+05:00"`
}

type DTOAuthorizationUserRequest struct {
	Login    string `json:"login"    example:"egor"`
	Password string `json:"password" example:"12345678"`
}

type DTOUserPasswordHash struct {
	ID           uuid.UUID
	PasswordHash string
}

type DTOToken struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTAyMzI5NzIsInVzZXJfaWQiOiI4ZjFWEJT578KTMDcxLTQ1YjUtODViMy0xYmY0ZmFhZTVjOGEifQ.faJA58348Qwr1z1G39aOgOywyxxNN1Cm4oXS5rSGSJ0"`
}
