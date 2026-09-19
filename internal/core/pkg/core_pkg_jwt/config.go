package core_pkg_jwt

import (
	"os"
)

type Config struct {
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	return &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}
