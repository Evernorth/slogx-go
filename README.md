# slogx-go
  
  [![Go Report Card](https://goreportcard.com/badge/github.com/Evernorth/slogx)](https://goreportcard.com/report/github.com/Evernorth/slogx)
  [![GoDoc](https://godoc.org/github.com/Evernorth/slogx?status.svg)](https://godoc.org/github.com/Evernorth/slogx)
  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
  [![Release](https://img.shields.io/github/v/release/Evernorth/slogx-go)](https://github.com/Evernorth/slogx-go/releases)

## Description

The `slogx` package is a collection of `slog` extensions. The [`slog`](https://pkg.go.dev/log/slog)  is a Go standard library package for a structured logger.

## Features
The `slogx` package provides the following features:
* `ContextHandler` allows you to add `slog` attributes (`slog.Attr` instances) to a `context.Context`.  These attributes are added to log records when the `*Context` function variants (`InfoContext`, `ErrorContext`, etc) on the logger are used.
* `LoggerBuilder` provides a simple way to build `slog.Logger` instances.
* `LevelManager` provides a way to manage `slog.LevelVar` instances from environment variables.
* Multiple loggers can be created with different log levels and formats. See [internal/examples](internal/examples) for more examples.

## Installation

``` go get -u github.com/Evernorth/slogx-go ```

## Usage

### Simple logging
The following example demonstrates how to create a logger and log a message.
```go
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
```
#### Simple log example
```text
{"time":"2024-10-21T12:03:41.103566-04:00","level":"INFO","msg":"Hello, World!"}
```
.
### Managing log levels
The following example demonstrates how to create a logger with a log level that can be changed at runtime.
```go
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
	logger2.Warn("logger2 warn message.")
	logger2.Error("logger2 error message.")

}

```

#### Log level example
```text
2024/10/21 12:05:40 INFO Logger initialized LOGGER1_LOG_LEVEL=INFO
2024/10/21 12:05:40 INFO Logger initialized LOGGER2_LOG_LEVEL=DEBUG
{"time":"2024-10-21T12:05:40.937543-04:00","level":"INFO","msg":"logger1 info message."}
{"time":"2024-10-21T12:05:40.937545-04:00","level":"WARN","msg":"logger1 warn message."}
{"time":"2024-10-21T12:05:40.937547-04:00","level":"ERROR","msg":"logger1 error message."}
{"time":"2024-10-21T12:05:40.937548-04:00","level":"DEBUG","msg":"logger2 debug message."}
{"time":"2024-10-21T12:05:40.93755-04:00","level":"INFO","msg":"logger2 info message."}
{"time":"2024-10-21T12:05:40.937551-04:00","level":"ERROR","msg":"logger2 error message."}

```

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
```
#### Context-aware log example
```text
{"time":"2024-10-21T12:09:44.301872-04:00","level":"INFO","msg":"Hello, World! Logging with Context."}
{"time":"2024-10-21T12:09:44.302091-04:00","level":"INFO","msg":"Context and some attributes","test1":"val1","test2":"val2"}
{"time":"2024-10-21T12:09:44.302095-04:00","level":"INFO","msg":"Context and some more attributes","test1":"val1","test2":"val2","test3":"val3"}
{"time":"2024-10-21T12:09:44.302098-04:00","level":"INFO","msg":"Context and update attributes","test1":"new-val1","test2":"val2","test3":"val3"}
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