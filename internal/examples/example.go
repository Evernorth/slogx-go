package main

import (
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"os"
)

const (
	logger1LevelEnvVar = "LOGGER1_LOG_LEVEL"
	logger2LevelEnvVar = "LOGGER2_LOG_LEVEL"
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

func manageLevelFromEnv(logLevel string, levelVar *slog.LevelVar) {
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

	// Enroll the levelVars with the LevelManager
	manageLevelFromEnv(logger1LevelEnvVar, levelVar1)
	manageLevelFromEnv(logger2LevelEnvVar, levelVar2)

	// Tell the LevelManager to update the levels
	slogx.GetLevelManager().UpdateLevels()

	// Log that the logger has been initialized
	slog.Info("logger1 initialized", slog.String(logger1LevelEnvVar, levelVar1.Level().String()))
	slog.Info("logger2 initialized", slog.String(logger2LevelEnvVar, levelVar2.Level().String()))

}

func main() {

	setup()

	// Log some test messages
	logger1.Debug("logger1 debug message.")
	logger1.Info("logger1 info message.")
	logger1.Warn("logger1 warn message.")
	logger1.Error("logger1 error message.")

	logger2.Debug("logger2 debug message.")
	logger2.Info("logger2 info message.")
	logger2.Warn("logger2 warn message.")
	logger2.Error("logger2 error message.")

}
