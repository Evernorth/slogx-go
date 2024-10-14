package slogx

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestLevelManager(t *testing.T) {

	levelManager := GetLevelManager()

	levelVar1 := slog.LevelVar{}
	levelVar1.Set(slog.LevelInfo)
	levelVar2 := slog.LevelVar{}
	levelVar2.Set(slog.LevelInfo)

	levelVar1Key := "LEVEL_VAR_1_LEVEL"
	levelVar2Key := "LEVEL_VAR_2_LEVEL"

	var err error
	err = levelManager.ManageLevelFromEnv(&levelVar1, levelVar1Key)
	assert.NoError(t, err)
	err = levelManager.ManageLevelFromEnv(&levelVar2, levelVar2Key)
	assert.NoError(t, err)

	levelManager.UpdateLevels()
	assert.Equal(t, slog.LevelInfo, levelVar1.Level())
	assert.Equal(t, slog.LevelInfo, levelVar2.Level())

	// Test that the level is updated when the environment variable is set
	require.NoError(t, os.Setenv(levelVar1Key, slog.LevelDebug.String()))
	require.NoError(t, os.Setenv(levelVar2Key, slog.LevelInfo.String()))

	levelManager.UpdateLevels()
	assert.Equal(t, slog.LevelDebug, levelVar1.Level())
	assert.Equal(t, slog.LevelInfo, levelVar2.Level())

	require.NoError(t, os.Setenv(levelVar2Key, slog.LevelDebug.String()))
	levelManager.UpdateLevels()
	assert.Equal(t, slog.LevelDebug, levelVar1.Level())
	assert.Equal(t, slog.LevelDebug, levelVar2.Level())
}

func TestLevelManagerSingleton(t *testing.T) {
	assert.Equal(t, GetLevelManager(), GetLevelManager())
}

func TestLevelManagerErrors(t *testing.T) {
	levelManager := GetLevelManager()

	var err error
	err = levelManager.ManageLevelFromEnv(nil, "LEVEL_VAR_1_LEVEL")
	assert.Error(t, err)
	err = nil
	err = levelManager.ManageLevelFromEnv(&slog.LevelVar{}, "")
	assert.Error(t, err)
}
