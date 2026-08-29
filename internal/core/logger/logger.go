package core_logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	file *os.File
}

func NewLogger(config Config) (*Logger, error) {
	level := zap.NewAtomicLevel()

	if err := level.UnmarshalText([]byte(config.LogLevel)); err != nil {
		return nil, err
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	if err := os.MkdirAll(config.LogFolder, 0755); err != nil {
		return nil, err
	}

	logFilePath := filepath.Join(config.LogFolder, "logger")

	logFile, err := os.OpenFile(
		logFilePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	fileCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(logFile),
		level,
	)

	core := zapcore.NewTee(
		consoleCore,
		fileCore,
	)

	logger := zap.New(core)

	return &Logger{
		Logger: logger,
		file:   logFile,
	}, nil
}

func (l *Logger) Close() {
	if err := l.Sync(); err != nil {
		fmt.Println("failed to sync logger:", err)
	}
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close log file:", err)
	}
}
