package main

import (
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
		Build()
)

func main() {
	logger.Info("Hello, World!")
}
