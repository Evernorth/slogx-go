package slogx

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"log/slog"
	"os"
	"testing"
)

func TestGetLevelFromEnv(t *testing.T) {
	testLevelDebug := "TEST_LEVEL_DEBUG"
	testLevelInfo := "TEST_LEVEL_INFO"
	testLevelWarn := "TEST_LEVEL_WARN"
	testLevelError := "TEST_LEVEL_ERROR"

	require.NoError(t, os.Setenv(testLevelDebug, slog.LevelDebug.String()))
	require.NoError(t, os.Setenv(testLevelInfo, slog.LevelInfo.String()))
	require.NoError(t, os.Setenv(testLevelWarn, slog.LevelWarn.String()))
	require.NoError(t, os.Setenv(testLevelError, slog.LevelError.String()))

	var actualLevel slog.Level
	actualLevel = GetLevelFromEnv(testLevelDebug, slog.LevelInfo)
	assert.Equal(t, slog.LevelDebug, actualLevel)
	actualLevel = GetLevelFromEnv(testLevelInfo, slog.LevelWarn)
	assert.Equal(t, slog.LevelInfo, actualLevel)
	actualLevel = GetLevelFromEnv(testLevelWarn, slog.LevelInfo)
	assert.Equal(t, slog.LevelWarn, actualLevel)
	actualLevel = GetLevelFromEnv(testLevelError, slog.LevelInfo)
	assert.Equal(t, slog.LevelError, actualLevel)
	actualLevel = GetLevelFromEnv("INVALID_ENV_VAR", slog.LevelInfo)
	assert.Equal(t, slog.LevelInfo, actualLevel)
}
