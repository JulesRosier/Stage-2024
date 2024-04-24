package helper

import (
	"log/slog"
	"stage2024/pkg/settings"
	"strings"
)

func GetLogLevel(set settings.Logger) slog.Level {
	switch strings.ToUpper(set.Level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
