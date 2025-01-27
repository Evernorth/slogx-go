package slogx

import (
	"errors"
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
		return nil, errors.New("invalid level name: " + levelName)
	}
}

// GetLevelFromEnv returns a slog.Level object for the provided environment variable key.
// If the environment variable is not set or the level name is not valid, the defaultLevel is returned.
// ToDo: Why is this public and how to deprecate?
func GetLevelFromEnv(key string, defaultLevel slog.Level) slog.Level {
	return GetLevelFromNameFunc(key, getEnvLevelNameFunc(), defaultLevel)
}

// GetLevelFromNameFunc returns a slog.Level object for the provided LevelNameFunc.
// If a key is not set or the level name is not valid, the defaultLevel is returned.
func GetLevelFromNameFunc(levelKey string, levelNameFunc LevelNameFunc, defaultLevel slog.Level) slog.Level {
	defaultLevelStr := defaultLevel.String()
	levelName := levelNameFunc(levelKey)
	if levelName == "" {
		if defaultLevelStr != "" {
			return defaultLevel
		} else {
			panic("required Level name variable " + levelKey + " not found")
		}
	}
	level, err := GetLevelByName(levelName)
	if err != nil {
		slog.Default().Warn("Key is an invalid logging level name.",
			slog.String("key", levelKey),
			slog.String("level", levelName))
		return defaultLevel
	}
	return *level
}

// getEnvLevelNameFunc gets the environment variable with the provided level name key.
func getEnvLevelNameFunc() LevelNameFunc {
	return func(key string) string {
		return os.Getenv(key)
	}
}
