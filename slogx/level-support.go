package slogx

import (
	"github.com/rotisserie/eris"
	"log/slog"
	"os"
	"strings"
)

// GetLevelByName returns a slog.Level object for the provided level name.
// If the level name is not valid, an error is returned
func GetLevelByName(levelName string) (*slog.Level, error) {
	if strings.EqualFold(levelName, slog.LevelDebug.String()) {
		levelDebug := slog.LevelDebug
		return &levelDebug, nil
	} else if strings.EqualFold(levelName, slog.LevelInfo.String()) {
		levelInfo := slog.LevelInfo
		return &levelInfo, nil
	} else if strings.EqualFold(levelName, slog.LevelWarn.String()) {
		levelWarn := slog.LevelWarn
		return &levelWarn, nil
	} else if strings.EqualFold(levelName, slog.LevelError.String()) {
		levelError := slog.LevelError
		return &levelError, nil
	} else {
		return nil, eris.New("invalid level name: " + levelName)
	}
}

// GetLevelFromEnv returns a slog.Level object for the provided environment variable key.
// If the environment variable is not set or the level name is not valid, the defaultLevel is returned.
func GetLevelFromEnv(key string, defaultLevel slog.Level) slog.Level {
	defaultLevelStr := defaultLevel.String()
	levelStr := getenv(key, &defaultLevelStr)
	level, err := GetLevelByName(levelStr)
	if err != nil {
		slog.Default().Warn("Environment variable has invalid logging level name.",
			slog.String("key", key),
			slog.String("level", levelStr))
		return defaultLevel
	}
	return *level
}

// getenv gets the environment variable with the provided key.  If the value is not empty, then it is returned.
// If it is empty, and a default value is provided, then the default value will be returned.  If the value is empty
// and no default value is provided, then a panic will occur.
func getenv(key string, defaultValue *string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue != nil {
			return *defaultValue
		} else {
			panic("required env variable " + key + " not found")
		}
	}
	return value
}
