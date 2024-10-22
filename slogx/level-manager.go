package slogx

import (
	"errors"
	"log/slog"
	"sync"
)

// LevelManager is an interface for managing slog.LevelVar objects from environment variables.  Call ManageLevelFromEnv to
// associate a slog.LevelVar with an environment variable key.  Call UpdateLevels to update the levels of all enrolled
// slog.LevelVar objects from their environment variables.
type LevelManager interface {
	ManageLevelFromEnv(levelVar *slog.LevelVar, key string) error
	UpdateLevels()
}

// defaultLevelManager is the default implementation of LevelManager.
type defaultLevelManager struct {
	levelVarMap *sync.Map
}

// defaultLevelManagerInstance is the singleton instance of defaultLevelManager.
var defaultLevelManagerInstance = &defaultLevelManager{
	levelVarMap: new(sync.Map),
}

// GetLevelManager returns the singleton instance of LevelManager.
func GetLevelManager() LevelManager {
	return defaultLevelManagerInstance
}

// ManageLevelFromEnv associates a slog.LevelVar with an environment variable key.  The level of the slog.LevelVar will be
// updated when UpdateLevels is called.
func (lm *defaultLevelManager) ManageLevelFromEnv(levelVar *slog.LevelVar, key string) error {
	if levelVar == nil {
		return errors.New("levelVar is required")
	}
	if key == "" {
		return errors.New("envVar is required")
	}
	lm.levelVarMap.Store(levelVar, key)
	return nil
}

// UpdateLevels updates the levels of all enrolled slog.LevelVar objects from their environment variables.
func (lm *defaultLevelManager) UpdateLevels() {
	lm.levelVarMap.Range(func(key, value interface{}) bool {
		var levelVar *slog.LevelVar
		var envVar string
		var ok bool

		// Cast the key and value
		levelVar, ok = key.(*slog.LevelVar)
		if !ok {
			panic("Could not cast key to *slog.LevelVar")
		}
		envVar, ok = value.(string)
		if !ok {
			panic("Could not cast value to string")
		}

		// Determine the level, defaulting to the current level if the environment variable is not set
		level := GetLevelFromEnv(envVar, levelVar.Level())

		// Update the level
		levelVar.Set(level)

		return true
	})
}
