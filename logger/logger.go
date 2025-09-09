package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	LogLvl string
}

func NewLoger(cnf Config) *slog.Logger {
	opt := &slog.HandlerOptions{}
	switch cnf.LogLvl {
	case "warn":
		opt.Level = slog.LevelWarn
	case "info":
		opt.Level = slog.LevelInfo
	case "debug":
		opt.Level = slog.LevelDebug
	case "error":
		opt.Level = slog.LevelError
	default:
		opt.Level = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opt))
	logger.With("component", "main")
	return logger
}
