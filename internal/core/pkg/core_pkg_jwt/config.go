package core_pkg_jwt

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}
