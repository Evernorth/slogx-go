package main

import (
	"context"
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"os"
)

// This gets us a slog.Logger with context support that logs in JSON format to stdout.
var (
	logger, _ = slogx.NewLoggerBuilder().
		WithWriter(os.Stdout).
		WithFormat(slogx.FormatJSON).
		WithLevel(slog.LevelInfo).
		WithContextHandler().
		Build()
)

// main This example demonstrates how to create a logger with context support that logs in JSON format to stdout.
// It also demonstrates how to log with context and attributes.
// The logger is configured to log in JSON format to stdout with a default log level of INFO.
func main() {
	ctx := context.Background()
	logger.InfoContext(ctx, "Hello, World! Logging with Context.")

	// Create a context with some attributes
	ctx = slogx.ContextWithAttrs(ctx, slog.String("test1", "val1"),
		slog.String("test2", "val2"))
	logger.InfoContext(ctx, "Context and some attributes")

	// Create a context with some more attributes
	ctx = slogx.ContextWithAttrs(ctx, slog.String("test3", "val3"))
	logger.InfoContext(ctx, "Context and some more attributes")

	// Create a context with updated attributes
	ctx = slogx.ContextWithAttrs(ctx, slog.String("test1", "new-val1"))
	logger.InfoContext(ctx, "Context and update attributes")

}
