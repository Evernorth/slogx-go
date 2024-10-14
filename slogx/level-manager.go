package slogx

import (
	"github.com/rotisserie/eris"
	"log/slog"
	"sync"
)

type LevelManager interface {
	ManageLevelFromEnv(levelVar *slog.LevelVar, key string) error
	UpdateLevels()
}

type defaultLevelManager struct {
	levelVarMap *sync.Map
}

var defaultLevelManagerInstance = &defaultLevelManager{
	levelVarMap: new(sync.Map),
}

func GetLevelManager() LevelManager {
	return defaultLevelManagerInstance
}

func (lm *defaultLevelManager) ManageLevelFromEnv(levelVar *slog.LevelVar, key string) error {
	if levelVar == nil {
		return eris.New("levelVar is required")
	}
	if key == "" {
		return eris.New("envVar is required")
	}
	lm.levelVarMap.Store(levelVar, key)
	return nil
}

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
