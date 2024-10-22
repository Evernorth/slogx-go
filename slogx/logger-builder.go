package slogx

import (
	"io"
	"log/slog"
	"os"
)

type Format int

const (
	FormatText Format = 0
	FormatJSON Format = 1
)

type LoggerBuilder interface {
	WithContextHandler() LoggerBuilder
	WithFormat(format Format) LoggerBuilder
	WithWriter(writer io.Writer) LoggerBuilder
	WithLevel(level slog.Level) LoggerBuilder
	WithLevelEnvVar(key string) LoggerBuilder
	Build() (*slog.Logger, *slog.LevelVar)
}

type defaultLoggerBuilder struct {
	writer            io.Writer
	format            Format
	level             slog.Level
	useContextHandler bool
	levelEnvVar       string
}

// NewLoggerBuilder creates a new LoggerBuilder with default values.  The default values are:  LevelInfo, FormatText,
// useContextHandler=false, levelEnvVar="" and writer=os.Stderr.
func NewLoggerBuilder() LoggerBuilder {
	return &defaultLoggerBuilder{
		level:             slog.LevelInfo,
		format:            FormatText,
		useContextHandler: false,
		levelEnvVar:       "",
		writer:            os.Stderr,
	}
}

// WithContextHandler enables the ContextHandler for the logger.
func (lb *defaultLoggerBuilder) WithContextHandler() LoggerBuilder {
	lb.useContextHandler = true
	return lb
}

// WithFormat sets the Format for the logger.
func (lb *defaultLoggerBuilder) WithFormat(format Format) LoggerBuilder {
	lb.format = format
	return lb
}

// WithWriter sets the io.Writer for the logger.
func (lb *defaultLoggerBuilder) WithWriter(writer io.Writer) LoggerBuilder {
	lb.writer = writer
	return lb
}

// WithLevel sets the slog.Level for the logger.
func (lb *defaultLoggerBuilder) WithLevel(level slog.Level) LoggerBuilder {
	lb.level = level
	return lb
}

// WithLevelEnvVar sets the environment variable key to use for the logger level.
func (lb *defaultLoggerBuilder) WithLevelEnvVar(key string) LoggerBuilder {
	lb.levelEnvVar = key
	return lb
}

// Build creates a new slog.Logger with the provided configuration. A slog.LevelVar to control the
// logger level is also returned.
func (lb *defaultLoggerBuilder) Build() (*slog.Logger, *slog.LevelVar) {

	// Set the default level
	levelVar := new(slog.LevelVar)
	levelVar.Set(lb.level)

	// If a level environment variable is provided, try to set the level from the environment
	if lb.levelEnvVar != "" {
		levelVar.Set(GetLevelFromEnv(lb.levelEnvVar, lb.level))
	}

	// Create the handler
	handlerOpts := &slog.HandlerOptions{
		Level: levelVar,
	}
	var handler slog.Handler
	if lb.format == FormatJSON {
		handler = slog.NewJSONHandler(lb.writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(lb.writer, handlerOpts)
	}

	// If the context handler is enabled, wrap the handler with a ContextHandler
	if lb.useContextHandler {
		handler = NewContextHandler(handler)
	}

	// Create the logger
	logger := slog.New(handler)

	slog.Default()

	return logger, levelVar
}
