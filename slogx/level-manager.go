package slogx

import (
	"errors"
	"log/slog"
	"os"
	"sync"
)

// LevelFunc is a function that returns the name of a level given a key.
type LevelFunc func(key string) string

// LevelManager is an interface for managing slog.LevelVar objects from environment variables.  Call ManageLevelFromEnv to
// associate a slog.LevelVar with an environment variable key.  Call UpdateLevels to update the levels of all enrolled
// slog.LevelVar objects from their environment variables.
type LevelManager interface {
	ManageLevelFromEnv(levelVar *slog.LevelVar, key string) error
	ManageLevelFromFunc(levelVar *slog.LevelVar, key string, levelFunc LevelFunc) error
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

type levelFuncHolder struct {
	levelKey  string
	levelFunc LevelFunc
}

// ManageLevelFromEnv associates a slog.LevelVar with an environment variable key.  The level of the slog.LevelVar will be
// updated when UpdateLevels is called.
func (lm *defaultLevelManager) ManageLevelFromEnv(defaultLevelVar *slog.LevelVar, key string) error {
	err := lm.ManageLevelFromFunc(defaultLevelVar, key, getEnvLevelFunc())
	if err != nil {
		return err
	}
	return nil
}

// ManageLevelFromFunc associates a slog.LevelVar with a key and a LevelFunc.  The level of the slog.LevelVar will be
// updated when UpdateLevels is called.
// A LevelFunc is useful for getting a level name using alternate sources, such as koanf, viper, etc.
func (lm *defaultLevelManager) ManageLevelFromFunc(defaultLevelVar *slog.LevelVar, key string, levelFunc LevelFunc) error {
	if defaultLevelVar == nil {
		return errors.New("defaultLevelVar is required")
	}
	if key == "" {
		return errors.New("key is required")
	}
	if levelFunc == nil {
		return errors.New("levelFunc is required")
	}

	lm.levelVarMap.Store(defaultLevelVar, levelFuncHolder{levelKey: key, levelFunc: levelFunc})
	return nil
}

// UpdateLevels updates the levels of all enrolled slog.LevelVar objects from their environment variables.
func (lm *defaultLevelManager) UpdateLevels() {
	lm.levelVarMap.Range(func(key, value interface{}) bool {
		var defaultLevelVar *slog.LevelVar
		var funcHolder levelFuncHolder
		var ok bool

		// Cast the key and value
		defaultLevelVar, ok = key.(*slog.LevelVar)
		if !ok {
			panic("Could not cast key to *slog.LevelVar")
		}
		funcHolder, ok = value.(levelFuncHolder)
		if !ok {
			panic("Could not cast value to levelFuncHolder")
		}

		// Determine the level, defaulting to the current level if the levelFunc does not return a level name
		level := GetLevelFromFunc(funcHolder.levelKey, funcHolder.levelFunc, defaultLevelVar.Level())

		// Update the level
		defaultLevelVar.Set(level)

		return true
	})
}

// getEnvLevelFunc gets the environment variable with the provided level name key.
func getEnvLevelFunc() LevelFunc {
	return func(key string) string {
		return os.Getenv(key)
	}
}
