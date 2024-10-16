package slogx

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestNewLoggerBuilder(t *testing.T) {
	builder := NewLoggerBuilder().(*defaultLoggerBuilder)
	assert.Equal(t, slog.LevelInfo, builder.level)
	assert.Equal(t, FormatText, builder.format)
	assert.False(t, builder.useContextHandler)
	assert.Equal(t, "", builder.levelEnvVar)
	assert.Equal(t, os.Stderr, builder.writer)
}

func TestWithContextHandler(t *testing.T) {
	builder := NewLoggerBuilder().WithContextHandler().(*defaultLoggerBuilder)
	assert.True(t, builder.useContextHandler)
}

func TestWithFormat(t *testing.T) {
	builder := NewLoggerBuilder().WithFormat(FormatJSON).(*defaultLoggerBuilder)
	assert.Equal(t, FormatJSON, builder.format)
}

func TestWithWriter(t *testing.T) {
	writer := os.Stdout
	builder := NewLoggerBuilder().WithWriter(writer).(*defaultLoggerBuilder)
	assert.Equal(t, writer, builder.writer)
}

func TestWithLevel(t *testing.T) {
	builder := NewLoggerBuilder().WithLevel(slog.LevelDebug).(*defaultLoggerBuilder)
	assert.Equal(t, slog.LevelDebug, builder.level)
}

func TestWithLevelEnvVar(t *testing.T) {
	envVar := "LOG_LEVEL"
	builder := NewLoggerBuilder().WithLevelEnvVar(envVar).(*defaultLoggerBuilder)
	assert.Equal(t, envVar, builder.levelEnvVar)
}

func TestBuild(t *testing.T) {
	builder := NewLoggerBuilder().
		WithWriter(os.Stdout).
		WithFormat(FormatJSON).
		WithLevel(slog.LevelDebug).
		WithContextHandler().
		WithLevelEnvVar("LOG_LEVEL")

	require.NoError(t, os.Setenv("LOG_LEVEL", "INFO"))

	logger, levelVar := builder.Build()
	assert.NotNil(t, logger)
	assert.NotNil(t, levelVar)
	assert.Equal(t, slog.LevelInfo, levelVar.Level())
}