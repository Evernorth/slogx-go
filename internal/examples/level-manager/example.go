package main

import (
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"os"
)

// Environment variables
// These can be set to change the log level at runtime.
// The log level can be set to one of the following values: DEBUG, INFO, WARN, ERROR, FATAL, PANIC
const (
	logger1LevelEnvVar = "LOGGER1_LOG_LEVEL"
	logger2LevelEnvVar = "LOGGER2_LOG_LEVEL"
)

// Loggers
// This gets us a slog.Logger with context support that logs in JSON format to stdout.
var (
	logger1, levelVar1 = slogx.NewLoggerBuilder().
				WithWriter(os.Stdout).
				WithFormat(slogx.FormatJSON).
				WithLevel(slog.LevelInfo).
				Build()

	logger2, levelVar2 = slogx.NewLoggerBuilder().
				WithWriter(os.Stdout).
				WithFormat(slogx.FormatJSON).
				WithLevel(slog.LevelDebug).
				Build()
)

// manageLevelFromEnv Manage the log level from an environment variable
func manageLevelFromEnv(logLevel string, levelVar *slog.LevelVar) {
	// Log the log level
	slog.Default().Info("", slog.String(logLevel, os.Getenv(logLevel)))

	// Set the log level
	err := slogx.GetLevelManager().ManageLevelFromEnv(levelVar, logLevel)
	if err != nil {
		panic(err)
	}

}

// setup Set the default level manager
func setup() {

	// Enroll the levelVars with the LevelManager
	manageLevelFromEnv(logger1LevelEnvVar, levelVar1)
	manageLevelFromEnv(logger2LevelEnvVar, levelVar2)

	// Tell the LevelManager to update the levels
	slogx.GetLevelManager().UpdateLevels()

	// Log that the logger has been initialized

	slog.Info("Logger initialized", slog.String(logger1LevelEnvVar, levelVar1.Level().String()))
	slog.Info("Logger initialized", slog.String(logger2LevelEnvVar, levelVar2.Level().String()))

}

// main This function demonstrates how to use the slogx package to create a logger to manage log levels.
// It also demonstrates how to manage the log level from an environment variable.
// It also demonstrates how to use multiple loggers with log levels.
// The logger is configured to log in JSON format to stdout with a default log level of INFO.
func main() {

	setup()

	// Log some test messages

	// This message will not be logged because the log level is INFO.
	// This can be changed by setting the environment variable LOGGER1_LOG_LEVEL to DEBUG
	logger1.Debug("logger1 debug message.")

	logger1.Info("logger1 info message.")
	logger1.Warn("logger1 warn message.")
	logger1.Error("logger1 error message.")

	// Log some test messages using the 2nd logger

	// This message will be logged because the log level is DEBUG.
	// This can be changed by setting the environment variable LOGGER2_LOG_LEVEL to INFO
	logger2.Debug("logger2 debug message.")
	logger2.Info("logger2 info message.")
	logger2.Error("logger2 error message.")

}
