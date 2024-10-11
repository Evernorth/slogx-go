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
