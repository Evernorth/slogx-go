package slogx

import (
	"errors"
	"log/slog"
	"sync"
)

// LevelNameFunc is a function that returns the name of a level given a key.
type LevelNameFunc func(key string) string

// LevelManager is an interface for managing slog.LevelVar objects from environment variables.  Call ManageLevelFromEnv to
// associate a slog.LevelVar with an environment variable key.  Call UpdateLevels to update the levels of all enrolled
// slog.LevelVar objects from their environment variables.
type LevelManager interface {
	ManageLevelFromEnv(levelVar *slog.LevelVar, key string) error
	ManageLevelFromFunc(levelVar *slog.LevelVar, key string, levelNameFunc LevelNameFunc) error
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
	levelKey      string
	levelNameFunc LevelNameFunc
}

// ManageLevelFromEnv associates a slog.LevelVar with an environment variable key.  The level of the slog.LevelVar will be
// updated when UpdateLevels is called.
func (lm *defaultLevelManager) ManageLevelFromEnv(defaultLevelVar *slog.LevelVar, key string) error {
	err := lm.ManageLevelFromFunc(defaultLevelVar, key, getEnvLevelNameFunc())
	if err != nil {
		return err
	}
	return nil
}

func (lm *defaultLevelManager) ManageLevelFromFunc(defaultLevelVar *slog.LevelVar, key string, levelNameFunc LevelNameFunc) error {
	if defaultLevelVar == nil {
		return errors.New("defaultLevelVar is required")
	}
	if key == "" {
		return errors.New("key is required")
	}
	if levelNameFunc == nil {
		return errors.New("levelNameFunc is required")
	}

	lm.levelVarMap.Store(defaultLevelVar, levelFuncHolder{levelKey: key, levelNameFunc: levelNameFunc})
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

		// Determine the level, defaulting to the current level if the levelNameFunc does not return a level name
		level := GetLevelFromNameFunc(funcHolder.levelKey, funcHolder.levelNameFunc, defaultLevelVar.Level())

		// Update the level
		defaultLevelVar.Set(level)

		return true
	})
}
