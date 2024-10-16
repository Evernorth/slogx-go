package slogx

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"strings"
	"testing"
)

func TestContextLogger(t *testing.T) {

	buffer := bytes.NewBufferString("")
	logger, _ := NewLoggerBuilder().
		WithWriter(buffer).
		WithFormat(FormatJSON).
		WithLevel(slog.LevelInfo).
		WithContextHandler().
		Build()

	ctx := context.Background()
	newCtx := ContextWithAttrs(ctx,
		slog.String("test1", "val1"),
		slog.String("test2", "val2"))

	assert.NotEqual(t, ctx, newCtx)

	logger.InfoContext(newCtx, "test msg")

	assert.Equal(t, true, strings.Contains(buffer.String(), "test msg"))
	assert.Equal(t, true, strings.Contains(buffer.String(), "\"test1\":\"val1\",\"test2\":\"val2\""))
}

func TestContextLoggerWithMultipleUpdates(t *testing.T) {
	buffer := bytes.NewBufferString("")
	logger, _ := NewLoggerBuilder().
		WithWriter(buffer).
		WithFormat(FormatJSON).
		WithLevel(slog.LevelInfo).
		WithContextHandler().
		Build()

	ctx := context.Background()
	newCtx := ContextWithAttrs(ctx,
		slog.String("test1", "val1"),
		slog.String("test2", "val2"))

	// Update the context with additional attributes
	updatedCtx := ContextWithAttrs(newCtx,
		slog.String("test3", "val3"),
		slog.String("test4", "val4"))

	assert.NotEqual(t, ctx, newCtx)
	assert.NotEqual(t, newCtx, updatedCtx)

	logger.InfoContext(updatedCtx, "test msg")

	logOutput := buffer.String()
	assert.Contains(t, logOutput, "test msg")
	assert.Contains(t, logOutput, "\"test1\":\"val1\"")
	assert.Contains(t, logOutput, "\"test2\":\"val2\"")
	assert.Contains(t, logOutput, "\"test3\":\"val3\"")
	assert.Contains(t, logOutput, "\"test4\":\"val4\"")
}

func TestContextLoggerWithUpdate(t *testing.T) {
	buffer := bytes.NewBufferString("")
	logger, _ := NewLoggerBuilder().
		WithWriter(buffer).
		WithFormat(FormatJSON).
		WithLevel(slog.LevelInfo).
		WithContextHandler().
		Build()

	ctx := context.Background()
	newCtx := ContextWithAttrs(ctx,
		slog.String("test1", "val1"),
		slog.String("test2", "val2"))

	// Update the context with additional attributes
	updatedCtx := ContextWithAttrs(newCtx,
		slog.String("test1", "new-val1"))

	assert.NotEqual(t, ctx, newCtx)
	assert.NotEqual(t, newCtx, updatedCtx)

	logger.InfoContext(updatedCtx, "test msg")

	logOutput := buffer.String()
	assert.Contains(t, logOutput, "test msg")
	assert.Contains(t, logOutput, "\"test1\":\"new-val1\"")
	assert.Contains(t, logOutput, "\"test2\":\"val2\"")
}
