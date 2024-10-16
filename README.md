# slogx-go
  
  [![Go Report Card](https://goreportcard.com/badge/github.com/Evernorth/slogx)](https://goreportcard.com/report/github.com/Evernorth/slogx)
  [![GoDoc](https://godoc.org/github.com/Evernorth/slogx?status.svg)](https://godoc.org/github.com/Evernorth/slogx)
  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
  [![Release](https://img.shields.io/github/v/release/Evernorth/slogx)](https://github.com/Evernorth/slogx-go/releases)

## Description

A collection of `slog` extensions.
* `ContextHandler` allows you to add `slog` attributes (`slog.Attr` instances) to a `context.Context`.  These attributes are added to log records when the `*Context` function variants (`InfoContext`, `ErrorContext`, etc) on the logger are used.
* `LoggerBuilder` provides a simple way to build `slog.Logger` instances.
* `LevelManager` provides a way to manage `slog.LevelVar` instances from environment variables.
* Multiple loggers can be created with different log levels and formats. See [internal/examples](internal/examples) for more examples.

## Installation

``` go get -u github.com/Evernorth/slogx-go ```

## Usage

### Context-aware logging
Create a context-aware `slog.Logger`, then use the logger making sure to use a _`Context`_ variant function.

```go
package main

import (
	"context"
	"github.com/Evernorth/slogx-go/slogx"
	"log/slog"
	"os"
)

// Environment variables
// These can be set to change the log level at runtime.
// The log level can be set to one of the following values: DEBUG, INFO, WARN, ERROR, FATAL, PANIC
const (
	logger1LevelEnvVar = "LOGGER1_LOG_LEVEL"
)

// Loggers
// This gets us a slog.Logger with context support that logs in JSON format to stdout.
var (
	logger1, levelVar1 = slogx.NewLoggerBuilder().
				WithWriter(os.Stdout).
				WithFormat(slogx.FormatJSON).
				WithLevel(slog.LevelInfo).
				WithContextHandler().
				Build()


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
func setup(ctx context.Context) {

	// Enroll the levelVars with the LevelManager
	manageLevelFromEnv(logger1LevelEnvVar, levelVar1)

	// Tell the LevelManager to update the levels
	slogx.GetLevelManager().UpdateLevels()

	// Log that the logger has been initialized
	slog.InfoContext(ctx, "Logger initialized", slog.String(logger1LevelEnvVar, levelVar1.Level().String()))

}

// main This function demonstrates how to use the slogx package to create a logger with context support.
// It also demonstrates how to manage the log level from an environment variable.
// The logger is configured to log in JSON format to stdout with a default log level of INFO.
func main() {
	ctx := context.Background()
	setup(ctx)

	// Log some test messages
	logger1.Debug("logger1 debug message.")
	logger1.Info("logger1 info message.")
	logger1.Warn("logger1 warn message.")
	logger1.Error("logger1 error message.")

	// Log some test messages with context
	logger1.DebugContext(ctx, "logger1 debug message.")
	logger1.InfoContext(ctx, "logger1 info message.")
	logger1.WarnContext(ctx, "logger1 warn message.")
	logger1.ErrorContext(ctx, "logger1 error message.")
  

}

```

## Dependencies
See the [go.mod](go.mod) file.

## Support
If you have questions, concerns, bug reports, etc. See [CONTRIBUTING](CONTRIBUTING.md).

## License
slogx is open source software released under the [Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0.html).

## Original Contributors
- Steve Sefton, Evernorth
- Shellee Stewart, Evernorth
- Neil Powell, Evernorth