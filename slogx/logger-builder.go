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
	WithLevelNameFunc(key string, levelNameFunc LevelNameFunc) LoggerBuilder
	Build() (*slog.Logger, *slog.LevelVar)
}

type defaultLoggerBuilder struct {
	writer            io.Writer
	format            Format
	level             slog.Level
	useContextHandler bool
	levelKey          string
	levelNameFunc     LevelNameFunc
}

// NewLoggerBuilder creates a new LoggerBuilder with default values.  The default values are:  LevelInfo, FormatText,
// useContextHandler=false, levelKey="", levelNameFunc=nil and writer=os.Stderr.
func NewLoggerBuilder() LoggerBuilder {
	return &defaultLoggerBuilder{
		level:             slog.LevelInfo,
		format:            FormatText,
		useContextHandler: false,
		levelKey:          "",
		levelNameFunc:     nil,
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

// WithLevelEnvVar sets the environment variable key to use to get the level.
func (lb *defaultLoggerBuilder) WithLevelEnvVar(key string) LoggerBuilder {
	lb.levelKey = key
	return lb
}

// WithLevelNameFunc sets the level variable key and function to use to get the level name.
func (lb *defaultLoggerBuilder) WithLevelNameFunc(key string, levelNameFunc LevelNameFunc) LoggerBuilder {
	lb.levelKey = key
	lb.levelNameFunc = levelNameFunc
	return lb
}

// Build creates a new slog.Logger with the provided configuration. A slog.LevelVar to control the
// logger level is also returned.
func (lb *defaultLoggerBuilder) Build() (*slog.Logger, *slog.LevelVar) {

	// Set the default level
	levelVar := new(slog.LevelVar)
	levelVar.Set(lb.level)

	// If a level variable key is provided and a level name function is provided, set the level from the level name function
	// Otherwise, try to set the level from the environment
	if lb.levelKey != "" {
		if lb.levelNameFunc != nil {
			levelVar.Set(GetLevelFromNameFunc(lb.levelKey, lb.levelNameFunc, lb.level))
		} else {
			levelVar.Set(GetLevelFromNameFunc(lb.levelKey, getEnvLevelNameFunc(), lb.level))
		}
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
