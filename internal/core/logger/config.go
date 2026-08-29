package core_logger

import "os"

type Config struct {
	LogLevel  string
	LogFolder string
}

func NewConfig() Config {
	return Config{
		LogLevel:  os.Getenv("LOG_LEVEL"),
		LogFolder: os.Getenv("LOG_FOLDER"),
	}
}
