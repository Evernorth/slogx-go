package main

import (
	"context"
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"os"
)

const (
	AppLogLevel    = "APP_LOG_LEVEL"
	SystemLogLevel = "SYSTEM_LOG_LEVEL"
)

var (
	// This gets us a slog.Logger with context support that logs in JSON format to stdout.
	logger1, levelVar1 = slogx.NewLoggerBuilder().
				WithWriter(os.Stdout).
				WithFormat(slogx.FormatJSON).
				WithLevel(slog.LevelInfo).
				WithContextHandler().
				Build()

	logger2, levelVar2 = slogx.NewLoggerBuilder().
				WithWriter(os.Stdout).
				WithFormat(slogx.FormatJSON).
				WithLevel(slog.LevelInfo).
				WithContextHandler().
				Build()
)

func initializeLogLevel(logLevel string, levelVar *slog.LevelVar) {
	// Log the log level
	slog.Default().Info("", slog.String(logLevel, os.Getenv(logLevel)))

	// Set the log level
	err := slogx.GetLevelManager().ManageLevelFromEnv(levelVar, logLevel)
	if err != nil {
		panic(err)
	}

}

// Set the default level manager
func setup() {

	// Initialize the log levels
	initializeLogLevel(AppLogLevel, levelVar1)
	initializeLogLevel(SystemLogLevel, levelVar2)
	slogx.GetLevelManager().UpdateLevels()

	// Log that the logger has been initialized
	slog.InfoContext(context.TODO(), "Logger initialized", slog.String(AppLogLevel, levelVar1.Level().String()))
	slog.InfoContext(context.TODO(), "Logger initialized", slog.String(SystemLogLevel, levelVar2.Level().String()))

}

func main() {

	setup()

	// Log some messages

	logger1.DebugContext(context.TODO(), "Test debug AppLogLevel message.") // This will log out due to the AppLogLevel being set to Debug
	logger1.ErrorContext(context.TODO(), "Test error AppLogLevel message.")
	logger1.WarnContext(context.TODO(), "Test warn AppLogLevel message.")
	logger1.InfoContext(context.TODO(), "Test info AppLogLevel message.")

	logger2.DebugContext(context.TODO(), "Test debug SystemLogLevel message.") // This will not log out due to the SystemLogLevel being set to Info
	logger2.ErrorContext(context.TODO(), "Test error SystemLogLevel message.")
	logger2.WarnContext(context.TODO(), "Test warn SystemLogLevel message.")
	logger2.InfoContext(context.TODO(), "Test info SystemLogLevel message.")

}
