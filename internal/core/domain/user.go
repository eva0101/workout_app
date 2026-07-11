package core_domain

import "time"

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
